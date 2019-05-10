package stats

const (
	Get = 0
	Set = 1
	Add = 2
	Sub = 3
)

type Metric struct {
	Metric string
	Op     int
	// value not used for "Get"
	Value int
	// channel only required for "Get" op
	Resp chan Value
}

type Value struct {
	Metric string
	Value  int
}

var metrics map[string]int

func ProcessMetrics(requests chan Metric) {
	metrics = make(map[string]int)
	for {
		request := <-requests
		switch request.Op {
		case Get:
			if request.Metric == "*" {
				for k, v := range metrics {
					request.Resp <- Value{k, v}
				}
				request.Resp <- Value{"", 0}
			} else {
				request.Resp <- Value{request.Metric, metrics[request.Metric]}
			}
		case Set:
			metrics[request.Metric] = request.Value
		case Add:
			metrics[request.Metric] = metrics[request.Metric] + request.Value
		case Sub:
			metrics[request.Metric] = metrics[request.Metric] - request.Value

		}
	}

}

func CreateMetricChannel() chan Metric {
	return make(chan Metric, 1000)
}

func CreateReplyChannel() chan Value {
	// this buffer should be large enough to hold a full set of metrics
	return make(chan Value, 10)
}
