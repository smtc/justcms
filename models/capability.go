package models

// 权限

var (
	js_roles = make(map[string]*Role)
)

type Role struct {
	Name         string          `json:"name"`
	DisplayName  string          `json:"display_name"`
	Capabilities map[string]bool `json:"capabilities"`
}

func (r *Role) AddCap(cap string, grant bool) {
	r.Capabilities[cap] = grant
}

func (r *Role) RemoveCap(cap string) {
	delete(r.Capabilities, cap)
}

func (r *Role) HasCap(cap string) bool {
	// TODO: apply filter
	// $capabilities = apply_filters( 'role_has_cap', $this->capabilities, $cap, $this->name );
	if c, ok := r.Capabilities[cap]; ok {
		return c
	}
	return false
}
