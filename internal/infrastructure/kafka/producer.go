package kafka

import (
	"encoding/json"
	"time"

	"asset-core/internal/config"
	"asset-core/internal/pkg/logger"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
)

type Producer interface {
	Publish(topic string, event Event) error
	Close() error
}

type Event struct {
	EventID    string      `json:"event_id"`
	EventType  string      `json:"event_type"`
	OccurredAt time.Time   `json:"occurred_at"`
	TraceID    string      `json:"trace_id,omitempty"`
	OperatorID string      `json:"operator_id,omitempty"`
	Payload    interface{} `json:"payload"`
}

type noopProducer struct{}

func (p noopProducer) Publish(string, Event) error { return nil }
func (p noopProducer) Close() error                { return nil }

type saramaProducer struct {
	producer sarama.SyncProducer
	log      *logger.Logger
}

func NewProducer(cfg config.KafkaConfig, log *logger.Logger) Producer {
	if !cfg.Enabled {
		return noopProducer{}
	}
	scfg := sarama.NewConfig()
	scfg.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(cfg.Brokers, scfg)
	if err != nil {
		log.Fatal("kafka producer init failed", logger.Error(err))
	}
	return &saramaProducer{producer: producer, log: log}
}

func NewEvent(eventType string, payload interface{}) Event {
	return Event{
		EventID:    uuid.NewString(),
		EventType:  eventType,
		OccurredAt: time.Now(),
		Payload:    payload,
	}
}

func (p *saramaProducer) Publish(topic string, event Event) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}
	_, _, err = p.producer.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(event.EventID),
		Value: sarama.ByteEncoder(body),
	})
	return err
}

func (p *saramaProducer) Close() error {
	return p.producer.Close()
}
