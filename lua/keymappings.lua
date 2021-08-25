vim.api.nvim_set_keymap('n', '<Space>', '<NOP>', { noremap = true, silent = true})

vim.g.mapleader = ' '

vim.api.nvim_set_keymap('n', '<Leader>h', ':set hlsearch!<CR>', { noremap = true, silent = true})

vim.api.nvim_set_keymap('n', '<Leader>e', ':NvimTreeToggle<CR>', { noremap = true, silent = true})

-- better window movement
vim.api.nvim_set_keymap('n', '<C-h>', '<C-w>h', { silent = true})
vim.api.nvim_set_keymap('n', '<C-j>', '<C-w>j', { silent = true})
vim.api.nvim_set_keymap('n', '<C-k>', '<C-w>k', { silent = true})
vim.api.nvim_set_keymap('n', '<C-l>', '<C-w>l', { silent = true})

-- better indenting
vim.api.nvim_set_keymap('v', '>', '>gv', { noremap = true, silent = true})
vim.api.nvim_set_keymap('v', '<', '<gv', { noremap = true, silent = true})

-- I hate escape
vim.api.nvim_set_keymap('i', 'kj', '<Esc>', { noremap = true, silent = true})
vim.api.nvim_set_keymap('i', 'jk', '<Esc>', { noremap = true, silent = true})

-- Tab switch buffer
vim.api.nvim_set_keymap('n', '<TAB>', ':bnext<CR>', { noremap = true, silent = true})
vim.api.nvim_set_keymap('n', '<S-TAB>', ':bprevious<CR>', { noremap = true, silent = true})

-- Move selected line / block of text in visual mode
vim.api.nvim_set_keymap('x', 'K', ':move \'<-2<CR>gv-gv\'', { noremap = true, silent = true})
vim.api.nvim_set_keymap('x', 'J', ':move \'<+1<CR>gv-gv\'', { noremap = true, silent = true})

vim.g.onedark_termcolors=256


vim.g.synmaxcol=128




-- Example config in lua
-- vim.g.moonlight_italic_comments = true
-- vim.g.moonlight_italic_keywords = true
-- vim.g.moonlight_italic_functions = true
-- vim.g.moonlight_italic_variables = false
-- vim.g.moonlight_contrast = true
-- vim.g.moonlight_borders = false 
-- vim.g.moonlight_disable_background = true

-- Load the colorscheme
-- require('moonlight').set()

-- TAB completion
-- vim.api.nvim_set_keymap('i', '<expr><TAB>', 'pumvisible() ? \"\\<C-n>\" : \"\\<TAB>\"', { noremap = true, silent = true})
