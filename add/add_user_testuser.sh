#!/bin/bash

rabbitmqctl add_user testuser testpassword || exit 1
rabbitmqctl set_user_tags testuser administrator || exit 1
rabbitmqctl set_permissions -p / testuser ".*" ".*" ".*" || exit 1
echo success
