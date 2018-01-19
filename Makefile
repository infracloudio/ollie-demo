###############################################################################
# The build architecture is select by setting the ARCH variable.
# For example: When building on ppc64le you could use ARCH=ppc64le make <....>.
# When ARCH is undefined it defaults to amd64.
ARCH?=arm

# Determine which OS.
OS?=$(shell uname -s | tr A-Z a-z)

# Get version from git.
GIT_VERSION?=$(shell git describe --tags --dirty)
SKILL_SERVER?=ollie-skill-server
CERT_DIR?=/tmp/$(SKILL_SERVER)

.PHONY: all binary clean help ssl-keys
default: help

## Make all targets
all: cmds

## Generate all cmds
cmds: $(SKILL_SERVER)

$(SKILL_SERVER):
	$(MAKE) OS=$(OS) ARCH=$(ARCH) BINARY=$(SKILL_SERVER) binary

## Run API Server
run-api-server: $(SKILL_SERVER)
	sudo ./dist/$(SKILL_SERVER)-$(OS)-$(ARCH)  --port=5000

ssl-keys:
	mkdir -p $(CERT_DIR)
	openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
		-keyout $(CERT_DIR)/dev.key -out $(CERT_DIR)/dev.crt \
		-subj "/C=IN/ST=MH/L=Pune/CN=ollie-demo.io/emailAddress=sanket@infracloud.com"

binary:
	GOOS=$(OS) GOARCH=$(ARCH) CGO_ENABLED=0 go build -o dist/$(BINARY)-$(OS)-$(ARCH) \
	-ldflags "-X main.VERSION=$(GIT_VERSION)" cmd/$(BINARY)/main.go

## Generate swagger code
swagger-codegen:
	swagger generate server -f config/swagger/swagger.yaml -t ./pkg/
	rm -r cmd
	mv pkg/cmd .

## Cleanup all build files
clean:
	rm -r dist

.PHONY: help
## Display this help text.
help: # Some kind of magic from https://gist.github.com/rcmachado/af3db315e31383502660
	$(info Available targets)
	@awk '/^[a-zA-Z\-\_0-9\/]+:/ {                                      \
		nb = sub( /^## /, "", helpMsg );                                \
		if(nb == 0) {                                                   \
			helpMsg = $$0;                                              \
			nb = sub( /^[^:]*:.* ## /, "", helpMsg );                   \
		}                                                               \
		if (nb)                                                         \
			printf "\033[1;31m%-" width "s\033[0m %s\n", $$1, helpMsg;  \
	}                                                                   \
	{ helpMsg = $$0 }'                                                  \
	width=20                                                            \
	$(MAKEFILE_LIST)
