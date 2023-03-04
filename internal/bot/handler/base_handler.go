package handler

import (
	"app/internal/model"
	"app/internal/utils"
)

type Nav struct {
	ValidSources []int
	ValidBranch  int
}

type BaseHandler struct {
}

var (
	MockEdgarSubs = map[string]string{
		"TSLA": "Tesla",
		"NFLX": "Netflix",
		"META": "Meta",
	}
)

func (h *BaseHandler) ValidState(user model.User, validNav Nav) bool {
	return utils.IntInSlice(validNav.ValidSources, user.State.NavCurrent) && validNav.ValidBranch == user.State.Branch
}

func (h *BaseHandler) ValidText(input string, validSources []string) bool {
	return utils.StrInSlice(validSources, input)
}
