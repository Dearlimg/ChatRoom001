// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package db

import (
	"context"
	"database/sql"
)

type Querier interface {
	CheckRelationTypeByID(ctx context.Context, id int64) (bool, error)
	CountAccountByUserID(ctx context.Context, userID int64) (int64, error)
	CreateAccount(ctx context.Context, arg *CreateAccountParams) error
	CreateApplication(ctx context.Context, arg *CreateApplicationParams) error
	// -- name: CreateFile :one
	// START TRANSACTION;
	// INSERT INTO files (
	//     file_name, file_type, file_size, `key`, url, relation_id, account_id
	// ) VALUES (
	//              ?, ?, ?, ?, ?, ?, ?
	//          );
	// SELECT * FROM files
	// WHERE file_id = LAST_INSERT_ID();
	// COMMIT;
	CreateFile(ctx context.Context, arg *CreateFileParams) error
	CreateFriendRelation(ctx context.Context, arg *CreateFriendRelationParams) error
	CreateGet(ctx context.Context, account1ID int64) (*Application, error)
	CreateGroupNotify(ctx context.Context, arg *CreateGroupNotifyParams) error
	// -- name: CreateGroupNotifyReturn :one
	// select *
	// from group_notify
	// where create_at=?
	// and relation_id=?
	// and account_id=?;
	CreateGroupNotifyReturn(ctx context.Context) (*GroupNotify, error)
	CreateGroupRelation(ctx context.Context, arg *CreateGroupRelationParams) error
	CreateGroupRelationReturn(ctx context.Context, arg *CreateGroupRelationReturnParams) (int64, error)
	CreateManySetting(ctx context.Context, arg *CreateManySettingParams) error
	CreateMessage(ctx context.Context, arg *CreateMessageParams) error
	CreateMessageReturn(ctx context.Context) (*CreateMessageReturnRow, error)
	// -- name: CreateRelationReturn :one
	// select id
	// from relations
	// where last_insert_id();
	CreateRelationReturn(ctx context.Context, arg *CreateRelationReturnParams) (int64, error)
	CreateSetting(ctx context.Context, arg *CreateSettingParams) error
	CreateUser(ctx context.Context, arg *CreateUserParams) error
	DeleteAccount(ctx context.Context, id int64) error
	DeleteAccountByUserID(ctx context.Context, userID int64) error
	DeleteApplication(ctx context.Context, arg *DeleteApplicationParams) error
	DeleteFileByID(ctx context.Context, id int64) error
	DeleteFriendRelation(ctx context.Context, arg *DeleteFriendRelationParams) error
	DeleteFriendRelationByAccountID(ctx context.Context, account1ID sql.NullInt64) error
	DeleteGroup(ctx context.Context, relationID int64) error
	DeleteGroupNotify(ctx context.Context, id int64) error
	DeleteRelation(ctx context.Context, id int64) error
	DeleteSetting(ctx context.Context, arg *DeleteSettingParams) error
	DeleteSettingsByAccountID(ctx context.Context, accountID int64) error
	DeleteUser(ctx context.Context, id int64) error
	ExistAccountByID(ctx context.Context, id int64) (bool, error)
	ExistEmail(ctx context.Context, email string) (bool, error)
	ExistGroupLeaderByAccountIDWithLock(ctx context.Context, accountID int64) (bool, error)
	ExistsAccountByNameAndUserID(ctx context.Context, arg *ExistsAccountByNameAndUserIDParams) (bool, error)
	ExistsApplicationByIDWithLock(ctx context.Context, arg *ExistsApplicationByIDWithLockParams) (bool, error)
	ExistsFriendRelation(ctx context.Context, arg *ExistsFriendRelationParams) (bool, error)
	ExistsFriendSetting(ctx context.Context, arg *ExistsFriendSettingParams) (bool, error)
	ExistsIsLeader(ctx context.Context, arg *ExistsIsLeaderParams) (bool, error)
	ExistsSetting(ctx context.Context, arg *ExistsSettingParams) (bool, error)
	ExistsUserByID(ctx context.Context, id int64) (bool, error)
	GetAccountByID(ctx context.Context, arg *GetAccountByIDParams) (*GetAccountByIDRow, error)
	GetAccountByUserID(ctx context.Context, userID int64) ([]*GetAccountByUserIDRow, error)
	GetAccountIDsByRelationID(ctx context.Context, relationID int64) ([]int64, error)
	// -- name: GetAccountsByName :many
	// SELECT
	//     a.*,
	//     r.id AS relation_id,
	//     (SELECT COUNT(*) FROM accounts WHERE name LIKE CONCAT('%', ?, '%')) AS total
	// FROM (
	//          SELECT id, name, avatar, gender
	//          FROM accounts
	//          WHERE name LIKE CONCAT('%', ?, '%')
	//      ) AS a
	//          LEFT JOIN relations r
	//                    ON r.relation_type = 'friend'
	//                        AND (
	//                           (r.account1_id = a.id AND r.account2_id = ?)
	//                               OR
	//                           (r.account1_id = ? AND r.account2_id = a.id)
	//                           )
	// LIMIT ? OFFSET ?;
	GetAccountNameByID(ctx context.Context, id int64) (string, error)
	GetAccountsByName(ctx context.Context, arg *GetAccountsByNameParams) ([]*GetAccountsByNameRow, error)
	GetAcountIDsByUserID(ctx context.Context, userID int64) ([]int64, error)
	GetAllEmail(ctx context.Context) ([]string, error)
	GetAllGroupRelation(ctx context.Context) ([]int64, error)
	GetAllRelationIDs(ctx context.Context) ([]int64, error)
	GetAllRelationOnRelation(ctx context.Context) ([]*Relation, error)
	// -- name: GetApplicationByID :one
	// select *
	// from applications
	// where account1_id = ?
	//   and account2_id = ?
	// limit  1;
	GetApplicationByID(ctx context.Context, arg *GetApplicationByIDParams) (*Application, error)
	GetApplications(ctx context.Context, arg *GetApplicationsParams) ([]*GetApplicationsRow, error)
	GetCreateFile(ctx context.Context, key string) (*File, error)
	GetFileByRelation(ctx context.Context, relationID sql.NullInt64) ([]*File, error)
	GetFileByRelationID(ctx context.Context, relationID sql.NullInt64) ([]*File, error)
	GetFileDetailsByID(ctx context.Context, id int64) (*File, error)
	GetFileKeyByID(ctx context.Context, id int64) (string, error)
	GetFriendPinSettingsOrderByPinTime(ctx context.Context, arg *GetFriendPinSettingsOrderByPinTimeParams) ([]*GetFriendPinSettingsOrderByPinTimeRow, error)
	GetFriendRelationByID(ctx context.Context, id int64) (*Relation, error)
	GetFriendSettingsByName(ctx context.Context, arg *GetFriendSettingsByNameParams) ([]*GetFriendSettingsByNameRow, error)
	// SELECT
	//     s.relation_id,
	//     s.nick_name,
	//     s.is_not_disturb,
	//     s.is_pin,
	//     s.pin_time,
	//     s.is_show,
	//     s.last_show,
	//     s.is_self,
	//     a.id AS account_id,
	//     a.name AS account_name,
	//     a.avatar AS account_avatar
	// FROM (
	//          SELECT
	//              st.relation_id,
	//              st.nick_name,
	//              st.is_not_disturb,
	//              st.is_pin,
	//              st.pin_time,
	//              st.is_show,
	//              st.last_show,
	//              st.is_self,
	//              CASE
	//                  WHEN rt.relation_type = 'friend' THEN rt.account2_id
	//                  ELSE NULL
	//                  END AS friend_account_id
	//          FROM settings st
	//                   INNER JOIN relations rt
	//                              ON st.relation_id = rt.id
	//                                  AND rt.relation_type = 'friend'
	//          WHERE
	//              st.account_id = ?
	//            AND st.is_show = true
	//      ) s
	//          INNER JOIN accounts a
	//                     ON a.id = s.friend_account_id
	// ORDER BY
	//     s.last_show DESC;
	//
	// -- name: GetFriendShowSettingsOrderByShowTime :many
	// SELECT
	//     s.relation_id,
	//     s.nick_name,
	//     s.is_not_disturb,
	//     s.is_pin,
	//     s.pin_time,
	//     s.is_show,
	//     s.last_show,
	//     s.is_self,
	//     a.id AS account_id,
	//     a.name AS account_name,
	//     a.avatar AS account_avatar
	// FROM (
	//          SELECT
	//              settings.relation_id,
	//              settings.nick_name,
	//              settings.is_not_disturb,
	//              settings.is_pin,
	//              settings.pin_time,
	//              settings.is_show,
	//              settings.last_show,
	//              settings.is_self
	//          FROM
	//              settings
	//                  JOIN
	//              relations ON settings.relation_id = relations.id
	//          WHERE
	//              settings.account_id =?  -- 这里的? 是占位符，实际使用时需要替换为具体的值
	//            AND settings.is_show = true
	//            AND relations.relation_type = 'friend'
	//      ) AS s
	//          JOIN
	//      accounts a ON a.id = (
	//          SELECT
	//              sub_settings.account_id  -- 给子查询中的 settings 表取别名 sub_settings
	//          FROM
	//              settings AS sub_settings  -- 给子查询中的 settings 表取别名 sub_settings
	//          WHERE
	//              sub_settings.relation_id = s.relation_id
	//            AND (sub_settings.account_id !=? OR sub_settings.is_self = true)  -- 这里的? 是占位符，实际使用时需要替换为具体的值
	//      )
	// ORDER BY
	//     s.last_show DESC;
	//
	// -- name: GetFriendShowSettingsOrderByShowTime :many
	// select s.*,
	//        r.id,r.group_avatar,r.group_name,r.group_description,r.created_at
	// from (select relation_id,
	//              nick_name,
	//              is_not_disturb,
	//              is_pin,
	//              pin_time,
	//              is_show,
	//              last_show,
	//              is_self
	//       from settings,
	//            relations
	//       where settings.account_id = ?
	//         and settings.relation_id = relations.id
	//         and settings.is_show = true
	//         and relations.relation_type = 'friend') as s,
	//      relations r
	// where r.id = (select relation_id from settings where relation_id = s.relation_id and settings.account_id = ? limit 1)
	// order by s.last_show desc;
	GetFriendSettingsOrderByName(ctx context.Context, arg *GetFriendSettingsOrderByNameParams) ([]*GetFriendSettingsOrderByNameRow, error)
	GetFriendShowSettingsOrderByShowTime(ctx context.Context, arg *GetFriendShowSettingsOrderByShowTimeParams) ([]*GetFriendShowSettingsOrderByShowTimeRow, error)
	GetGroupAvatar(ctx context.Context, relationID sql.NullInt64) (*File, error)
	GetGroupList(ctx context.Context, arg *GetGroupListParams) ([]*GetGroupListRow, error)
	GetGroupMembers(ctx context.Context, relationID int64) ([]int64, error)
	GetGroupMembersByID(ctx context.Context, arg *GetGroupMembersByIDParams) ([]*GetGroupMembersByIDRow, error)
	GetGroupNotifyByID(ctx context.Context, relationID sql.NullInt64) ([]*GetGroupNotifyByIDRow, error)
	// select s.*,
	//        a.id as account_id,
	//        a.name as account_name,
	//        a.avatar as account_avatar
	// from (select settings.relation_id, settings.nick_name, settings.pin_time,settings.is_pin,settings.is_show,settings.is_not_disturb
	//       from settings,
	//            relations
	//       where settings.account_id = ?
	//         and settings.is_pin = true
	//         and settings.relation_id = relations.id
	//         and relations.relation_type = 'friend') as s,
	//      accounts a
	// where a.id = (select account_id from settings where relation_id = s.relation_id and (settings.account_id != ? or is_self = true))
	// order by s.pin_time;
	GetGroupPinSettingsOrderByPinTime(ctx context.Context, arg *GetGroupPinSettingsOrderByPinTimeParams) ([]*GetGroupPinSettingsOrderByPinTimeRow, error)
	GetGroupRelationByID(ctx context.Context, id int64) (*GetGroupRelationByIDRow, error)
	// SELECT
	//     s.*,
	//     a.id AS account_id,
	//     a.name AS account_name,
	//     a.avatar AS account_avatar,
	//     COUNT(*) OVER () AS total
	// FROM (
	//          SELECT
	//              s.relation_id,
	//              s.nick_name,
	//              s.is_not_disturb,
	//              s.is_pin,
	//              s.pin_time,
	//              s.is_show,
	//              s.last_show,
	//              s.is_self
	//          FROM settings s
	//                   INNER JOIN relations r ON s.relation_id = r.id
	//          WHERE
	//              s.account_id = ?
	//            AND r.relation_type = 'friend'
	//      ) AS s
	//          CROSS JOIN accounts a
	// WHERE
	//     a.id = (
	//         SELECT account_id
	//         FROM settings
	//         WHERE
	//             relation_id = s.relation_id
	//           AND (settings.account_id != ? OR s.is_self = 1)
	//     )
	//   AND (
	//     a.name LIKE CONCAT('%', ?, '%')
	//         OR s.nick_name LIKE CONCAT('%', ?, '%')
	//     )
	// ORDER BY a.name
	// LIMIT ? OFFSET ?;
	GetGroupSettingsByName(ctx context.Context, arg *GetGroupSettingsByNameParams) ([]*GetGroupSettingsByNameRow, error)
	GetGroupShowSettingsOrderByShowTime(ctx context.Context, arg *GetGroupShowSettingsOrderByShowTimeParams) ([]*GetGroupShowSettingsOrderByShowTimeRow, error)
	GetLastMessageByRelation(ctx context.Context, relationID int64) (*GetLastMessageByRelationRow, error)
	GetMessageByID(ctx context.Context, id int64) (*GetMessageByIDRow, error)
	GetMsgByRelationIDAndTime(ctx context.Context, arg *GetMsgByRelationIDAndTimeParams) ([]*GetMsgByRelationIDAndTimeRow, error)
	GetMsgsByContent(ctx context.Context, arg *GetMsgsByContentParams) ([]*GetMsgsByContentRow, error)
	GetMsgsByContentAndRelation(ctx context.Context, arg *GetMsgsByContentAndRelationParams) ([]*GetMsgsByContentAndRelationRow, error)
	GetPinMsgsByRelationID(ctx context.Context, arg *GetPinMsgsByRelationIDParams) ([]*GetPinMsgsByRelationIDRow, error)
	GetRelationIDByAccountID(ctx context.Context, arg *GetRelationIDByAccountIDParams) (int64, error)
	GetRelationIDByAccountIDFromSettings(ctx context.Context, accountID int64) ([]int64, error)
	GetRelationIDByInfo(ctx context.Context, arg *GetRelationIDByInfoParams) (int64, error)
	GetRlyMsgsInfoByMsgID(ctx context.Context, arg *GetRlyMsgsInfoByMsgIDParams) ([]*GetRlyMsgsInfoByMsgIDRow, error)
	GetSettingByID(ctx context.Context, arg *GetSettingByIDParams) (*Setting, error)
	GetTopMsgByRelationID(ctx context.Context, arg *GetTopMsgByRelationIDParams) (*GetTopMsgByRelationIDRow, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByID(ctx context.Context, id int64) (*User, error)
	OfferMsgsByAccountIDAndTime(ctx context.Context, arg *OfferMsgsByAccountIDAndTimeParams) ([]*OfferMsgsByAccountIDAndTimeRow, error)
	TransferIsLeaderFalse(ctx context.Context, arg *TransferIsLeaderFalseParams) error
	// select s.*,
	//        r.id as realtion_id,
	//        r.group_name AS group_name,
	//        r.group_avatar AS group_avatar,
	//        r.group_description AS description,
	//        count(*) over () as total
	// from (select relation_id,
	//     nick_name,
	//     is_not_disturb,
	//     is_pin,
	//     pin_time,
	//     is_show,
	//     last_show,
	//     is_self
	//     from settings,
	//     relations
	//     where settings.account_id = ?
	//     and settings.relation_id = relations.id
	//     and relations.relation_type = 'group') as s,
	//     relations r
	// where r.id = (select s.relation_id from settings where s.relation_id=s.relation_id and (settings.account_id=?))
	//     and ((r.group_name like ('%' || ? || '%')))
	// order by (r.group_name)
	// limit ? offset ?;
	TransferIsLeaderTrue(ctx context.Context, arg *TransferIsLeaderTrueParams) error
	UpdateAccount(ctx context.Context, arg *UpdateAccountParams) error
	UpdateAccountAvatar(ctx context.Context, arg *UpdateAccountAvatarParams) error
	UpdateApplication(ctx context.Context, arg *UpdateApplicationParams) error
	UpdateGroupAvatar(ctx context.Context, arg *UpdateGroupAvatarParams) error
	UpdateGroupNotify(ctx context.Context, arg *UpdateGroupNotifyParams) error
	UpdateGroupNotifyReturn(ctx context.Context, id int64) (*GroupNotify, error)
	UpdateGroupRelation(ctx context.Context, arg *UpdateGroupRelationParams) error
	UpdateMsgPin(ctx context.Context, arg *UpdateMsgPinParams) error
	// -- name: UpdateMsgReads :exec
	// UPDATE messages AS m
	// SET read_ids =
	//         CASE
	//             WHEN JSON_CONTAINS(read_ids, CAST(@accountID AS JSON)) = 0
	//                 THEN JSON_ARRAY_APPEND(read_ids, '$', @accountID)
	//             ELSE read_ids
	//             END
	// WHERE relation_id = ?
	//   AND JSON_CONTAINS(@target_ids, CAST(m.id AS JSON));
	// 先执行更新操作
	UpdateMsgReads(ctx context.Context, relationID int64) error
	// 再查询受影响的行
	UpdateMsgReadsReturn(ctx context.Context, relationID int64) ([]*UpdateMsgReadsReturnRow, error)
	UpdateMsgRevoke(ctx context.Context, arg *UpdateMsgRevokeParams) error
	UpdateMsgTop(ctx context.Context, arg *UpdateMsgTopParams) error
	UpdateSettingDisturb(ctx context.Context, arg *UpdateSettingDisturbParams) error
	UpdateSettingNickName(ctx context.Context, arg *UpdateSettingNickNameParams) error
	UpdateSettingPin(ctx context.Context, arg *UpdateSettingPinParams) error
	UpdateSettingShow(ctx context.Context, arg *UpdateSettingShowParams) error
	UpdateUser(ctx context.Context, arg *UpdateUserParams) error
}

var _ Querier = (*Queries)(nil)
