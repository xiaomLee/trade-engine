package hub

import (
	"testing"
	"time"
)

func TestNewMemoryHub(t *testing.T) {
	h := NewMemoryHub()
	h.Start()

	// add 1
	if resp, err := h.DoTask(NewTask(TASKTYPEADD, &TaskAdd{Value: "1"})); err != nil {
		t.Error(err)
	} else {
		t.Logf("add response:%s", resp)
	}

	// get all
	if resp, err := h.DoTask(NewTask(TASKTYPEGETALL, &TaskGetAll{})); err != nil {
		t.Error(err)
	} else {
		t.Logf("get all response:%s", resp)
	}

	// add 2
	if resp, err := h.DoTask(NewTask(TASKTYPEADD, &TaskAdd{Value: "2"})); err != nil {
		t.Error(err)
	} else {
		t.Logf("add response:%s", resp)
	}

	// get all
	if resp, err := h.DoTask(NewTask(TASKTYPEGETALL, &TaskGetAll{})); err != nil {
		t.Error(err)
	} else {
		t.Logf("get all response:%s", resp)
	}

}

func TestNewRaftHubSingle(t *testing.T) {
	h, err := NewRaftHub("./", "127.0.0.1:1234", "test1", true)
	if err != nil {
		t.Error(err)
	}
	println("init success")
	h.Start()
	println("start hub success")

	time.Sleep(10 * time.Second)
	println(h.LeaderAddress())

	// add 1
	if resp, err := h.DoTask(NewTask(TASKTYPEADD, &TaskAdd{Value: "1"})); err != nil {
		t.Error(err)
	} else {
		t.Logf("add response:%s", resp)
	}

	// get all
	if resp, err := h.DoTask(NewTask(TASKTYPEGETALL, &TaskGetAll{})); err != nil {
		t.Error(err)
	} else {
		t.Logf("get all response:%s", resp)
	}

	// add 2
	if resp, err := h.DoTask(NewTask(TASKTYPEADD, &TaskAdd{Value: "2"})); err != nil {
		t.Error(err)
	} else {
		t.Logf("add response:%s", resp)
	}

	// get all
	if resp, err := h.DoTask(NewTask(TASKTYPEGETALL, &TaskGetAll{})); err != nil {
		t.Error(err)
	} else {
		t.Logf("get all response:%s", resp)
	}
}

func TestNewRaftHubCluster(t *testing.T) {
	h1, err := NewRaftHub("./1/", "127.0.0.1:1234", "1", true)
	if err != nil {
		t.Error(err)
	}
	h1.Start()
	time.Sleep(time.Second * 5)

	h2, err := NewRaftHub("./2/", "127.0.0.1:1235", "2", false)
	if err != nil {
		t.Error(err)
	}
	h2.Start()

	// add follower
	//if resp, err := h1.DoTask(NewTask(TASKTYPEADDFOLLOWER, &TaskAddFollower{ServerId: "2", Addr: "127.0.0.1:1235"})); err != nil {
	//	t.Error(err)
	//} else {
	//	t.Logf("add response:%s", resp)
	//}
	//time.Sleep(time.Second * 3)
	println("h2 leader:", h2.LeaderAddress())

	h3, err := NewRaftHub("./3/", "127.0.0.1:1236", "3", false)
	if err != nil {
		t.Error(err)
	}
	h3.Start()

	// add follower
	//if resp, err := h1.DoTask(NewTask(TASKTYPEADDFOLLOWER, &TaskAddFollower{ServerId: "3", Addr: "127.0.0.1:1236"})); err != nil {
	//	t.Error(err)
	//} else {
	//	t.Logf("add response:%s", resp)
	//}
	println("h3 leader:", h3.LeaderAddress())

	go func() {
		time.Sleep(3 * time.Second)
		h1.Close()
	}()

	select {}

}
