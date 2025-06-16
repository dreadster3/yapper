.PHONY: build
build:
	mkdir -p out
	go build -o out/server ./cmd/server

.PHONY: test
test:
	go test ./...

.PHONY: docker
docker:
	docker build -t yapper/server:latest .

.PHONY: run
run:
	go run ./cmd/server


