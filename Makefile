.PHONY: proto

proto: proto/*.proto
	protoc -I=. proto/*.proto --go_out=plugins=grpc:$(GOPATH)/src/github.com/jackgardner/go-ledger
