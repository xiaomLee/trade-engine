package raftcache

import (
	"log"

	"github.com/xiaomLee/trade-engine/entrust/mcache"

	"github.com/hashicorp/raft"
)

type CacheFSMSnapShot struct {
	cache  *mcache.Cache
	logger *log.Logger
}

// Persist data in specific type
// kv item serialize in google protubuf
func (f *CacheFSMSnapShot) Persist(sink raft.SnapshotSink) error {
	f.logger.Printf("Persist action in fsmSnapshot")
	defer sink.Close()

	var i int
	data := f.cache.Data()
	for i = 0; i < len(data); i++ {
		item := data[i]
		if _, err := sink.Write([]byte(item)); err != nil {
			return err
		}
	}
	f.logger.Printf("Persist total %d keys", i)

	return nil
}

func (f *CacheFSMSnapShot) Release() {
	f.logger.Printf("Release action in fsmSnapshot")
}
