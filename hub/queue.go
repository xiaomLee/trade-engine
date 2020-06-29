package hub

import "errors"

const defaultQueueCap = 200

// 分段队列
// order by price desc. buckets[0].max is the max price of all buy buckets
// eg. buckets[0].min > buckets[1].max
type BuyQueue struct {
	Side    int8
	buckets []*bucket
}

// order by price asc. buckets[0].min is the min price of all sell buckets
// eg. buckets[0].max < buckets[1].min
type SellQueue struct {
	Side    int8
	buckets []*bucket
}

// 委托队列
// buy buckets orders order by price desc. if price equal, order by time asc
// sell buckets orders order by price asc. if price equal, order by time asc
// n > defaultQueueCap will expansion.
type bucket struct {
	max    float64
	min    float64
	orders []*Order
	n      int
}

func (q *BuyQueue) AddOrder(o *Order) error {
	if o.Side != BuyOrderSide {
		return errors.New("order side invalid")
	}
	if q.buckets == nil {
		q.buckets = make([]*bucket, 0)
	}
	if len(q.buckets) == 0 {
		es := &bucket{
			max:    o.Price,
			min:    o.Price,
			orders: make([]*Order, 0),
			n:      1,
		}
		es.orders[0] = o
		q.buckets[0] = es
		return nil
	}
	var i int
	for i = 0; i <= len(q.buckets); i++ {
		if q.buckets[i].min <= o.Price {
			break
		}
	}
	return nil
}
