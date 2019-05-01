package go_promise

import "sync"

type Promise struct {
	Wg  sync.WaitGroup
	Res string
	Err error
}

func NewPromise(fun func() (string, error)) *Promise {
	p := &Promise{}
	p.Wg.Add(1)
	go func() {
		p.Res, p.Err = fun()
		p.Wg.Done()
	}()
	return p
}

func (p *Promise) Then(r func(string), e func(error)) {
	go func() {
		p.Wg.Wait()
		if p.Err != nil {
			e(p.Err)
			return
		}
		r(p.Res)
	}()
	return
}
