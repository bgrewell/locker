.PHONY: all clean

all: locker

locker:
	go build -o bin/locker cmd/main.go

clean:
	rm -f bin/locker