### Step 1: Container
Run Container
```
docker run -itp 9001:9001 --name go_temp -v /usr/local/project/temp/go_amqp/:/home/ -d golang:1.16.6
```

Enter Container
```
docker exec -it go_temp bash
```

### Step 2: Build
```
cd /home
go env -w GO111MODULE=on
```

```
go build mq.go
```

### Step 3: Config

create file by config/*_demo.yaml
- config.yaml
- env.yaml


config.yaml's field | Description
---|---
http_port | [producer only]
amqp.username | RabbitMQ's
amqp.password |
amqp.host |
amqp.port |
amqp.vhost | need "/"
amqp.queue_prefix |
amqp.queue_count |
amqp.queue_start_no |
amqp.transfer_url | [comsumer only]

`transfer_url`

consumer service will http post `transfer_url` with parameters

`queue_prefix` `queue_count` `queue_start_no`

case 1:
- queue_start_no=1
- queue_count=2
- queue_prefix="seck"

service wille produce/consume
- exchange
    - seck.ex
- queue
    - seck.q.1
    - seck.q.2
- routingKey
    - seck.rk.1
    - seck.rk.2

case 2:
- queue_start_no=3
- queue_count=1
- queue_prefix="seck"

service wille produce/consume
- exchange
    - seck.ex
- queue
    - seck.q.3
- routingKey
    - seck.rk.3

### Step 4: Run

```
sudo chmod -R 777 ./mq && ./mq
```