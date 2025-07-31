package schemes


type RegisterScheme struct {
	Email string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=5,max=64"`
	Name string `json:"name" validate:"required,max=128"`
}

type LoginScheme struct {
	Email string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=5,max=64"`
	AppId int32 `json:"app_id" validate:"required"`
}
