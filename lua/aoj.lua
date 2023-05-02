function GetWindowList()
    local wls = {}
    local maxWinNr = vim.api.nvim_eval("winnr('$')")
    for i = 1, maxWinNr, 1 do
        local bufnr = vim.api.nvim_eval("winbufnr(" .. i .. ")")
        wls[tostring(i)] = bufnr
    end
    return wls
end
