mq:
	docker run --rm -it -p 5672:5672 -p 15672:15672 -v lavinmq_data:/tmp/amqp cloudamqp/lavinmq

mq-def:
    docker run -d --hostname my-lavinmq --name my-lavinmq \
    -p 5672:5672 -p 15672:15672 \
    -v definitions_mq.json:/etc/rabbitmq/definitions.json \
    -e RABBITMQ_DEFAULT_USER=guest -e RABBITMQ_DEFAULT_PASS=guest \
  lavinmq/lavinmq