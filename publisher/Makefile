

build: compressor.go config/config.go main.go mqtt-server.go
	@echo "----------- building mqtt publish server------------"

	@go mod tidy
	@go get
	@go fmt
	@go vet
	@go build -o publisher.bin

	@echo "BUILD SUCCESS, you can type 'make run' to run server"

run: publisher.bin
	@go run publisher

publisher.bin: build

clean:
	@rm *.bin




