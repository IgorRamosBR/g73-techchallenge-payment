#!/bin/bash

# Wait for RabbitMQ server to start
sleep 10

# Define RabbitMQ variables
RABBITMQ_USER="techchallenge"
RABBITMQ_PASS="admin123"
RABBITMQ_HOST="localhost"
RABBITMQ_PORT="15672"

# Create exchanges
curl -i -u $RABBITMQ_USER:$RABBITMQ_PASS -H "content-type:application/json" \
  -XPUT http://$RABBITMQ_HOST:$RABBITMQ_PORT/api/exchanges/%2f/order_events \
  -d'{"type":"topic","durable":true}'

# Create queues
curl -i -u $RABBITMQ_USER:$RABBITMQ_PASS -H "content-type:application/json" \
  -XPUT http://$RABBITMQ_HOST:$RABBITMQ_PORT/api/queues/%2f/orders_paid_queue \
  -d'{"durable":true}'
curl -i -u $RABBITMQ_USER:$RABBITMQ_PASS -H "content-type:application/json" \
  -XPUT http://$RABBITMQ_HOST:$RABBITMQ_PORT/api/queues/%2f/orders_ready_queue \
  -d'{"durable":true}'

# Create bindings
curl -i -u $RABBITMQ_USER:$RABBITMQ_PASS -H "content-type:application/json" \
  -XPOST http://$RABBITMQ_HOST:$RABBITMQ_PORT/api/bindings/%2f/e/order_events/q/orders_paid_queue \
  -d'{"routing_key":"orders.paid"}'
curl -i -u $RABBITMQ_USER:$RABBITMQ_PASS -H "content-type:application/json" \
  -XPOST http://$RABBITMQ_HOST:$RABBITMQ_PORT/api/bindings/%2f/e/order_events/q/orders_ready_queue \
  -d'{"routing_key":"orders.ready"}'
