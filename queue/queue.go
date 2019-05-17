package queue

const (
	Ok         = 0
	QueueEmpty = 1
	NotFound   = 2
)

type Queue interface {
	PeekNext() (string, int)
	GetNext() (string, int)
	GetScore(itemId string) (int, int)
	Update(itemId string, score int) (int, int)
}
