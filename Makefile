.PHONY: proto

proto: proto/*.proto
	protoc -I . proto/*.proto --go_out=.