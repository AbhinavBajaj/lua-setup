Vim�UnDo� �#5�N���h���|%���I�h��4���}?                                     `�~�    _�                             ����                                                                                                                                                                                                                                                                                                                                       ;           V        `�~�     �              ;   let mapleader = "\<space>"   let maplocalleader = ","    let g:loaded_python_provider = 0   let g:loaded_perl_provider = 0       2if empty(glob('~/.config/nvim/autoload/plug.vim'))   D    silent !curl -fLo ~/.config/nvim/autoload/plug.vim --create-dirs   M        \ https://raw.githubusercontent.com/junegunn/vim-plug/master/plug.vim   ;    autocmd VimEnter * PlugInstall --sync | source $MYVIMRC   endif   !" Specify a directory for plugins   ," - For Neovim: stdpath('data') . '/plugged'   :" - Avoid using standard Vim directory names like 'plugin'   !call plug#begin('~/.vim/plugged')   "   " adding vim-fugitive   ,Plug 'https://github.com/tpope/vim-fugitive'       $" adding emmet vim for html editing    Plug 'mattn/emmet-vim'       " adding color scheme   Plug 'chriskempson/base16-vim'   "   " adding vim surround   Plug 'tpope/vim-surround'       " startup screen   Plug 'mhinz/vim-startify'   " Autocompletion plugin    W Plug 'neoclide/coc.nvim', {'branch': 'master', 'do': 'yarn install --frozen-lockfile'}   "   " adding commonly used commands   3Plug 'junegunn/fzf', { 'do': { -> fzf#install() } }   Plug 'junegunn/fzf.vim'   "   "vim status bar   Plug 'itchyny/lightline.vim'   "   "add support for go   2Plug 'fatih/vim-go', { 'do': ':GoUpdateBinaries' }       " faster key press   Plug 'rhysd/accelerated-jk'   "   " Initialize plugin system       call plug#end()       $" vim-colors-solarized-ours settings   syntax enable   set background=dark   set termguicolors   *colorscheme base16-tomorrow-night-eighties   "   " remap keys   inoremap kj <Esc>   "    5�_�                             ����                                                                                                                                                                                                                                                                                                                                                  V        `�~�    �                   �               5��