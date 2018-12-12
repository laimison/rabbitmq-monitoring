#!/bin/bash

rabbitmqctl list_queues -q name messages -p Some_Virtual_Host

printf "\nMore verbose version:\n"
rabbitmqctl list_queues -p Some_Virtual_Host name messages messages_ready state consumer_utilisation

