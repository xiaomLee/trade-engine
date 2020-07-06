package raftcache

import (
	"fmt"
	"testing"
)

func TestRaftCache_Start(t *testing.T) {
	notify := make(chan bool, 1)

	cache := NewRaftCache("./", ":1234")
	cache.Start(true, "1", notify)
	select {
	case leader := <-notify:
		fmt.Println(leader)
		fmt.Println(cache.raft.Leader())
		if leader {
			cache.Add("a")
			cache.Add("b")
			cache.Add("c")

			values := cache.GetAll()
			fmt.Println(values)
			cache.Remove(0)
			cache.Remove(1)

		}
	}
}
