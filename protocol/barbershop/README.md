#Barbershop protocol
This protocol is based loosely on the Redis protocol specification.
[https://redis.io/topics/protocol](https://redis.io/topics/protocol)

Every barbershop command or data transmitted by the client and the server is terminated by "\r\n" (CRLF).

The simplest commands are the inline commands. This is an example of a server/client chat (the server chat starts with S:, the client chat with C:).

    C: UPDATE some_string_key_without_whitespace 5\r\n
    S: +OK\r\n
    C: NEXT\r\n
    S: +some_string_key_without_whitespace\r\n

An inline command is a CRLF-terminated string sent to the client. The server can reply to commands in different ways:

    With an error message (the first byte of the reply will be "-")
    With a single line reply (the first byte of the reply will be "+)

This service does not support bulk commands or responses.

## Commands

There are only a handful of commands supported at this point.

'UPDATE <item key> <value>'

Update the priority of a given item by X.

    C: UPDATE item_key 5\r\n
    S: +OK\r\n

'NEXT'

Return the next item in the queue.

    C: NEXT\r\n
    S: +item_key\r\n

When there are no more items to return a '-1' is returned.

    C: NEXT\r\n
    S: +-1\r\n

'PEEK'

Return the next item in the queue without removing it from the queue.

    C: PEEK\r\n
    S: +item_key\r\n

When there are no more items to return a '-1' is returned.

    C: PEEK\r\n
    S: +-1\r\n

'SCORE <item key>'

Return the score of a given item.

    C: SCORE item_key\r\n
    S: +5\r\n

If the item does not exist or have a score then a '-1' is returned.

    C: SCORE item_key\r\n
    S: +-1\r\n

'INFO'

Return some server stats. This command deviates from the standard response
format by returning a list of key-value pairs separated by a ':'.

*Note: the exact items that are available here are currently in progress and may change*
* 'uptime' (32u) Number of seconds this server has been running.
* 'version' (string) Version string of this server.
* 'updates' (32u) Number of update commands received by this server.
* 'items' (32u) Number of items.
* 'pools' (32u) Number of pools.

    C: INFO\r\n
    S: uptime:60000\r\n
    S: version:0.2.1\r\n
    S: updates:9742851\r\n
    S: items:2132931\r\n
    S: pools:47831\r\n
