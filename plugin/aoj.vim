if exists('g:loaded_aoj')
  finish
endif
let g:loaded_aoj = 1

function! s:RequireAOJ(host) abort
  return jobstart(['aoj.nvim'], {
        \ 'rpc': v:true,
        \ 'on_stderr': {->execute('echom string(a:000)', 1)},
        \ })
endfunction


call remote#host#Register('aoj.nvim', '0', function('s:RequireAOJ'))
call remote#host#RegisterPlugin('aoj.nvim', '0', [
\ {'type': 'command', 'name': 'AojSession', 'sync': 1, 'opts': {}},
\ {'type': 'function', 'name': 'AojDescription', 'sync': 0, 'opts': {}},
\ {'type': 'function', 'name': 'AojRunSample', 'sync': 0, 'opts': {}},
\ {'type': 'function', 'name': 'AojSubmit', 'sync': 0, 'opts': {}},
\ ])

nnoremap <silent><C-d>s :<C-u>call AojSubmit(input("problem id(submit): "))<CR>
nnoremap <silent><C-d>t :<C-u>call AojRunSample(input("problem id(run sample): "))<CR>
nnoremap <silent><C-d>d :<C-u>call AojDescription(input("problem id(description): "))<CR>

