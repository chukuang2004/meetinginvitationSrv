package utils

import (
	"sync"
	"time"
)

type Schedule struct {
	Done   map[string]chan bool
	Locker sync.RWMutex
}

type Task func()

const (
	SUCCESS = 0
	IDEXIST = 1
)

var schInst *Schedule = nil

func GetScheduleInstance() *Schedule {
	if schInst == nil {
		schInst = &Schedule{
			Done: make(map[string]chan bool),
		}
	}

	return schInst
}
func (s *Schedule) Start(id string, task Task, sec time.Duration) int {

	s.Locker.RLock()
	_, ok := s.Done[id]
	s.Locker.RUnlock()

	if ok {
		return IDEXIST
	}

	go func() {
		ch := make(chan bool)

		s.Locker.Lock()
		s.Done[id] = ch
		s.Locker.Unlock()

		select {
		case <-ch:

		case <-time.After(sec):
			task()
		}

		close(ch)
		s.Locker.Lock()
		delete(s.Done, id)
		s.Locker.Unlock()
	}()

	return SUCCESS
}

func (s *Schedule) Quit(id string) {

	s.Locker.RLock()
	ch, ok := s.Done[id]
	s.Locker.RUnlock()

	if !ok {
		return
	}

	ch <- true
}
