#!/bin/bash

rabbitmqctl add_vhost Some_Virtual_Host || exit 1
rabbitmqctl set_permissions -p Some_Virtual_Host guest ".*" ".*" ".*" || exit 1
rabbitmqctl set_permissions -p Some_Virtual_Host mon ".*" ".*" ".*" || exit 1
rabbitmqctl list_vhosts
