set encoding=utf-8
scriptencoding utf-8
setglobal fileencoding=utf-8

set autoindent                        " maintain indent of current line
set smartindent                       " Smart indent
set backspace=indent,eol,start

if exists('$SUDO_USER')
  set nobackup                        " don't create root-owned files
  set nowritebackup
else
  set backupdir=$HOME/.config/nvim/tmp/backup " keep backup files out of the way
endif

if exists('&belloff')
  set belloff=all                     " never ring the bell for any reason
endif

set cursorline                        " Highlight current line
set diffopt+=foldcolumn:0             " don't show fold column in diff view

if exists('$SUDO_USER')
  set noswapfile                      " don't create root-owned files
else
  set directory=$HOME/.config/nvim/tmp/swap/  " keep swap files out of the way
endif

set noemoji                           " don't assume all emoji are double width
set expandtab                         " always use spaces instead of tabs
set tabstop=2 " spaces per tab
set softtabstop=-1                    " use 'shiftwidth' for tab/bs at end of line
set shiftwidth=2
set noshowmode

if has('folding')
  if has('windows')
    set fillchars=diff:∙               " BULLET OPERATOR (U+2219, UTF-8: E2 88 99)
    set fillchars+=fold:·              " MIDDLE DOT (U+00B7, UTF-8: C2 B7)
    set fillchars+=vert:┃              " BOX DRAWINGS HEAVY VERTICAL (U+2503, UTF-8: E2 94 83)
  endif

  if has('nvim-0.3.1')
    set fillchars+=eob:\              " suppress ~ at EndOfBuffer
  endif

  set foldmethod=indent               " not as cool as syntax, but faster
  set foldlevelstart=99               " start unfolded
  set foldtext=chand#settings#foldtext()
endif

if v:version > 703 || v:version == 703 && has('patch541')
  set formatoptions+=j                " remove comment leader when joining comment lines
endif

set formatoptions+=n                  " smart auto-indenting inside numbered lists
set hidden                            " allows you to hide buffers with unsaved changes without being prompted

if exists('&inccommand')
  set inccommand=split                " live preview of :s results
endif

set laststatus=2                      " always show status line
set lazyredraw                        " don't bother updating screen during macro playback

if has('linebreak')
  set linebreak                       " wrap long lines at characters in 'breakat'
endif

set list                              " show whitespace
set listchars=nbsp:⦸                  " CIRCLED REVERSE SOLIDUS (U+29B8, UTF-8: E2 A6 B8)
set listchars+=tab:››                 " RIGHT-POINTING SINGLE ANGLE QUOTATION MARK (U+203A, UTF-8: E2 80 BA)
set listchars+=extends:»              " RIGHT-POINTING DOUBLE ANGLE QUOTATION MARK (U+00BB, UTF-8: C2 BB)
set listchars+=precedes:«             " LEFT-POINTING DOUBLE ANGLE QUOTATION MARK (U+00AB, UTF-8: C2 AB)
set listchars+=trail:•                " BULLET (U+2022, UTF-8: E2 80 A2)
set modelines=5                       " scan this many lines looking for modeline
set nojoinspaces                      " don't autoinsert two spaces after '.', '?', '!' for join command
set number                            " show line numbers in gutter

if exists('+relativenumber')
  set relativenumber                  " show relative numbers in gutter
endif

set scrolloff=3                       " start scrolling 3 lines before edge of viewport
set shell=sh                          " shell to use for `!`, `:!`, `system()` etc.
set noshiftround                      " don't always indent by multiple of shiftwidth
set shiftwidth=2                      " spaces per tab (when shifting)
set shortmess+=A                      " ignore annoying swapfile messages
set shortmess+=I                      " no splash screen
set shortmess+=O                      " file-read message overwrites previous
set shortmess+=T                      " truncate non-file messages in middle
set shortmess+=W                      " don't echo "[w]"/"[written]" when writing
set shortmess+=a                      " use abbreviations in messages eg. `[RO]` instead of `[readonly]`
if has('patch-7.4.314')
  set shortmess+=c                    " completion messages
endif
set shortmess+=o                      " overwrite file-written messages
set shortmess+=t                      " truncate file messages at start
set mouse=a                           " enable mouse in vim

if has('linebreak')
  let &showbreak='↳ '                 " DOWNWARDS ARROW WITH TIP RIGHTWARDS (U+21B3, UTF-8: E2 86 B3)
endif

if has('showcmd')
  set noshowcmd                       " don't show extra info at end of command line
endif

set sidescroll=0                      " sidescroll in jumps because terminals are slow
set sidescrolloff=3                   " same as scrolloff, but for columns
set smarttab                          " <tab>/<BS> indent/dedent in leading whitespace

set incsearch             " Show matches While searching
set ignorecase            " ignore case on search
set smartcase             " Ignores case if search is all lower, case sensitive otherwise
set hlsearch              " Highlight Search

set grepprg=rg\ --vimgrep " Use RipGrep for grepping
set clipboard+=unnamedplus " Use system Clipboard

if (has("nvim"))
  let $NVIM_TUI_ENABLE_TRUE_COLOR=1
endif

if has('syntax')
  set spellcapcheck=                  " don't check for capital letters at start of sentence
endif

if has('windows')
  set splitbelow                      " open horizontal splits below current window
endif

if has('vertsplit')
  set splitright                      " open vertical splits to the right of the current window
endif

if exists('&swapsync')
  set swapsync=                       " let OS sync swapfiles lazily
endif
set switchbuf=usetab                  " try to reuse windows/tabs when switching buffers

if has('syntax')
  set synmaxcol=150                   " don't bother syntax highlighting long lines
endif

if has('termguicolors')
  set termguicolors                   " use guifg/guibg instead of ctermfg/ctermbg in terminal
endif

if has('persistent_undo')
  if exists('$SUDO_USER')
    set noundofile                    " don't create root-owned files
  else
    set undodir=~/.config/nvim/tmp/undo       " keep undo files out of the way
    set undofile                      " actually use undo files
  endif
endif

set updatecount=80                    " update swapfiles every 80 typed chars
set updatetime=2000                   " CursorHold interval

if has('viminfo') " ie. Vim.
  let s:viminfo='viminfo'
elseif has('shada') " ie. Neovim.
  let s:viminfo='shada'
endif

if exists('s:viminfo')
  if exists('$SUDO_USER')
    " Don't create root-owned files.
    execute 'set ' . s:viminfo . '='
  else
    " Defaults:
    "   Neovim: !,'100,<50,s10,h
    "   Vim:    '100,<50,s10,h
    "
    " - ! save/restore global variables (only all-uppercase variables)
    " - '100 save/restore marks from last 100 files
    " - <50 save/restore 50 lines from each register
    " - s10 max item size 10KB
    " - h do not save/restore 'hlsearch' setting
    "
    " Our overrides:
    " - '0 store marks for 0 files
    " - <0 don't save registers
    " - f0 don't store file marks
    " - n: store in ~/.config/nvim/tmp
    "
    execute 'set ' . s:viminfo . "=!,'100,<0,s10,n~/.config/nvim/tmp/" . s:viminfo

    if !empty(glob('~/.config/nvim/tmp/' . s:viminfo))
      if !filereadable(expand('~/.config/nvim/tmp/' . s:viminfo))
        echoerr 'warning: ~/.config/nvim/tmp/' . s:viminfo . ' exists but is not readable'
      endif
    endif
  endif
endif

if has('mksession')
  set viewdir=$HOME/.config/nvim/tmp/sessions         " override ~/.vim/view default
  set viewoptions=cursor,folds        " save/restore just these (with `:{mk,load}view`)
endif

if has('virtualedit')
  set virtualedit=block               " allow cursor to move where there is no text in visual block mode
endif
set visualbell t_vb=                  " stop annoying beeping for non-error errors
set whichwrap=b,h,l,s,<,>,[,],~       " allow <BS>/h/l/<Left>/<Right>/<Space>, ~ to cross line boundaries
set wildcharm=<C-z>                   " substitute for 'wildchar' (<Tab>) in macros
if has('wildignore')
  set wildignore+=*.o,*.rej           " patterns to ignore during file-navigation
endif
if has('wildmenu')
  set wildmenu                        " show options as list when switching buffers etc
endif
set wildmode=longest:full,full        " shell-like autocomplete to unambiguous portion
set wildignore+=*.so,*.pyc,*.png,*.jpg,*.gif,*.jpeg,*.ico,*.pdf
set wildignore+=*.wav,*.mp4,*.mp3
set wildignore+=*.o,*.out,*.obj,.git,*.rbc,*.rbo,*.class,.svn,*.gem
set wildignore+=*.zip,*.tar.gz,*.tar.bz2,*.rar,*.tar.xz
set wildignore+=*/vendor/gems/*,*/vendor/cache/*,*/.bundle/*,*/.sass-cache/*
set wildignore+=*.swp,*~,._*
set wildignore+=_pycache_,.DS_Store,.vscode,.localized
set wildignore+=.cache,node_modules,package-lock.json,yarn.lock,dist,.git
set wildignore+=.vimruncmd

