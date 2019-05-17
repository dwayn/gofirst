# RESP - Redis Serialization Protocol
This protocol is based strictly on the Redis protocol specification
[https://redis.io/topics/protocol](https://redis.io/topics/protocol)



# Commands

## General
* Commands are case insensitive
* Item keys are case sensitive
* Item keys cannot contain any spaces or newlines (bulk string is not supported at this time for keys)



## UPDATE
Adds an item with given score to the queue or updates the item with the given score. 

    UPDATE item_key [score]\r\n

Args:

        item_key  - string identifier for the item in the queue
        score     - integer score to queue an item with (optional)

Returns an integer response on success with the item's current score (this may change or be dependent on exact queue implementation):

    UPDATE foo 123\r\n
    :123\r\n
    UPDATE foo -12\r\n
    :111\r\n


## NEXT
Removes the next item from the queue and returns it

    NEXT\r\n

Returns string response prefixed with a `+` or a nil string identifier of `$-1\r\n` if the queue is empty

    UPDATE foo 123\r\n
    :123\r\n
    NEXT\r\n
    +foo\r\n
    NEXT\r\n
    $-1\r\n
    


## PEEK
Fetch the next item off the queue without removing it

    PEEK item_key\r\n

Args:

        item_key  - string identifier for the item in the queue
    
Returns string response prefixed with a `+` or a nil string identifier of `$-1\r\n` if the queue is empty

    UPDATE foo 123\r\n
    :123\r\n
    PEEK\r\n
    +foo\r\n
    NEXT\r\n
    +foo\r\n
    PEEK\r\n
    $-1\r\n


## INFO
Fetches info on the current running server.

    INFO\r\n

Returns an array response that represents a map of stat name to value

    INFO\r\n
    *6\r\n
    $5\r\n
    pools\r\n
    :3\r\n
    $7\r\n
    updates\r\n
    :11\r\n
    $9\r\n
    connected\r\n
    :1\r\n
    $11\r\n
    connections\r\n
    :1\r\n
    $3\r\n
    ops\r\n
    :21\r\n
    $5\r\n
    items\r\n
    :2\r\n

The above response represents an equivalent json dictionary with the following value

    {
        "pools":        3
        "updates":      11
        "connected":    1
        "connections":  1
        "ops":          21
        "items":        2
    }

Stats currently collected/exposed
* pools - number of score buckets currently in the queue
* updates - number of update operations that have come into the queue
* connected - number of currently connected sockets
* connections - number of connections made during the full lifetime of the process
* items - number of items in the queue






