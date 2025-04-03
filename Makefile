#
# Simple Makefile for Golang based Projects.
#
PROJECT = newt

GIT_GROUP = caltechlibrary

RELEASE_DATE=$(shell date +'%Y-%m-%d')

RELEASE_HASH=$(shell git log --pretty=format:'%h' -n 1)

GO_PROGRAMS = newt ndr nte # $(shell ls -1 cmd)

TS_PROGRAMS =

MAN_PAGES = newt.1 ndr.1 nte.1 # $(shell ls -1 *.1.md | sed -E 's/\.1.md/.1/g')

HTML_PAGES = newt.1.html ndr.1.html nte.1.html # $(shell find . -type f | grep -E '\.html')

VERSION = $(shell grep '"version":' codemeta.json | cut -d\"  -f 4)

BRANCH = $(shell git branch | grep '* ' | cut -d\  -f 2)

PACKAGE = $(shell ls -1 *.go | grep -v 'version.go')

SUBPACKAGES = $(shell ls -1 */*.go)

OS = $(shell uname)

#PREFIX = /usr/local/bin
PREFIX = $(HOME)

ifneq ($(prefix),)
	PREFIX = $(prefix)
endif

EXT =
ifeq ($(OS), Windows)
	EXT = .exe
endif

DIST_FOLDERS = bin/*

build: version.go version.ts $(GO_PROGRAMS) $(TS_PROGRAMS) man CITATION.cff about.md installer.sh installer.ps1

version.ts: .FORCE
	@echo '' | pandoc --from t2t --to plain \
		--metadata-file codemeta.json \
		--metadata package=$(PROJECT) \
		--metadata version=$(VERSION) \
		--metadata release_date=$(RELEASE_DATE) \
		--metadata release_hash=$(RELEASE_HASH) \
		--template codemeta-version-ts.tmpl \
		LICENSE >version.ts

version.go: .FORCE
	@echo '' | pandoc --from t2t --to plain \
		--metadata-file codemeta.json \
		--metadata package=$(PROJECT) \
		--metadata version=$(VERSION) \
		--metadata release_date=$(RELEASE_DATE) \
		--metadata release_hash=$(RELEASE_HASH) \
		--template codemeta-version-go.tmpl \
		LICENSE >version.go

$(GO_PROGRAMS): $(PACKAGE)
	@mkdir -p bin
	go build -o "bin/$@$(EXT)" cmd/$@/*.go
	./bin/$@ -help >$@.1.md

$(TS_PROGRAMS): $(TS_PACKAGE)
	@mkdir -p bin
	env EXT=$(EXT) deno task compile_$@
	./bin/$@$(EXT) -help >$@.1.md

CITATION.cff: .FORCE
	@cat codemeta.json | sed -E   's/"@context"/"at__context"/g;s/"@type"/"at__type"/g;s/"@id"/"at__id"/g' >_codemeta.json
	@echo '' | pandoc --metadata title="Cite $(PROJECT)" --metadata-file=_codemeta.json --template=codemeta-cff.tmpl >CITATION.cff

about.md: .FORCE
	@cat codemeta.json | sed -E 's/"@context"/"at__context"/g;s/"@type"/"at__type"/g;s/"@id"/"at__id"/g' >_codemeta.json
	@echo "" | pandoc --metadata-file=_codemeta.json --template codemeta-about.tmpl >about.md 2>/dev/null;
	@if [ -f _codemeta.json ]; then rm _codemeta.json; fi

installer.sh: .FORCE
	@echo '' | pandoc --metadata title="Installer" --metadata git_org_or_person="$(GIT_GROUP)" --metadata-file codemeta.json --template codemeta-bash-installer.tmpl >installer.sh
	@chmod 775 installer.sh
	@git add -f installer.sh

installer.ps1: .FORCE
	@echo '' | pandoc --metadata title="Windows Powershell Installer" --metadata git_org_or_person="$(GIT_GROUP)" --metadata-file codemeta.json --template codemeta-ps1-installer.tmpl >installer.ps1
	@chmod 775 installer.ps1
	@git add -f installer.ps1

presentations: .FORCE
	- make -f presentation1.mak
	- make -f presentation2.mak
	- make -f presentation3.mak

clean-website:
	make -f website.mak clean

website: clean-website presentations .FORCE
	make -f website.mak

man: $(MAN_PAGES_1)

$(MAN_PAGES_1): .FORCE
	@mkdir -p man/man1
	pandoc $@.md --from markdown --to man -s >man/man1/$@

# NOTE: on macOS you must use "mv" instead of "cp" to avoid problems
install: build man .FORCE
	@if [ ! -d $(PREFIX)/bin ]; then mkdir -p $(PREFIX)/bin; fi
	@echo "Installing programs in $(PREFIX)/bin"
	@for FNAME in $(GO_PROGRAMS); do if [ -f ./bin/$$FNAME ]; then mv -v ./bin/$$FNAME $(PREFIX)/bin/$$FNAME; fi; done
	@echo ""
	@echo "Make sure $(PREFIX)/bin is in your PATH"
	@echo ""
	@echo "Installing manpages in $(PREFIX)/man"
	@if [ ! -d $(PREFIX)/man/man1 ]; then mkdir -p $(PREFIX)/man/man1; fi
	@for MAN_PAGE in $(MAN_PAGES); do cp -v man/man1/$$MAN_PAGE $(PREFIX)/man/man1/;done
	@echo ""
	@echo "Make sure $(PREFIX)/man is in your MANPATH"
	@echo ""

uninstall: .FORCE
	@echo "Removing programs in $(PREFIX)/bin"
	-for FNAME in $(GO_PROGRAMS); do if [ -f $(PREFIX)/bin/$$FNAME ]; then rm -v $(PREFIX)/bin/$$FNAME; fi; done
	-for MAN_PAGE in $(MAN_PAGES); do if [ -f "$(PREFIX)/man/man1/$$MAN_PAGE" ]; then rm "$(PREFIX)/man/man1/$$MAN_PAGE"; fi; done


hash: .FORCE
	git log --pretty=format:'%h' -n 1

check: .FORCE
	for FNAME in $(shell ls -1 *.go); do go fmt $$FNAME; done
	go vet *.go
	for FNAME in $(shell ls -1 *.ts); do deno fmt $$FNAME; deno check $$FNAME; done

test: .FORCE
	#dropdb birds
	#createdb birds
	#cd testdata && psql -d birds -c '\i birds-setup.sql'
	@echo 'NOTE: You need to run testdata/setup-for_tests.bash for test to succeed'
	go test #-test.v
	for FNAME in $(shell ls -1 *_test.ts); do deno test $$FNAME; done

clean:
	-if [ -d bin ]; then rm -fR bin; fi
	-if [ -d dist ]; then rm -fR dist; fi
	-if [ -d testout ]; then rm -fR testout; fi
	-for MAN_PAGE in $(MAN_PAGES); do if [ -f man/man1/$$MAN_PAGE.1 ]; then rm man/man1/$$MAN_PAGE.1; fi;done
	-make -f website.mak clean
	-make -f presentation1.mak clean
	-make -f presentation2.mak clean


status:
	git status

save:
	@if [ "$(msg)" != "" ]; then git commit -am "$(msg)"; else git commit -am "Quick Save"; fi
	git push origin $(BRANCH)

refresh:
	git fetch origin
	git pull origin $(BRANCH)

publish: build website save .FORCE
	./publish.bash


dist/Linux-x86_64: $(GO_PROGRAMS)
	@mkdir -p dist/bin
	@for FNAME in $(GO_PROGRAMS); do env  GOOS=linux GOARCH=amd64 go build -o "dist/bin/$${FNAME}" cmd/$${FNAME}/*.go; done
	@for FNAME in $(TS_PROGRAMS); do env EXT= TARGET="--target x86_64-unknown-linux-gnu" deno task compile_$$FNAME; done
	@cd dist && zip -r $(PROJECT)-v$(VERSION)-Linux-x86_64.zip LICENSE codemeta.json CITATION.cff *.md bin/* man/*
	@rm -fR dist/bin

dist/Linux-aarch64: $(GO_PROGRAMS)
	@mkdir -p dist/bin
	@for FNAME in $(GO_PROGRAMS); do env  GOOS=linux GOARCH=arm64 go build -o "dist/bin/$${FNAME}" cmd/$${FNAME}/*.go; done
	@for FNAME in $(TS_PROGRAMS); do env EXT= TARGET="--target aarch64-unknown-linux-gnu" deno task compile_$$FNAME; done
	@cd dist && zip -r $(PROJECT)-v$(VERSION)-Linux-aarch64.zip LICENSE codemeta.json CITATION.cff *.md bin/* man/*
	@rm -fR dist/bin

dist/macOS-x86_64: $(GO_PROGRAMS)
	@mkdir -p dist/bin
	@for FNAME in $(GO_PROGRAMS); do env GOOS=darwin GOARCH=amd64 go build -o "dist/bin/$${FNAME}" cmd/$${FNAME}/*.go; done
	@for FNAME in $(TS_PROGRAMS); do env EXT= TARGET="--target x86_64-apple-darwin" deno task compile_$$FNAME; done
	@cd dist && zip -r $(PROJECT)-v$(VERSION)-macOS-x86_64.zip LICENSE codemeta.json CITATION.cff *.md bin/* man/*
	@rm -fR dist/bin


dist/macOS-arm64: $(GO_PROGRAMS)
	@mkdir -p dist/bin
	@for FNAME in $(GO_PROGRAMS); do env GOOS=darwin GOARCH=arm64 go build -o "dist/bin/$${FNAME}" cmd/$${FNAME}/*.go; done
	@for FNAME in $(TS_PROGRAMS); do env EXT= TARGET="--target aarch64-apple-darwin" deno task compile_$$FNAME; done
	@cd dist && zip -r $(PROJECT)-v$(VERSION)-macOS-arm64.zip LICENSE codemeta.json CITATION.cff *.md bin/* man/*
	@rm -fR dist/bin


dist/Windows-x86_64: $(GO_PROGRAMS)
	@mkdir -p dist/bin
	@for FNAME in $(GO_PROGRAMS); do env GOOS=windows GOARCH=amd64 go build -o "dist/bin/$${FNAME}.exe" cmd/$${FNAME}/*.go; done
	@for FNAME in $(TS_PROGRAMS); do env EXT= TARGET="--target x86_64-pc-windows-msvc" deno task compile_$$FNAME; done
	@cd dist && zip -r $(PROJECT)-v$(VERSION)-Windows-x86_64.zip LICENSE codemeta.json CITATION.cff *.md bin/* man/*
	@rm -fR dist/bin


distribute_docs:
	@mkdir -p dist/
	@cp -v codemeta.json dist/
	@cp -v CITATION.cff dist/
	@cp -v README.md dist/
	@cp -v LICENSE dist/
	@cp -v INSTALL.md dist/
	@cp -vR man dist/

release: .FORCE installer.sh build distribute_docs dist/Linux-x86_64 dist/Linux-aarch64 \
	dist/macOS-x86_64 dist/macOS-arm64 \
	dist/Windows-x86_64

.FORCE:
