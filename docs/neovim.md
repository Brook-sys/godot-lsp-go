# Neovim

Example with `nvim-lspconfig`:

```lua
local lspconfig = require('lspconfig')
local configs = require('lspconfig.configs')

configs.godot_lsp_go = {
  default_config = {
    cmd = { 'godot-lsp-go' },
    filetypes = { 'gdscript' },
    root_dir = lspconfig.util.root_pattern('project.godot'),
  },
}

lspconfig.godot_lsp_go.setup{}
```

For auto-launch:

```lua
cmd = { 'godot-lsp-go', '--launch' }
```
