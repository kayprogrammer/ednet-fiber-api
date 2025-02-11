package base

type StatusData struct {
	Status string `json:"status" example:"failure"`
}

type FieldData struct {
	Field string `json:"field" example:"This field is required"`
}

type ValidationErrorExample struct {
	StatusData
	Message string `json:"message" example:"Invalid Entry"`
	Data    FieldData
}

type NotFoundErrorExample struct {
	StatusData
	Message string `json:"message" example:"The item was not found"`
}

type UnauthorizedErrorExample struct {
	StatusData
	Message string `json:"message" example:"Unauthorized user/Invalid credentials/Invalid Token"`
}

type InvalidErrorExample struct {
	StatusData
	Message string `json:"message" example:"Request was invalid due to ..."`
}
