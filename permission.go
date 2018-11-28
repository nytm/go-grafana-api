package gapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

type PermissionType int8

const (
	None            PermissionType = 0
	PermissionView  PermissionType = 1
	PermissionEdit  PermissionType = 2
	PermissionAdmin PermissionType = 4
)

func NewPermissionType(p int) (PermissionType, error) {
	switch p {
	case 1:
		return PermissionView, nil
	case 2:
		return PermissionEdit, nil
	case 4:
		return PermissionAdmin, nil
	}
	return None, errors.New(fmt.Sprintf("Unknow permission: %d", p))
}

func (p *PermissionType) Value() int8 {
	switch *p {
	case PermissionView:
		return 1
	case PermissionEdit:
		return 2
	case PermissionAdmin:
		return 4
	}
	return 0
}

func (d *PermissionType) UnmarshalJSON(b []byte) error {
	var s int64
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	var err error
	*d, err = NewPermissionType(int(s))
	return err
}

func (d *PermissionType) MarshalJSON() ([]byte, error) {
	return []byte(d.String()), nil
}

func (p *PermissionType) String() string {
	return strconv.FormatInt(int64(p.Value()), 10)
}

type Permission struct {
	Id         int64  `json:"id,omitempty"`
	FolderUid  string `json:"folderUid,omitempty"`
	UserId     int64  `json:"userId,omitempty"`
	TeamId     int64  `json:"teamId,omitempty"`
	Role       string `json:"role,omitempty"`
	Permission int    `json:"permission,omitempty"`
	IsFolder   bool   `json:"isFolder,omitempty"`
}
