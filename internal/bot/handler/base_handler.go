package handler

type BaseHandler struct {
}

func (h *BaseHandler) CheckState(previousNav string, validSources []string) bool {
	for _, v := range validSources {
		if v == previousNav {
			return true
		}
	}

	return false
}
