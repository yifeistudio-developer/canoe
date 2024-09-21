package route

import . "canoe/internal/model"

type rootController struct {
}

func (r *rootController) Get() *Result {
	return Success(nil)
}
