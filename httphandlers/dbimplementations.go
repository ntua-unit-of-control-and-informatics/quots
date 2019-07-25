package httphandlers

import (
	db "quots/database"
)

var IAppdb db.IApplicationDao
var dbAppImpl = &db.ApplicationDao{}

var IUsersdb db.IUsersDao
var dbUserImpl = &db.UsersDao{}
