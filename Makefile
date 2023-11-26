VERSION=$(strip $(subst v, , $(shell git describe --tags --abbrev=0)))
HASH=$(shell git rev-parse --short HEAD)
TIME=$(shell git log --pretty=format:"%ad" --date=format:'%Y%m%d%H%M%S' $(HASH) -1)
BUILD=$(VERSION).$(HASH).$(TIME)
BUILDTIME=$(shell date "+%Y-%m-%dT%H:%M:%S%z")

GO_BUILD_LDFLAGS = -X 'github.com/lawyzheng/kragins/pkg/buildinfo.version=$(VERSION)' \
  				   -X 'github.com/lawyzheng/kragins/pkg/buildinfo.build=$(BUILD)' \
  				   -X 'github.com/lawyzheng/kragins/pkg/buildinfo.buildTime=$(BUILDTIME)'
GO_BUILD_FLAGS	= -ldflags "-w -s $(GO_BUILD_LDFLAGS)"

OS = darwin linux windows
ARCH = arm64 amd64

.PHONY: generate
generate:
	cd cmd/kragins/internal/cli; wire .

.PHONY: build
build: generate
	go build -a -tags prod -o "cmd/kragins/dist/$(shell go env GOOS)/$(shell go env GOARCH)/" $(GO_BUILD_FLAGS) \
		cmd/kragins/kragins.go;

.PHONY: prod
prod: generate
	rm -rf cmd/kragins/dist
	for arch in $(ARCH); do \
  		for os in $(OS); do \
			echo "release [$$os]-[$$arch] production"; \
			GOOS=$$os \
			GOARCH=$$arch \
			go build -a -tags prod -o "cmd/kragins/dist/$$os/$$arch/" $(GO_BUILD_FLAGS) \
				cmd/kragins/kragins.go; \
			cd cmd/kragins/dist/$$os/$$arch/; \
			zip -r kragins-$$os-$$arch-v$(VERSION).zip . ; \
			cd -; \
		done \
	done
