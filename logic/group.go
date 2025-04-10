package logic

import (
	"github.com/Dearlimg/Goutils/pkg/app/errcode"
	"github.com/gin-gonic/gin"
)

type group struct{}

func (group) CreateGroup(ctx *gin.Context, accountID int64, name, description string) (relationID int64, err errcode.Err) {

	return 0, nil
}
