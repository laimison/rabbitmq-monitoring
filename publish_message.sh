#!/bin/bash

rabbitmqadmin publish --vhost=Some_Virtual_Host exchange=some_exchange routing_key=outgoing_routing_key payload="hello"
