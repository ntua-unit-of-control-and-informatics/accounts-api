package httphandlers

import (
	"os"

	db "euclia.xyz/accounts-api/database"
)

var datastore db.IDatastore = &db.Datastore{}
var notHandler db.INotificationsDao = &db.NotsDao{}
var orgHandler db.IOrganizationsDao = &db.OrgsDao{}
var usersHandler db.IUsersDao = &db.UsersDao{}
var invitHandler db.IInvitationsDao = &db.InvitationsDao{}
var reqHandler db.IRequestDao = &db.RequestsDao{}

var EmailServer string
var EmailUsername string
var EmailPassword string

func Init() {
	EmailServer = os.Getenv("EMAIL_SERVER")
	if EmailServer == "" {
		EmailServer = "smtp.gmail.com"
	}
	EmailUsername = os.Getenv("EMAIL_USERNAME")
	if EmailUsername == "" {
		EmailUsername = "contact@jaqpot.org"
	}
	EmailPassword = os.Getenv(" EmailUsername")
	if EmailPassword == "" {
		EmailPassword = "8HOu8muFJc"
	}

}
