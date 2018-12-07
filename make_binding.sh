#!/bin/bash

rabbitmqadmin --vhost="Some_Virtual_Host" declare binding source="some_exchange" destination_type="queue" destination="some_incoming_queue" routing_key="some_routing_key"
