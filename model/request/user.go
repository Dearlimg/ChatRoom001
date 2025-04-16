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
	Email string `json:"email" binding:"required,email,lte=50"` // 邮箱
	Code  string `json:"code" binding:"required,gte=6,lte=6"`   // 验证码
}
type ParamUpdateUserPassword struct {
	Code        string `json:"code" binding:"required,gte=6,lte=50"`        //验证码
	NewPassword string `json:"newPassword" binding:"required,gte=6,lte=50"` //新的密码
}
