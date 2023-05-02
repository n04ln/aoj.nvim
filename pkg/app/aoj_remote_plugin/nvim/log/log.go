package log

import (
	"fmt"

	"github.com/neovim/go-client/nvim"
)

const (
	pluginName = "aoj.nvim"
)

var (
	opt = map[string]interface{}{}
)

func Log(v *nvim.Nvim, message string) {
	sanitizeMessage := func(message string) string {
		return fmt.Sprintf("[%s] %s", pluginName, message)
	}
	_ = v.Echo(
		[]nvim.TextChunk{{Text: sanitizeMessage(message), HLGroup: ""}}, true, opt)
	return
}
