.DEFAULT_GOAL := default
go_mod:
	go mod download
imac-install: darwin
	cp ./darwin/amd64/fr $(shell echo $$PATH | cut -d ':' -f 1)
mac-install: darwin
	cp ./darwin/arm64/fr $(shell echo $$PATH | cut -d ':' -f 1)
install: linux
	cp ./linux/adm64/fr /usr/local/bin/
clean:
	rm -rf ./{linux,darwin}

default: linux darwin

linux: go_mod
	mkdir -p linux/amd64
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o linux/amd64/fr -a -ldflags '-extldflags "-static"' .
darwin: go_mod
	mkdir -p darwin/amd64
	mkdir -p darwin/arm64
	env GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o darwin/amd64/fr -a -ldflags '-extldflags "-static"' .
	env GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -o darwin/arm64/fr -a -ldflags '-extldflags "-static"' .
