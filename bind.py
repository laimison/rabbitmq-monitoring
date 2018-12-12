#!/usr/bin/env python
import pika

rabbitmq_host = "127.0.0.1"
rabbitmq_port = 5672
rabbitmq_virtual_host = "Some_Virtual_Host"
rabbitmq_send_exchange = "some_exchange"
rabbitmq_rcv_exchange = "some_exchange"

rabbitmq_rcv_queue = "some_incoming_queue"
rabbitmq_rcv_key = "some_routing_key"

# rabbitmq_user = "testuser"
# rabbitmq_password = "testpassword"

rabbitmq_user = "guest"
rabbitmq_password = "guest"

outgoingRoutingKeys = ["outgoing_routing_key"]
outgoingQueues = ["some_outgoing_queue"]

# The binding area
credentials = pika.PlainCredentials(rabbitmq_user, rabbitmq_password)
connection = pika.BlockingConnection(pika.ConnectionParameters(rabbitmq_host, rabbitmq_port, rabbitmq_virtual_host, credentials))
channel = connection.channel()
channel.queue_bind(exchange=rabbitmq_rcv_exchange, queue=rabbitmq_rcv_queue, routing_key=rabbitmq_rcv_key)

for index in range(len(outgoingRoutingKeys)):
    channel.queue_bind(exchange=rabbitmq_send_exchange, queue=outgoingQueues[index], routing_key=outgoingRoutingKeys[index])
