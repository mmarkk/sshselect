# Installation paths
PREFIX=/usr/local
BINDIR=$(PREFIX)/bin
MANDIR=$(PREFIX)/share/man
MAN1DIR=$(MANDIR)/man1

# User-local installation paths
USER_PREFIX=$(HOME)/.local
USER_BINDIR=$(USER_PREFIX)/bin
USER_MANDIR=$(USER_PREFIX)/share/man
USER_MAN1DIR=$(USER_MANDIR)/man1

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
BINARY_NAME=sshselect

all: build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(MAN1DIR)/$(BINARY_NAME).1.gz
	rm -f $(USER_MAN1DIR)/$(BINARY_NAME).1.gz

# System-wide installation (requires sudo)
install: build
	@echo "Installing system-wide (requires sudo)..."
	sudo install -d $(DESTDIR)$(BINDIR)
	sudo install -m 755 $(BINARY_NAME) $(DESTDIR)$(BINDIR)/
	sudo install -d $(DESTDIR)$(MAN1DIR)
	sudo install -m 644 $(BINARY_NAME).1 $(DESTDIR)$(MAN1DIR)/
	sudo gzip -f $(DESTDIR)$(MAN1DIR)/$(BINARY_NAME).1

# User-local installation
install-user: build
	@echo "Installing for current user..."
	install -d $(USER_BINDIR)
	install -m 755 $(BINARY_NAME) $(USER_BINDIR)/
	install -d $(USER_MAN1DIR)
	install -m 644 $(BINARY_NAME).1 $(USER_MAN1DIR)/
	gzip -f $(USER_MAN1DIR)/$(BINARY_NAME).1 || true
	@if [ ! -f $(USER_MAN1DIR)/$(BINARY_NAME).1.gz ]; then \
		cp $(BINARY_NAME).1 $(USER_MAN1DIR)/$(BINARY_NAME).1 && \
		gzip -f $(USER_MAN1DIR)/$(BINARY_NAME).1; \
	fi
	@echo "Add $(USER_BINDIR) to your PATH if not already done"
	@echo "Add $(USER_MANDIR) to your MANPATH if not already done"

uninstall:
	sudo rm -f $(DESTDIR)$(BINDIR)/$(BINARY_NAME)
	sudo rm -f $(DESTDIR)$(MAN1DIR)/$(BINARY_NAME).1.gz

uninstall-user:
	rm -f $(USER_BINDIR)/$(BINARY_NAME)
	rm -f $(USER_MAN1DIR)/$(BINARY_NAME).1.gz

.PHONY: all build clean install install-user uninstall uninstall-user
