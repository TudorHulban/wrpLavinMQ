# Advanced Message Queuing Protocol (AMQP)

Integration example based on Makefile target, see docker invocation with definitions.

## Users

An admin user (`gtest/gtest`) is created through the definitions file to run the tests. This user should have access to the definitions virtual host.

## Management UI

Available at `localhost:15672` credentials are `gtest/gtest`.

## Messaging entities

Messaging entities are created through the defintions file in ops folder.

## Troubleshooting

### amqp-tools

Install on Alma Liux:

#### Install dependencies

```sh
sudo dnf groupinstall "Development Tools"
sudo dnf install cmake libuuid-devel
```

#### Clone and build

```sh
git clone https://github.com/rabbitmq/rabbitmq-tools.git
cd rabbitmq-tools
cmake .
make
sudo make install
```

#### Verify

```sh
amqp-publish --help
amqp-consume --help
```
