package main

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/brandaogabriel7/full-cycle-go-intensivo/internal/infra/database"
	"github.com/brandaogabriel7/full-cycle-go-intensivo/internal/usecase"
	"github.com/brandaogabriel7/full-cycle-go-intensivo/pkg/kafka"
	"github.com/brandaogabriel7/full-cycle-go-intensivo/pkg/rabbitmq"
	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	amqp "github.com/rabbitmq/amqp091-go"

	// sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./orders.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	repository := database.NewOrderRepository(db)
	usecase := usecase.NewCalculateFinalPrice(repository)

	msgChanKafka := make(chan *ckafka.Message)
	topics := []string{"orders"}
	servers := "host.docker.internal:29092"
	fmt.Println("Kafka consumer has started")
	go kafka.Consume(topics, servers, msgChanKafka)
	go kafkaWorker(msgChanKafka, usecase)

	ch, err := rabbitmq.OpenChannel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	msgRabbitmqChannel := make(chan amqp.Delivery)
	go rabbitmq.Consume(ch, msgRabbitmqChannel)
	rabbitmqWorker(msgRabbitmqChannel, usecase)
}

func kafkaWorker(msgChan chan *ckafka.Message, uc *usecase.CalculateFinalPrice) {
	fmt.Println("Kafka worker has started")
	for msg := range msgChan {
		var orderInputDTO usecase.OrderInputDTO
		err := json.Unmarshal(msg.Value, &orderInputDTO)
		if err != nil {
			panic(err)
		}
		outputDto, err := uc.Execute(orderInputDTO)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Kafka has processed order %s\n", outputDto.ID)
	}
}

func rabbitmqWorker(msgChan chan amqp.Delivery, uc *usecase.CalculateFinalPrice) {
	fmt.Println("RabbitMQ worker has started")
	for msg := range msgChan {
		var orderInputDTO usecase.OrderInputDTO
		err := json.Unmarshal(msg.Body, &orderInputDTO)
		if err != nil {
			panic(err)
		}
		outputDto, err := uc.Execute(orderInputDTO)
		if err != nil {
			panic(err)
		}
		msg.Ack(false)
		fmt.Printf("RabbitMQ has processed order %s\n", outputDto.ID)
	}
}
