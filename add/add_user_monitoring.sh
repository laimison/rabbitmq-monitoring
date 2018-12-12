#!/bin/bash

rabbitmqctl add_user monitoring password || exit 1
rabbitmqctl set_user_tags monitoring monitoring || exit 1
rabbitmqctl set_permissions -p / mon "" "" "" || exit 1
echo success
