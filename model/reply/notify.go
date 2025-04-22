package reply

import (
	"encoding/json"
	"time"
)

type ParamGroupNotify struct {
	ID         int64           `json:"id"`
	RelationID int64           `json:"relation_id"`
	MsgContent string          `json:"msg_content"`
	MsgExtend  json.RawMessage `json:"msg_extend"`
	AccountID  int64           `json:"account_id"`
	CreateAt   time.Time       `json:"create_at"`
	ReadIDs    json.RawMessage `json:"read_ids"`
}

type ParamGetNotifyByID struct {
	List  []ParamGroupNotify `json:"list"`
	Total int64              `json:"total"`
}
