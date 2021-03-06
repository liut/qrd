.SILENT :
.PHONY : main clean dist

NAME:=qrd
TAG:=`git describe --tags`
LDFLAGS:=-X main.version=$(TAG)


all: main

main:
	echo "Building $(NAME)"
	go build -ldflags "$(LDFLAGS)"


clean:
	rm -rf dist
	rm -f $(NAME)
	rm -f $(NAME)-*.tar.?z

dist: clean
	mkdir -p dist/linux/amd64 && GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/linux/amd64/$(NAME)
	mkdir -p dist/darwin/amd64 && GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/darwin/amd64/$(NAME)

release: dist
	tar -cvJf $(NAME)-linux-amd64-$(TAG).tar.xz -C dist/linux/amd64 $(NAME)
	tar -cvJf $(NAME)-darwin-amd64-$(TAG).tar.xz -C dist/darwin/amd64 $(NAME)
