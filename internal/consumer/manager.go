package consumer

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/segmentio/kafka-go"
	"github.com/qnhqn1/file-validator/config"
	"github.com/qnhqn1/file-validator/internal/metrics"
	"github.com/qnhqn1/file-validator/internal/producer"
	"github.com/qnhqn1/file-validator/internal/services/validator"
)


type Manager struct {
	reader    *kafka.Reader
	producers *producer.Manager
	svc       validator.Service
	cfg       *config.Config
	collector *metrics.Collector
}


func New(cfg *config.Config, svc validator.Service, producers *producer.Manager, collector *metrics.Collector) *Manager {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        cfg.Kafka.Brokers,
		GroupID:        cfg.Kafka.GroupID,
		Topic:          cfg.Topics.Input,
		MinBytes:       1,
		MaxBytes:       10e6,
		CommitInterval: 0,
	})
	return &Manager{reader: r, producers: producers, svc: svc, cfg: cfg, collector: collector}
}




func (m *Manager) Run(ctx context.Context) error {
	log.Printf("file-validator: консьюмер работает для топика %s", m.cfg.Topics.Input)
	for {
		msg, err := m.reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				log.Printf("file-validator: консьюмер выключен")
				return nil
			}
			log.Printf("file-validator: ошибка получения: %v", err)
			return err
		}

		m.collector.RecordReceived(ctx)


		var ev map[string]interface{}
		if err := json.Unmarshal(msg.Value, &ev); err != nil {
			log.Printf("file-validator: недопустимый payload: %v", err)
			m.collector.RecordError(ctx, metrics.CategoryInvalidFile)

			_ = m.reader.CommitMessages(ctx, msg)
			continue
		}


		reqID, _ := ev["request_id"].(string)
		objName, _ := ev["object_name"].(string)
		docID, _ := ev["document_id"].(string)
		if objName == "" {
			log.Printf("file-validator: отсутствует object_name в payload: %v", ev)
			m.collector.RecordError(ctx, metrics.CategoryInvalidFile)


			if reqID != "" {
				resp := map[string]interface{}{"request_id": reqID, "status": "invalid", "error": "missing_object_name"}
				b, _ := json.Marshal(resp)
				_ = m.producers.SendValidated(ctx, msg.Key, b)
			}
			_ = m.reader.CommitMessages(ctx, msg)
			continue
		}


		u, err := url.Parse(m.cfg.Minio.Endpoint)
		if err != nil {
			log.Printf("file-validator: недопустимый endpoint minio: %v", err)
			m.collector.RecordError(ctx, metrics.CategoryUnknown)
			_ = m.reader.CommitMessages(ctx, msg)
			continue
		}

		u.Path = path.Join(u.Path, m.cfg.Minio.Bucket, objName)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
		if err != nil {
			log.Printf("file-validator: ошибка построения http запроса: %v", err)
			m.collector.RecordError(ctx, metrics.CategoryUnknown)
			_ = m.reader.CommitMessages(ctx, msg)
			continue
		}

		if ak := os.Getenv("MINIO_ACCESS_KEY"); ak != "" {
			req.SetBasicAuth(ak, os.Getenv("MINIO_SECRET_KEY"))
		}
		respHTTP, err := http.DefaultClient.Do(req)
		if err != nil || respHTTP.StatusCode != 200 {
			log.Printf("file-validator: ошибка http получения объекта: %v status=%v", err, func() interface{} {
				if respHTTP != nil {
					return respHTTP.Status
				}
				return nil
			}())
			m.collector.RecordError(ctx, metrics.CategoryCorruptFile)
			if reqID != "" {
				resp := map[string]interface{}{"request_id": reqID, "status": "invalid", "error": "object_fetch_failed"}
				b, _ := json.Marshal(resp)
				_ = m.producers.SendValidated(ctx, msg.Key, b)
			}
			if respHTTP != nil && respHTTP.Body != nil {
				_ = respHTTP.Body.Close()
			}
			_ = m.reader.CommitMessages(ctx, msg)
			continue
		}
		data, err := io.ReadAll(respHTTP.Body)
		_ = respHTTP.Body.Close()
		if err != nil {
			log.Printf("file-validator: ошибка чтения объекта: %v", err)
			m.collector.RecordError(ctx, metrics.CategoryCorruptFile)
			if reqID != "" {
				resp := map[string]interface{}{"request_id": reqID, "status": "invalid", "error": "object_read_failed"}
				b, _ := json.Marshal(resp)
				_ = m.producers.SendValidated(ctx, msg.Key, b)
			}
			_ = m.reader.CommitMessages(ctx, msg)
			continue
		}


		if err := m.svc.ValidateAndStore(ctx, docID, data); err != nil {
			log.Printf("file-validator: валидация/сохранение не удались для id=%s: %v", docID, err)
			m.collector.RecordError(ctx, metrics.CategoryInvalidFile)

			if reqID != "" {
				resp := map[string]interface{}{"request_id": reqID, "status": "invalid", "error": err.Error()}
				b, _ := json.Marshal(resp)
				_ = m.producers.SendValidated(ctx, msg.Key, b)
			}
			_ = m.reader.CommitMessages(ctx, msg)
			continue
		}

		m.collector.RecordProcessed(ctx)


		if reqID != "" {
			resp := map[string]interface{}{"request_id": reqID, "status": "valid"}
			b, _ := json.Marshal(resp)
			if err := m.producers.SendValidated(ctx, msg.Key, b); err != nil {
				log.Printf("file-validator: ошибка отправки ответа: %v", err)
			}
		}

		if err := m.reader.CommitMessages(ctx, msg); err != nil {
			log.Printf("file-validator: ошибка коммита: %v", err)
			return err
		}
	}
}


func (m *Manager) Close() {
	if m.reader != nil {
		_ = m.reader.Close()
	}
}


