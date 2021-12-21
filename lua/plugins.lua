local execute = vim.api.nvim_command
local fn = vim.fn

local install_path = fn.stdpath('data')..'/site/pack/packer/start/packer.nvim'

if fn.empty(fn.glob(install_path)) > 0 then
  fn.system({'git', 'clone', 'https://github.com/wbthomason/packer.nvim', install_path})
  execute 'packadd packer.nvim'
end

return require('packer').startup(function()
   use 'wbthomason/packer.nvim'                                 -- packer manager
   use 'mhinz/vim-startify'                                     -- startup screen
   use {'junegunn/fzf', run = ":call fzf#install()"}                                           -- adding commonly used commands
   use 'junegunn/fzf.vim'
   use 'terryma/vim-multiple-cursors'                           -- use multiple cursors
   use 'tpope/vim-surround'                                     -- adding vim surround
   use 'jiangmiao/auto-pairs'                                   -- autopair for syntax completion
   use 'tomtom/tcomment_vim'                                    -- comment lines
   use 'honza/vim-snippets'                                     -- add snippets 
   use 'hrsh7th/nvim-compe'                                     -- autocompletion
   use {'nvim-treesitter/nvim-treesitter', run = ":TSUpdate"}   -- tree sitter
   use 'neovim/nvim-lspconfig'                                  -- common lsp configs       
   use {
       'kyazdani42/nvim-tree.lua',
       requires = {
         'kyazdani42/nvim-web-devicons', -- optional, for file icon
       },
       config = function() require'nvim-tree'.setup {} end
   }
   use 'marko-cerovac/material.nvim'
   use 'nvim-lua/lsp-status.nvim'
   use 'sheerun/vim-polyglot'
   use 'xiyaowong/accelerated-jk.nvim'
   use 'ojroques/vim-oscyank'
   use 'hoob3rt/lualine.nvim'

   use 'airblade/vim-gitgutter'                               -- gutter for git status
   use 'Mofiqul/vscode.nvim' -- color schemes
   use {'akinsho/bufferline.nvim', requires = 'kyazdani42/nvim-web-devicons'}
end)

