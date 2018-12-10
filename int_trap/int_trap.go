package int_trap

import (
	"os"
	"os/signal"
)

type Callback interface {
	Call()
}

type Hub struct {
	// Registered callbacks.
	callbacks map[*Callback]bool
}

var hub *Hub = &Hub{callbacks: make(map[*Callback]bool)}

func init() {
	signalReceiver()
}

func signalReceiver() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	go func() {
		select {
		case <-c:
			println("[keyboard interrupt]")
			for fn := range hub.callbacks {
				(*fn).Call()
			}
			os.Exit(0)
		}
	}()
}

func Register(fn *Callback) {
	hub.callbacks[fn] = true
}

func Unregister(fn *Callback) {
	if _, ok := hub.callbacks[fn]; ok {
		delete(hub.callbacks, fn)
	}
}

func Call(fn *Callback) {
	(*fn).Call()
	Unregister(fn)
}

type simpleCallback struct {
	fn func()
}

func (c *simpleCallback) Call() {
	(*c).fn()
}

func RegisterFunc(fn func()) *Callback {
	var cb Callback = &simpleCallback{fn: fn}
	Register(&cb)
	return &cb
}
