package request

type ParamRegister struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required,gte=6,lte=255"` //长度介于6到255
	Code     string `json:"code" binding:"required,gte=6,lte=255"`
}

type ParamLogin struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required,gte=6,lte=255"`
}

type ParamUpdateUserEmail struct {
}
