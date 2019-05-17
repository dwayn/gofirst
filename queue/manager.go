package queue

import (
	"fmt"

	"github.com/dwayn/gofirst/command"
)

func RunQueue(commandQueue chan command.Request, internalQueue Queue) {

	for {
		c := <-commandQueue
		var reply command.Response
		switch c.OpType {
		case command.Update:
			newScore, err := internalQueue.Update(c.ItemKey, c.ItemValue)
			switch err {
			case Ok:
				reply.ErrorCode = command.Ok
				reply.ResponseBody = fmt.Sprintf("%d", newScore)

			default:
				reply.ErrorCode = command.ErrUnknownError
			}
		case command.Next:
			val, err := internalQueue.GetNext()
			switch err {
			case Ok:
				reply.ErrorCode = command.Ok
				reply.ResponseBody = val
			case QueueEmpty:
				reply.ErrorCode = command.ErrQueueEmpty
			default:
				reply.ErrorCode = command.ErrUnknownError
			}
		case command.Peek:
			val, err := internalQueue.PeekNext()
			switch err {
			case Ok:
				reply.ErrorCode = command.Ok
				reply.ResponseBody = val
			case QueueEmpty:
				reply.ErrorCode = command.ErrQueueEmpty
			default:
				reply.ErrorCode = command.ErrUnknownError
			}
		case command.Score:
			val, err := internalQueue.GetScore(c.ItemKey)
			switch err {
			case Ok:
				reply.ErrorCode = command.Ok
				reply.ResponseBody = fmt.Sprintf("%d", val)
			case NotFound:
				reply.ErrorCode = command.ErrNotFound
			default:
				reply.ErrorCode = command.ErrUnknownError
			}
		default:
			reply.ErrorCode = command.ErrUnknownCommand
			reply.ErrorMessage = fmt.Sprintf("ERROR Unknown command: %s", c.OpType)
		}
		c.ResponseChannel <- reply
	}
}

func CreateChannel() chan command.Request {
	queueChannel := make(chan command.Request)
	return queueChannel
}
