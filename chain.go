package eventd

import (
	"sync"
)

type Chain[T any] struct {
	len int

	header *ChainNode[T]
	tail   *ChainNode[T]

	mu sync.RWMutex
}

func (c *Chain[T]) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	length := c.len

	return length
}

type ChainNode[T any] struct {
	elem T

	weight int

	prev *ChainNode[T]
	next *ChainNode[T]
}

type ChainIterator[T any] struct {
	list []T

	cursor int
}

func (c *Chain[T]) AddElem(elem T, weight int) *ChainNode[T] {
	c.mu.Lock()
	defer c.mu.Unlock()

	node := new(ChainNode[T])
	node.elem = elem
	node.weight = weight

	if c.header == nil {
		c.header = node
		c.tail = node
		c.len = 1

		return node
	}

	if node.weight > c.header.weight {
		node.next = c.header
		c.header.prev = node
		c.header = node
	} else if node.weight <= c.tail.weight {
		node.prev = c.tail
		c.tail.next = node
		c.tail = node
	} else {
		prev := c.tail.prev

		for node.weight > prev.weight {
			prev = prev.prev
		}

		node.prev = prev
		node.next = prev.next
		node.prev.next = node
		node.next.prev = node
	}

	c.len += 1
	return node
}

func (c *Chain[T]) RemoveNode(node *ChainNode[T]) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.len == 0 || node == nil {
		return
	}

	target := c.header
	for target != nil && target != node {
		target = target.next
	}

	if target != nil {
		if target != c.header && target != c.tail {
			target.prev.next = target.next
			target.next.prev = target.prev
		} else {
			if target == c.header {
				c.header = target.next
				if c.header != nil {
					c.header.prev = nil
				}
			}

			if target == c.tail {
				c.tail = target.prev
				if c.tail != nil {
					c.tail.next = nil
				}
			}
		}
	}

	c.len -= 1
	target = nil
}

func (c *Chain[T]) Iteration() *ChainIterator[T] {
	c.mu.RLock()
	defer c.mu.RUnlock()

	iterator := new(ChainIterator[T])
	iterator.list = make([]T, 0, c.len)
	iterator.cursor = -1

	p := c.header

	for p != nil {
		iterator.list = append(iterator.list, p.elem)
		p = p.next
	}

	return iterator
}

func (iterator *ChainIterator[T]) Next() bool {
	iterator.cursor++

	return iterator.cursor < len(iterator.list)
}

func (iterator *ChainIterator[T]) Get() T {
	if iterator.cursor < len(iterator.list) {
		return iterator.list[iterator.cursor]
	} else {
		return *new(T)
	}
}
