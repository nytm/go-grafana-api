package gapi

import (
	"errors"
	"fmt"
	"strconv"
)

type PermissionType int8

const (
	None  PermissionType = 0
	View  PermissionType = 1
	Edit  PermissionType = 2
	Admin PermissionType = 4
)

func NewPermissionType(p int) (PermissionType, error) {
	switch p {
	case 1:
		return View, nil
	case 2:
		return Edit, nil
	case 4:
		return Admin, nil
	}
	return None, errors.New(fmt.Sprintf("Unknow permission: %d", p))
}

func (p *PermissionType) Value() int8 {
	switch *p {
	case View:
		return 1
	case Edit:
		return 2
	case Admin:
		return 4
	}
	return 0
}

func (p *PermissionType) String() string {
	return strconv.FormatInt(int64(p.Value()), 10)
}
