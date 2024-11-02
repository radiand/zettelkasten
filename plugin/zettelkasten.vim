" Vim integration with zettelkasten.go script.
" Maintainer: radiand@protonmail.com

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

function! s:complete_workspaces(arglead, cmdline, cursorpos)
    " Returns newline delimited string with available workspaces.
    " To be used with :command-completion-custom
    "
    " Args:
    "   ...unused, enforced by Vim.

    return system('zettelkasten get workspace')
endfunction

function! s:new(workspace=v:null)
    " Create new zettelkasten note.
    "
    " Args:
    "   workspace: workspace to create note in; use default if not specified.

    let system_cmd = 'zettelkasten new'
    if !empty(a:workspace)
        system_cmd = system_cmd .. ' ' .. a:workspace
    endif
    let new_note_path = trim(system(system_cmd))

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
        echoerr "No valid UID under cursor"
        return
    endif
    execute ":edit " .. s:get_path_of_note(uid)
endfunction

command! -nargs=? -complete=custom,s:complete_workspaces ZkNew :call s:new(<f-args>)
command! ZkGoto :call s:goto()

nnoremap <leader>zk :ZkNew<Space>
nnoremap <leader>zg :ZkGoto<CR>
