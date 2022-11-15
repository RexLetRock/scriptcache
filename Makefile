.DEFAULT_GOAL := bench
bench:
	go test -bench=. ./ztcp/ztcpclient/*_test.go

zbuffer-bench:
	go test -v ./zbuffer/ -bench=. -run=xxx -benchmem

zbufferv2-bench:
	go test -v ./zbufferv2/ -bench=. -run=xxx -benchmem
