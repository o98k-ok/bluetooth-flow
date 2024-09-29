all: build

build:
	mkdir -p dist
	go build -o dist/bluetooth main.go
	cp -r icons dist/icons
	cp blueutil dist/blueutil
	cp grid.sh dist/grid.sh

run:
	go run ./dist/bluetooth

clean:
	rm -rf dist