.PHONY: install_dependencies run clean deck test

SOURCES = $(shell find . -name \*.go)
BIN = pongo

build: $(BIN)

$(BIN): $(SOURCES)
	go build

run: $(BIN)
	./$(BIN)

deck:
	./deck/build.sh

install_dependencies:
	sudo dnf install mesa-libGLU-devel mesa-libGLES-devel libXrandr-devel \
		libXcursor-devel libXinerama-devel libXi-devel libXxf86vm-devel \
		alsa-lib-devel pkg-config

clean:
	rm -f $(BIN)

test:
	go test
