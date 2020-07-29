package bls

/*
#cgo CFLAGS:-fopenmp
#cgo LDFLAGS:-fopenmp
#include <omp.h>
#include <string.h>
#include <mcl/bn_c384_256.h>
#include <bls/bls.h>
inline void blsAggregateSignatureMT(blsSignature *aggSig, const blsSignature *sigVec, mclSize n)
{
	blsSignature ret;
	memset(&ret, 0, sizeof(ret));
	#pragma omp declare reduction(Agg: blsSignature: blsSignatureAdd(&omp_out, &omp_in)) initializer(omp_priv=omp_orig)
	#pragma omp parallel for reduction(Agg:ret)
	for (mclSize i = 0; i < n; i++) {
		blsSignatureAdd(&ret, &sigVec[i]);
	}
	*aggSig = ret;
}
*/
import "C"
import (
	"runtime"
)

func (sig *Sign) AggregateMT2(sigVec []Sign) {
	C.blsAggregateSignatureMT(&sig.v, &sigVec[0].v, C.mclSize(len(sigVec)))
}

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
