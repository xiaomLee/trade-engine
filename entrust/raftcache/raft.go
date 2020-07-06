package raftcache

import (
	"encoding/json"
	"errors"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/xiaomLee/trade-engine/entrust/queue"

	raftboltdb "github.com/hashicorp/raft-boltdb"

	"github.com/hashicorp/raft"
)

type RaftCache struct {
	RaftDir  string
	RaftBind string

	raft *raft.Raft // The consensus mechanism
	fsm  *CacheFSM

	logger *log.Logger
}

func NewRaftCache(dir, bind string) *RaftCache {
	fsm := NewCacheFSM()
	return &RaftCache{
		RaftDir:  dir,
		RaftBind: bind,
		fsm:      fsm,
		logger:   log.New(os.Stderr, "[RaftCache]", log.LstdFlags),
	}
}

// NotifyCh is used to provide a channel that will be notified of leadership changes.
func (c *RaftCache) Start(bootstrap bool, localID string, notify chan bool) error {
	conf := raft.DefaultConfig()
	conf.LocalID = raft.ServerID(localID)
	conf.SnapshotInterval = 5 * time.Second
	conf.SnapshotThreshold = 2
	conf.NotifyCh = notify

	addr, err := net.ResolveTCPAddr("tcp", c.RaftBind)
	if err != nil {
		return err
	}

	transport, err := raft.NewTCPTransport(c.RaftBind, addr, 3, 10*time.Second, os.Stderr)
	if err != nil {
		return err
	}

	ss, err := raft.NewFileSnapshotStore(c.RaftDir, 2, os.Stderr)
	if err != nil {
		return err
	}

	// boltDB implement log store and stable store interface
	logStore, err := raftboltdb.NewBoltStore(filepath.Join(c.RaftDir, "raftcache-log.db"))
	if err != nil {
		return err
	}
	stableStore, err := raftboltdb.NewBoltStore(filepath.Join(c.RaftDir, "raftcache-stable.db"))
	if err != nil {
		return err
	}

	// raftcache system
	r, err := raft.NewRaft(conf, c.fsm, logStore, stableStore, ss, transport)
	if err != nil {
		return err
	}
	c.raft = r

	if bootstrap {
		configuration := raft.Configuration{
			Servers: []raft.Server{
				{
					ID:      conf.LocalID,
					Address: transport.LocalAddr(),
				},
			},
		}
		c.raft.BootstrapCluster(configuration)
	}
	return nil
}

func (c *RaftCache) GetBuyList(coinType string) []queue.Item {
	return c.fsm.GetBuyList(coinType)
}

func (c *RaftCache) GetSellList(coinType string) []queue.Item {
	return c.fsm.GetSellList(coinType)
}

func (c *RaftCache) AddBuy(coinType string, item queue.Item) error {
	cmd := ApplyCommand{
		Op:       "AddBuy",
		CoinType: coinType,
		Msg:      item,
	}

	return c.apply(cmd)
}

func (c *RaftCache) AddSell(coinType string, item queue.Item) error {
	cmd := ApplyCommand{
		Op:       "AddSell",
		CoinType: coinType,
		Msg:      item,
	}

	return c.apply(cmd)
}

func (c *RaftCache) RemoveBuy(coinType string, index int) error {
	cmd := ApplyCommand{
		Op:       "RemoveBuy",
		CoinType: coinType,
		Msg:      strconv.Itoa(index),
	}

	return c.apply(cmd)
}

func (c *RaftCache) RemoveSell(coinType string, index int) error {
	cmd := ApplyCommand{
		Op:       "RemoveSell",
		CoinType: coinType,
		Msg:      strconv.Itoa(index),
	}

	return c.apply(cmd)
}

func (c *RaftCache) apply(cmd ApplyCommand) error {
	if c.raft.State() != raft.Leader {
		c.logger.Println("not leader")
		return errors.New("not leader")
	}

	msg, err := json.Marshal(cmd)
	if err != nil {
		return err
	}

	return c.raft.Apply(msg, 10*time.Second).Error()
}

func (c *RaftCache) Leader() string {
	return string(c.raft.Leader())
}

func (c *RaftCache) AddFollower(serverId, addr string) error {
	future := c.raft.AddVoter(raft.ServerID(serverId), raft.ServerAddress(addr), 0, 0)
	if err := future.Error(); err != nil {
		return err
	}
	return nil
}

func (c *RaftCache) Close() error {
	future := c.raft.Shutdown()
	return future.Error()
}
