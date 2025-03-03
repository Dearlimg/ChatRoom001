// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package db

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type AccountsGender string

const (
	AccountsGenderValue0 AccountsGender = "男"
	AccountsGenderValue1 AccountsGender = "女"
	AccountsGenderValue2 AccountsGender = "未知"
)

func (e *AccountsGender) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = AccountsGender(s)
	case string:
		*e = AccountsGender(s)
	default:
		return fmt.Errorf("unsupported scan type for AccountsGender: %T", src)
	}
	return nil
}

type NullAccountsGender struct {
	AccountsGender AccountsGender `json:"accounts_gender"`
	Valid          bool           `json:"valid"` // Valid is true if AccountsGender is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullAccountsGender) Scan(value interface{}) error {
	if value == nil {
		ns.AccountsGender, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.AccountsGender.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullAccountsGender) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.AccountsGender), nil
}

type ApplicationsStatus string

const (
	ApplicationsStatusValue0 ApplicationsStatus = "已申请"
	ApplicationsStatusValue1 ApplicationsStatus = "已同意"
	ApplicationsStatusValue2 ApplicationsStatus = "已拒绝"
)

func (e *ApplicationsStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = ApplicationsStatus(s)
	case string:
		*e = ApplicationsStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for ApplicationsStatus: %T", src)
	}
	return nil
}

type NullApplicationsStatus struct {
	ApplicationsStatus ApplicationsStatus `json:"applications_status"`
	Valid              bool               `json:"valid"` // Valid is true if ApplicationsStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullApplicationsStatus) Scan(value interface{}) error {
	if value == nil {
		ns.ApplicationsStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.ApplicationsStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullApplicationsStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.ApplicationsStatus), nil
}

type FilesFileType string

const (
	FilesFileTypeImg  FilesFileType = "img"
	FilesFileTypeFile FilesFileType = "file"
)

func (e *FilesFileType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = FilesFileType(s)
	case string:
		*e = FilesFileType(s)
	default:
		return fmt.Errorf("unsupported scan type for FilesFileType: %T", src)
	}
	return nil
}

type NullFilesFileType struct {
	FilesFileType FilesFileType `json:"files_file_type"`
	Valid         bool          `json:"valid"` // Valid is true if FilesFileType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullFilesFileType) Scan(value interface{}) error {
	if value == nil {
		ns.FilesFileType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.FilesFileType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullFilesFileType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.FilesFileType), nil
}

type MessagesMsgType string

const (
	MessagesMsgTypeText MessagesMsgType = "text"
	MessagesMsgTypeFile MessagesMsgType = "file"
)

func (e *MessagesMsgType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = MessagesMsgType(s)
	case string:
		*e = MessagesMsgType(s)
	default:
		return fmt.Errorf("unsupported scan type for MessagesMsgType: %T", src)
	}
	return nil
}

type NullMessagesMsgType struct {
	MessagesMsgType MessagesMsgType `json:"messages_msg_type"`
	Valid           bool            `json:"valid"` // Valid is true if MessagesMsgType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullMessagesMsgType) Scan(value interface{}) error {
	if value == nil {
		ns.MessagesMsgType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.MessagesMsgType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullMessagesMsgType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.MessagesMsgType), nil
}

type MessagesNotifyType string

const (
	MessagesNotifyTypeSystem MessagesNotifyType = "system"
	MessagesNotifyTypeCommon MessagesNotifyType = "common"
)

func (e *MessagesNotifyType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = MessagesNotifyType(s)
	case string:
		*e = MessagesNotifyType(s)
	default:
		return fmt.Errorf("unsupported scan type for MessagesNotifyType: %T", src)
	}
	return nil
}

type NullMessagesNotifyType struct {
	MessagesNotifyType MessagesNotifyType `json:"messages_notify_type"`
	Valid              bool               `json:"valid"` // Valid is true if MessagesNotifyType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullMessagesNotifyType) Scan(value interface{}) error {
	if value == nil {
		ns.MessagesNotifyType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.MessagesNotifyType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullMessagesNotifyType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.MessagesNotifyType), nil
}

type MsgNotificationsMsgType string

const (
	MsgNotificationsMsgTypeSystem MsgNotificationsMsgType = "system"
	MsgNotificationsMsgTypeCommon MsgNotificationsMsgType = "common"
)

func (e *MsgNotificationsMsgType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = MsgNotificationsMsgType(s)
	case string:
		*e = MsgNotificationsMsgType(s)
	default:
		return fmt.Errorf("unsupported scan type for MsgNotificationsMsgType: %T", src)
	}
	return nil
}

type NullMsgNotificationsMsgType struct {
	MsgNotificationsMsgType MsgNotificationsMsgType `json:"msg_notifications_msg_type"`
	Valid                   bool                    `json:"valid"` // Valid is true if MsgNotificationsMsgType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullMsgNotificationsMsgType) Scan(value interface{}) error {
	if value == nil {
		ns.MsgNotificationsMsgType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.MsgNotificationsMsgType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullMsgNotificationsMsgType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.MsgNotificationsMsgType), nil
}

type RelationsRelationType string

const (
	RelationsRelationTypeGroup  RelationsRelationType = "group"
	RelationsRelationTypeFriend RelationsRelationType = "friend"
)

func (e *RelationsRelationType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = RelationsRelationType(s)
	case string:
		*e = RelationsRelationType(s)
	default:
		return fmt.Errorf("unsupported scan type for RelationsRelationType: %T", src)
	}
	return nil
}

type NullRelationsRelationType struct {
	RelationsRelationType RelationsRelationType `json:"relations_relation_type"`
	Valid                 bool                  `json:"valid"` // Valid is true if RelationsRelationType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullRelationsRelationType) Scan(value interface{}) error {
	if value == nil {
		ns.RelationsRelationType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.RelationsRelationType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullRelationsRelationType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.RelationsRelationType), nil
}

type Account struct {
	ID        int64          `json:"id"`
	UserID    int64          `json:"user_id"`
	Name      string         `json:"name"`
	Avatar    string         `json:"avatar"`
	Gender    AccountsGender `json:"gender"`
	Signature string         `json:"signature"`
	CreateAt  time.Time      `json:"create_at"`
}

type Application struct {
	ID       int64              `json:"id"`
	Status   ApplicationsStatus `json:"status"`
	CreateAt time.Time          `json:"create_at"`
}

type File struct {
	ID       int64         `json:"id"`
	FileType FilesFileType `json:"file_type"`
	FilePath string        `json:"file_path"`
	CreateAt time.Time     `json:"create_at"`
}

type Group struct {
	ID          int32          `json:"id"`
	Name        string         `json:"name"`
	Description sql.NullString `json:"description"`
	Avatar      sql.NullString `json:"avatar"`
}

type GroupNotify struct {
	ID            int64           `json:"id"`
	RelationID    sql.NullInt64   `json:"relation_id"`
	MsgContent    string          `json:"msg_content"`
	MsgExpand     json.RawMessage `json:"msg_expand"`
	AccountID     sql.NullInt64   `json:"account_id"`
	CreateAt      time.Time       `json:"create_at"`
	ReadIds       json.RawMessage `json:"read_ids"`
	MsgContentTsv sql.NullString  `json:"msg_content_tsv"`
}

type Message struct {
	ID            int64              `json:"id"`
	NotifyType    MessagesNotifyType `json:"notify_type"`
	MsgType       MessagesMsgType    `json:"msg_type"`
	MsgContent    string             `json:"msg_content"`
	MsgExtend     json.RawMessage    `json:"msg_extend"`
	FileID        sql.NullInt64      `json:"file_id"`
	AccountID     sql.NullInt64      `json:"account_id"`
	RlyMsgID      sql.NullInt64      `json:"rly_msg_id"`
	RelationID    int64              `json:"relation_id"`
	CreateAt      time.Time          `json:"create_at"`
	IsRevoke      bool               `json:"is_revoke"`
	IsTop         bool               `json:"is_top"`
	IsPin         bool               `json:"is_pin"`
	PinTime       time.Time          `json:"pin_time"`
	ReadIds       json.RawMessage    `json:"read_ids"`
	MsgContentTsy sql.NullString     `json:"msg_content_tsy"`
}

type MsgNotification struct {
	ID       int64                   `json:"id"`
	MsgType  MsgNotificationsMsgType `json:"msg_type"`
	Content  string                  `json:"content"`
	CreateAt time.Time               `json:"create_at"`
}

type Relation struct {
	ID               int64                 `json:"id"`
	RelationType     RelationsRelationType `json:"relation_type"`
	GroupName        sql.NullString        `json:"group_name"`
	GroupDescription sql.NullString        `json:"group_description"`
	Account1ID       sql.NullInt64         `json:"account1_id"`
	Account2ID       sql.NullInt64         `json:"account2_id"`
	CreateAt         time.Time             `json:"create_at"`
}

type Setting struct {
	AccountID    int64     `json:"account_id"`
	RelationID   int64     `json:"relation_id"`
	NickName     string    `json:"nick_name"`
	IsNotDisturb bool      `json:"is_not_disturb"`
	IsPin        bool      `json:"is_pin"`
	PinTime      time.Time `json:"pin_time"`
	IsShow       bool      `json:"is_show"`
	LastShow     time.Time `json:"last_show"`
	IsLeader     bool      `json:"is_leader"`
	IsSelf       bool      `json:"is_self"`
}

type User struct {
	ID       int64     `json:"id"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
	CreateAt time.Time `json:"create_at"`
}
