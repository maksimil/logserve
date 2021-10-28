run:
	$(MAKE) build-web
	go run . serve

build:
	$(MAKE) build-web
	go build .

install: 
	$(MAKE) build-web
	go install .

test:
	go test ./cmd

build-web:
	cgscript ./web/build.go