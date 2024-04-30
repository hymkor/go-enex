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

NAME=unenex
VERSION=$(shell git describe --tags 2>$(NUL) || echo v0.0.0)
GOOPT=-ldflags "-s -w -X main.version=$(VERSION)"
EXT=$(shell go env GOEXE)

all:
	go fmt
	$(SET) "CGO_ENABLED=0" && go build $(GOOPT)
	$(foreach I,$(wildcard cmd/*),go fmt -C $(I) && $(SET) "CGO_ENABLED=0" && go build -C $(I) -o $(CURDIR) $(GOOPT) && ) echo OK

_dist:
	$(MAKE) all
	zip $(NAME)-$(VERSION)-$(GOOS)-$(GOARCH).zip $(NAME)$(EXT)

dist:
	$(SET) "GOOS=linux"   && $(SET) "GOARCH=386"   && $(MAKE) _dist
	$(SET) "GOOS=linux"   && $(SET) "GOARCH=amd64" && $(MAKE) _dist
	$(SET) "GOOS=windows" && $(SET) "GOARCH=386"   && $(MAKE) _dist
	$(SET) "GOOS=windows" && $(SET) "GOARCH=amd64" && $(MAKE) _dist

release:
	gh release create -d -t $(VERSION) $(VERSION) $(wildcard enexToHtml-$(VERSION)-*.zip)
manifest:
	make-scoop-manifest *-windows-*.zip > $(NAME).json

clean:
	$(RM) *.png *.html

test-html:
	$(foreach I,$(wildcard *.enex),"./$(NAME)" $(I) && ) echo OK

test-md:
	$(foreach I,$(wildcard *.enex),"./$(NAME)" -markdown $(I) && ) echo OK

.PHONY: dist clean _dist all release manifest test-html test-md
