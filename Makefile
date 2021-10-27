run:
	go run . serve

build:
	go build .

install: 
	go install .

test:
	go test ./cmd

build-web:
	cgscript ./web/build.go