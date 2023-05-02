package handler

import (
	"context"
	"log"
	"time"

	nvimutil "github.com/n04ln/aoj.nvim/pkg/app/aoj_remote_plugin/nvim"
	"github.com/neovim/go-client/nvim"
)

type cycle uint64

type result string

var (
	loadingDone chan result
	cycleChars  = []string{
		"\\", "|", "/", "-", "\\", "|", "/", "-",
	}
)

func (c *cycle) String() string {
	tmp := cycle(uint64(*c) + 1)
	*c = tmp
	return cycleChars[uint64(*c)%uint64(len(cycleChars))]
}

func (r result) String() string {
	return string(r)
}

func (a *Aoj) drawLoading(v *nvim.Nvim) {
	c := new(cycle)
	loadingDone = make(chan result)
	err := nvimutil.ShowScratchBuffer(v, *a.scratchBuffer, c)
	if err != nil {
		return
	}
	for {
		select {
		case s := <-loadingDone:
			if s == "" {
				return
			}
			err := nvimutil.ShowScratchBuffer(v, *a.scratchBuffer, s)
			if err != nil {
				return
			}
			return
		default:
			time.Sleep(100 * time.Millisecond)
			err := nvimutil.SetContentToBuffer(v, *a.scratchBuffer, c)
			if err != nil {
				return
			}
		}
	}
}

func (a *Aoj) Submit(v *nvim.Nvim, args []string) (err error) {
	problemID, err := sanitizeProblemID(args)
	if err != nil {
		return err
	}

	if err := a.switchScratchBuffer(v); err != nil {
		return err
	}
	log.Println("id is " + problemID)

	go a.drawLoading(v)
	defer func() {
		if err != nil {
			loadingDone <- result("")
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	language, err := checkCurrentLanguage(v)
	if err != nil {
		return err
	}

	sourceCode, err := readCurrentBuffer(v)
	if err != nil {
		return err
	}

	token, err := a.aojClient.Submit(ctx, problemID, language, sourceCode)
	if err != nil {
		return err
	}

	stat, err := a.aojClient.Status(ctx, token, problemID)
	if err != nil {
		return err
	}

	loadingDone <- result(stat.String())

	return nil
}
