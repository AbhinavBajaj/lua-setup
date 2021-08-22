let mapleader = "\<space>"
let maplocalleader = ","
let g:loaded_python_provider = 0
let g:loaded_perl_provider = 0

if empty(glob('~/.config/nvim/autoload/plug.vim'))
    silent !curl -fLo ~/.config/nvim/autoload/plug.vim --create-dirs
        \ https://raw.githubusercontent.com/junegunn/vim-plug/master/plug.vim
    autocmd VimEnter * PlugInstall --sync | source $MYVIMRC
endif

call plug#begin('~/.vim/plugged')

Plug 'tpope/vim-fugitive' " adding vim-fugitive
Plug 'tomasiser/vim-code-dark' " vscode like dark mode
Plug 'tpope/vim-surround' " adding vim surround
Plug 'fatih/vim-go', { 'do': ':GoUpdateBinaries' } " adding go support
Plug 'mhinz/vim-startify' " startup screen
Plug 'junegunn/fzf', { 'do': { -> fzf#install() } } " adding commonly used commands
Plug 'neoclide/coc.nvim', {'branch': 'master', 'do': 'yarn install --frozen-lockfile'} " Autocompletion plugin 
Plug 'junegunn/fzf.vim'
Plug 'vim-airline/vim-airline'
Plug 'vim-airline/vim-airline-themes'
Plug 'rhysd/accelerated-jk' " faster key press
Plug 'jiangmiao/auto-pairs' " autopair for syntax completion
Plug 'airblade/vim-gitgutter' " gutter for git status
Plug 'tomtom/tcomment_vim' " comment lines
Plug 'terryma/vim-multiple-cursors' " use multiple cursors
Plug 'voldikss/vim-floaterm' " add terminal to vim
Plug 'honza/vim-snippets'

call plug#end() " vim-colors-solarized-ours settings syntax enable set t_Co=256
set t_ut=
colorscheme codedark

" beautify using vim-go
let g:go_highlight_types = 1
let g:go_highlight_fields = 1
let g:go_highlight_functions = 1
let g:go_highlight_function_calls = 1
let g:go_highlight_operators = 1

" vim multiple cursor settings
let g:multi_cursor_select_all_word_key = '<C-a>'

"remap keys
inoremap kj <Esc>

nnoremap  <silent>   <tab>  :if &modifiable && !&readonly && &modified <CR> :write<CR> :endif<CR>:bnext<CR>
nnoremap  <silent> <s-tab>  :if &modifiable && !&readonly && &modified <CR> :write<CR> :endif<CR>:bprevious<CR>

set clipboard=unnamed "use system clipboard when yanking
set lazyredraw

function! Osc52Yank()
    let buffer=system('base64 -w0', @0)
    let buffer=substitute(buffer, "\n$", "", "")
    let buffer='\e]52;c;'.buffer.'\x07'
    silent exe "!echo -ne ".shellescape(buffer)." > ".shellescape(g:tty)
endfunction
augroup Yank
    autocmd!
    autocmd TextYankPost * if v:event.operator ==# 'y' | call Osc52Yank() | endif
augroup END

unmap <C-i>
