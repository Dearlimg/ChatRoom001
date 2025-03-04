package model

import (
	"encoding/json"
	"github.com/Dearlimg/Goutils/pkg/token"
)

type TokenType string

type Content struct {
	TokenType TokenType `json:"token_type,omitempty"`
	ID        int64     `json:"id,omitempty"`
}

type Token struct {
	AccessToken string
	PayLoad     *token.Payload
	Content     *Content
}

const (
	UserToken    TokenType = "user"
	AccountToken TokenType = "account"
)

func NewTokenContent(t TokenType, id int64) *Content {
	return &Content{
		TokenType: t,
		ID:        id,
	}
}

func (c *Content) Marshal() ([]byte, error) {
	return json.Marshal(c)
}

func (c *Content) Unmarshal(data []byte) error {
	if err := json.Unmarshal(data, &c); err != nil {
		return err
	}
	return nil
}
