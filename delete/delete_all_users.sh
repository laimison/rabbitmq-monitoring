#!/bin/bash

echo "deleting users"
users=`rabbitmqctl list_users -q | awk -F ' ' '{print $1}'`
for u in $users; do rabbitmqctl delete_user $u; done
