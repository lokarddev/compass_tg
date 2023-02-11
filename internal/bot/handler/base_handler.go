package handler

import "app/internal/bot/utils"

type BaseHandler struct{}

var (
	MockEdgarSubs = map[string]string{
		"TSLA": "Tesla",
		"NFLX": "Netflix",
		"META": "Meta",
	}
)

func (h *BaseHandler) ValidState(currentNav string, validSources []string) bool {
	return utils.StrInSlice(validSources, currentNav)
}

func (h *BaseHandler) ValidText(input string, validSources []string) bool {
	return utils.StrInSlice(validSources, input)
}
