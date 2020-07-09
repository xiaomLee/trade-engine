package queue

import "testing"

type TestItem struct {
	Id    string
	Price float64
}

func (v TestItem) Key() string {
	return v.Id
}

func (v TestItem) Compare(item Item) int {
	v2 := item.(TestItem)
	if v.Price == v2.Price {
		return 0
	}
	if v.Price < v2.Price {
		return -1
	}
	return 1
}

func TestQueue_AddItem(t *testing.T) {
	descQueue := &Queue{Sort: SortDesc}
	descQueue.SetBucketSize(3)

	// queue list: 1
	insertData(t, descQueue, "1", 1, 0, 0)
	// queue list: 2, 1
	insertData(t, descQueue, "2", 2, 0, 0)
	// queue list: 2, 1, 3
	insertData(t, descQueue, "3", 1, 0, 2)
	// queue list: 2, 4; 1, 3
	insertData(t, descQueue, "4", 1.5, 0, 1)
	// queue list: 2, 4, 5; 1, 3
	insertData(t, descQueue, "5", 1.5, 0, 2)
	// queue list: 6, 2, 4, 5; 1, 3 --> 6, 2; 4, 5; 1, 3
	insertData(t, descQueue, "6", 3, 0, 0)

	ascQueue := &Queue{Sort: SortAsc}
	ascQueue.SetBucketSize(3)

	// queue list: 1
	insertData(t, ascQueue, "1", 1, 0, 0)
	// queue list: 1, 2
	insertData(t, ascQueue, "2", 2, 0, 1)
	// queue list: 1, 3, 2
	insertData(t, ascQueue, "3", 1, 0, 1)
	// queue list: 1, 3, 4, 2 --> 1, 3; 4, 2
	insertData(t, ascQueue, "4", 1.5, 1, 0)
	// queue list: 1, 3; 4, 5, 2
	insertData(t, ascQueue, "5", 1.5, 1, 1)
	// queue list: 1, 3; 4, 5; 2, 6
	insertData(t, ascQueue, "6", 3, 2, 1)
	// queue list: 1, 3; 4, 5, 7; 2, 6
	insertData(t, ascQueue, "7", 1.5, 1, 2)

}

func TestQueue_FindPosition(t *testing.T) {
	queue := &Queue{Sort: SortDesc}
	queue.SetBucketSize(3)

	insertData(t, queue, "1", 1, 0, 0)
	insertData(t, queue, "2", 2, 0, 0)
	insertData(t, queue, "3", 1, 0, 2)
	insertData(t, queue, "4", 1.5, 0, 1)

	item3 := TestItem{
		Id:    "3",
		Price: 1,
	}
	// after expand find item3 again
	i, j := queue.FindPosition(item3)
	t.Logf("bucketPos:%d, dataPos:%d \n", i, j)

	if i != 1 || j != 1 {
		t.Error("item position not expect, expect 1, 1")
	}
}

func TestQueue_Remove(t *testing.T) {
	queue := &Queue{Sort: SortDesc}
	queue.SetBucketSize(3)

	insertData(t, queue, "1", 1, 0, 0)
	insertData(t, queue, "2", 2, 0, 0)
	insertData(t, queue, "3", 1, 0, 2)
	insertData(t, queue, "4", 1.5, 0, 1)

	// the queue list is: 2, 4; 1, 3
	if err := queue.Remove(TestItem{Id: "1", Price: 1}); err != nil {
		t.Error(err)
	}

	// after remove, the queue list is: 2, 4; 3
	i, j := queue.FindPosition(TestItem{"3", 1})
	t.Logf("bucketPos:%d, dataPos:%d \n", i, j)
	if i != 1 || j != 0 {
		t.Error("item position not expect, expect 1, 0")
	}

	// remove 2, 4
	if err := queue.Remove(TestItem{Id: "2", Price: 2}); err != nil {
		t.Error(err)
	}
	if err := queue.Remove(TestItem{Id: "4", Price: 1.5}); err != nil {
		t.Error(err)
	}

	// after remove, the queue list is: 3
	i, j = queue.FindPosition(TestItem{"3", 1})
	t.Logf("bucketPos:%d, dataPos:%d \n", i, j)
	if i != 0 || j != 0 {
		t.Error("item position not expect, expect 0, 0")
	}

	// insert data
	insertData(t, queue, "5", 5, 0, 0)
}

func insertData(t *testing.T, queue *Queue, id string, price float64, bucketPosExpect, dataPosExpect int) {
	item := TestItem{
		Id:    id,
		Price: price,
	}

	if err := queue.AddItem(item); err != nil {
		t.Error(err)
	}

	i, j := queue.FindPosition(item)

	if i != bucketPosExpect || j != dataPosExpect {
		t.Errorf("item position not expect. bucketPosExpect:%d dataPosExpect:%d \n", bucketPosExpect, dataPosExpect)
		t.Logf("sortType:%s, id:%s, price:%f, bucketPos:%d, dataPos:%d \n", queue.Sort, id, price, i, j)
	}
}
