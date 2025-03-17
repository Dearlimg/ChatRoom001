package reply

import (
	"ChatRoom001/model/common"
	"time"
)

type ParamCreateAccount struct {
	ParamAccountInfo     ParamAccountInfo     `json:"param_account_info"`
	ParamGetAccountToken ParamGetAccountToken `json:"param_get_account_token"`
}

type ParamAccountInfo struct {
	ID        int64  `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Avatar    string `json:"avatar,omitempty"`
	Gender    string `json:"gender,omitempty"`
	Signature string `json:"signature,omitempty"`
}

type ParamGetAccountToken struct {
	AccountToken common.Token `json:"account_token"`
}

type ParamGetAccountByUserID struct {
	List  []ParamAccountInfo `json:"list,omitempty"`
	Total int64              `json:"total,omitempty"`
}

type ParamFriendInfo struct {
	ParamAccountInfo
	RelationID int64 `json:"relation_id,omitempty"`
}

type ParamGetAccountByID struct {
	Info       ParamAccountInfo `json:"info"`
	Signature  string           `json:"signature,omitempty"`
	CreateAt   time.Time        `json:"create_at,omitempty"`
	RelationID int64            `json:"relation_id"`
}

type ParamGetAccountsByName struct {
	List  []*ParamFriendInfo `json:"list,omitempty"`  // 账号列表
	Total int64              `json:"total,omitempty"` // 总数
}
