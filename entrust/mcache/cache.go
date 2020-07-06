package mcache

import (
	"errors"
	"sync"

	"github.com/xiaomLee/trade-engine/entrust/queue"
)

type Cache struct {
	buyQueue  map[string]*queue.Queue
	sellQueue map[string]*queue.Queue
	sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{
		buyQueue:  make(map[string]*queue.Queue),
		sellQueue: make(map[string]*queue.Queue),
		RWMutex:   sync.RWMutex{},
	}
}

func (c *Cache) GetBuyList(coinType string) []queue.Item {
	list := make([]queue.Item, 0)
	c.RLock()
	defer c.RUnlock()
	q, ok := c.buyQueue[coinType]
	if !ok {
		return list
	}
	buckets := q.Buckets()
	for i := 0; i < len(buckets); i++ {
		list = append(list, buckets[i].Items()...)
	}
	return list
}

func (c *Cache) GetSellList(coinType string) []queue.Item {
	list := make([]queue.Item, 0)
	c.RLock()
	defer c.RUnlock()
	q, ok := c.buyQueue[coinType]
	if !ok {
		return list
	}
	buckets := q.Buckets()
	for i := 0; i < len(buckets); i++ {
		copy(list, buckets[i].Items())
	}
	return list
}

func (c *Cache) AddBuy(coinType string, item queue.Item) error {
	c.Lock()
	defer c.Unlock()
	if _, ok := c.buyQueue[coinType]; !ok {
		c.buyQueue[coinType] = &queue.Queue{Sort: queue.SortDesc}
	}

	return c.buyQueue[coinType].AddItem(item)
}

func (c *Cache) AddSell(coinType string, item queue.Item) error {
	c.Lock()
	defer c.Unlock()
	if _, ok := c.buyQueue[coinType]; !ok {
		c.sellQueue[coinType] = &queue.Queue{Sort: queue.SortAsc}
	}

	return c.sellQueue[coinType].AddItem(item)
}

func (c *Cache) RemoveBuy(coinType string, index int) error {
	c.RLock()
	defer c.RUnlock()
	q, ok := c.buyQueue[coinType]
	if !ok {
		return errors.New("no found queue " + coinType)
	}

	return q.Remove(index)
}

func (c *Cache) RemoveSell(coinType string, index int) error {
	c.RLock()
	defer c.RUnlock()
	q, ok := c.sellQueue[coinType]
	if !ok {
		return errors.New("no found queue " + coinType)
	}

	return q.Remove(index)
}

func (c *Cache) AddFollower(serverId, addr string) error {
	return nil
}

func (c *Cache) Close() error {
	return nil
}
