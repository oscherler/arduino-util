BINARY = arduino-util
PLATFORMS = darwin,arm64 darwin,amd64 linux,arm64 linux,amd64

default: build

build:
	go build -o $(BINARY)

build_all: clean build
	./build_all.sh $(PLATFORMS)

build_dist:
	@mkdir -p dist
	go build -o dist/$(BINARY)_$(GOOS)_$(GOARCH)

clean:
	rm -f $(BINARY)
	find dist -name '$(BINARY)*' -delete
