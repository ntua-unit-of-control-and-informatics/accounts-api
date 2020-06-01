package models

type Meta struct {
	Description string   `json:"description,omitempty" bson:"description"`
	Picture     string   `json:"picture,omitempty" bson:"picture"`
	Date        int64    `json:"date,omitempty" bson:"date"`
	Subjects    []string `json:"subjects,omitempty" bson:"subjects"`
	Audiences   []string `json:"audiences,omitempty" bson:"audiences"`
	About       string   `json:"about,omitempty" bson:"about"`
	City        string   `json:"city,omitempty" bson:"city"`
	Country     string   `json:"country,omitempty" bson:"country"`
	Github      string   `json:"github,omitempty" bson:"github"`
	Website     string   `json:"website,omitempty" bson:"website"`
	Linkedin    string   `json:"linkedin,omitempty" bson:"linkedin"`
}

type User struct {
	Id               string   `json:"_id" bson:"_id"`
	Occupation       string   `json:"occupation,omitempty" bson:"occupation"`
	OccupationAt     string   `json:"occupationAt,omitempty" bson:"occupation_at"`
	MetaInfo         Meta     `json:"meta,omitempty" bson:"meta"`
	EmailVerified    bool     `json:"emailVerified,omitempty" bson:"email_verified"`
	Name             string   `json:"name" bson:"name"`
	PreferedUsername string   `json:"preferredUsername" bson:"preferred_username"`
	GivenName        string   `json:"givenName" bson:"given_name"`
	FamilyName       string   `json:"familyName" bson:"family_name"`
	Email            string   `json:"email" bson:"email"`
	Credits          float64  `json:"credits" bson:"credits"`
	Spenton          []Spent  `json:"spenton,omitempty" bson:"spenton"`
	Organizations    []string `json:"organizations,omitempty" bson:"organizations"`
}

type Organization struct {
	Id           string            `json:"_id" bson:"_id"`
	Title        string            `json:"title" bson:"title"`
	MetaInfo     Meta              `json:"meta,omitempty" bson:"meta"`
	Creator      string            `json:"creator" bson:"creator"`
	Admins       []string          `json:"admins,omitempty" bson:"admins"`
	Users        []string          `json:"users,omitempty" bson:"users"`
	Markdown     string            `json:"markdown,omitempty" bson:"markdown"`
	Contact      map[string]string `json:"contact,omitempty" bson:"contact"`
	Credits      float64           `json:"credits,omitempty" bson:"credits"`
	ContactTypes []string          `json:"contacttypes,omitempty" bson:"contact_types"`
}

type ErrorReport struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

type QuotsUser struct {
	Id       string  `json:"id" bson:"_id"`
	Email    string  `json:"email" bson:"email"`
	Username string  `json:"username" bson:"username"`
	Credits  float64 `json:"credits" bson:"credits"`
	Spenton  []Spent `json:"spenton" bson:"spenton"`
}

type CanProceed struct {
	Userid  string `json:"userid" bson:"userid"`
	Proceed bool   `json:"proceed"`
}

type Spent struct {
	Appid string                 `json:"appid" bson:"appid"`
	Usage map[string]interface{} `json:"usage"  bson:"usage"`
}

type Usage struct {
	Usagetype string  `json:"usagetype" bson:"usagetype"`
	Used      float64 `json:"used" bson:"used"`
}

type Notification struct {
	Id                string `json:"_id" bson:"_id"`
	From              string `json:"from" bson:"from"`
	To                string `json:"to" bson:"to"`
	EmailTo           string `json:"email_to" bson:"email_to"`
	Body              string `json:"body" bson:"body"`
	Type              string `json:"type" bson:"type"`
	Viewed            string `json:"viewed" bson:"viewed"`
	InvitationTo      string `json:"invitationTo" bson:"invitation_to"`
	AffiliationFromTo string `json:"affiliationFromTo" bson:"affiliation_from_to"`
	CreditsFrom       string `json:"creditsFrom" bson:"credits_from"`
}

type Invitation struct {
	Id                string `json:"_id" bson:"_id"`
	MetaInfo          Meta   `json:"meta,omitempty" bson:"meta"`
	From              string `json:"from" bson:"from"`
	FromId            string `json:"fromId" bson:"from_id"`
	To                string `json:"to" bson:"to"`
	EmailTo           string `json:"emailTo" bson:"email_to"`
	Body              string `json:"body" bson:"body"`
	Type              string `json:"type" bson:"type"`
	Viewed            bool   `json:"viewed" bson:"viewed"`
	InvitationTo      string `json:"invitationTo" bson:"invitation_to"`
	InvitationToId    string `json:"invitationToId" bson:"invitation_to_id"`
	AffiliationFromTo string `json:"affiliationFromTo" bson:"affiliation_from_to"`
	Action            string `json:"action" bson:"action"`
	// CreditsFrom       string `json:"creditsFrom" bson:"credits_from"`
}

type Request struct {
	Id        string `json:"_id" bson:"_id"`
	FromUser  string `json:"fromUser" bson:"from_user"`
	UserEmail string `json:"userEmail" bson:"user_email"`
	Message   string `json:"message" bson:"message"`
}
