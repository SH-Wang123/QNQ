package common

const (
	qPoolReady = iota
	qPoolRunning
	qPoolStopping
	qPoolStop
)

var globalQPool qPool

type qPool struct {
	status    int
	sub       chan Task
	scheduler qPoolScheduler
	actuators []qPoolActuator
}

// actuator 执行器
type qPoolActuator struct {
}

type qPoolScheduler struct {
}

func (q *qPool) Start() {
	//初始化执行器

	//初始化调度器

	q.status = qPoolRunning
}

func (q *qPool) Shutdown() {
	q.status = qPoolStopping

	q.status = qPoolStop
}

func (q *qPool) Submit() {

}

func createNewQPool() {
	globalQPool = qPool{
		status: qPoolReady,
	}
}
