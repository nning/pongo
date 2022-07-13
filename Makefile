.PHONY: install_dependencies run

SOURCES = $(shell find . -name \*.go)
BIN = pongo

build: $(BIN)

$(BIN): $(SOURCES)
	go build

run:
	go run .

install_dependencies:
	sudo dnf install mesa-libGLU-devel mesa-libGLES-devel libXrandr-devel \
		libXcursor-devel libXinerama-devel libXi-devel libXxf86vm-devel \
		alsa-lib-devel pkg-config

clean:
	rm -f $(BIN)