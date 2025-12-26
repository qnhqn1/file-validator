package metrics

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)


const (
	CategoryInvalidFile  = "invalid_file"
	CategoryMissingParts = "missing_parts"
	CategoryCorruptFile  = "corrupt_file"
	CategoryUnknown      = "unknown"
)


type Snapshot struct {
	Received        int64            `json:"received"`
	Processed       int64            `json:"processed"`
	Errors          int64            `json:"errors"`
	ErrorCategories map[string]int64 `json:"error_categories"`
}


type Collector struct {
	receivedCounter      metric.Int64Counter
	processedCounter     metric.Int64Counter
	errorsCounter        metric.Int64Counter
	errorCategoryCounter metric.Int64Counter
	received             int64
	processed            int64
	errors               int64
	errorCategories      map[string]int64
	mu                   sync.Mutex
}


func New() (*Collector, error) {
	reader := sdkmetric.NewManualReader()
	meterProvider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))
	otel.SetMeterProvider(meterProvider)

	meter := meterProvider.Meter("file-validator")

	receivedCounter, err := meter.Int64Counter("kafka_messages_received_total")
	if err != nil {
		return nil, err
	}
	processedCounter, err := meter.Int64Counter("kafka_messages_processed_total")
	if err != nil {
		return nil, err
	}
	errorsCounter, err := meter.Int64Counter("kafka_messages_errors_total")
	if err != nil {
		return nil, err
	}
	errorCategoryCounter, err := meter.Int64Counter("kafka_messages_errors_by_category_total")
	if err != nil {
		return nil, err
	}

	categories := map[string]int64{
		CategoryInvalidFile:  0,
		CategoryMissingParts: 0,
		CategoryCorruptFile:  0,
		CategoryUnknown:      0,
	}

	return &Collector{
		receivedCounter:      receivedCounter,
		processedCounter:     processedCounter,
		errorsCounter:        errorsCounter,
		errorCategoryCounter: errorCategoryCounter,
		errorCategories:      categories,
	}, nil
}


func (c *Collector) RecordReceived(ctx context.Context) {
	c.mu.Lock()
	c.received++
	c.mu.Unlock()
	c.receivedCounter.Add(ctx, 1)
}


func (c *Collector) RecordProcessed(ctx context.Context) {
	c.mu.Lock()
	c.processed++
	c.mu.Unlock()
	c.processedCounter.Add(ctx, 1)
}


func (c *Collector) RecordError(ctx context.Context, category string) {
	if category == "" {
		category = CategoryUnknown
	}

	c.mu.Lock()
	c.errors++
	c.errorCategories[category]++
	c.mu.Unlock()

	c.errorsCounter.Add(ctx, 1)
	c.errorCategoryCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("category", category)))
}


func (c *Collector) Snapshot() Snapshot {
	c.mu.Lock()
	defer c.mu.Unlock()

	categorySnapshot := make(map[string]int64, len(c.errorCategories))
	for name, count := range c.errorCategories {
		categorySnapshot[name] = count
	}

	return Snapshot{
		Received:        c.received,
		Processed:       c.processed,
		Errors:          c.errors,
		ErrorCategories: categorySnapshot,
	}
}


func (c *Collector) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(c.Snapshot())
	}
}


