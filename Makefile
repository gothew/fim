make:
	go run main.go

test:
	go test ./... -short

build:
	go build -o fim main.go

install:
	go install
