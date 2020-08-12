package queue

import (
	"errors"
)

const (
	SortAsc           = SortType("asc")
	SortDesc          = SortType("desc")
	defaultBucketSize = 3
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
	// this bucket`s last item
	tail Item
	// record if item in this bucket, for quick find and avoid deduplication.
	// key is item`s unique key.
	// notice: value is not the index, it`s just 1.
	dataMap map[string]uint8
	items   []Item
	n       int
}

type Item interface {
	Key() string
	// this > item return 1
	// this = item return 0
	// this < item return -1
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
		return q.addItemAsc(item)
	default:
		return errors.New("queue sort type error " + string(q.Sort))
	}
}

func (q *Queue) FindPosition(item Item) (bucketPos, dataPos int) {
	if len(q.buckets) == 0 {
		return -1, -1
	}
	for index, b := range q.buckets {
		bucketPos = index
		if _, ok := b.dataMap[item.Key()]; ok {
			break
		}
		if index == len(q.buckets)-1 {
			bucketPos = -1
		}
	}
	if bucketPos == -1 {
		return -1, -1
	}

	for index, data := range q.buckets[bucketPos].items {
		dataPos = index
		if data.Key() == item.Key() {
			break
		}
		if index == len(q.buckets[bucketPos].items)-1 {
			panic("queue data err. item find in data map, but not find in items")
		}
	}

	return bucketPos, dataPos
}

func (q *Queue) Remove(item Item) error {
	i, j := q.FindPosition(item)
	if i < 0 {
		return nil
	}

	b := q.buckets[i]
	// this bucket has only one item, remove this bucket
	if b.n == 1 {
		dst := make([]*bucket, i)
		copy(dst, q.buckets[:i])
		q.buckets = append(dst, q.buckets[i+1:]...)
	}

	dst := make([]Item, j)
	copy(dst, b.items[:j])
	b.items = append(dst, b.items[j+1:]...)
	delete(b.dataMap, item.Key())
	b.n--

	return nil
}

func (q *Queue) addItemDesc(item Item) error {
	if q.bucketSize == 0 {
		q.bucketSize = defaultBucketSize
	}
	if q.buckets == nil {
		q.buckets = make([]*bucket, 0)
	}
	if len(q.buckets) == 0 {
		b := &bucket{
			dataMap: make(map[string]uint8),
			items:   make([]Item, 0),
			n:       0,
		}
		q.buckets = append(q.buckets, b)
	}

	// find bucket
	var i int
	for index, b := range q.buckets {
		i = index
		if b.n == 0 {
			break
		}
		ret := item.Compare(b.tail)
		if ret > 0 {
			break
		}
		if ret == 0 && b.n < q.bucketSize {
			break
		}
	}

	// current bucket is null, just append and return
	if q.buckets[i].n == 0 {
		q.buckets[i].tail = item
		q.buckets[i].items = append(q.buckets[i].items, item)
		q.buckets[i].dataMap[item.Key()] = 1
		q.buckets[i].n++
		return nil
	}

	// compare and insert into the right index of items
	var j int
	for index, data := range q.buckets[i].items {
		j = index
		ret := item.Compare(data)
		if ret > 0 {
			break
		}
	}

	if j == q.buckets[i].n-1 && item.Compare(q.buckets[i].items[j]) <= 0 {
		q.buckets[i].tail = item
		q.buckets[i].items = append(q.buckets[i].items, item)
	} else {
		old := q.buckets[i].items[j:]

		dst := make([]Item, q.buckets[i].n-len(old))
		copy(dst, q.buckets[i].items[:len(dst)])

		q.buckets[i].items = dst
		q.buckets[i].items = append(q.buckets[i].items, item)
		q.buckets[i].items = append(q.buckets[i].items, old...)
	}
	q.buckets[i].dataMap[item.Key()] = 1
	q.buckets[i].n++

	if q.buckets[i].n > q.bucketSize {
		b := q.buckets[i]
		b1 := &bucket{
			dataMap: make(map[string]uint8),
			items:   make([]Item, 0),
			n:       0,
		}
		b2 := &bucket{
			dataMap: make(map[string]uint8),
			items:   make([]Item, 0),
			n:       0,
		}

		// copy data
		for index, data := range b.items {
			if index < b.n/2 {
				b1.dataMap[data.Key()] = 1
				b1.items = append(b1.items, data)
				b1.n++
				if index == b.n/2-1 {
					b1.tail = data
				}
				continue
			}

			b2.dataMap[data.Key()] = 1
			b2.items = append(b2.items, data)
			b2.n++
			if index == b.n-1 {
				b2.tail = data
			}
		}

		if i == len(q.buckets)-1 {
			q.buckets[i] = b1
			q.buckets = append(q.buckets, b2)
			return nil
		}

		oldBuckets := q.buckets[i+1:]
		dst := make([]*bucket, i)
		copy(dst, q.buckets)

		q.buckets = append(dst, b1)
		q.buckets = append(q.buckets, b2)
		q.buckets = append(q.buckets, oldBuckets...)
	}

	return nil
}

func (q *Queue) addItemAsc(item Item) error {
	if q.bucketSize == 0 {
		q.bucketSize = defaultBucketSize
	}
	if q.buckets == nil {
		q.buckets = make([]*bucket, 0)
	}
	if len(q.buckets) == 0 {
		b := &bucket{
			dataMap: make(map[string]uint8),
			items:   make([]Item, 0),
			n:       0,
		}
		q.buckets = append(q.buckets, b)
	}

	// find bucket
	var i int
	for index, b := range q.buckets {
		i = index
		if b.n == 0 {
			break
		}
		ret := item.Compare(b.tail)
		if ret < 0 {
			break
		}
		if ret == 0 && b.n < q.bucketSize {
			break
		}
	}

	// current bucket is null, just append and return
	if q.buckets[i].n == 0 {
		q.buckets[i].tail = item
		q.buckets[i].items = append(q.buckets[i].items, item)
		q.buckets[i].dataMap[item.Key()] = 1
		q.buckets[i].n++
		return nil
	}

	// compare and insert into the right index of items
	var j int
	for index, data := range q.buckets[i].items {
		j = index
		ret := item.Compare(data)
		if ret < 0 {
			break
		}
	}

	if j == q.buckets[i].n-1 && item.Compare(q.buckets[i].items[j]) >= 0 {
		q.buckets[i].tail = item
		q.buckets[i].items = append(q.buckets[i].items, item)
	} else {
		old := q.buckets[i].items[j:]

		dst := make([]Item, q.buckets[i].n-len(old))
		copy(dst, q.buckets[i].items[:len(dst)])

		q.buckets[i].items = dst
		q.buckets[i].items = append(q.buckets[i].items, item)
		q.buckets[i].items = append(q.buckets[i].items, old...)
	}
	q.buckets[i].dataMap[item.Key()] = 1
	q.buckets[i].n++

	if q.buckets[i].n > q.bucketSize {
		b := q.buckets[i]
		b1 := &bucket{
			dataMap: make(map[string]uint8),
			items:   make([]Item, 0),
			n:       0,
		}
		b2 := &bucket{
			dataMap: make(map[string]uint8),
			items:   make([]Item, 0),
			n:       0,
		}

		// copy data
		for index, data := range b.items {
			if index < b.n/2 {
				b1.dataMap[data.Key()] = 1
				b1.items = append(b1.items, data)
				b1.n++
				if index == b.n/2-1 {
					b1.tail = data
				}
				continue
			}

			b2.dataMap[data.Key()] = 1
			b2.items = append(b2.items, data)
			b2.n++
			if index == b.n-1 {
				b2.tail = data
			}
		}

		if i == len(q.buckets)-1 {
			q.buckets[i] = b1
			q.buckets = append(q.buckets, b2)
			return nil
		}

		oldBuckets := q.buckets[i+1:]
		dst := make([]*bucket, i)
		copy(dst, q.buckets)

		q.buckets = append(dst, b1)
		q.buckets = append(q.buckets, b2)
		q.buckets = append(q.buckets, oldBuckets...)
	}

	return nil
}
