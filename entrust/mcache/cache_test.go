package mcache

import (
	"testing"

	"github.com/xiaomLee/trade-engine/entrust/queue"
)

type TestStruct struct {
	price     float64
	timeStamp int64
}

func (v TestStruct) Compare(item queue.Item) int {
	v2 := item.(TestStruct)
	if v.price == v2.price {
		return 0
	}
	if v.price < v2.price {
		return -1
	}
	return 1
}

func TestCache_AddBuy(t *testing.T) {
	coinType := "BTC/USDT"
	cache := NewCache()
	item1 := TestStruct{
		price:     1,
		timeStamp: 111,
	}
	item2 := TestStruct{
		price:     1,
		timeStamp: 222,
	}
	item3 := TestStruct{
		price:     3,
		timeStamp: 333,
	}
	if err := cache.AddBuy(coinType, item1); err != nil {
		t.Error(err)
	}
	if err := cache.AddBuy(coinType, item2); err != nil {
		t.Error(err)
	}
	if err := cache.AddBuy(coinType, item3); err != nil {
		t.Error(err)
	}

	ret := cache.GetBuyList(coinType)
	for i, r := range ret {
		t.Logf("index:%d value:%+v \n", i, r.(TestStruct))
	}

	if len(ret) != 3 {
		t.Fatalf("length is not 3")
	}
	if ret[0].(TestStruct).price != 1 {
		t.Fatalf("ret first is not 1 %v", ret[0].(TestStruct))
	}
	if ret[1].(TestStruct).price != 2 {
		t.Fatalf("ret second is not 2 %v", ret[1].(TestStruct))
	}
	if ret[2].(TestStruct).price != 3 {
		t.Fatalf("ret thrid is not 1 %v", ret[2].(TestStruct))
	}
}
