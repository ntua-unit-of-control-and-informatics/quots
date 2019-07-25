package database

import (
	models "quots/models"
)

type IApplicationDao interface {
	CreateApplication(applicationToCreate models.Application) (applicationCreated models.Application, err error)
	GetApplicationBiId(id string) (appFound models.Application, err error)
	UpdateApp(applicationToUpdate models.Application) (applicationUpdated models.Application, err error)
	GetAllApps(min int64, max int64) (counted int64, applications []*models.Application, err error)
}

type IUsersDao interface {
	CreateUser(userGot models.User) (userCreated models.User, err error)
	GetUserById(id string) (userFound models.User, err error)
	UpdateUserCredits(userToUpdate models.User) (userUpdated models.User, err error)
	UpdateUsersSpentOn(userToUpdate models.User) (userUpdated models.User, err error)
	GetUsersPaginated(min int64, max int64) (counted int64, usersFound []*models.User, err error)
	FindUserByEmail(email string) (userFound models.User, err error)
}
