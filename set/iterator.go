package set

import "sync"

// Iterator 表示集合元素迭代器。
type Iterator[T comparable] struct {
	// C 是输出集合元素的只读通道。
	C <-chan T

	stop     chan struct{}
	stopOnce sync.Once
}

// newIterator 根据元素快照创建迭代器。
func newIterator[T comparable](values []T) *Iterator[T] {
	channel := make(chan T)
	stop := make(chan struct{})
	iterator := &Iterator[T]{C: channel, stop: stop}

	go func() {
		defer close(channel)

		for _, value := range values {
			select {
			case <-stop:
				return
			case channel <- value:
			}
		}
	}()

	return iterator
}

// newIterChannel 根据元素快照创建不会阻塞发送方的遍历通道。
func newIterChannel[T comparable](values []T) <-chan T {
	channel := make(chan T, len(values))
	for _, value := range values {
		channel <- value
	}
	close(channel)
	return channel
}

// Stop 停止迭代器并关闭输出通道。
func (i *Iterator[T]) Stop() {
	if i == nil {
		return
	}

	i.stopOnce.Do(func() {
		if i.stop != nil {
			close(i.stop)
		}
	})
}
