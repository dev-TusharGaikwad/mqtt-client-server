

build: decompressor.go ./config/config.go main.go mqtt-client.go
	@echo "----------- building mqtt subscriber client------------"

	@go mod tidy
	@go get
	@go fmt
	@go vet
	@go build -o subscriber.bin

	@echo "BUILD SUCCESS, you can now type 'make run' to run subscriber"

run: subscriber.bin
	@go run subscriber

subscriber.bin: build

clean:
	@rm *.bin




