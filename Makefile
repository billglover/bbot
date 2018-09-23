build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/actionHandler cmd/actionHandler/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/authHandler cmd/authHandler/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/msgFlagger cmd/msgFlagger/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/msgSender cmd/msgSender/main.go

.PHONY: clean
clean:
	rm -rf ./bin ./vendor

.PHONY: deploy
deploy: clean build
	sls deploy --verbose
