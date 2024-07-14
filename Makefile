NAME = moji-code
PREFIX = /usr/local/bin

LDFLAGS =-w -s

$(NAME): clean
	@go build -ldflags="$(LDFLAGS)" -o $(NAME) ./*.go

.PHONY: clean install uninstall

clean:
	@$(RM) $(NAME)

install:
	cp -i $(NAME) $(PREFIX)

uninstall:
	$(RM) -i $(PREFIX)/$(NAME)
