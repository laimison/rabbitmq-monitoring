# rabbitmq-monitoring

This repo helps to install and configure RabbitMQ. The main goal is to prepare monitoring for the queues then.

## Usage

Run everything in the order

`./install_dependencies.sh`

`./install_and_start_rabbitmq.sh`

`./add_user_perm.sh` - add a user and permissions

`./make_vhost_and_set_perm.sh` - make a virtual host and set permissions for guest

`./make_exchange.sh` - make an exchange for virtual host

`./make_queue.sh` - make a queue for virtual host

`./make_binding.py` or `./make_binding.sh` - make a binding - additionally created some_incoming_queue in vhost so script is able to reach that`

`./enable_rabbitmq_management.sh`

`./list_queues.sh` - list queue names and messages count

`./list_queues_builtin_api.sh` - list information about the queues from API in json

`./publish_message.sh` - publish message in a queue by hand

## Requirements

Scripts are designed for Mac OS.

## References

[https://stackoverflow.com/questions/4545660/rabbitmq-creating-queues-and-bindings-from-a-command-line](https://stackoverflow.com/questions/4545660/rabbitmq-creating-queues-and-bindings-from-a-command-line)
