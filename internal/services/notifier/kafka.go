package notifier

import (
	"auth_service/internal/config"
	"auth_service/internal/domain/models"
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

type Notification struct {
	Message string `json:"message"`
}

type Notifier interface {
	Send(*models.User) error
}

type KafkaNotifier struct {
	producer sarama.SyncProducer
	topic    string
}

func setupProducer(cfg config.KafkaConfig) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer([]string{cfg.Host},
		config)
	if err != nil {
		return nil, fmt.Errorf("failed to setup producer: %w", err)
	}
	return producer, nil
}

func New(cfg config.KafkaConfig) (*KafkaNotifier, error) {

	producer, err := setupProducer(cfg)
	if err != nil {
		log.Fatalf("failed start notifier %v", err.Error())
	}

	return &KafkaNotifier{
		producer: producer,
		topic:    cfg.Topic,
	}, nil
}

func (k *KafkaNotifier) Send(usr *models.User) error {

	message := fmt.Sprintf("User %d Email %s", usr.ID, usr.Email)

	notification := Notification{
		Message: message,
	}

	notificationJSON, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: k.topic,
		Key:   sarama.StringEncoder(fmt.Sprintf("%d", usr.ID)),
		Value: sarama.StringEncoder(notificationJSON),
	}

	_, _, err = k.producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}
