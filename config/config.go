package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)


type Config struct {
	ServiceName   string         `yaml:"serviceName"`
	Port          int            `yaml:"port"`
	EnableSwagger bool           `yaml:"enableSwagger"`
	Database      DatabaseConfig `yaml:"database"`
	Kafka         KafkaConfig    `yaml:"kafka"`
	Topics        TopicsConfig   `yaml:"topics"`
	Redis         RedisConfig    `yaml:"redis"`
	Minio         MinioConfig    `yaml:"minio"`
}


type DatabaseConfig struct {
	Shards []ShardConfig `yaml:"shards"`
}


type ShardConfig struct {
	Name     string `yaml:"name"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	DBName   string `yaml:"dbname"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	ConnURL  string `yaml:"connURL"`
}


type KafkaConfig struct {
	Brokers         []string `yaml:"brokers"`
	GroupID         string   `yaml:"groupId"`
	MaxMessageBytes int      `yaml:"maxMessageBytes"`
}


type TopicsConfig struct {
	Input    string `yaml:"input"`
	Output   string `yaml:"output"`
	Response string `yaml:"response"`
}


type RedisConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	DB   int    `yaml:"db"`
	TTL  int    `yaml:"ttlSeconds"`
}


type MinioConfig struct {
	Endpoint string `yaml:"endpoint"`
	Bucket   string `yaml:"bucket"`
}


func LoadConfig(filename string) (*Config, error) {
	cfg := &Config{}

	if strings.TrimSpace(filename) != "" {
		data, err := os.ReadFile(filename)
		if err != nil {
			return nil, fmt.Errorf("не удалось прочитать файл конфигурации: %w", err)
		}

		if err := yaml.Unmarshal(data, cfg); err != nil {
			return nil, fmt.Errorf("не удалось разобрать YAML: %w", err)
		}
	}
	cfg.applyEnv()
	return cfg, nil
}

func (c *Config) applyEnv() {
	if env := strings.TrimSpace(os.Getenv("SERVICE_NAME")); env != "" {
		c.ServiceName = env
	}
	if env := strings.TrimSpace(os.Getenv("PORT")); env != "" {
		if port, err := strconv.Atoi(env); err == nil {
			c.Port = port
		}
	}

	if env := strings.TrimSpace(os.Getenv("ENABLE_SWAGGER")); env != "" {
		if enabled, err := strconv.ParseBool(env); err == nil {
			c.EnableSwagger = enabled
		}
	}
	if env := strings.TrimSpace(os.Getenv("KAFKA_BOOTSTRAP_SERVERS")); env != "" {
		parts := strings.Split(env, ",")
		for i := range parts {
			parts[i] = strings.TrimSpace(parts[i])
		}
		c.Kafka.Brokers = parts
	}
	if env := strings.TrimSpace(os.Getenv("KAFKA_GROUP_ID")); env != "" {
		c.Kafka.GroupID = env
	}
	if env := strings.TrimSpace(os.Getenv("KAFKA_MAX_MESSAGE_BYTES")); env != "" {
		if maxBytes, err := strconv.Atoi(env); err == nil {
			c.Kafka.MaxMessageBytes = maxBytes
		}
	}

	if env := strings.TrimSpace(os.Getenv("VALIDATOR_INPUT_TOPIC")); env != "" {
		c.Topics.Input = env
	}
	if env := strings.TrimSpace(os.Getenv("VALIDATOR_OUTPUT_TOPIC")); env != "" {
		c.Topics.Output = env
	}
	if env := strings.TrimSpace(os.Getenv("VALIDATOR_RESPONSE_TOPIC")); env != "" {
		c.Topics.Response = env
	}

	if env := strings.TrimSpace(os.Getenv("REDIS_ADDR")); env != "" {

		if strings.Contains(env, ":") {
			parts := strings.Split(env, ":")
			c.Redis.Host = parts[0]
			if port, err := strconv.Atoi(parts[1]); err == nil {
				c.Redis.Port = port
			}
		} else {
			c.Redis.Host = env
		}
	}
	if env := strings.TrimSpace(os.Getenv("REDIS_DB")); env != "" {
		if db, err := strconv.Atoi(env); err == nil {
			c.Redis.DB = db
		}
	}


	if env := strings.TrimSpace(os.Getenv("MINIO_ENDPOINT")); env != "" {
		c.Minio.Endpoint = env
	}
	if env := strings.TrimSpace(os.Getenv("MINIO_BUCKET")); env != "" {
		c.Minio.Bucket = env
	}


	if env := strings.TrimSpace(os.Getenv("DB_SHARDS")); env != "" {
		parts := strings.Split(env, ",")
		shards := make([]ShardConfig, 0, len(parts))
		for i, p := range parts {
			shards = append(shards, ShardConfig{
				Name:    fmt.Sprintf("shard%d", i),
				ConnURL: strings.TrimSpace(p),
			})
		}
		c.Database.Shards = shards
	}
}


