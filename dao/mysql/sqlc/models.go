// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package db

import (
	"database/sql"
	"database/sql/driver"
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
	AccountsGender AccountsGender
	Valid          bool // Valid is true if AccountsGender is not NULL
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
	ApplicationsStatus ApplicationsStatus
	Valid              bool // Valid is true if ApplicationsStatus is not NULL
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
	FilesFileType FilesFileType
	Valid         bool // Valid is true if FilesFileType is not NULL
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
	MsgNotificationsMsgType MsgNotificationsMsgType
	Valid                   bool // Valid is true if MsgNotificationsMsgType is not NULL
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

type RelationsGroupType string

const (
	RelationsGroupTypePublic  RelationsGroupType = "public"
	RelationsGroupTypePrivate RelationsGroupType = "private"
)

func (e *RelationsGroupType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = RelationsGroupType(s)
	case string:
		*e = RelationsGroupType(s)
	default:
		return fmt.Errorf("unsupported scan type for RelationsGroupType: %T", src)
	}
	return nil
}

type NullRelationsGroupType struct {
	RelationsGroupType RelationsGroupType
	Valid              bool // Valid is true if RelationsGroupType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullRelationsGroupType) Scan(value interface{}) error {
	if value == nil {
		ns.RelationsGroupType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.RelationsGroupType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullRelationsGroupType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.RelationsGroupType), nil
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
	RelationsRelationType RelationsRelationType
	Valid                 bool // Valid is true if RelationsRelationType is not NULL
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
	ID        int64
	UserID    int64
	Name      string
	Avatar    string
	Gender    AccountsGender
	Signature string
	CreateAt  time.Time
}

type Application struct {
	ID       int64
	Status   ApplicationsStatus
	CreateAt time.Time
}

type File struct {
	ID       int64
	FileType FilesFileType
	FilePath string
	CreateAt time.Time
}

type Group struct {
	ID          int32
	Name        string
	Description sql.NullString
	Avatar      sql.NullString
}

type MsgNotification struct {
	ID       int64
	MsgType  MsgNotificationsMsgType
	Content  string
	CreateAt time.Time
}

type Relation struct {
	ID           int64
	RelationType RelationsRelationType
	// 群组类型，仅 relation_type=group 时有效
	GroupType NullRelationsGroupType
	// 好友账号1 ID，仅 relation_type=friend 时有效
	FriendAccount1ID sql.NullInt64
	// 好友账号2 ID，仅 relation_type=friend 时有效
	FriendAccount2ID sql.NullInt64
	CreatedAt        sql.NullTime
}

type User struct {
	ID       int64
	Email    string
	Password string
	CreateAt time.Time
}
