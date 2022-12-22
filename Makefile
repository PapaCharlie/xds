export GRPC_VERBOSITY=DEBUG
export GRPC_GO_LOG_VERBOSITY_LEVEL=99
export GRPC_GO_LOG_SEVERITY_LEVEL=info
export GRPC_XDS_BOOTSTRAP=./bootstrap.json

all: format
	go run -v server.go

format:
	goimports -w .