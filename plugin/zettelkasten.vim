" Vim integration with zettelkasten.go script.
" Maintainer: radiand@protonmail.com
" License: MIT

if exists('g:loaded_zettelkasten')
  finish
endif
let g:loaded_zettelkasten = 1

function! s:get_root_dir()
    return trim(system('zettelkasten get config ZettelkastenDir'))
endfunction

function! s:get_path_of_note(uid)
    return trim(system('zettelkasten get -p note ' .. a:uid))
endfunction

function! s:new(workspace=v:null)
    " Create new zettelkasten note.

    let system_cmd = 'zettelkasten new'
    if !empty(a:workspace)
        system_cmd = system_cmd .. ' ' .. a:workspace
    endif
    let new_note_path = system(system_cmd)

    " Open buffer with new note.
    execute ":edit " . new_note_path

    " Jump cursor to the bottom.
    normal G
endfunction

function! s:goto()
    " Get UID under cursor and open the note in new buffer.

    let cur_word = expand("<cword>")
    let uid_pattern = "\\d\\{4}\\d\\{2}\\d\\{2}T\\d\\{2}\\d\\{2}\\d\\{2}Z"
    let uid = matchstr(cur_word, uid_pattern)
    if empty(uid)
        echo "No valid UID under cursor"
        return
    endif
    execute ":edit " .. s:get_path_of_note(uid)
endfunction

command! -nargs=? ZkNew :call s:new(<f-args>)
command! ZkGoto :call s:goto()

nnoremap <leader>zk :ZkNew<Space>
nnoremap <leader>zg :ZkGoto<CR>

" For telescope.nvim
" nnoremap <leader>zf :lua require("telescope.builtin").live_grep({cwd=vim.fn.ZkGetRootDir(), additional_args={"--no-ignore"}})<CR>

" For fzf.vim
nnoremap <leader>zf :call RGPath(0, s:get_root_dir())<CR>
