GO_NAME=metainfo.go
BINARY_NAME=metainfo

build:
	go build -o bin/main metainfo.go

compile:
	echo "Compiling for every OS"
	GOOS=freebsd GOARCH=386 go build -o bin/freebsd/${BINARY_NAME} ${GO_NAME}
	GOOS=linux GOARCH=386 go build -o bin/linux/86/${BINARY_NAME} ${GO_NAME}
	GOOS=windows GOARCH=386 go build -o bin/windows/${BINARY_NAME} ${GO_NAME}
	GOOS=linux GOARCH=amd64 go build -o bin/linux/amd64/${BINARY_NAME} ${GO_NAME}

clean:
	go clean
	rm -rf bin/*
run:
	go run ${GO_NAME} test 32



