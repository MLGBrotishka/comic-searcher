package entity

type User struct {
	ID       int    `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"` //TODO: add salt
	Salt     []byte `json:"salt"`
	Admin    bool   `json:"admin"` //TODO: role model
}

func (u *User) IsAdmin() bool {
	return u.Admin
}
