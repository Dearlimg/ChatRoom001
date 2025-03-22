package server

import "ChatRoom001/model/reply"

type SendMsg struct {
	EnToken string                    `json:"en_token,omitempty"` // 加密后的 token
	MsgInfo reply.ParamMsgInfoWithRly `json:"msg_info"`
}

type ReadMsg struct {
	EnToken  string  `json:"en_token,omitempty"` // 加密后的 Token
	MsgIDs   []int64 `json:"msg_ids"`            // 已读消息 IDs
	ReaderID int64   `json:"reader_id"`          // 读者账号 ID
}
