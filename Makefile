.PHONY: mq

mq:
	docker run --rm -it -p 5672:5672 -p 15672:15672 -v lavinmq_data:/tmp/amqp cloudamqp/lavinmq