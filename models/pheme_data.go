package models

// PhemeParamsPost params
// @Description post params
type PhemeParamsPost struct {
	Visibilty byte   `json:"visibility" validate:"min=0,max=255"`
	Category  string `json:"category" validate:"required"`
	Text      string `json:"text" validate:"required"`
	UserID    uint   `json:"userID" validate:"required"`
}

// PhemeParamsID param
// @Description id param
type PhemeParamsID struct {
	ID uint `json:"id" query:"id" validate:"required"`
}
