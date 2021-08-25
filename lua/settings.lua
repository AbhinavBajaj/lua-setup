vim.o.encoding="utf-8"
vim.o.fileencoding="utf-8"

vim.o.autoindent=true        
vim.o.smartindent=true
vim.o.backspace="indent,eol,start"

vim.o.expandtab=true                    --always use spaces instead of tabs
vim.o.tabstop=2                         -- spaces per tab
vim.o.softtabstop=-1                    --use 'shiftwidth' for tab/bs at end of line
vim.o.shiftwidth=2

vim.o.hidden=true                            --allows you to hide buffers with unsaved changes without being prompted
vim.o.laststatus=2                      --always show status line

vim.o.list=true                              --show whitespace
-- vim.o.relativenumber=true                   --ow relative numbers in gutter

vim.o.scrolloff=3                       --start scrolling 3 lines before edge of viewport
vim.o.shell="sh"                          --shell to use for `!`, `:!`, `system()` etc.
vim.o.shiftwidth=2                      --spaces per tab (when shifting)
vim.o.mouse="a"                           --enable mouse in vim

vim.o.incsearch=true                         --Show matches While searching
vim.o.ignorecase=true                        --ignore case on search
vim.o.smartcase=true                         --Ignores case if search is all lower, case sensitive otherwise
vim.o.hlsearch=true                          --Highlight Search

vim.o.wildmode="longest:full,full"        -- shell-like autocomplete to unambiguous portion
vim.o.wildignore="+=*.so,*.pyc,*.png,*.jpg,*.gif,*.jpeg,*.ico,*.pdf"
vim.o.wildignore="+=*.wav,*.mp4,*.mp3"
vim.o.wildignore="+=*.o,*.out,*.obj,.git,*.rbc,*.rbo,*.class,.svn,*.gem"
vim.o.wildignore="+=*.zip,*.tar.gz,*.tar.bz2,*.rar,*.tar.xz"
vim.o.wildignore="+=*/vendor/gems/*,*/vendor/cache/*,*/.bundle/*,*/.sass-cache/*"
vim.o.wildignore="+=*.swp,*~,._*"
vim.o.wildignore="+=_pycache_,.DS_Store,.vscode,.localized"
vim.o.wildignore="+=.cache,node_modules,package-lock.json,yarn.lock,dist,.git"
vim.o.wildignore="+=.vimruncmd"

vim.o.wrap=false
vim.o.showmode=false
vim.o.list=false

vim.o.lazyredraw=true
vim.o.modifiable=true

-- vim.lsp.set_log_level("debug")


-- vim.cmd('colorscheme codedark')
