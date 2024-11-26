.PHONY: clean impl

.DEFAULT_GOAL := impl

GIT_COMMIT=$(shell git rev-parse HEAD)
GIT_DATE=$(shell date "+%Y-%m-%dT%H:%M:%S")
GIT_TAG=$(shell git describe --tags --abbrev=0 --always)
BUILD_DIR=$(CURDIR)/build/bin

LDFLAGS:=-X 'main.gitCommit=$(GIT_COMMIT)' -X 'main.gitDate=$(GIT_DATE)' -X 'main.gitTag=$(GIT_TAG)'

ifeq ($(DEBUG),1)
	GCFLAGS:=all=-N -l
else 
	LDFLAGS:=$(LDFLAGS) -w -s
endif

LDFLAGS:=-ldflags="$(LDFLAGS)"
GCFLAGS:=-gcflags="$(GCFLAGS)"

impl:
	@echo "Building target: $@" 
	go build $(LDFLAGS) $(GCFLAGS) -o $(BUILD_DIR)/$@ $(CURDIR)/cmd/$@
	@echo "Done building."

clean:
	@rm -rf $(BUILD_DIR)/*

all: impl
