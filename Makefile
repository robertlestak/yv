bin: bin/yv_darwin bin/yv_linux bin/yv_windows

bin/yv_darwin:
	@echo "Building yv_darwin"
	GOOS=darwin GOARCH=amd64 go build -o bin/yv_darwin cmd/yv/*.go

bin/yv_linux:
	@echo "Building yv_linux"
	GOOS=linux GOARCH=amd64 go build -o bin/yv_linux cmd/yv/*.go

bin/yv_windows:
	@echo "Building yv_windows"
	GOOS=windows GOARCH=amd64 go build -o bin/yv_windows cmd/yv/*.go

docker:
	@echo "Building docker image"
	docker build -t yv .