NAME     := wordle-search-api
VERSION  := v0.0.0
REVISION := $(shell git rev-parse --short HEAD)

SRCS    := $(shell find . -type f -name '*.go')
DSTDIR  := /srv/http/bin
USER    := http
GROUP   := http
LDFLAGS := -ldflags="-s -w -X \"main.Version=$(VERSION)\" -X \"main.Revision=$(REVISION)\" -extldflags \"-static\""

GOVER     := $(shell go version | awk '{ print substr($$3, 3) }' | tr "." " ")
VER_JUDGE := $(shell if [ $(word 1,$(GOVER)) -eq 1 ] && [ $(word 2,$(GOVER)) -le 10 ]; then echo 0; else echo 1; fi)

run:
	@go run *.go

init:
ifeq ($(VER_JUDGE),1)
	@go mod init $(NAME) && go mod tidy
else
	@echo "Packageの取得は手動で行ってください"
endif

build: $(SRCS)
	@go build -a -tags netgo -installsuffix netgo $(LDFLAGS) -o bin/$(NAME)

install:
	@command cp -r bin/$(NAME) $(DSTDIR)/
	@chown $(USER):$(GROUP) $(DSTDIR)/$(NAME)

uninstall: revoke_service
	@rm -f $(DSTDIR)/$(NAME)

create_service:
	@echo -e "[Unit]\nDescription=$(NAME)(Golang App)\n\n[Service]\nEnvironment=\"GIN_MODE=release\"\nWorkingDirectory=$(DSTDIR)/\n\nExecStart=$(DSTDIR)/$(NAME)\nExecStop=/bin/kill -HUP $MAINPID\nExecReload=/bin/kill -HUP $MAINPID && $(DSTDIR)/$(NAME)\n\nRestart=always\nType=simple\nUser=$(USER)\nGroup=$(GROUP)\n\n[Install]\nWantedBy=multi-user.target" | tee /etc/systemd/system/$(NAME).service
	@systemctl enable $(NAME).service

start: create_service
	@systemctl start $(NAME).service

revoke_service: /etc/systemd/system/$(NAME).service
	@systemctl stop $(NAME).service
	@systemctl disable $(NAME).service
	@rm -f /etc/systemd/system/$(NAME).service

clean:
	@rm -rf bin/*
	@rm -rf vendor/*

.PHONY: test
test:
	@go test

