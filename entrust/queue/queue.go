package queue

import "errors"

const (
	SortAsc           = SortType("asc")
	SortDesc          = SortType("desc")
	defaultBucketSize = 200
)

type SortType string

// Segment Queue
// desc. buckets[0].max is the max of all buy buckets
// eg. buckets[0].min > buckets[1].max
// asc. buckets[0].min is the min of all sell buckets
// eg. buckets[0].max < buckets[1].min
type Queue struct {
	Sort       SortType
	bucketSize int
	buckets    []*bucket
}

// desc. if compare return equal, insert into the last equal item tail.
// asc. if compare equal, order by time asc
// n > defaultBucketSize will expansion.
type bucket struct {
	max   Item
	min   Item
	items []Item
	n     int
}

type Item interface {
	Compare(Item) int
}

func (b *bucket) Items() []Item {
	return b.items
}

// call immediately after the Queue init is recommend
func (q *Queue) SetBucketSize(size int) {
	if len(q.buckets) > 0 {
		return
	}
	q.bucketSize = size
}

func (q *Queue) Buckets() []*bucket {
	return q.buckets
}

func (q *Queue) AddItem(item Item) error {
	switch q.Sort {
	case SortDesc:
		return q.addItemDesc(item)
	case SortAsc:
		return q.addItemDesc(item)
	default:
		return errors.New("queue sort type error " + string(q.Sort))
	}
}

func (q *Queue) Remove(index int) error {
	return nil
}

func (q *Queue) addItemDesc(item Item) error {
	if q.buckets == nil {
		q.buckets = make([]*bucket, 0)
	}
	if len(q.buckets) == 0 {
		es := &bucket{
			max:   item,
			min:   item,
			items: make([]Item, 0),
			n:     1,
		}
		es.items = append(es.items, item)
		q.buckets = append(q.buckets, es)
		return nil
	}

	// find bucket
	var i int
	for i = 0; i < len(q.buckets); i++ {
		ret := item.Compare(q.buckets[i].min)
		if ret <= 0 {
			break
		}
	}

	// todo compare and insert into the right index
	var j int
	for j = 0; j < len(q.buckets[i].items); j++ {
		ret := item.Compare(q.buckets[i].items[j])
		if ret <= 0 {
			break
		}
	}
	if j == 0 {
		q.buckets[i].max = item
		q.buckets[i].items = append([]Item{item}, q.buckets[i].items...)
	} else if j == len(q.buckets[i].items) {
		q.buckets[i].min = item
		q.buckets[i].items = append(q.buckets[i].items, item)
	} else {
		old := q.buckets[i].items[j+1 : len(q.buckets[i].items)]
		q.buckets[i].items = append(q.buckets[i].items[:j], item)
		q.buckets[i].items = append(q.buckets[i].items, old...)
	}

	return nil
}
