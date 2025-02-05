package base

type ResponseSchema struct {
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"Data fetched/created/updated/deleted"`
}

func (obj ResponseSchema) Init() ResponseSchema {
	if obj.Status == "" {
		obj.Status = "success"
	}
	return obj
}

func ResponseMessage(message string) ResponseSchema {
	return ResponseSchema{Status: "success", Message: message}
}

type PaginatedResponseDataSchema struct {
	PerPage     uint `json:"per_page" example:"100"`
	CurrentPage uint `json:"current_page" example:"1"`
	LastPage    uint `json:"last_page" example:"100"`
}

type UserDataSchema struct {
	Name     string  `json:"name" example:"John Doe"`
	Username string  `json:"username" example:"john-doe"`
	Avatar   *string `json:"avatar" example:"https://img.url"`
}

// func (u UserDataSchema) Init(user *ent.User) UserDataSchema {
// 	user.Name = user.Name
// 	user.Username = user.Username
// 	user.Avatar = user.Avatar
// 	return user
// }