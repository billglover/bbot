build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/acceptRequest cmd/acceptRequest/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/takeAction cmd/takeAction/main.go

.PHONY: clean
clean:
	rm -rf ./bin ./vendor

.PHONY: deploy
deploy: clean build
	sls deploy --verbose
