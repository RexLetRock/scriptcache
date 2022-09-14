.DEFAULT_GOAL := bench
bench:
	go test -bench=. ./ztcp/ztcpclient/*_test.go

