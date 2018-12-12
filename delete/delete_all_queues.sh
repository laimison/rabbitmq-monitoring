#!/bin/bash

echo "deleting queues"
for v in `rabbitmqctl list_vhosts -q`; do echo $v; for q in `rabbitmqctl list_queues -p $v -q`; do rabbitmqctl delete_queue $q -p $v; done ; done
