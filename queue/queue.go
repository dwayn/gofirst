package queue

const (
	Ok         = 0
	QueueEmpty = 1
	NotFound   = 2
	AllocError = 3 // not sure this can actually happen in go memory management, but leaving for now until sure
)

type Queue interface {
	PeekNext() (string, int)
	GetNext() (string, int)
	GetScore(itemId string) (int, int)
	Update(itemId string, score int) (int, int)
}
