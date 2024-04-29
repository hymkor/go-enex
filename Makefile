ifeq ($(OS),Windows_NT)
    SHELL=CMD.EXE
    SET=SET
    RM=del
    NUL=nul
else
    SET=export
    RM=rm
    NUL=/dev/null
endif

NAME=$(lastword $(subst /, ,$(abspath .)))
VERSION=$(shell git.exe describe --tags 2>$(NUL) || echo v0.0.0)
GOOPT=-ldflags "-s -w -X main.version=$(VERSION)"
EXT=$(shell go env GOEXE)

all:
	go fmt
	$(SET) "CGO_ENABLED=0" && go build $(GOOPT)
	$(foreach I,$(wildcard cmd/*),go fmt -C $(I) && $(SET) "CGO_ENABLED=0" && go build -C $(I) -o $(CURDIR) $(GOOPT) && ) echo OK

_package:
	$(MAKE) all
	zip enexToHtml-$(VERSION)-$(GOOS)-$(GOARCH).zip enexToHtml$(EXT)

package:
	$(SET) "GOOS=linux"   && $(SET) "GOARCH=386"   && $(MAKE) _package
	$(SET) "GOOS=linux"   && $(SET) "GOARCH=amd64" && $(MAKE) _package
	$(SET) "GOOS=windows" && $(SET) "GOARCH=386"   && $(MAKE) _package
	$(SET) "GOOS=windows" && $(SET) "GOARCH=amd64" && $(MAKE) _package

release:
	gh release create -d -t $(VERSION) $(VERSION) $(wildcard enexToHtml-$(VERSION)-*.zip)
manifest:
	make-scoop-manifest *-windows-*.zip > enexToHtml.json

clean:
	$(RM) *.png *.html

test-html:
	$(foreach I,$(wildcard *.enex),"./enexToHtml" $(I) && ) echo OK

test-md:
	$(foreach I,$(wildcard *.enex),"./enexToHtml" -markdown $(I) && ) echo OK

test-embed:
	$(foreach I,$(wildcard *.enex),"./enexToHtml" -embed $(I) && ) echo OK
