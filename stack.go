package calculator

import (
	"sync"
)

type calcItem calcMinor

type calcItemStack struct {
	items []*calcMinor
	lock  sync.RWMutex
}

// 创建栈
func (s *calcItemStack) New() *calcItemStack {
	s.items = []*calcMinor{}
	return s
}

// 入栈
func (s *calcItemStack) Push(t *calcMinor) {
	s.lock.Lock()
	s.items = append(s.items, t)
	s.lock.Unlock()
}

// 出栈
func (s *calcItemStack) Pop() *calcMinor {
	s.lock.Lock()
	if s.IsEmpty() {
		s.lock.Unlock()
		return nil
	}
	item := s.items[len(s.items)-1]
	s.items = s.items[0 : len(s.items)-1]
	s.lock.Unlock()
	return item
}

// 取栈顶
func (s *calcItemStack) Top() *calcMinor {
	if s.IsEmpty() {
		return nil
	}
	return s.items[len(s.items)-1]
}

// 判空
func (s *calcItemStack) IsEmpty() bool {
	return len(s.items) == 0
}
