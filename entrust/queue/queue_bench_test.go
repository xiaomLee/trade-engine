package queue

import (
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func BenchmarkQueue_AddItem50(b *testing.B) {
	benchmarkQueueAddItem(50, b)
}

func BenchmarkQueue_AddItem100(b *testing.B) {
	benchmarkQueueAddItem(100, b)
}

func BenchmarkQueue_AddItem200(b *testing.B) {
	benchmarkQueueAddItem(200, b)
}

func BenchmarkQueue_AddItem500(b *testing.B) {
	benchmarkQueueAddItem(500, b)
}

func BenchmarkQueue_AddItem1000(b *testing.B) {
	benchmarkQueueAddItem(1000, b)
}

func BenchmarkQueue_AddItem100000(b *testing.B) {
	benchmarkQueueAddItem(100000, b)
}

func benchmarkQueueAddItem(size int, b *testing.B) {
	queue := Queue{Sort: SortDesc}
	queue.SetBucketSize(size)
	rand.NewSource(time.Now().UnixNano())

	//b.N = 1000000
	data := make(map[int]TestItem)
	for i := 0; i < b.N; i++ {
		data[i] = TestItem{
			Id:    strconv.Itoa(i),
			Price: rand.Float64(),
		}
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		queue.AddItem(data[i])
	}
}

func BenchmarkQueue_Get50_2000(b *testing.B) {
	benchmarkQueueGet(50, 5000, b)
}

func BenchmarkQueue_Get100_5000(b *testing.B) {
	benchmarkQueueGet(100, 5000, b)
}

func BenchmarkQueue_Get200_5000(b *testing.B) {
	benchmarkQueueGet(200, 5000, b)
}

func BenchmarkQueue_Get500_5000(b *testing.B) {
	benchmarkQueueGet(500, 5000, b)
}

func BenchmarkQueue_Get50_100000(b *testing.B) {
	benchmarkQueueGet(50, 100000, b)
}

func BenchmarkQueue_Get100_100000(b *testing.B) {
	benchmarkQueueGet(100, 100000, b)
}

func BenchmarkQueue_Get200_100000(b *testing.B) {
	benchmarkQueueGet(200, 100000, b)
}

func BenchmarkQueue_Get500_100000(b *testing.B) {
	benchmarkQueueGet(500, 100000, b)
}

func BenchmarkQueue_Get1000_100000(b *testing.B) {
	benchmarkQueueGet(1000, 100000, b)
}

func benchmarkQueueGet(size, num int, b *testing.B) {
	queue := &Queue{Sort: SortDesc}
	queue.SetBucketSize(size)
	rand.NewSource(time.Now().UnixNano())

	data := initData(queue, num)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		queue.FindPosition(data[rand.Intn(num)])
	}

}

func BenchmarkQueue_Remove50_2000(b *testing.B) {
	benchmarkQueueRemove(50, 5000, b)
}

func BenchmarkQueue_Remove100_5000(b *testing.B) {
	benchmarkQueueRemove(100, 5000, b)
}

func BenchmarkQueue_Remove200_5000(b *testing.B) {
	benchmarkQueueRemove(200, 5000, b)
}

func BenchmarkQueue_Remove500_5000(b *testing.B) {
	benchmarkQueueRemove(500, 5000, b)
}

func BenchmarkQueue_Remove50_100000(b *testing.B) {
	benchmarkQueueRemove(50, 100000, b)
}

func BenchmarkQueue_Remove100_100000(b *testing.B) {
	benchmarkQueueRemove(100, 100000, b)
}

func BenchmarkQueue_Remove200_100000(b *testing.B) {
	benchmarkQueueRemove(200, 100000, b)
}

func BenchmarkQueue_Remove500_100000(b *testing.B) {
	benchmarkQueueRemove(500, 100000, b)
}

func BenchmarkQueue_Remove1000_100000(b *testing.B) {
	benchmarkQueueRemove(1000, 100000, b)
}

func benchmarkQueueRemove(size, num int, b *testing.B) {
	queue := &Queue{Sort: SortDesc}
	queue.SetBucketSize(size)
	rand.NewSource(time.Now().UnixNano())

	data := initData(queue, num)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		queue.Remove(data[rand.Intn(num)])
	}
}

func initData(queue *Queue, num int) map[int]TestItem {
	data := make(map[int]TestItem)
	for i := 0; i < num; i++ {
		item := TestItem{
			Id:    strconv.Itoa(i),
			Price: rand.Float64(),
		}
		data[i] = item
		queue.AddItem(item)
	}

	return data
}
