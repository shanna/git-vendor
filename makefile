PREFIX ?= /usr/local
INSTALL_FILES = `find bin -type f 2>/dev/null`

all:

test:

install:
	for file in $(INSTALL_FILES); do cp $$file $(PREFIX)/$$file; done

uninstall:
	for file in $(INSTALL_FILES); do rm -f $(PREFIX)/$$file; done

.PHONY: install uninstall all
