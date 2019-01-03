# rabbitmq-monitoring

This repo helps to install, configure RabbitMQ and also monitor your queues using RabbitMQ API or Prometheus.

## Usage

* Firstly, install dependencies and RabbitMQ. Check this directory for the scripts:

`install/`

* If you want to do some clean up in your development environment check scripts in a directory - this is destructive and can delete all your users and queues!

`delete/`

* Create users, vhosts, exchanges and queues

`add/` - you need monitoring user if you are looking for monitoring solution

* Make a binding

`./bind.py` or `./bind.sh`

* To enable plugin, check this directory

`enable/` - rabbitmq_management is usually enabled so you can do monitoring GET queries

* There are more commands in a root dir to check the queues, publish messages, check some statuses, utilizations, etc.

* Main monitoring script is written in Go:

`./monitoring.go --help`

## Requirements

Scripts are tested on Mac OS, but could be suitable for Linux.

## References

[https://stackoverflow.com/questions/4545660/rabbitmq-creating-queues-and-bindings-from-a-command-line](https://stackoverflow.com/questions/4545660/rabbitmq-creating-queues-and-bindings-from-a-command-line)
[https://stackoverflow.com/questions/52010390/permissions-that-need-to-be-assigned-to-a-monitoring-software](https://stackoverflow.com/questions/52010390/permissions-that-need-to-be-assigned-to-a-monitoring-software)
