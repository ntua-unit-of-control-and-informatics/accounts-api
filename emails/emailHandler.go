package email

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strconv"
)

var _emailUsername string
var _emailServer string
var _emailPass string
var _emailPort int
var _emailServerAll string
var _smtpServer smtpServer
var _accountsDeployment string

func Init() {
	log.Println("Starting Email")

	_accountsDeployment = os.Getenv("ACCOUNTS_DEPLOYMENT")
	if _accountsDeployment == "" {
		_accountsDeployment = "jaqpot ntua server"
	}

	_emailUsername = os.Getenv("EMAIL")
	if _emailUsername == "" {
		_emailUsername = "jaqpot@jaqpot.org"
	}
	_emailPass = os.Getenv("EMAIL_PASS")
	if _emailPass == "" {
		_emailPass = "8HOu8muFJc"
	}

	_emailServer = os.Getenv("EMAIL_SERVER")
	if _emailServer == "" {
		_emailServer = "smtp.gmail.com"
	}
	_emailPortS := os.Getenv("EMAIL_PORT")
	if _emailPortS == "" {
		_emailPortS = "587"
	}
	_emailPortInt, err := strconv.Atoi(_emailPortS)
	_emailPort = _emailPortInt
	_emailServerAll = _emailPass + ":" + _emailPortS
	_smtpServer = smtpServer{host: _emailServer, port: _emailPortS}

	if err != nil {
		log.Fatal(err.Error())
	}
	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()
}

type smtpServer struct {
	host string
	port string
}

// Address URI to smtp server
func (s *smtpServer) Address() string {
	return s.host + ":" + s.port
}

func SendInvitation(userGivenName string, toEmail string, fromBody string, orgTitle string) {
	to := []string{
		toEmail,
	}
	message := []byte(
		"Subject:" + userGivenName + " invited you to join his / her organization\r\n" +
			"\r\n" +
			userGivenName +
			" sent you an invitation to join " + orgTitle + ". \n \n" +
			fromBody +
			"\n\n The invitation will be available on https://accounts.jaqpot.org \r\n \n \n" +
			"Thank you for using Jaqpot services!")
	auth := smtp.PlainAuth("Jaqpot", _emailUsername, _emailPass, _smtpServer.host)
	// Sending email.
	err := smtp.SendMail(_smtpServer.Address(), auth, _emailUsername, to, message)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func SendCreditRequest(fromId string, fromEmail string, messageG string) {
	// from := "jaqpot@jaqpot.org"
	// password := "8HOu8muFJc"
	// Receiver email address.
	to := []string{
		"jaqpot@jaqpot.org",
		"pantelispanka@gmail.com",
	}
	// smtp server configuration.
	// Message.
	message := []byte(
		"Subject: Credit Request!\r\n" +
			"\r\n" +
			"User " + fromEmail + " with id " + fromId +
			" requested for credits for the deployemnt at " + _accountsDeployment + "\r\n" +
			"Message: " + messageG + "")
	// Authentication.
	auth := smtp.PlainAuth("Jaqpot", _emailUsername, _emailPass, _smtpServer.host)
	// Sending email.
	err := smtp.SendMail(_smtpServer.Address(), auth, _emailUsername, to, message)
	if err != nil {
		fmt.Println(err)
		return
	}

	to2 := []string{
		fromEmail,
	}
	message2 := []byte(
		"Subject: Credit Request\r\n" +
			"\r\n" +
			"Your request for more credits have been received. It will be validated and we will come back at you! \n \n \n" +
			"Thank you for using our services! \r\n")
	// Authentication.
	// Sending email.
	err2 := smtp.SendMail(_smtpServer.Address(), auth, _emailUsername, to2, message2)
	if err2 != nil {
		fmt.Println(err)
		return
	}

}
