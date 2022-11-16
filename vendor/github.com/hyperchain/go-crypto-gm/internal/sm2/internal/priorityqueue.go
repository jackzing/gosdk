package internal

import "sync"

//Item item
type Item struct {
	value    *point // 优先级队列中的数据
	priority uint64 // 优先级队列中节点的优先级
	index    int    // index是该节点在堆中的位置
}

//PriorityQueue 优先级队列需要实现heap的interface
type PriorityQueue []*Item

var queuePool = sync.Pool{
	New: func() interface{} {
		return make(PriorityQueue, 0, maxBatchSize)
	},
}

func getPriorityQueue() PriorityQueue {
	return queuePool.Get().(PriorityQueue)
}

func closePriorityQueue(in PriorityQueue) {
	queuePool.Put(in)
}

var itemPool = sync.Pool{
	New: func() interface{} {
		return new(Item)
	},
}

//GetItem 获取Item实例
func (pq PriorityQueue) GetItem() *Item {
	return itemPool.Get().(*Item)
}

//PutItem 归还Item实例
func (pq PriorityQueue) PutItem(item *Item) {
	itemPool.Put(item)
}

// 绑定Len方法
func (pq PriorityQueue) Len() int {
	return len(pq)
}

// 大根堆
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].priority > pq[j].priority
}

// 绑定swap方法
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index, pq[j].index = i, j
}

//Pop 绑定put方法，将index置为-1是为了标识该数据已经出了优先级队列了
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	*pq = old[0 : n-1]
	item.index = -1
	return item
}

//Push 绑定push方法
func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}
