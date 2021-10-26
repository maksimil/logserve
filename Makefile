run:
	go run .

build:
	go build -ldflags "-s -w" .

install: 
	go install .

test:
	go test ./cmd

build-web:
	cd web; cgscript ./build.go