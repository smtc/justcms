package auth

import (
	drv "database/sql/driver"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/smtc/justcms/meta"
)

// 权限

var (
	_authRoles = struct {
		sync.Mutex `json:"-"`
		Roles      map[string]*Role `json:"roles"`
	}{Roles: make(map[string]*Role)}
)

type Role struct {
	Name         string          `json:"name"`
	DisplayName  string          `json:"display_name"`
	Capabilities map[string]bool `json:"capabilities"`
}

// 用户的权限分为两个部分：一部分是由用户的role继承的权限，一部分是单独赋给用户的权限Caps
// AllCaps是两种权限的结合, Caps中的权限可以覆盖role中的权限
type UserCap struct {
	Roles    []string        `json:"roles"`
	Caps     map[string]bool `json:"caps"`
	AllCaps  map[string]bool `json:"-"`
	computed bool            `json:"-"`
}

// 更新role到数据库中
func updateRoleOption() error {
	buf, err := json.Marshal(_authRoles)
	if err != nil {
		return err
	}

	return meta.UpdateOptions("_capabilities", string(buf))
}

// 创建一个新角色
func CreateRole(name, displayName string, caps map[string]bool) (err error) {
	_authRoles.Lock()
	defer _authRoles.Unlock()

	if _, ok := _authRoles.Roles[name]; ok {
		return fmt.Errorf("Role %s has exist.", name)
	}

	_authRoles.Roles[name] = &Role{name, displayName, caps}

	return updateRoleOption()
}

// 删除角色
func RemoveRole(name string) (err error) {
	_authRoles.Lock()
	delete(_authRoles.Roles, name)
	defer _authRoles.Unlock()

	return updateRoleOption()
}

// 返回所有角色的名字
func AllRoleNames() []string {
	rnames := make([]string, len(_authRoles.Roles))
	i := 0
	for name, _ := range _authRoles.Roles {
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

func (ac *UserCap) Scan(src interface{}) error {
	buf, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("UserCap Scan: cannot convert src(%v) to byte", src)
	}
	err := json.Unmarshal([]byte(buf), ac)
	return err
}

func (ac UserCap) Value() (drv.Value, error) {
	buf, err := json.Marshal(ac)
	if err != nil {
		return nil, err
	}
	return string(buf), nil
}

// User role interface

func (a *User) GetCaps() error {
	if a.Capability == nil {
		a.Capability = &UserCap{}
	}
	err := meta.ScanMetaData(a.ObjectId, "capability", a.Capability)
	if err != nil {
		return err
	}
	return nil
}

func (a *User) SetCaps() error {
	if a.Capability == nil {
		return fmt.Errorf("capability is nil")
	}
	val, err := a.Capability.Value()
	if err != nil {
		return err
	}

	return meta.UpdateMetaData(a.ObjectId, "users", "capability", val.(string), true)
}

func mergeRoleCaps(rname string, caps map[string]bool) {
	_authRoles.Lock()
	defer _authRoles.Unlock()

	role, ok := _authRoles.Roles[rname]
	if !ok {
		log.Printf("mergeRoleCaps: Not found role by name " + rname)
		return
	}
	for name, grant := range role.Capabilities {
		caps[name] = grant
	}
}

// 计算account role中的权限，加入到AllCap中
func (a *User) ComputeCaps() error {
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

func (a *User) getRoleCaps(rname string) {
	mergeRoleCaps(rname, a.Capability.AllCaps)
}

func (a *User) AddRole(rname string) error {
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

func (a *User) DelRole(rname string) error {
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

func (a *User) SetRole(rname string) error {
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

func (a *User) AddCap(cap string, grant bool) error {
	if a.Capability == nil {
		if err := a.GetCaps(); err != nil {
			return err
		}
	}
	a.Capability.Caps[cap] = grant
	a.Capability.AllCaps[cap] = grant

	return a.SetCaps()

}

func (a *User) RemoveCap(cap string) error {
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

func (a *User) RemoveAllRole() error {
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

func (a *User) RemoveAllCap() error {
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

func (a *User) HasCap(cap string) bool {
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
