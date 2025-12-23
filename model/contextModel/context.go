package contextModel

import (
	"jiyu/model/adminUserModel"
	"jiyu/model/userModel"
)

type Context struct {
	User *userModel.User
}

type AdminContext struct {
	User *adminUserModel.User
}
