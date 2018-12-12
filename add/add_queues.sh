#!/bin/bash

rabbitmqadmin declare queue --vhost=Some_Virtual_Host name=some_outgoing_queue durable=true || exit 1
rabbitmqadmin declare queue --vhost=Some_Virtual_Host name=some_incoming_queue || exit 1
echo success
