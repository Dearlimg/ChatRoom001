package request

type ParamExistEmail struct {
	Email string `json:"email" binding:"required,email,lte=50"`
}

type ParamSendEmail struct {
	Email string `json:"email" binding:"required,lte=50"` // 邮箱
}
