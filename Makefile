all: clear install gen-example

.PHONY: proto
proto:
	protoc --go_out=paths=source_relative:. proto/marshaler.proto

install: proto
	go install
gen-example:
	protoc -I $(GOPATH)/src:. --roles_out=. --go_out=. example/example.proto
clear:
	rm -f proto/marshaler.pb.go
	rm -f example/example.ex.go
	rm -f example/example.pb.go