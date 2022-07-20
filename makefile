.PHONY: clean
clean:
	rm -rf ./build/
.PHONY: build #打包字体 证书
build:
	go build -o=./build/dbMigrate cmd/main.go
	CGO_ENABLED=0  GOOS=linux  GOARCH=amd64  go build  -o=./build/dbMigrate_linux cmd/main.go
	CGO_ENABLED=0 GOOS=windows  GOARCH=amd64  go  build -o=./build/dbMigrate_windows cmd/main.go
