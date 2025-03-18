package request

type ParamCreateApplication struct {
	AccountID      int64  `json:"account_id" binding:"required,gte=1"`
	ApplicationMsg string `json:"application_msg" binding:"lte=200"`
}

type ParamDeleteApplication struct {
	AccountID int64 `json:"account_id" binding:"required,gte=1"`
}

type ParamRefuseApplication struct {
	AccountID int64  `json:"account_id" binding:"required,gte=1"`
	RefuseMsg string `json:"refuse_msg" binding:"lte=200"`
}

type ParamAcceptApplication struct {
	AccountID int64 `json:"account_id" binding:"required,gte=1"`
}
