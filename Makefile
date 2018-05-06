include Makefile.waterlog

PKGNAME  = eopkg-deps
SUBPKGS  = cli index storage
PROJREPO = github.com/DataDrake

GOBIN    = build/bin
GOSRC    = build/src
PROJROOT = $(GOSRC)/$(PROJREPO)

LDFLAGS   = -ldflags "-s -w"
TAGS      = --tags "libsqlite3 linux"
GOCC      = GOPATH=$(shell pwd)/build go
GOFMT     = $(GOCC) fmt -x
GOGET     = $(GOCC) get $(LDFLAGS)
GOINSTALL = $(GOCC) install -v $(LDFLAGS) $(TAGS)
GOTEST    = $(GOCC) test -x
GOVET     = $(GOCC) vet -x

MEGACHECK = $(GOBIN)/megacheck
GOLINT    = $(GOBIN)/golint -set_exit_status

DESTDIR ?=
PREFIX  ?= /usr
BINDIR   = $(PREFIX)/bin

all: build

build: setup
	@$(call stage,BUILD)
	@$(GOINSTALL) $(PROJREPO)/$(PKGNAME)
	@$(call pass,BUILD)

setup:
	@$(call stage,SETUP)
	@$(call task,Setting up GOPATH...)
	@mkdir -p $(GOPATH)
	@$(call task,Setting up src/...)
	@mkdir -p $(GOSRC)
	@$(call task,Setting up project root...)
	@mkdir -p $(PROJROOT)
	@$(call task,Setting up symlinks...)
	@if [ ! -d $(PROJROOT)/$(PKGNAME) ]; then ln -s $(shell pwd) $(PROJROOT)/$(PKGNAME); fi
	@$(call task,Getting dependencies...)
	@if [ ! -e $(GOBIN)/glide ]; then $(GOGET) github.com/Masterminds/glide; rm -rf build/src/github.com/Masterminds; fi
	@$(GOBIN)/glide install
	@$(call pass,SETUP)

test: build
	@$(call stage,TEST)
	@for d in $(SUBPKGS); do $(GOTEST) ./$$d/... || exit 1; done
	@$(call pass,TEST)

validate: setup-validate
	@$(call stage,FORMAT)
	@for d in $(SUBPKGS); do $(GOFMT) ./$$d/...|| exit 1; done || $(GOFMT) $(PKGNAME).go
	@$(call pass,FORMAT)
	@$(call stage,VET)
	@for d in $(SUBPKGS); do $(MEGACHECK) ./$$d || exit 1; done || $(MEGACHECK) $(PKGNAME).go || exit 1
	@$(call pass,VET)
	@$(call stage,LINT)
	@for d in $(SUBPKGS); do $(GOLINT) ./$$d/... || exit 1; done || $(GOLINT) $(PKGNAME).go || exit 1;
	@$(call pass,LINT)

setup-validate:
	@if [ ! -e $(GOBIN)/megacheck ]; then \
	    printf "Installing megacheck..."; \
	    $(GOGET) honnef.co/go/tools/cmd/megacheck; \
	    printf "DONE\n\n"; \
	fi
	@if [ -d build/src/honnef.co ]; then rm -rf build/src/honnef.co; fi
	@if [ ! -e $(GOBIN)/golint ]; then \
	    printf "Installing golint..."; \
	    $(GOGET) github.com/golang/lint/golint; \
	    printf "DONE\n\n"; \
	fi
	@if [ -d build/src/golang.org ]; then rm -rf build/src/golang.org; fi
	@if [ -d build/src/github.com/golang ]; then rm -rf build/src/github.com/golang; fi

install:
	@$(call stage,INSTALL)
	install -D -m 00755 $(GOBIN)/$(PKGNAME) $(DESTDIR)$(BINDIR)/$(PKGNAME)
	@$(call pass,INSTALL)

uninstall:
	@$(call stage,UNINSTALL)
	rm -f $(DESTDIR)$(BINDIR)/$(PKGNAME)
	@$(call pass,UNINSTALL)

clean:
	@$(call stage,CLEAN)
	@$(call task,Removing symlinks...)
	@unlink $(PROJROOT)/$(PKGNAME)
	@$(call task,Removing build directory...)
	@rm -rf build
	@$(call pass,CLEAN)
