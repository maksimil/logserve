run:
	go run .

build:
	go build -ldflags "-s -w" .

install: 
	go install .