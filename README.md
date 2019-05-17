# gofirst

Gofirst is a lightweight in memory queue. My initial goal was to rewrite [Barbershop](https://github.com/ngerakines/babershop) in Go to familiarize myself with the language, but in the process, I wanted to see if I can bring some improvements along the way. Please excuse the mess, as this is part of a learning process for me, but feel free to log a PR or file an issue if there is anything that can be done better.

# Installation
```go get -u -v github.com/dwayn/gofirst```

# Usage

Currently there is only a priority queue implementation, but as this moves forward, the plan is to add more queue types as well as named queues with the ability to target them at query time.

## Command Line Options

    -h  --help      Print help information
    -l  --listen    Interface to listen on. Default: localhost
    -p  --port      Port to listen on. Default: 3333
    -r  --protocol  Network protocol: resp, barbershop. Default: resp


# Protocols 

## Barbershop Protocol
An equivalent implementation of the protocol defined in [Barbershop](https://github.com/ngerakines/barbershop). Since this project was defined as a clone of Barbershop, it seemed only fitting to implement the protocol.
Barbershop protocol [README is here](protocol/barbershop/README.md)

## Redis Serialization Protocol
In an effort to standardize the protocol (and build a modularized protocol handler), a modified version of the Barbershop protocol that is completely, rather than loosely, [RESP](https://redis.io/topics/protocol) compliant is available.
RESP protocol [README is here](protocol/resp/README.md)


