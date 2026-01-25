package models

// UserRight представляет права пользователя
type UserRight struct {
	Read    bool `json:"read"`
	Write   bool `json:"write"`
	Execute bool `json:"execute"`
	Admin   bool `json:"admin"`
}

// GetFromString создает UserRight из строки (например, "rwxa")
func GetUserRightFromString(right string) UserRight {
	ur := UserRight{}
	for _, char := range right {
		switch char {
		case 'r', 'R':
			ur.Read = true
		case 'w', 'W':
			ur.Write = true
		case 'x', 'X':
			ur.Execute = true
		case 'a', 'A':
			ur.Admin = true
		}
	}
	return ur
}

// ToString преобразует UserRight в строку
func (ur UserRight) ToString() string {
	result := ""
	if ur.Read {
		result += "r"
	}
	if ur.Write {
		result += "w"
	}
	if ur.Execute {
		result += "x"
	}
	if ur.Admin {
		result += "a"
	}
	return result
}
