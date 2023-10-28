LAMBDA_OUTPUT_DIR=./tmp/lambda

.PHONY: clean build

build: build-for-lambda

clean:
	rm -rf ./tmp

build-for-lambda:
	@mkdir -p $(LAMBDA_OUTPUT_DIR)
	GOOS=linux GOARCH=amd64 go build -o $(LAMBDA_OUTPUT_DIR)/handler ./cmd/handler.go
	zip -j ./tmp/handler.zip $(LAMBDA_OUTPUT_DIR)/handler

deploy: clean build
	pulumi up --yes
