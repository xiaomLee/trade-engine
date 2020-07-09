package entrust

import "github.com/xiaomLee/trade-engine/entrust/queue"

type Entrust struct {
	Id          string
	EntrustType string
	CoinType    string
	BidType     uint8 // 买卖方向
	UserId      int64
	Num         float64
	Price       float64
}

func (e *Entrust) Key() string {
	return e.Id
}

func (e *Entrust) Compare(item queue.Item) int {
	bf := item.(*Entrust)
	if e.Price == bf.Price {
		return 0
	}
	if e.Price < bf.Price {
		return -1
	}

	return 1
}
