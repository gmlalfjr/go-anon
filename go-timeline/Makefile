.PHONY: build clean deploy
build:
    env GOOS=linux go build -ldflags="-s -w" -o <your bunary path>
clean:
    rm -rf ./bin
deploy: clean build
    sls deploy --verbose