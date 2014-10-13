package models

import (
	drv "database/sql/driver"
	"encoding/json"
	"fmt"
	"log"
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

// 用户的权限分为两个部分：一部分是由用户的role继承的权限，一部分是单独赋给用户的权限Caps
// AllCaps是两种权限的结合, Caps中的权限可以覆盖role中的权限
type AccountCap struct {
	Roles    []string        `json:"roles"`
	Caps     map[string]bool `json:"caps"`
	AllCaps  map[string]bool `json:"-"`
	computed bool            `json:"-"`
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
	err := ScanMetaData(a.ObjectId, "capability", a.Capability)
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

	return UpdateMetaData(a.ObjectId, "users", "capability", val.(string), true)
}

func mergeRoleCaps(rname string, caps map[string]bool) {
	jcRoles.lock.Lock()
	defer jcRoles.lock.Unlock()

	role, ok := jcRoles.Roles[rname]
	if !ok {
		log.Printf("mergeRoleCaps: Not found role by name " + rname)
		return
	}
	for name, grant := range role.Capabilities {
		caps[name] = grant
	}
}

// 计算account role中的权限，加入到AllCap中
func (a *Account) ComputeCaps() error {
	if a.Capability == nil {
		if err := a.GetCaps(); err != nil {
			return err
		}
	}

	for _, role := range a.Capability.Roles {
		a.getRoleCaps(role)
	}

	for name, grant := range a.Capability.Caps {
		a.Capability.AllCaps[name] = grant
	}
	a.Capability.computed = true

	return nil
}

func (a *Account) getRoleCaps(rname string) {
	mergeRoleCaps(rname, a.Capability.AllCaps)
}

func (a *Account) AddRole(rname string) error {
	if a.Capability == nil {
		if err := a.GetCaps(); err != nil {
			return err
		}
	}

	for _, name := range a.Capability.Roles {
		if name == rname {
			return nil
		}
	}
	a.Capability.Roles = append(a.Capability.Roles, rname)
	a.getRoleCaps(rname)

	return a.SetCaps()
}

func (a *Account) DelRole(rname string) error {
	if a.Capability == nil {
		if err := a.GetCaps(); err != nil {
			return err
		}
	}
	for i, name := range a.Capability.Roles {
		if name == rname {
			a.Capability.Roles = append(a.Capability.Roles[:i], a.Capability.Roles[i+1:]...)
			return a.SetCaps()
		}
	}
	// 重新计算AllCaps
	a.Capability.AllCaps = map[string]bool{}
	a.ComputeCaps()
	return nil
}

func (a *Account) SetRole(rname string) error {
	if a.Capability == nil {
		if err := a.GetCaps(); err != nil {
			return err
		}
	}

	a.Capability.Roles = []string{rname}
	// 重新计算AllCaps
	a.Capability.AllCaps = map[string]bool{}
	a.ComputeCaps()

	return a.SetCaps()

}

func (a *Account) AddCap(cap string, grant bool) error {
	if a.Capability == nil {
		if err := a.GetCaps(); err != nil {
			return err
		}
	}
	a.Capability.Caps[cap] = grant
	a.Capability.AllCaps[cap] = grant

	return a.SetCaps()

}

func (a *Account) RemoveCap(cap string) error {
	if a.Capability == nil {
		if err := a.GetCaps(); err != nil {
			return err
		}
	}
	delete(a.Capability.Caps, cap)
	// 重新计算AllCaps
	a.Capability.AllCaps = map[string]bool{}
	a.ComputeCaps()

	return a.SetCaps()

}

func (a *Account) RemoveAllRole() error {
	if a.Capability == nil {
		if err := a.GetCaps(); err != nil {
			return err
		}
	}

	a.Capability.Roles = []string{}
	// 重新计算AllCaps
	a.Capability.AllCaps = map[string]bool{}
	for name, grant := range a.Capability.Caps {
		a.Capability.AllCaps[name] = grant
	}

	return a.SetCaps()
}

func (a *Account) RemoveAllCap() error {
	if a.Capability == nil {
		if err := a.GetCaps(); err != nil {
			return err
		}
	}

	a.Capability.Caps = map[string]bool{}
	// 重新计算AllCaps
	a.Capability.AllCaps = map[string]bool{}
	a.ComputeCaps()

	return a.SetCaps()

}

func (a *Account) HasCap(cap string) bool {
	if a.Capability == nil {
		if err := a.GetCaps(); err != nil {
			return false
		}
	}
	if a.Capability.computed == false {
		a.ComputeCaps()
	}

	return a.Capability.AllCaps[cap]
}
