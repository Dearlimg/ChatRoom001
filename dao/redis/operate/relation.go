package operate

import (
	"context"
	"encoding/json"
	"github.com/Dearlimg/Goutils/pkg/utils"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

const KeyGroup = "KeyGroup"

// DelAllPrefixLua 这段代码是一个 Lua 脚本，用于在 Redis 中批量删除指定前缀的键。
const DelAllPrefixLua = "local redisKeys = redis.call('keys', KEYS[1] .. '*');for i, k in pairs(redisKeys) do redis.call('expire', k, 0);end"

func (r *RDB) ReloadRelationIDs(ctx context.Context, groupMap map[int64][]int64) error {
	groupMapJson, err := json.Marshal(groupMap)
	if err != nil {
		return err
	}
	// Lua 脚本功能，允许开发者在服务器直接执行脚本，这样可以减少网络延迟和往返次数，提高操作的原子性和效率
	if err = r.rdb.Eval(ctx, DelAllPrefixLua, []string{KeyGroup}, string(groupMapJson)).Err(); err != nil && !errors.Is(err, redis.Nil) {
		return err
	}

	pipe := r.rdb.Pipeline()
	for relationID, ids := range groupMap {
		data := make([]interface{}, len(ids))
		for i, id := range ids {
			data[i] = utils.IDToString(id)
		}
		r.rdb.SAdd(ctx, utils.LinkStr(KeyGroup, utils.IDToString(relationID)), data...)

	}
	_, err = pipe.Exec(ctx)
	return err
}

func (r *RDB) AddRelationAccount(ctx context.Context, relationID int64, accountIDs ...int64) (err error) {
	if len(accountIDs) == 0 {
		return nil
	}
	data := make([]interface{}, len(accountIDs))
	for i, id := range accountIDs {
		data[i] = utils.IDToString(id)
	}
	return r.rdb.SAdd(ctx, utils.LinkStr(KeyGroup, utils.IDToString(relationID)), data...).Err()
}

func (r *RDB) DeleteRelations(ctx context.Context, relationIDs ...int64) (err error) {
	if len(relationIDs) == 0 {
		return nil
	}
	pipe := r.rdb.TxPipeline()
	for _, relationID := range relationIDs {
		pipe.Del(ctx, utils.LinkStr(KeyGroup, utils.IDToString(relationID)))
	}
	_, err = pipe.Exec(ctx)
	return
}

func (r *RDB) DeleteAccountFromRelations(ctx context.Context, accountID int64, relationIDs ...int64) error {
	if len(relationIDs) == 0 {
		return nil
	}
	pipe := r.rdb.TxPipeline()
	for _, relationID := range relationIDs {
		pipe.SRem(ctx, utils.LinkStr(KeyGroup, utils.IDToString(relationID)), utils.IDToString(accountID))
	}
	_, err := pipe.Exec(ctx)
	return err
}

// DeleteRelationAccount 从一个群聊中删除多个成员
func (r *RDB) DeleteRelationAccount(ctx context.Context, relationID int64, accountIDs ...int64) error {
	if len(accountIDs) == 0 {
		return nil
	}
	data := make([]interface{}, len(accountIDs))
	for i, v := range accountIDs {
		data[i] = utils.IDToString(v)
	}
	return r.rdb.SRem(ctx, utils.LinkStr(KeyGroup, utils.IDToString(relationID)), data...).Err()
}

func (r *RDB) GetAllAccountsByRelationID(ctx context.Context, relationID int64) ([]int64, error) {
	id := utils.IDToString(relationID)
	key := utils.LinkStr(KeyGroup, id)
	accountsIDStr, err := r.rdb.SMembers(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	accountIDs := make([]int64, 0, len(accountsIDStr))
	for _, id := range accountsIDStr {
		accountID := utils.StringToIDMust(id)
		accountIDs = append(accountIDs, accountID)
	}
	return accountIDs, nil
}
