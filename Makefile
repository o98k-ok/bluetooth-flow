all: build

build:
	mkdir -p dist
	go build -o dist/bluetooth main.go
	CCO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -o dist/bluetooth_arm main.go
	CCO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -o dist/bluetooth_amd main.go
	makefat ./dist/bluetooth ./dist/bluetooth_*
	rm -rf ./dist/bluetooth_*
	cp -r icons dist/icons
	cp blueutil dist/blueutil
	cp grid.sh dist/grid.sh

run:
	go run ./dist/bluetooth

clean:
	rm -rf dist