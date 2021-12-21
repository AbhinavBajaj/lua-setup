vim.api.nvim_set_keymap('n', '<Leader>p', ':Files<CR>', { noremap = true, silent = true})
vim.api.nvim_set_keymap('n', '<Leader>f', ':Ag<CR>', { noremap = true, silent = true})
vim.api.nvim_set_keymap('n', '<Leader>P', ':Files ~/go-code/idl/<CR>', { noremap = true, silent = true})
vim.api.nvim_set_keymap('n', '<Leader>b', ':Buffers<CR>', { noremap = true, silent = true})
vim.api.nvim_set_keymap('n', '<Leader>t', ':Tags<CR>', { noremap = true, silent = true})
vim.api.nvim_set_keymap('n', '<Leader>x', ':Helptags<CR>', { noremap = true, silent = true})
