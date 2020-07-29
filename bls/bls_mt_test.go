package bls

import (
	"testing"
)

func prepareAggregate(n int) (sigs []Sign, sigs2 chan *Sign) {
	Init(BLS12_381)
	var sec SecretKey
	sec.SetByCSPRNG()
	sigs = make([]Sign, n)
	sigs2 = make(chan *Sign, n)
	msg := make([]byte, 1)
	for i := 0; i < n; i++ {
		msg[0] = byte(i)
		sigs[i] = *sec.SignByte(msg)
		sigs2 <- &sigs[i]
	}
	return sigs, sigs2
}

const N = 1000

func TestAggregate(t *testing.T) {
	sigs, sigs2 := prepareAggregate(N)
	var aggSig Sign
	var aggSig2 Sign
	var aggSig3 Sign
	aggSig.Aggregate(sigs)
	aggSig3.AggregateMT2(sigs)
	aggSig2.AggregateMT(sigs2)
	if !aggSig.IsEqual(&aggSig2) {
		t.Error("AggregateMT")
	}
	if !aggSig.IsEqual(&aggSig3) {
		t.Error("AggregateMT2")
	}
}

func BenchmarkAggregate(b *testing.B) {
	sigs, _ := prepareAggregate(N)
	b.ResetTimer()
	var aggSig Sign
	aggSig.Aggregate(sigs)
}

func BenchmarkAggregateMT(b *testing.B) {
	_, sigs := prepareAggregate(N)
	b.ResetTimer()
	var aggSig Sign
	aggSig.AggregateMT(sigs)
}

func BenchmarkAggregateMT2(b *testing.B) {
	sigs, _ := prepareAggregate(N)
	b.ResetTimer()
	var aggSig Sign
	aggSig.AggregateMT2(sigs)
}
