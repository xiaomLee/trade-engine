package raftcache

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"

	"github.com/xiaomLee/trade-engine/entrust/queue"

	"github.com/xiaomLee/trade-engine/entrust/mcache"

	"github.com/hashicorp/raft"
)

type CacheFSM struct {
	cache  *mcache.Cache
	logger *log.Logger
}

type ApplyCommand struct {
	Op       string
	CoinType string
	Msg      interface{}
}

func NewCacheFSM() *CacheFSM {
	return &CacheFSM{
		cache:  mcache.NewCache(),
		logger: log.New(os.Stderr, "[CacheFSM]", log.LstdFlags),
	}
}

func (f *CacheFSM) GetBuyList(coinType string) []queue.Item {
	return f.cache.GetBuyList(coinType)
}

func (f *CacheFSM) GetSellList(coinType string) []queue.Item {
	return f.cache.GetSellList(coinType)
}

// apply the add & remove from app mcache
func (f *CacheFSM) Apply(l *raft.Log) interface{} {
	var (
		cmd ApplyCommand
	)
	if err := json.Unmarshal(l.Data, &cmd); err != nil {
		f.logger.Println(err)
		return err
	}
	switch cmd.Op {
	case "AddBuy":
		return f.applyAddBuy(cmd)
	case "AddSell":
		return f.applyAddSell(cmd)
	case "RemoveBuy":
		return f.applyRemoveBuy(cmd)
	case "RemoveSell":
		return f.applyRemoveSell(cmd)
	}
	return errors.New("op type invalid:" + cmd.Op)
}

func (f *CacheFSM) Snapshot() (raft.FSMSnapshot, error) {
	f.logger.Printf("Generate FSMSnapshot")
	return &CacheFSMSnapShot{
		cache:  f.cache,
		logger: log.New(os.Stderr, "[fsmSnapshot] ", log.LstdFlags),
	}, nil
}

func (f *CacheFSM) Restore(rc io.ReadCloser) error {

	return nil
}

func (f *CacheFSM) applyAddBuy(cmd ApplyCommand) interface{} {
	if cmd.Op != "AddBuy" {
		return errors.New("op type invalid:" + cmd.Op)
	}

	return f.cache.AddBuy(cmd.CoinType, cmd.Msg.(queue.Item))
}

func (f *CacheFSM) applyAddSell(cmd ApplyCommand) interface{} {
	if cmd.Op != "AddSell" {
		return errors.New("op type invalid:" + cmd.Op)
	}

	return f.cache.AddSell(cmd.CoinType, cmd.Msg.(queue.Item))
}

func (f *CacheFSM) applyRemoveBuy(cmd ApplyCommand) interface{} {
	if cmd.Op != "RemoveBuy" {
		return errors.New("op type invalid:" + cmd.Op)
	}

	return f.cache.RemoveBuy(cmd.CoinType, cmd.Msg.(int))
}

func (f *CacheFSM) applyRemoveSell(cmd ApplyCommand) interface{} {
	if cmd.Op != "RemoveSell" {
		return errors.New("op type invalid:" + cmd.Op)
	}

	return f.cache.RemoveSell(cmd.CoinType, cmd.Msg.(int))
}
