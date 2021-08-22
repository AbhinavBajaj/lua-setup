" Normal mode mappings.

" Multi-mode mappings (Normal, Visual, Operating-pending modes).
noremap Y y$

" Store relative line number jumps in the jumplist if they exceed a threshold.
nnoremap <expr> k (v:count > 5 ? "m'" . v:count : '') . 'k'
nnoremap <expr> j (v:count > 5 ? "m'" . v:count : '') . 'j'

" Close pane using c-w
noremap <silent> <C-w> :bd<Cr>

nnoremap <c-p> :bprevious<cr>
nnoremap <c-n> :bnext<cr>
nnoremap <f7> :tabprevious<cr>
nnoremap <f8> :tabnext<cr>

" Change text without putting the text into register,
nnoremap c "_c
nnoremap C "_C
nnoremap cc "_cc

" Make s/S/ss behave like d/D/dd without saving to register
nnoremap s  "_d
nnoremap S  "_D
nnoremap ss "_dd

" Quicker window movement
nnoremap <C-j> <C-w>j
nnoremap <C-k> <C-w>k
nnoremap <C-h> <C-w>h
nnoremap <C-l> <C-w>l

" Moves selected Lines up and Down with alt-j/k
nnoremap ∆ :m .+1<CR>==
nnoremap ˚ :m .-2<CR>==
inoremap ∆ <Esc>:m .+1<CR>==gi
inoremap ˚ <Esc>:m .-2<CR>==gi
vnoremap ∆ :m '>+1<CR>gv=gv
vnoremap ˚ :m '<-2<CR>gv=gv

" faster key press
nmap j <Plug>(accelerated_jk_gj)
nmap k <Plug>(accelerated_jk_gk)

