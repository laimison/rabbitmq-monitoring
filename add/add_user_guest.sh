#!/bin/bash

rabbitmqctl add_user guest guest || exit 1
rabbitmqctl set_user_tags guest administrator || exit 1
rabbitmqctl set_permissions -p / guest ".*" ".*" ".*" || exit 1
rabbitmqctl change_password guest guest || exit 1
echo success
