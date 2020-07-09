package mcache

import (
	"testing"

	"github.com/xiaomLee/trade-engine/entrust/queue"
)

type TestStruct struct {
	id        string
	price     float64
	timeStamp int64
}

func (v TestStruct) Key() string {
	return v.id
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
		id:        "1",
		price:     1,
		timeStamp: 111,
	}
	item2 := TestStruct{
		id:        "2",
		price:     1,
		timeStamp: 222,
	}
	item3 := TestStruct{
		id:        "3",
		price:     3,
		timeStamp: 333,
	}
	item4 := TestStruct{
		id:        "4",
		price:     4,
		timeStamp: 444,
	}
	item5 := TestStruct{
		id:        "5",
		price:     0.5,
		timeStamp: 555,
	}
	item6 := TestStruct{
		id:        "6",
		price:     1,
		timeStamp: 666,
	}
	item7 := TestStruct{
		id:        "7",
		price:     2,
		timeStamp: 777,
	}
	item8 := TestStruct{
		id:        "8",
		price:     2,
		timeStamp: 888,
	}
	if err := cache.AddBuy(coinType, item5); err != nil {
		t.Error(err)
	}
	t.Logf("queue list:%v \n", cache.GetBuyList(coinType))
	if err := cache.AddBuy(coinType, item3); err != nil {
		t.Error(err)
	}
	t.Logf("queue list:%v \n", cache.GetBuyList(coinType))
	if err := cache.AddBuy(coinType, item4); err != nil {
		t.Error(err)
	}
	t.Logf("queue list:%v \n", cache.GetBuyList(coinType))
	if err := cache.AddBuy(coinType, item1); err != nil {
		t.Error(err)
	}
	t.Logf("queue list:%v \n", cache.GetBuyList(coinType))
	if err := cache.AddBuy(coinType, item2); err != nil {
		t.Error(err)
	}
	t.Logf("queue list:%v \n", cache.GetBuyList(coinType))
	if err := cache.AddBuy(coinType, item6); err != nil {
		t.Error(err)
	}
	t.Logf("queue list:%v \n", cache.GetBuyList(coinType))
	if err := cache.AddBuy(coinType, item7); err != nil {
		t.Error(err)
	}
	t.Logf("queue list:%v \n", cache.GetBuyList(coinType))
	if err := cache.AddBuy(coinType, item8); err != nil {
		t.Error(err)
	}
	t.Logf("queue list:%v \n", cache.GetBuyList(coinType))

	ret := cache.GetBuyList(coinType)
	//for i, r := range ret {
	//	t.Logf("index:%d value:%+v \n", i, r.(TestStruct))
	//}

	if len(ret) != 6 {
		t.Fatalf("length is not 6")
	}

	if ret[0].(TestStruct).price != 3 {
		t.Fatalf("ret first is not 1 %v", ret[0].(TestStruct))
	}
	if v := ret[1].(TestStruct); v.price != 1 && v.id != "1" {
		t.Fatalf("ret second is not 1 %v", ret[1].(TestStruct))
	}
	if v := ret[2].(TestStruct); v.price != 2 && v.id != "2" {
		t.Fatalf("ret thrid is not 1 %v", ret[2].(TestStruct))
	}
}
