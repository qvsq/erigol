client: cmd/client/main.go
	@mkdir -p bin/
	@go build -o bin/client cmd/client/main.go
server: cmd/client/main.go
	@mkdir -p bin/
	@go build -o bin/server cmd/server/main.go
clean:
	@rm -rf bin/
remove_output:
	@rm -f output.txt
deps:
	@go mod download