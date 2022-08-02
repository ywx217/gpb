bench_tmp_file=/tmp/gpb-bench.txt

gen: gen-proto

gen-proto:
	@protoc -I=./internal/testprotos --go_out=. ./internal/testprotos/test.proto

coverage:
	@go test -gcflags=all=-l -coverprofile="coverage.txt" -covermode=atomic -count=1 -race ./...
	@go tool cover -func=coverage.txt -o coverage.out

bench:
	@go test . -test.bench='^Benchmark.*$$' -benchmem

bench-compare:
	@go test . -test.bench='^Benchmark(Gpb|GoProtobuf).*$$' -benchmem -count=5 > $(bench_tmp_file) && benchstat $(bench_tmp_file)

bench-optimized-cpu:
	@go test . -test.bench='^BenchmarkGpb.*$$' -benchtime=10s -cpuprofile profile.out
	@go tool pprof profile.out

bench-optimized-mem:
	@go test . -test.bench='^BenchmarkGpb.*$$' -benchtime=10s -benchmem -memprofile memprofile.out
	@go tool pprof memprofile.out

clean:
	@rm -f profile.out memprofile.out
