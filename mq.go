package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
	"log"
	"mgp"
	"os"
	"strconv"
	"strings"
	"time"
)

func init() {
	now := time.Now()
	hour, minute, second := now.Clock()

	// go_2021-11-25_171957.log
	filepath := "./runtime/log/go_" +
		strings.Replace(now.String()[:10], ":", "_", 3) + "_" +
		strconv.Itoa(hour) + strconv.Itoa(minute) + strconv.Itoa(second) +
		".log"

	logFile, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	// set logFile storage location
	log.SetOutput(logFile)
}

func main() {
	env := mgp.ReadEnv()

	if env.ServiceRole == "producer" {
		runAsProducer()
	} else if env.ServiceRole == "consumer" {
		runAsConsumer()
	} else {
		fmt.Println("ServiceRole Error")
	}
}

func runAsProducer() {
	fmt.Println("runAsProducer")
	log.Println("runAsProducer")

	config := mgp.ReadConfig()
	httpPort := config.HttpPort
	amqpConfig := config.Amqp

	mgp.InitChannel(amqpConfig)

	r := gin.Default()

	r.POST("/gateway", func(c *gin.Context) {
		requestHandler(c)
	})

	r.GET("/gateway", func(c *gin.Context) {
		requestHandler(c)
	})

	r.Run(":" + strconv.Itoa(httpPort))
}

func runAsConsumer() {
	fmt.Println("runAsConsumer")
	config := mgp.ReadConfig()
	amqpConfig := config.Amqp
	transferUrl := amqpConfig.TransferUrl
	queueStartNo := amqpConfig.QueueStartNo

	qCount := amqpConfig.QueueCount
	var msgPool [] <-chan amqp.Delivery

	for i := 0; i < qCount; i++ {
		qNo := i + queueStartNo

		queue := amqpConfig.QueuePrefix + ".q." + strconv.Itoa(qNo)
		ch := mgp.GetChannel(amqpConfig, qNo)

		messages, err := ch.Consume(
			queue, // queue
			queue,       // consumer
			true,         // auto-ack
			false,        // exclusive
			false,        // no-local
			false,        // no-wait
			nil,          // args
		)
		failOnError(err, "Failed to register a consumer")

		msgPool = append(msgPool, messages)
	}
	forever := make(chan bool)
	for i := 0; i < qCount; i++ {
		no := i
		qNo := no + queueStartNo

		go func() {
			fmt.Println("go handle msg queue pool:" + strconv.Itoa(qNo))

			for d := range msgPool[no] {
				msg := d.Body

				fmt.Printf("Queue Pool "+strconv.Itoa(qNo)+" Received a message: %s\n", msg)

				var msgStruct mgp.ProduceData

				json.Unmarshal(msg, &msgStruct)
				uuid := msgStruct.Uuid
				mgp.Info("Queue Pool "+strconv.Itoa(qNo)+"Received a message: "+string(msg), uuid)

				_ = mgp.HttpPostForm(transferUrl, msgStruct, uuid)
			}
		}()
	}

	fmt.Println("Waiting for messages")

	<-forever
}

func requestHandler(c *gin.Context) {
	uuid := mgp.InitUuid()
	argJson := c.Request.FormValue("json")
	mgp.Info("args:" + argJson, uuid)

	i, ch, mqPro := mgp.ChBean()
	go mgp.Produce(i, mqPro, ch, argJson, uuid)

	c.JSON(200, gin.H{
		"code": 200,
		"msg": "success",
	})
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}