# sdl2-config --cflags --static-libs
# CGO_CFLAGS="-I/usr/include/SDL2 -D_REENTRANT -pthread -lSDL2 -lglib-2.0 -lgobject-2.0 -lgio-2.0 -libus-1.0 -ldbus-1 -ldl -lm -Wl,--no-undefined -Wl,-z,relro -Wl,--as-needed -Wl,-z,now -specs=/usr/lib/rpm/redhat/redhat-hardened-ld -pthread -lSDL2" \
# go build -tags static -ldflags "-s -w" # static
# sdl2-config --cflags --libs

build:
	CGO_ENABLED=1 \
	CC=gcc \
	CGO_CFLAGS="-I/usr/include/SDL2 -Llib/ -D_REENTRANT -pthread -lglib-2.0 -lgobject-2.0 -lgio-2.0 -libus-1.0 -ldbus-1 -ldl -lm -Wl,--no-undefined -Wl,-z,relro -Wl,--as-needed -Wl,-z,now" \
	go build -ldflags "-w" -o dist/boing_ball #dynamic
	#go build -tags static -ldflags "-s -w" #static
	patchelf --set-rpath '$$ORIGIN/lib' dist/boing_ball

run: build
	dist/boing_ball

dist: build
	-mkdir dist || true
	cp -r lib dist/
	tar czvf boing_ball.tgz -C dist/ .

install-deps:
	sudo dnf install -y SDL2{,_image,_mixer,_ttf,_gfx}-devel
	sudo dnf install -y libXxf86vm-devel libXScrnSaver-devel patchelf

remove-deps:
	sudo dnf remove -y SDL2{,_image,_mixer,_ttf,_gfx}-devel
	sudo dnf remove -y libXxf86vm-devel libXScrnSaver-devel patchelf
