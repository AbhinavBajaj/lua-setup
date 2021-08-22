function! GetSourceGraphLink()
  let l:remote = trim(system("git remote get-url origin"))
  let l:remote = substitute(l:remote, "gitolite@", "", "")
  let l:remote = substitute(l:remote, ":", "/", "")
  let l:remote = substitute(l:remote, "ssh://", "", "")
  let l:remote = substitute(l:remote, "\.git$", "", "")
  if l:remote == "go-code"
    let l:commit = trim(system("git merge-base @ origin/main"))
  else
    let l:commit = trim(system("git merge-base @ origin/master"))
  endif
  let l:filename = trim(system("git ls-files --full-name " . expand("%")))
  let l:url = "https://sourcegraph.uberinternal.com/" . l:remote . "@" . l:commit . "/-/blob/" . l:filename . "#L" . line(".") . ":" . col(".")
  return l:url
endfunction

function! CopySourceGraphLinkToClipboard()
  let l:url = GetSourceGraphLink()
  call setreg("a", l:url)
  echo "Copied SourceGraph URL to clipboard"
endfunction

function! CopySourceGraphLinkForPhab()
  let l:url = GetSourceGraphLink()
  call setreg("a", "[[ " . l:url . " | " . expand("%:t") . ":" . line(".") . " ]]")
  echo "Copied SourceGraph URL to clipboard with [[ ... | ... ]] wrapper"
endfunction

nmap <silent> <leader><C-g> :call CopySourceGraphLinkToClipboard()<cr>
nmap <silent> <leader><C-p> :call CopySourceGraphLinkForPhab()<cr>
