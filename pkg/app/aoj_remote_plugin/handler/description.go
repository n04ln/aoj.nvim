package handler

import (
	"context"
	"log"

	nvimutil "github.com/n04ln/aoj.nvim/pkg/app/aoj_remote_plugin/nvim"
	"github.com/neovim/go-client/nvim"
)

func (a *Aoj) Description(v *nvim.Nvim, args []string) (err error) {
	problemID, err := sanitizeProblemID(args)
	if err != nil {
		return err
	}

	if err := a.switchScratchBuffer(v); err != nil {
		return err
	}
	log.Println("id is " + problemID)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	description, err := a.aojClient.GetDescription(ctx, problemID)
	if err != nil {
		return err
	}

	if err := nvimutil.ShowScratchBuffer(v, *a.scratchBuffer, description); err != nil {
		return err
	}

	return nil
}
