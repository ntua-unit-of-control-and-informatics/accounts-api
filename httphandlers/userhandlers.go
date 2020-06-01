package httphandlers

import (
	"encoding/json"
	"net/http"

	// db "quots/database"
	"euclia.xyz/accounts-api/models"

	"log"
	"strconv"
	"strings"

	auth "euclia.xyz/accounts-api/authentication"
	goquots "github.com/euclia/goquots"
	"github.com/gorilla/mux"
	// jwt "gopkg.in/square/go-jose.v2/jwt"
)

// var quo = quots.InitQuots("http://localhost:8000", )

func GetUser(w http.ResponseWriter, r *http.Request) {
	min := r.FormValue("min")
	max := r.FormValue("max")
	email := r.FormValue("email")
	if min == "" && max == "" {
		apiAll := r.Header.Get("Authorization")
		apiKeyAr := strings.Split(apiAll, " ")
		token := apiKeyAr[1]
		claims, err := auth.GetClaims(token)
		userG, err := usersHandler.GetUserById(claims.Subject)
		// usersHandler.UpdateUser()
		if err != nil {
			if err.Error() == "mongo: no documents in result" {
				createUserC := make(chan models.User)
				quotsUserC := make(chan goquots.QuotsUser)
				// quotsU, err := QuotsClient.CreateUser(claims.Subject, claims.PreferedUsername, claims.Email)
				go createUser(claims, createUserC)
				go createQuotsUser(claims, quotsUserC)
				user := <-createUserC
				quotsU := <-quotsUserC
				userCreated := models.User(user)
				for _, s := range quotsU.SpentOn {
					var sp models.Spent
					sp.Appid = s.AppId
					sp.Usage = s.Usage
					userCreated.Spenton = append(userCreated.Spenton, sp)
				}
				json.NewEncoder(w).Encode(userCreated)
			} else {
				err := models.ErrorReport{
					Message: err.Error(),
					Status:  http.StatusNotFound,
				}
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(err)
			}
		}
		quotsU, err := QuotsClient.GetUser(claims.Subject)
		userG.Credits = quotsU.Credits
		for _, s := range quotsU.SpentOn {
			var sp models.Spent
			sp.Appid = s.AppId
			sp.Usage = s.Usage
			userG.Spenton = append(userG.Spenton, sp)
		}
		if userG.Name != claims.Name || userG.GivenName != claims.GivenName || userG.FamilyName != claims.FamilyName {
			var userToUp models.User
			userToUp.Name = claims.Name
			userToUp.FamilyName = claims.FamilyName
			userToUp.GivenName = claims.GivenName
			usersHandler.UpdateUserNames(userToUp)
			userGUpdatd, _ := usersHandler.GetUserById(claims.Subject)
			var ar [1]models.User
			ar[0] = userGUpdatd
			json.NewEncoder(w).Encode(ar)
		}
		var ar [1]models.User
		ar[0] = userG

		json.NewEncoder(w).Encode(ar)
	} else if email != "" {
		var user models.User
		user, err := usersHandler.FindUserByEmail(email)
		if err != nil {
			err := models.ErrorReport{
				Message: err.Error(),
				Status:  http.StatusNotFound,
			}
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(err)
		} else {
			var ar [1]models.User
			ar[0] = user
			json.NewEncoder(w).Encode(ar)
		}
	} else {
		minval, err := strconv.ParseInt(min, 10, 64)
		maxval, err := strconv.ParseInt(max, 10, 64)
		total, users, err := usersHandler.GetUsersPaginated(minval, maxval)
		if err != nil {
			err := models.ErrorReport{
				Message: err.Error(),
				Status:  http.StatusNotFound,
			}
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(err)
		} else {
			t := strconv.FormatInt(total, 10)
			w.Header().Add("total", t)
			json.NewEncoder(w).Encode(users)
		}
	}
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	_ = json.NewDecoder(r.Body).Decode(&user)
	apiAll := r.Header.Get("Authorization")
	apiKeyAr := strings.Split(apiAll, " ")
	token := apiKeyAr[1]
	claims, err := auth.GetClaims(token)
	if err != nil {
		err := models.ErrorReport{
			Message: err.Error(),
			Status:  http.StatusBadRequest,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
	} else if claims.Subject != user.Id {
		err := models.ErrorReport{
			Message: err.Error(),
			Status:  http.StatusForbidden,
		}
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(err)
	} else {
		usersHandler.UpdateUser(user)
		json.NewEncoder(w).Encode(user)
	}
}

func UpdateUserCredits(w http.ResponseWriter, r *http.Request) {
	// vars := mux.Vars(r)
	// id := vars["id"]
	var user models.User
	_ = json.NewDecoder(r.Body).Decode(&user)
	apiAll := r.Header.Get("Authorization")
	apiKeyAr := strings.Split(apiAll, " ")
	token := apiKeyAr[1]
	claims, err := auth.GetClaims(token)
	if err != nil {
		err := models.ErrorReport{
			Message: err.Error(),
			Status:  http.StatusBadRequest,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
	} else if indexOf("/Administrator", claims.Groups) == -1 {
		err := models.ErrorReport{
			Message: err.Error(),
			Status:  http.StatusForbidden,
		}
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(err)
	} else {
		var qu goquots.QuotsUser
		qu.Id = user.Id
		qu.Email = user.Email
		qu.Credits = user.Credits
		qugot, err := QuotsClient.UpdateUserCredits(qu)
		if err != nil {
			err := models.ErrorReport{
				Message: err.Error(),
				Status:  http.StatusBadRequest,
			}
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(err)
		}
		user.Credits = qugot.Credits
		usersHandler.UpdateUser(user)
		json.NewEncoder(w).Encode(user)
	}
}

func UpdateUserOrganizations(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	orgid := vars["orgid"]
	var user models.User
	_ = json.NewDecoder(r.Body).Decode(&user)
	apiAll := r.Header.Get("Authorization")
	apiKeyAr := strings.Split(apiAll, " ")
	token := apiKeyAr[1]
	claims, err := auth.GetClaims(token)
	if err != nil {
		err := models.ErrorReport{
			Message: err.Error(),
			Status:  http.StatusBadRequest,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
	}
	userEmail := claims.Email
	userFound, err := usersHandler.GetUserById(id)
	if err != nil {
		err := models.ErrorReport{
			Message: err.Error(),
			Status:  http.StatusBadRequest,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
	}
	_, err = notHandler.FindInvitation(id, userEmail, orgid)
	if err != nil {
		err := models.ErrorReport{
			Message: err.Error(),
			Status:  http.StatusBadRequest,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
	}
	if indexOf(orgid, userFound.Organizations) > -1 {
		org, err := orgHandler.GetOrgById(orgid)
		if err != nil {
			err := models.ErrorReport{
				Message: err.Error(),
				Status:  http.StatusBadRequest,
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err)
		}
		org.Users = append(org.Users, id)
		org, err = orgHandler.UpdateOrg(org)
		if err != nil {
			err := models.ErrorReport{
				Message: err.Error(),
				Status:  http.StatusBadRequest,
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err)
		}
		user, err = usersHandler.UpdateUser(user)
		json.NewEncoder(w).Encode(user)
	} else if indexOf(orgid, userFound.Organizations) == -1 {
		org, err := orgHandler.GetOrgById(orgid)
		if err != nil {
			err := models.ErrorReport{
				Message: err.Error(),
				Status:  http.StatusBadRequest,
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err)
		}
		org.Users = append(org.Users[:indexOf(orgid, userFound.Organizations)], org.Users[indexOf(orgid, userFound.Organizations)+1:]...)
		org.Users = append(org.Users, id)
		org, err = orgHandler.UpdateOrg(org)
		if err != nil {
			err := models.ErrorReport{
				Message: err.Error(),
				Status:  http.StatusBadRequest,
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err)
		}
		user, err = usersHandler.UpdateUser(user)
		json.NewEncoder(w).Encode(user)
	}

	json.NewEncoder(w).Encode(user)
}

// func UpdateUserAndOrg(w http.ResponseWriter, r *http.Request){
// 	id := r.FormValue("id")
// 	orgid := r.FormValue("orgid")

// 	json.NewEncoder(w).Encode(r.Body)
// }

func createUser(claims auth.OidcClaims, c chan models.User) {
	var user models.User
	user.Id = claims.Subject
	user.Name = claims.Name
	user.PreferedUsername = claims.PreferedUsername
	user.Email = claims.Email
	user.FamilyName = claims.FamilyName
	user.GivenName = claims.GivenName
	userC, err := usersHandler.CreateUser(user)
	if err != nil {
		log.Println(err.Error())
	}
	c <- userC
	close(c)
}

func createQuotsUser(claims auth.OidcClaims, c chan goquots.QuotsUser) {
	quo, err := QuotsClient.GetUser(claims.Subject)
	if err != nil {
		quo, err := QuotsClient.CreateUser(claims.Subject, claims.PreferedUsername, claims.Email)
		if err != nil {
			log.Println(err.Error())
		}
		c <- quo
		close(c)
	}
	c <- quo
	close(c)
}

func getQuotsUser(claims auth.OidcClaims, id string, c chan goquots.QuotsUser) {
	quotsUser, err := QuotsClient.GetUser(id)
	if err != nil {
		quo, err := QuotsClient.CreateUser(claims.Subject, claims.PreferedUsername, claims.Email)
		if err != nil {
			log.Println(err.Error())
		}
		c <- quo
		close(c)
	}
	c <- quotsUser
	close(c)
}

func indexOf(element string, data []string) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1 //not found.
}
