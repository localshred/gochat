HOST="localhost"
PORT=5555
PROGRAM="gochat"

all: build

build:
	go install

client:
	nc $(HOST) $(PORT)

server:
	$(PROGRAM)

watch-build:
	watchman-make -p '*.go' '**/*.go' '**/*.json' 'Makefile*' -t build

.PHONY: client server watch-build watch-client watch-server