package entrust

import (
	"fmt"
	"testing"

	"github.com/xiaomLee/trade-engine/entrust/mcache"
)

func TestMcache(t *testing.T) {
	cache := mcache.NewCache()

	e := &Entrust{
		EntrustType: "1",
		CoinType:    "BTC/USDT",
		BidType:     1,
		UserId:      1111,
		Num:         1,
		Price:       100,
	}

	if err := cache.AddBuy(e.CoinType, e); err != nil {
		t.Error(err)
	}

	e2 := &Entrust{
		EntrustType: "1",
		CoinType:    "BTC/USDT",
		BidType:     1,
		UserId:      222,
		Num:         1,
		Price:       102,
	}
	if err := cache.AddBuy(e.CoinType, e2); err != nil {
		t.Error(err)
	}

	list := cache.GetBuyList(e.CoinType)
	for i, l := range list {
		fmt.Printf("index:%d value: %+v \n", i, l)
	}
}
