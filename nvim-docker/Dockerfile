FROM alpine:latest

RUN apk add neovim git build-base
RUN mkdir -p /home/neovim
RUN adduser -h /home/neovim neovim -D -u 1000
COPY --chown=1000:1000 config /home/neovim/.config/nvim
USER neovim
WORKDIR /home/neovim
RUN nvim --headless +"Lazy! sync" \
    +"TSInstallSync \
        astro \
        bash \
        c \
        css \
        csv \
        diff \
        dockerfile \
        fish \
        git_config \
        git_rebase \
        gitattributes \
        gitcommit \
        gitignore \
        go \
        gomod \
        gosum \
        gowork \
        hcl \
        hocon \
        html \
        javascript \
        jsdoc \
        json \
        jsonc \
        kdl \
        kotlin \
        lua \
        luadoc \
        luap \
        markdown \
        markdown_inline \
        proto \
        python \
        query \
        regex \
        ruby \
        rust \
        scss \
        sql \
        ssh_config \
        svelte \
        templ \
        terraform \
        toml \
        tsx \
        typescript \
        vim \
        vimdoc \
        yaml \
        zig \
    " \
    +qa
ENTRYPOINT [ "/usr/bin/nvim" ]
