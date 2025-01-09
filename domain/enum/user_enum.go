package enum

type UserRole string

const (
	UserRole_User  = UserRole("user")
	UserRole_Admin = UserRole("admin")
)

var ValidRoles = []UserRole{UserRole_User, UserRole_Admin}

func (role UserRole) String() string {
	return string(role)
}

func (r UserRole) IsValid() bool {
	for _, validRole := range ValidRoles {
		if r == validRole {
			return true
		}
	}
	return false
}
