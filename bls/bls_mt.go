package bls

import (
	"runtime"
)

// AggregateMT (Multi thread aggregation using up to max threads) --
func (sig *Sign) AggregateMT(sigChan chan *Sign) {
	MTAdd := func(res *Sign, op *Sign, sigChan chan *Sign, workerChan chan int) {
		res.Add(op)
		sigChan <- res
		<-workerChan
	}
	c := len(sigChan)
	maxThread := runtime.NumCPU() * 100
	workerChan := make(chan int, maxThread)
	for {
		if len(workerChan) < maxThread && len(sigChan) > 1 {
			workerChan <- 1
			go MTAdd((<-sigChan), (<-sigChan), sigChan, workerChan)
			c--
		}
		if c <= 1 && len(workerChan) == 0 {
			sig.v = (<-sigChan).v
			break
		}
	}
}
