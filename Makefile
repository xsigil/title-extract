PREFIX ?= $(HOME)/.local
BINDIR = $(PREFIX)/bin
MANDIR = $(PREFIX)/share/man/man1

BIN = build/title-extract
SRC = src/main.go
MAN = doc/title-extract.1

.PHONY: all build install clean test

all: build

build: $(BIN)

$(BIN): $(SRC)
	@mkdir -p build
	go build -ldflags="-s -w" -o $(BIN) ./src

install: build
	@mkdir -p $(BINDIR) $(MANDIR)
	install -m 755 $(BIN) $(BINDIR)/title-extract
	install -m 644 $(MAN) $(MANDIR)/title-extract.1

clean:
	rm -rf build

test: build
	@echo "==> Running functional stream test..."
	@echo "<html><head><title>Makefile Test</title></head></html>" | ./$(BIN) | grep -q "Makefile Test" && echo "[✓] Pass"
