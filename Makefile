.PHONY: proto

proto: proto/*.proto
	protoc -I proto proto/*.proto --go_out=plugins=grpc:proto
