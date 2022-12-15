package main

import (
	"sync"
)

type Node struct {
	next *Node
	val  int
}

type Queue struct {
	head *Node
	tail *Node
}

func (q *Queue) Pop() *Node {
	if q.head == nil {
		return nil
	}

	if q.head == q.tail {
		head := q.head
		q.head = nil
		q.tail = nil
		return head
	}

	head := q.head
	q.head = q.head.next
	return head
}

func (q *Queue) Push(elem int) {
	newNode := &Node{
		next: nil,
		val:  elem,
	}

	if q.head == nil {
		q.head = newNode
		q.tail = newNode
		return
	} else if q.head == q.tail {
		q.tail = newNode
		q.head.next = q.tail
		return
	}

	q.tail.next = newNode
	q.tail = newNode
}

type BlockingQueue struct {
	q    *Queue
	cap  int
	cond *sync.Cond
	size int
}

func NewBlockingQueue(cap int) *BlockingQueue {
	q := &Queue{}
	l := &sync.Mutex{}
	cond := sync.NewCond(l)

	return &BlockingQueue{
		q:    q,
		cap:  cap,
		size: 0,
		cond: cond,
	}
}

func (b *BlockingQueue) Enqueue(elem int) {

	b.cond.L.Lock()
	defer b.cond.L.Unlock()
	for b.size == b.cap {
		b.cond.Wait()
	}

	b.q.Push(elem)
	b.size += 1
	b.cond.Broadcast()
}

func (b *BlockingQueue) Dequeue() int {

	b.cond.L.Lock()
	defer b.cond.L.Unlock()
	for b.size == 0 {
		b.cond.Wait()
	}

	val := b.q.Pop()
	b.size -= 1
	b.cond.Broadcast()
	return val.val

}

// func main() {
// 	bq := NewBlockingQueue(3)

// 	wg := &sync.WaitGroup{}

// 	// Enqueue and dequeue more elements than capacity
// 	for i := 0; i < 6; i++ {
// 		i := i
// 		wg.Add(1)
// 		go func() {
// 			defer wg.Done()
// 			bq.Enqueue(i)
// 		}()
// 	}

// 	for i := 0; i < 6; i++ {
// 		wg.Add(1)
// 		go func() {
// 			defer wg.Done()
// 			val := bq.Dequeue()
// 			fmt.Println(val)
// 		}()
// 	}

// 	// Sleep a bit and then try dequeuing from an empty queue
// 	time.Sleep(time.Second * 5)

// 	for i := 0; i < 6; i++ {
// 		wg.Add(1)
// 		go func() {
// 			defer wg.Done()
// 			val := bq.Dequeue()
// 			fmt.Println(val)
// 		}()
// 	}

// 	// Sleep a bit and then Enqueue enough elements to unblock all blocked threads
// 	time.Sleep(time.Second * 5)

// 	for i := 0; i < 6; i++ {
// 		i := i
// 		wg.Add(1)
// 		go func() {
// 			defer wg.Done()
// 			bq.Enqueue(i)
// 		}()
// 	}

// 	// Sleep a bit and then Enqueue more elements than capacity
// 	time.Sleep(time.Second * 5)

// 	for i := 0; i < 6; i++ {
// 		i := i
// 		wg.Add(1)
// 		go func() {
// 			defer wg.Done()
// 			bq.Enqueue(i)
// 		}()
// 	}

// 	// Sleep a bit and then dequeue all elements, unblocking all threads
// 	time.Sleep(time.Second * 5)

// 	for i := 0; i < 6; i++ {
// 		wg.Add(1)
// 		go func() {
// 			defer wg.Done()
// 			val := bq.Dequeue()
// 			fmt.Println(val)
// 		}()
// 	}

// 	wg.Wait()
// }
