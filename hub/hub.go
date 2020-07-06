package hub

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/xiaomLee/trade-engine/entrust"
	"github.com/xiaomLee/trade-engine/entrust/mcache"
	"github.com/xiaomLee/trade-engine/entrust/raftcache"
)

type Hub struct {
	requestChan   chan *Task
	leaderFlag    *uint32
	leaderAddress string
	cache         entrust.Cache
	exit          chan bool
	sync.RWMutex
}

func (h *Hub) workLoop() {
	for {
		select {
		case req := <-h.requestChan:
			if !h.IsLeader() {
				h.redirectToLeader(req)
				break
			}
			h.handleRequest(req)
		case <-h.exit:
			return
		}
	}
}

func (h *Hub) redirectToLeader(t *Task) {
	if h.leaderAddress == "" {
		t.ErrResponse(errors.New("no leader"))
	}
	go func() {
		t.SuccessResponse("ok")
	}()
}

// 串行
func (h *Hub) handleRequest(t *Task) {
	var resp interface{}
	var err error
	switch t.TaskType {
	case TASKTYPEADD:
		err = h.doBuyEntrust(t)
	case TASKTYPEREMOVE:
		err = h.doBuyEntrust(t)
	case TASKTYPEGETALL:
		resp = h.doGetBuyList(t)
	case TASKTYPEADDFOLLOWER:
		detail := t.Detail.(*TaskAddFollower)
		err = h.cache.AddFollower(detail.ServerId, detail.Addr)
	default:
		err = errors.New("invalid task type")
	}

	if err != nil {
		t.ErrResponse(err)
		return
	}
	if resp == nil {
		resp = "success"
	}
	t.SuccessResponse(resp)
}

func (h *Hub) DoTask(t *Task) (interface{}, error) {
	// no leader
	if !h.IsLeader() && h.LeaderAddress() == "" {
		return nil, errors.New("service busy: no leader")
	}

	select {
	case <-h.exit:
		return nil, errors.New("service busy")

	case h.requestChan <- t:
		return t.End()

	case <-time.After(3 * time.Second):
		return nil, errors.New("service busy")

	}
}

func (h *Hub) Start() error {
	go func() {
		h.workLoop()
	}()
	return nil
}

func (h *Hub) Close() {
	if err := h.cache.Close(); err != nil {
		println(err)
	}
	close(h.exit)
}

func NewMemoryHub() *Hub {
	var leader uint32
	h := &Hub{
		requestChan: make(chan *Task, 1024),
		leaderFlag:  &leader,
		cache:       mcache.NewCache(),
		exit:        make(chan bool),
	}
	h.setLeader(true)
	h.setLeaderAddress("127.0.0.1:0")
	return h
}

func NewRaftHub(dir, bind, serverId string, bootstrap bool) (*Hub, error) {
	var leader uint32
	notifyCh := make(chan bool, 1)
	cache := raftcache.NewRaftCache(dir, bind)
	h := &Hub{
		requestChan: make(chan *Task, 1024),
		leaderFlag:  &leader,
		cache:       cache,
		exit:        make(chan bool),
	}
	h.setLeader(true)

	err := cache.Start(bootstrap, serverId, notifyCh)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			select {
			case leader := <-notifyCh:
				println("leader changes:" + cache.Leader())
				h.setLeader(leader)
				h.setLeaderAddress(cache.Leader())
			case <-h.exit:
				println("hub close")
				return
			}
		}
	}()

	return h, nil
}

func (h *Hub) IsLeader() bool {
	return atomic.LoadUint32(h.leaderFlag) == 1
}

func (h *Hub) setLeader(leader bool) {
	if leader {
		atomic.SwapUint32(h.leaderFlag, 1)
		return
	}
	atomic.SwapUint32(h.leaderFlag, 0)
}

func (h *Hub) LeaderAddress() string {
	h.RLock()
	defer h.RUnlock()
	return h.leaderAddress
}

func (h *Hub) setLeaderAddress(addr string) {
	h.Lock()
	defer h.Unlock()
	h.leaderAddress = addr
}
