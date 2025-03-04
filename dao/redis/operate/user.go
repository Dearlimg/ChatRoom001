package operate

import (
	"ChatRoom001/global"
	"github.com/Dearlimg/Goutils/pkg/utils"
	"github.com/gin-gonic/gin"
)

var UserKey = "user"

func (r *RDB) SaveUserToken(ctx *gin.Context, userID int64, tokens []string) error {
	key := utils.LinkStr(UserKey, utils.IDToString(userID))
	for _, token := range tokens {
		if err := r.rdb.SAdd(ctx, key, token).Err(); err != nil {
			return err
		}
		r.rdb.Expire(ctx, key, global.PrivateSetting.Token.AccessTokenExpire)
	}
	return nil
}
