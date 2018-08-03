PROTOFILES=\
	apns/apns.proto\
	gcm/gcm.proto\
	webpush/webpush.proto\
	adm/adm.proto\
	master.proto\

GOFILES=$(PROTOFILES:%.proto=%.pb.go)

all: $(GOFILES)

%.pb.go: %.proto
	protoc -I. --go_out=plugins=grpc:$(shell go env GOPATH)/src $<
