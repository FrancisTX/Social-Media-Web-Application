package auth

import "main/server/db"

func CheckPassword(user *db.User, p2 string) (*db.User, string) {
	if user.Password == p2 {
		return user, ""
	} else {
		return user, "password is not correct"
	}
}

func Auth(username string, password string) (*db.User, string) {
	if user, err := db.QueryUser(username); err == nil {
		return CheckPassword(user, password)
	} else {
		return user, "user is not exist"
	}
}