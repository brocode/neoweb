FROM alpine:latest

RUN apk add neovim
RUN mkdir -p /home/neovim
RUN adduser -h /home/neovim neovim -D -u 1000
COPY config /home/neovim/.config/nvim
USER neovim
WORKDIR /home/neovim
ENTRYPOINT [ "/usr/bin/nvim" ]
