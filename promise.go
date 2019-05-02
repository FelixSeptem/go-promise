// Package go_promise implement a simple mode of promise
package go_promise

import (
	"sync"
)

const (
	DEFAULT_CAPACITY      = 1024
	DEFAULT_MAXCONCURRENT = 1024
)

// define the necessary information of promise pattern
type promise struct {
	wg     sync.WaitGroup
	result chan struct {
		res interface{}
		err error
	}
	done            chan struct{}
	concurrentMutex chan struct{}
	whenSuccess     func(interface{}, error)
	whenFailure     func(interface{}, error)
	whenComplete    func()
}

// init a new promise use factory pattern
func NewPromise(capacity int, maxConcurrent int) *promise {
	if capacity <= 0 {
		capacity = DEFAULT_CAPACITY
	}
	if maxConcurrent <= 0 {
		maxConcurrent = DEFAULT_MAXCONCURRENT
	}
	wg := sync.WaitGroup{}
	result := make(chan struct {
		res interface{}
		err error
	}, capacity)
	cm := make(chan struct{}, maxConcurrent)
	for i := 0; i < maxConcurrent; i++ {
		cm <- struct{}{}
	}
	return &promise{
		wg:              wg,
		result:          result,
		concurrentMutex: cm,
		done:            make(chan struct{}, 0),
	}
}

// return the promise's result capacity, max concurrent, length of result, concurrent
func (p *promise) Info() (int, int, int, int) {
	return cap(p.result), cap(p.concurrentMutex), len(p.result), cap(p.concurrentMutex) - len(p.concurrentMutex)
}

// start to hold and wait promise end
func (p *promise) Start() {
	go func() {
		p.wg.Wait()
		p.whenComplete()
		close(p.done)
	}()
	for {
		select {
		case res := <-p.result:
			if res.err != nil {
				p.whenFailure(res.res, res.err)
			} else {
				p.whenSuccess(res.res, res.err)
			}
		case <-p.done:
			return
		}
	}
}

// add new task to the given promise
func (p *promise) Add(fun func() (res interface{}, err error)) {
	p.wg.Add(1)
	go func() {
		<-p.concurrentMutex
		res, err := fun()
		p.result <- struct {
			res interface{}
			err error
		}{res: res, err: err}
		p.concurrentMutex <- struct{}{}
		p.wg.Done()
	}()
}

// register Success callback on promise
func (p *promise) OnSuccess(fun func(res interface{}, err error)) {
	p.whenSuccess = fun
}

// register Failure callback on promise
func (p *promise) OnFailure(fun func(res interface{}, err error)) {
	p.whenFailure = fun
}

// register Complete callback on promise
func (p *promise) OnComplete(fun func()) {
	p.whenComplete = fun
}
