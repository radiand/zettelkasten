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
        let system_cmd = system_cmd ..  ' ' .. a:workspace
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

function! s:fzf_sink_from_ripgrep(result)
    " Extract absolute path from ripgrep oneline output.
    "
    " Args:
    "   result: string like 'path/to/file.md:10:20:match'

    let root_dir = s:get_root_dir()
    let selected_path = split(a:result, ':')[0]
    execute ":edit " .. root_dir .. '/' .. selected_path
endfunction

function! s:fzf_find()
    " Spawn FZF window with search-as-you-type functionality.

    let command_fmt = 'rg --column --line-number --no-heading --color=always --smart-case -- %s || true'
    let preview_cmd = 'rg --passthru --no-line-number --color=always --smart-case --no-messages {q} {1}'
    let opts = {
    \   'source': printf(command_fmt, ''),
    \   'dir': s:get_root_dir(),
    \   'sink': function('s:fzf_sink_from_ripgrep'),
    \   'options': [
    \       '--disabled',
    \       '--ansi',
    \       '--bind', 'change:reload:sleep 0.1; ' .. printf(command_fmt, '{q}'),
    \       '--delimiter', ':',
    \       '--preview-window', 'right:60%:<80(up:60%)'
    \       '--preview', preview_cmd,
    \   ],
    \ }
    call fzf#run(fzf#wrap(opts))
endfunction

command! -nargs=? -complete=custom,s:complete_workspaces ZkNew :call s:new(<f-args>)
command! ZkGoto :call s:goto()
command! ZkFZF :call s:fzf_find()

nnoremap <leader>zk :ZkNew<Space>
nnoremap <leader>zg :ZkGoto<CR>
nnoremap <leader>zf :ZkFZF<CR>
