package auth

func CheckPassword(userinfo *Userinfo, p2 string) (*Userinfo, string) {
	if userinfo.Password == p2 {
		return userinfo, ""
	} else {
		return userinfo, "password is not correct"
	}
}

func Auth(username string, password string, userinfo *Userinfo) (*Userinfo, string) {
	if userinfo != nil {
		return CheckPassword(userinfo, password)
	} else {
		return userinfo, "user does not exist"
	}
}