package utils

import (
	"math"
	"sync"
)

type WaitGroupPool struct {
	pool chan struct{}
	wg   *sync.WaitGroup
}

func NewWaitGroupPool(size int) *WaitGroupPool {
	if size <= 0 {
		size = math.MaxInt32
	}
	// 让线程必须要有两个以上才可以
	if size == 1 {
		size = 2
	}
	return &WaitGroupPool{
		pool: make(chan struct{}, size),
		wg:   &sync.WaitGroup{},
	}
}

func (p *WaitGroupPool) Add() {
	p.pool <- struct{}{}
	p.wg.Add(1)
}

func (p *WaitGroupPool) Done() {
	<-p.pool
	p.wg.Done()
}

func (p *WaitGroupPool) Wait() {
	p.wg.Wait()
}

func (p *WaitGroupPool) Size() int {
	return len(p.pool)
}
