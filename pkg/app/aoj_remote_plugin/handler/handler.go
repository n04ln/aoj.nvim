package handler

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/n04ln/aoj.nvim/pkg/app/aoj_remote_plugin/aoj_client"
	nvimutil "github.com/n04ln/aoj.nvim/pkg/app/aoj_remote_plugin/nvim"
	"github.com/neovim/go-client/nvim"
)

type Aoj struct {
	aojClient *aoj_client.AojClient

	config struct {
		id, password string
	}

	scratchBuffer *nvim.Buffer
}

const (
	timeout = 10 * time.Second
)

func NewAoj(id, password string) *Aoj {
	aojClient := aoj_client.NewAojClient()

	// get cookie
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	aojClient.Session(ctx, id, password)

	return &Aoj{
		config: struct {
			id, password string
		}{
			id, password,
		},
		aojClient: aojClient,
	}
}

func (a *Aoj) switchScratchBuffer(v *nvim.Nvim) error {
	if a.scratchBuffer == nil {
		var err error
		a.scratchBuffer, err = newScratchBuffer(v)
		if err != nil {
			return err
		}
	}
	return nil
}

func newScratchBuffer(v *nvim.Nvim) (*nvim.Buffer, error) {
	return nvimutil.NewScratchBuffer(v, "AOJ")
}

func sanitizeProblemID(args []string) (string, error) {
	if len(args) != 1 {
		return "", errors.New("invalid args")
	}

	input := args[0]
	var problemID string
	if u, err := url.ParseRequestURI(input); err != nil {
		problemID = input
	} else {
		ids, ok := u.Query()["id"]
		if !ok || len(ids) == 0 {
			return "", errors.New("no such id")
		}
		problemID = ids[0]
	}
	return problemID, nil
}

var languageMap = map[string]string{
	"c":     "C",
	"hs":    "Haskell",
	"go":    "Go",
	"cpp":   "C++14",
	"java":  "JAVA",
	"cs":    "C#",
	"d":     "D",
	"rs":    "Rust",
	"rb":    "Ruby",
	"py":    "Python3",
	"js":    "JavaScript",
	"scala": "Scala",
	"php":   "PHP",
	"ml":    "OCaml",
	"kt":    "Kotlin",
}

func checkCurrentLanguage(v *nvim.Nvim) (string, error) {
	buf, err := v.CurrentBuffer()
	if err != nil {
		return "", err
	}

	bufferName, err := v.BufferName(buf)
	if err != nil {
		return "", err
	}

	ext := strings.Split(bufferName, ".")[len(strings.Split(bufferName, "."))-1]

	language, ok := languageMap[ext]
	if !ok {
		return "", fmt.Errorf("cannot submit this file: .%s", ext)
	}

	return language, nil
}

func readCurrentBuffer(v *nvim.Nvim) (string, error) {
	buf, err := v.CurrentBuffer()
	if err != nil {
		return "", err
	}

	lines, err := v.BufferLines(buf, 0, -1, true)
	if err != nil {
		return "", err
	}

	var content string
	for i, c := range lines {
		content += string(c)
		if i < len(lines)-1 {
			content += "\n"
		}
	}

	return content, nil
}
