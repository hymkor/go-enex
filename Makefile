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

debug:
	go fmt ./...
	$(SET) "CGO_ENABLED=0" && go build $(GOOPT) -tags debug ./cmd/unenex
	$(SET) "CGO_ENABLED=0" && go build $(GOOPT) -tags debug ./cmd/exstyle

_dist:
	$(SET) "CGO_ENABLED=0" && go build $(GOOPT) ./cmd/unenex
	$(SET) "CGO_ENABLED=0" && go build $(GOOPT) ./cmd/exstyle
	zip $(NAME)-$(VERSION)-$(GOOS)-$(GOARCH).zip unenex$(EXT) exstyle$(EXT)

dist:
	$(SET) "GOOS=linux"   && $(SET) "GOARCH=386"   && $(MAKE) _dist
	$(SET) "GOOS=linux"   && $(SET) "GOARCH=amd64" && $(MAKE) _dist
	$(SET) "GOOS=windows" && $(SET) "GOARCH=386"   && $(MAKE) _dist
	$(SET) "GOOS=windows" && $(SET) "GOARCH=amd64" && $(MAKE) _dist

release:
	gh release create -d --notes "" -t $(VERSION) $(VERSION) $(wildcard $(NAME)-$(VERSION)-*.zip)

manifest:
	make-scoop-manifest *-windows-*.zip > $(NAME).json

clean:
	$(RM) *.png *.html

test-html:
	$(foreach I,$(wildcard *.enex),"./$(NAME)" $(I) && ) echo OK

test-md:
	$(foreach I,$(wildcard *.enex),"./$(NAME)" -markdown $(I) && ) echo OK

.PHONY: dist clean _dist all release manifest test-html test-md
