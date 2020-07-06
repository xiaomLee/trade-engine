package hub

import "github.com/xiaomLee/trade-engine/entrust"

const (
	TASKTYPEADD         = 1
	TASKTYPEREMOVE      = 2
	TASKTYPEGETALL      = 3
	TASKTYPEADDFOLLOWER = 4
)

type Task struct {
	errChan  chan error
	respChan chan interface{}
	TaskType uint8
	Detail   interface{}
}

func NewTask(taskType uint8, detail interface{}) *Task {
	return &Task{
		errChan:  make(chan error),
		respChan: make(chan interface{}),
		TaskType: taskType,
		Detail:   detail,
	}
}

func (t *Task) ErrResponse(err error) {
	t.errChan <- err
}

func (t *Task) SuccessResponse(v interface{}) {
	t.errChan <- nil
	t.respChan <- v
}

func (t *Task) End() (interface{}, error) {
	e := <-t.errChan
	if e != nil {
		return nil, e
	}
	return <-t.respChan, nil
}

type TaskBuyEntrust struct {
	Entrust *entrust.Entrust
}

type TaskGetBuyList struct {
	CoinType string
}

type TaskRemove struct {
	Index int
}

type TaskGetAll struct {
}

type TaskAddFollower struct {
	ServerId string
	Addr     string
}
