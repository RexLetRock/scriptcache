package ztcpclient

import "testing"

func BenchmarkUnbufferedChannelEmptyStruct(b *testing.B) {
	ch := make(chan struct{})
	go func() {
		for {
			<-ch
		}
	}()
	for i := 0; i < b.N; i++ {
		ch <- struct{}{}
	}
}

func BenchmarkBufferedChannelEmptyStruct(b *testing.B) {
	ch := make(chan struct{}, 1<<63-1)
	go func() {
		for {
			<-ch
		}
	}()
	for i := 0; i < b.N; i++ {
		ch <- struct{}{}
	}
}
