return {
    {
        "nvim-treesitter/nvim-treesitter",
        build = ":TSUpdate",
        enabled = not vim.g.vscode,
        config = function()
            local configs = require("nvim-treesitter.configs")
            vim.opt.foldmethod = "expr"
            vim.opt.foldexpr = "nvim_treesitter#foldexpr()"

            configs.setup({
                highlight = { enable = true },
                indent = { enable = true },
            })
        end,
    },
}
