# rabbitmq-monitoring

This repo helps to install and configure RabbitMQ. The main goal is to prepare monitoring for the queues then.

## Usage

Run everything in the order

`./install_dependencies.sh`

`./install_and_start_rabbitmq.sh`

`./add_user_perm.sh`

`./make_vhost_and_set_perm.sh`

`./make_exchange.sh`

`./make_queue.sh`

`./make_binding.py` or `./make_binding.sh` - got access denied so fix is needed

## Requirements

Scripts are designed for Mac OS.

## References

[https://stackoverflow.com/questions/4545660/rabbitmq-creating-queues-and-bindings-from-a-command-line](https://stackoverflow.com/questions/4545660/rabbitmq-creating-queues-and-bindings-from-a-command-line)
