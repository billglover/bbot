GO111MODULE=on
export GOOS GO111MODULE

GOOS=linux
export GOOS

test:
	go vet ./...
	go fmt ./...
	go test ./...

.PHONY: build
build:
	go build -ldflags="-s -w" -o bin/actionHandler cmd/actionHandler/main.go
	go build -ldflags="-s -w" -o bin/authHandler cmd/authHandler/main.go
	go build -ldflags="-s -w" -o bin/msgFlagger cmd/msgFlagger/main.go
	go build -ldflags="-s -w" -o bin/msgSender cmd/msgSender/main.go

.PHONY: clean
clean:
	rm -rf ./bin ./vendor

.PHONY: stage
stage: clean build
	sls create_domain -s dev
	sls deploy --verbose --stage dev

.PHONY: release
release: clean build
	sls create_domain -s prod
	sls deploy --verbose --stage prod
