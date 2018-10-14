HOST="localhost"
PORT=5555
PROGRAM="gochat"

all: build

build:
	go install

client:
	telnet $(HOST) $(PORT)

delay:
	sleep 1s

server:
	$(PROGRAM)

watch-build:
	watchman-make -p '*.go' '**/*.go' '**/*.json' 'Makefile*' -t build

watch-server:
	watchman-make -p '*.go' '**/*.go' '**/*.json' 'Makefile*' -t delay server

watch-client:
	watchman-make -p '*.go' '**/*.go' '**/*.json' 'Makefile*' -t delay client

.PHONY: client server watch-build watch-client watch-server