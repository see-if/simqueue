package simqueue

import (
	"container/heap"
	"sync"
	"time"
)

var entry *DelayedQueue

type DelayedFunc struct {
	F     func()
	RunAt time.Time
	ID    string
}
type DelayedFuncs []*DelayedFunc

func (df DelayedFuncs) Len() int           { return len(df) }
func (df DelayedFuncs) Less(i, j int) bool { return df[i].RunAt.Before(df[j].RunAt) }
func (df DelayedFuncs) Swap(i, j int)      { df[i], df[j] = df[j], df[i] }

func (df *DelayedFuncs) Push(x interface{}) {
	*df = append(*df, x.(*DelayedFunc))
}

func (df *DelayedFuncs) Pop() interface{} {
	old := *df
	n := len(old)
	x := old[n-1]
	*df = old[0 : n-1]
	return x
}

type DelayedQueue struct {
	Funcs DelayedFuncs
}

var lock sync.Mutex

var delayOnce = &sync.Once{}

func GetEntry() *DelayedQueue {
	delayOnce.Do(func() {
		entry = new(DelayedQueue)
	})
	return entry
}

func (dq *DelayedQueue) Schedule(f func(), delay time.Duration, ID string) {
	df := &DelayedFunc{
		F:     f,
		RunAt: time.Now().Add(delay),
		ID:    ID,
	}
	lock.Lock()
	heap.Push(&dq.Funcs, df)
	lock.Unlock()
}

func (dq *DelayedQueue) Cancel(ID string) bool {
	lock.Lock()
	for i, df := range dq.Funcs {
		if df.ID == ID {
			heap.Remove(&dq.Funcs, i)
			lock.Unlock()
			return true
		}
	}
	lock.Unlock()
	return false
}

func (dq *DelayedQueue) Run() {
	for {
		lock.Lock()
		if dq.Funcs.Len() == 0 {
			lock.Unlock()
			time.Sleep(time.Second)
			continue
		}
		df := heap.Pop(&dq.Funcs).(*DelayedFunc)
		if df.RunAt.After(time.Now()) {
			heap.Push(&dq.Funcs, df)
			lock.Unlock()
			time.Sleep(df.RunAt.Sub(time.Now()))
			continue
		}
		lock.Unlock()
		df.F()
	}
}
