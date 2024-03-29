.PHONY: build clean deploy

build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/read-business-card read-business-card/main.go

clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose
