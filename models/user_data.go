package models

// UserParamsName user name param.
type UserParamsName struct {
	Name string `query:"name" validate:"required"`
}

// UserParamsID user id param.
// @Description name param
type UserParamsID struct {
	ID uint `query:"id" validate:"required"`
}

// UserParamsNew new user params.
type UserParamsNew struct {
	Name     string `query:"name" validate:"required"`
	Email    string `query:"email" validate:"required"`
	Password string `query:"password" validate:"required"`
}

// UserParamsLogin login user params.
type UserParamsLogin struct {
	Email    string `query:"email" validate:"required"`
	Password string `query:"password" validate:"required"`
}
