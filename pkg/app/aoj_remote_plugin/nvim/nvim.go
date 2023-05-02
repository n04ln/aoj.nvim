package nvim

import (
	"fmt"
	"strings"

	"github.com/neovim/go-client/nvim"
)

func NewScratchBuffer(v *nvim.Nvim, bufferName string) (*nvim.Buffer, error) {
	var scratchBuf nvim.Buffer
	var bwin nvim.Window
	var win nvim.Window

	b := v.NewBatch()
	b.CurrentWindow(&bwin)
	b.Command("silent! execute 'new' '" + bufferName + "'")
	b.CurrentBuffer(&scratchBuf)
	b.SetBufferOption(scratchBuf, "buftype", "nofile")
	b.SetBufferOption(scratchBuf, "bufhidden", "hide")
	b.Command("setlocal noswapfile")
	b.Command("setlocal nobuflisted")
	b.SetBufferOption(scratchBuf, "undolevels", -1)
	b.CurrentWindow(&win)
	b.SetWindowHeight(win, 15)
	if err := b.Execute(); err != nil {
		return nil, err
	}

	if err := v.SetCurrentWindow(bwin); err != nil {
		return nil, err
	}

	return &scratchBuf, nil
}

func SetContentToBuffer(v *nvim.Nvim, scratch nvim.Buffer, str fmt.Stringer) error {
	content := fmt.Sprint(str)
	lines := strings.Split(content, "\n")

	var byteContent [][]byte
	for _, c := range lines {
		byteContent = append(byteContent, []byte(c))
	}

	err := v.SetBufferLines(scratch, 0, -1, true, byteContent)
	if err != nil {
		return err
	}

	return nil
}

func ShowScratchBuffer(v *nvim.Nvim, scratch nvim.Buffer, str fmt.Stringer) error {
	if err := SetContentToBuffer(v, scratch, str); err != nil {
		return err
	}

	var winls map[string]int
	if err := v.ExecLua("return GetWindowList()", &winls); err != nil {
		return err
	}

	var opened bool
	for _, bufname := range winls {
		if nvim.Buffer(bufname) == scratch {
			opened = true
			break
		}
	}

	if !opened {
		var bwin nvim.Window
		var win nvim.Window

		b := v.NewBatch()
		b.CurrentWindow(&bwin)
		b.Command(fmt.Sprintf("sb %d", scratch))
		b.CurrentWindow(&win)
		b.SetWindowHeight(win, 15)
		if err := b.Execute(); err != nil {
			return err
		}

		return v.SetCurrentWindow(bwin)
	}

	return nil
}
