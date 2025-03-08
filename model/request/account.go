package request

import "ChatRoom001/model/common"

type ParamsCreateAccount struct {
	Name      string `json:"name" binding:"required,gte=1,lte=20"`       // 账户名（唯一）
	Gender    string `json:"gender" binding:"required,oneof= 男 女 未知"`    // 性别
	Signature string `json:"signature" binding:"required,gte=0,lte=100"` // 签名
}

type ParamGetAccountToken struct {
	AccountID int64 `json:"account_id" form:"account_id" binding:"required,gte=1"` // 账号 ID
}

type ParamDeleteAccount struct {
	AccountID int64 `json:"account_id" binding:"required,gte=1"`
}

type ParamUpdateAccount struct {
	Name      string `json:"name" binding:"required,gte=1,lte=20"`
	Gender    string `json:"gender" binding:"required,oneof= 男 女 未知"`
	Signature string `json:"signature" binding:"required,gte=0,lte=100"`
}

type ParamGetAccountByName struct {
	Name string `json:"Name" form:"Name" binding:"required,gte=1,lte=20"`
	common.Page
}

type ParamGetAccountByID struct {
	AccountID int64 `json:"account_id" form:"account_id" binding:"required,gte=1"`
}
