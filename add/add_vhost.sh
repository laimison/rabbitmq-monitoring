#!/bin/bash

rabbitmqctl add_vhost Some_Virtual_Host || exit 1
rabbitmqctl list_users | grep -q ^guest && rabbitmqctl set_permissions -p Some_Virtual_Host guest ".*" ".*" ".*" || exit 1
rabbitmqctl list_users | grep -q ^testuser && rabbitmqctl set_permissions -p Some_Virtual_Host testuser ".*" ".*" ".*" || exit 1
rabbitmqctl list_users | grep -q ^monitoring && rabbitmqctl set_permissions -p Some_Virtual_Host monitoring "" "" "" || exit 1
rabbitmqctl list_vhosts
