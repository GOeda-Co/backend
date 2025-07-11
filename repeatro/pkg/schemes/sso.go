package schemes


type RegisterScheme struct {
	Email string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=5,max=64"`
}

type LoginScheme struct {
	Email string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=5,max=64"`
	AppId int32 `json:"app_id" validate:"required"`
}

