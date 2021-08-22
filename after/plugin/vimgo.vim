nmap <unique> <leader>gt :GoTest<cr>
nmap <unique> <leader>gf :GoTestFunc<cr>
nmap <unique> <leader>gr :GoReferrers<cr>
nmap <leader> <leader>gp :GoDebugBreakpoint<cr>
nmap <leader> <leader>gc :GoDebugContinue<cr>
" autoimport 
let g:go_fmt_command = "goimports"

" run :GoBuild or :GoTestCompile based on the go file
function! s:build_go_files()
  let l:file = expand('%')
  if l:file =~# '^\f\+_test\.go$'
    call go#test#Test(0, 1)
  elseif l:file =~# '^\f\+\.go$'
    call go#cmd#Build(0)
  endif
endfunction

" build and coverage commangs 
autocmd FileType go nmap <leader>b :<C-u>call <SID>build_go_files()<CR>
autocmd FileType go nmap <Leader>c <Plug>(go-coverage-toggle)

" remove invisible characters
set nolist