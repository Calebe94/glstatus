##
# glstatus
#
# @file
# @version 0.1

name = glstatus
PREFIX = /usr/local

all: build

build:
	go build -o $(name)

install: all
	cp -f $(name) "$(DESTDIR)$(PREFIX)/bin/"
	chmod 755 "$(DESTDIR)$(PREFIX)/bin/glstatus"

uninstall:
	rm -f "$(DESTDIR)$(PREFIX)/bin/glstatus"

clean:
	rm -f $(name)

# end
