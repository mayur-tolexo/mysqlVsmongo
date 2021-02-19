GO	= go
GODOC	= godoc


BIN = $(CURDIR)/build
BIN_NAME = comparator
MAIN_APP = $(CURDIR)/main.go


M = $(shell printf "\033[34;1m>>>\033[0m")

build: ; $(info $(M) Building binary...) @
	env GO111MODULE=on go build -v -o $(BIN)/$(BIN_NAME) $(MAIN_APP)