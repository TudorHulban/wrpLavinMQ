all: mq-def

.default_goal := all

mq:
	docker run --rm -it -p 5672:5672 -p 15672:15672 -v lavinmq_data:/tmp/amqp cloudamqp/lavinmq

mq-def:
	docker run -d --hostname test-mq --name test-mq \
	-p 5672:5672 -p 15672:15672 \
	-v $(PWD)/ops/definitions_mq.json:/etc/rabbitmq/definitions.json \
	-v $(PWD)/ops/rabbitmq.conf:/etc/rabbitmq/rabbitmq.conf \
	rabbitmq:management-alpine
