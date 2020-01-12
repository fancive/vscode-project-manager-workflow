GO		:= go
GOBUILD	:= $(GO) build

all: clean build

build:
	$(GOBUILD) -o vscode *.go
	chmod -R 0755 ./*

clean:
	rm -rf ./vscode
	# rm -rf ./tmp/*