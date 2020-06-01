package database

import (
	models "euclia.xyz/accounts-api/models"
)

type IUsersDao interface {
	CreateUser(userGot models.User) (userCreated models.User, err error)
	GetUserById(id string) (userFound models.User, err error)
	DeleteUser(id string) (usersDeleted int64, err error)
	UpdateUser(userToUpdate models.User) (userUpdated models.User, err error)
	// UpdateUsersSpentOn(userToUpdate models.User) (userUpdated models.User, err error)
	// UpdateUserEmailAndUsername(userToUpdate models.User) (userUpdated models.User, err error)
	GetUsersPaginated(min int64, max int64) (counted int64, usersFound []*models.User, err error)
	FindUserByEmail(email string) (userFound models.User, err error)
	SearchUserByNameAndMail(name string) (counted int, orgsFound []*models.User, err error)
	UpdateUserNames(user models.User) (userUpdated models.User, err error)
	FindOrgsUsers(org string) (counted int, orgsFound []*models.User, err error)
	UpdateUserOrganizations(user models.User) (userUpdated models.User, err error)
}

type IOrganizationsDao interface {
	CreateOrganization(orgGot models.Organization) (orgCreated models.Organization, err error)
	GetOrgById(id string) (orgGot models.Organization, err error)
	DeleteOrg(id string) (orgsDeleted int64, err error)
	UpdateOrg(orgToUpdate models.Organization) (orgUpdated models.Organization, err error)
	SearchOrgByName(name string) (counted int, orgsFound []*models.Organization, err error)
	GetOrgByAdmin(id string) (counted int, orgsFound []*models.Organization, err error)
}

type INotificationsDao interface {
	CreateNotification(not models.Notification) (notCreated models.Notification, err error)
	GetNotByToIDAndMail(id string, email string) (notFound models.Notification, err error)
	UpdateNotification(not models.Notification) (updatedNot models.Notification, err error)
	DeleteNotification(id string) (deletedNot int64, err error)
	FindInvitation(id string, email string, invitationto string) (notFound models.Notification, err error)
}

type IInvitationsDao interface {
	CreateInvitation(invGot models.Invitation) (invCreated models.Invitation, err error)
	GetByEmail(email string, start int64, skip int64, viewed bool) (invFound []models.Invitation, err error)
	GetReceivedById(id string) (invFound models.Invitation, err error)
	GetSentById(id string) (invFound models.Invitation, err error)
	UpdateInvitation(inv models.Invitation) (invUpdated models.Invitation, err error)
	DeleteInvitation(id string) (deletedNot int64, err error)
	CountWithEmail(email string, viewed bool) (total int64, err error)
}

type IRequestDao interface {
	CreateRequest(reqGot models.Request) (reqCreated models.Request, err error)
	GetRequest(userId string) (reqCreated models.Request, err error)
}
