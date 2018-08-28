build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/acceptRequest functions/acceptRequest/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/flagMessage functions/flagMessage/main.go

.PHONY: clean
clean:
	rm -rf ./bin ./vendor

.PHONY: deploy
deploy: clean build
	sls deploy --verbose
