package mgp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"math/rand"
	"strconv"
	"time"
)

type ProduceData struct {
	Ct      int
	Uuid    string
	BizJson string
}

type MqProperties struct {
	prefix   string
	qCount   int
	qStartNo int
}

var mqPool [] *amqp.Channel

var mqPro MqProperties


func ChBean() (int, *amqp.Channel, MqProperties) {
	/*
		This is because by default the seed is always the same, the number 1.
		To actually get a random number, you need to provide a unique seed for your program.
		You really want to not forget seeding, and instead properly seed our pseudonumber generator. How?

		Use rand.Seed() before calling any math/rand method, passing an int64 value.
		You just need to seed once in your program, not every time you need a random number.
		The most used seed is the current time, converted to int64 by UnixNano with rand.Seed(time.Now().UnixNano()):
	*/
	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(mqPro.qCount)

	return r, mqPool[r], mqPro
}

func InitChannel(amqpConfig AmqpCfg) {
	prefix := amqpConfig.QueuePrefix
	qCount := amqpConfig.QueueCount

	mqPro.prefix = prefix
	mqPro.qCount = qCount
	mqPro.qStartNo = amqpConfig.QueueStartNo

	for i := 0; i < qCount; i++ {
		no := i + amqpConfig.QueueStartNo

		mqPool = append(mqPool, GetChannel(amqpConfig, no))
	}
}

func getEx() string{
	return mqPro.prefix + ".ex"
}

func getQueue(no int) string{
	return mqPro.prefix + ".q." + strconv.Itoa(no)
}

func getRk(no int) string{
	return mqPro.prefix + ".rk." + strconv.Itoa(no)
}

// GetChannel int | start from 1
func GetChannel(amqpConfig AmqpCfg, no int) *amqp.Channel {
	ex := getEx()
	queue := getQueue(no)
	rk := getRk(no)

	var buf bytes.Buffer

	userName := amqpConfig.Username
	password := amqpConfig.Password
	host := amqpConfig.Host
	port := amqpConfig.Port
	vHost := amqpConfig.VHost

	buf.WriteString("amqp://")
	buf.WriteString(userName)
	buf.WriteString(":")
	buf.WriteString(password)

	buf.WriteString("@" + host + ":" + strconv.Itoa(port) + "/" + vHost)

	url := buf.String()

	fmt.Printf("rabbitmq connect: %s\n", url)

	conn, err := amqp.Dial(url)
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	err = ch.ExchangeDeclare(ex, "direct", true, false, false, false, nil)
	failOnError(err, "Failed to Declare a exchange")

	q, err := ch.QueueDeclare(
		queue, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue"+q.Name)

	err = ch.QueueBind(queue, rk, ex, false, nil)
	failOnError(err, "Failed to bind a queue")

	return ch
}

func pkgProduceData(bizJson string, uuid string) string {
	var produceData ProduceData

	produceData.Uuid = uuid
	produceData.Ct = int(time.Now().Unix())
	produceData.BizJson = bizJson

	jsonBytes, err := json.Marshal(produceData)
	if err != nil {
		fmt.Println(err)
	}

	return string(jsonBytes)
}

func Produce(i int, mqPro MqProperties, ch *amqp.Channel, bizJson string, uuid string) {
	msg := pkgProduceData(bizJson, uuid)

	ex := mqPro.prefix + ".ex"
	rk := mqPro.prefix + ".rk." + strconv.Itoa(i+mqPro.qStartNo)
	Info("ex:"+ex+" | rk:"+rk+" | body:"+msg, uuid)

	err := ch.Publish(
		ex,    // exchange
		rk,    // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		},
	)

	failOnError(err, "Failed to publish a message")
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
