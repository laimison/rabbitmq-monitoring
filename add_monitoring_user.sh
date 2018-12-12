#!/bin/bash

rabbitmqctl add_user mon password123 || exit 1
rabbitmqctl set_user_tags mon monitoring || exit 1
rabbitmqctl set_permissions -p / mon ".*" ".*" ".*" || exit 1
echo success
