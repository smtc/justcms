package models

import (
	drv "database/sql/driver"
	"encoding/json"
	"fmt"
	"sync"
)

// 权限

var (
	jcRoles = struct {
		lock  sync.Mutex       `json:"-"`
		Roles map[string]*Role `json:"roles"`
	}{Roles: make(map[string]*Role)}
)

type Role struct {
	Name         string          `json:"name"`
	DisplayName  string          `json:"display_name"`
	Capabilities map[string]bool `json:"capabilities"`
}

type AccountCap struct {
	Roles   []string        `json:"roles"`
	Caps    map[string]bool `json:"caps"`
	AllCaps map[string]bool `json:"-"`
}

// 更新role到数据库中
func updateRoleOption() error {
	jcRoles.lock.Lock()
	defer jcRoles.lock.Unlock()

	buf, err := json.Marshal(jcRoles)
	if err != nil {
		return err
	}

	return updateOptions("_capabilities", string(buf))
}

// 创建一个新角色
func CreateRole(name, displayName string, caps map[string]bool) (err error) {
	jcRoles.lock.Lock()

	if _, ok := jcRoles.Roles[name]; ok {
		jcRoles.lock.Unlock()
		return fmt.Errorf("Role %s has exist.", name)
	}
	jcRoles.lock.Unlock()

	jcRoles.Roles[name] = &Role{name, displayName, caps}

	return updateRoleOption()
}

// 删除角色
func RemoveRole(name string) (err error) {
	jcRoles.lock.Lock()
	delete(jcRoles.Roles, name)
	jcRoles.lock.Unlock()

	return updateRoleOption()
}

// 返回所有角色的名字
func AllRoleNames() []string {
	rnames := make([]string, len(jcRoles.Roles))
	i := 0
	for name, _ := range jcRoles.Roles {
		rnames[i] = name
	}
	return rnames
}

func (r *Role) addCap(cap string, grant bool) {
	r.Capabilities[cap] = grant
}

func (r *Role) delCap(cap string) {
	delete(r.Capabilities, cap)
}

// role是否有cap权限
func (r *Role) hasCap(cap string) bool {
	// TODO: apply filter
	// $capabilities = apply_filters( 'role_has_cap', $this->capabilities, $cap, $this->name );
	if c, ok := r.Capabilities[cap]; ok {
		return c
	}
	return false
}

func (ac *AccountCap) Scan(src interface{}) error {
	buf, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("AccountCap Scan: cannot convert src(%v) to byte", src)
	}
	err := json.Unmarshal([]byte(buf), ac)
	return err
}

func (ac AccountCap) Value() (drv.Value, error) {
	buf, err := json.Marshal(ac)
	if err != nil {
		return nil, err
	}
	return string(buf), nil
}

// Account role interface

func (a *Account) GetCaps() error {
	if a.Capability == nil {
		a.Capability = &AccountCap{}
	}
	err := ScanMetaData(a.Id, "users", "capability", a.Capability)
	if err != nil {
		return err
	}
	return nil
}

func (a *Account) SetCaps() error {
	if a.Capability == nil {
		return fmt.Errorf("capability is nil")
	}
	val, err := a.Capability.Value()
	if err != nil {
		return err
	}

	return UpdateMetaData(a.Id, "users", "capability", val.(string), true)
}

func (a *Account) AddRole(rname string) {

}

func (a *Account) DelRole(rname string) {
}

func (a *Account) SetRole(rname string) {

}

func (a *Account) AddCap(cap string) {

}

func (a *Account) RemoveCap(cap string) {

}

func (a *Account) RemoveAllRole() {

}

func (a *Account) RemoveAllCap() {

}

func (a *Account) HasCap(cap string) bool {
	return false
}
