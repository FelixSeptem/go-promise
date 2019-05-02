package go_promise

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"fmt"
	"time"
	"errors"
)

func TestNewPromise(t *testing.T) {
	p1 := NewPromise(128, 128)
	cr, cc, lr, lc := p1.Info()
	assert.Equal(t, 128, cr)
	assert.Equal(t, 128, cc)
	assert.Equal(t, 0, lr)
	assert.Equal(t, 0, lc)

	p2 := NewPromise(0, 0)
	cr, cc, lr, lc = p2.Info()
	assert.Equal(t, 1024, cr)
	assert.Equal(t, 1024, cc)
	assert.Equal(t, 0, lr)
	assert.Equal(t, 0, lc)
}

func TestPromise_Start(t *testing.T) {
	p := NewPromise(128, 128)
	p.OnSuccess(func(res interface{}, err error) {
		fmt.Printf("Success:%v\n", res)
	})
	p.OnFailure(func(res interface{}, err error) {
		fmt.Printf("Failure:%v\n", res)
	})
	p.OnComplete(func() {
		fmt.Println("End")
	})
	// Success Case
	p.Add(func() (res interface{}, err error) {
		return "Ok", nil
	})
	// Failed Case
	p.Add(func() (res interface{}, err error) {
		return "Failed", errors.New("failed")
	})
	// blocking success
	p.Add(func() (res interface{}, err error) {
		time.Sleep(time.Millisecond*500)
		return "Ok", nil
	})
	// blocking failed
	p.Add(func() (res interface{}, err error) {
		time.Sleep(time.Millisecond*500)
		return "Failed", errors.New("failed")
	})

	p.Start()
}