package handler

import (
	"context"

	"github.com/neovim/go-client/nvim"
)

func (a *Aoj) Session(v *nvim.Nvim, args []string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return a.aojClient.Session(ctx, a.config.id, a.config.password)
}
