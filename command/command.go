package command

const (
	Update = "update"
	Peek   = "peek"
	Score  = "score"
	Next   = "next"
	Info   = "info"
)

const (
	Ok                = 0
	ErrUnknownCommand = 1
	ErrQueueEmpty     = 2
	ErrNotFound       = 3
	ErrUnknownError   = 4
)

type Response struct {
	ErrorCode    int
	ErrorMessage string
	ResponseBody string
}

type Request struct {
	OpType          string
	ItemKey         string
	ItemValue       int
	ResponseChannel chan Response
}
