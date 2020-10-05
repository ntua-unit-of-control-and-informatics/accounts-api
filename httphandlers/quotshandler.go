package httphandlers

import (
	"log"
	"os"

	quots "github.com/euclia/goquots"
)

var QuotsClient quots.IQuots

func InitQuots() {
	quotsBase := os.Getenv("QUOTS_URL")
	appId := os.Getenv("QUOTS_APP")
	appSecret := os.Getenv("QUOTS_SECRET")
	if quotsBase == "" {
		quotsBase = "http://192.168.10.100:30180"
		// quotsBase = "http://localhost:8000"
	}
	if appId == "" {
		// appId = "GOQUOTS"
		appId = "JAQPOT-ACCOUNTS"
	}
	if appSecret == "" {
		// appSecret = "IlFELGMLf^BmJg2MVV"
		appSecret = "LqnYUFPm*sDnGAXiGl"
	}
	log.Println("Starting quots at " + quotsBase + " with quots app " + appId)

	var q quots.IQuots
	q = quots.InitQuots(quotsBase, appId, appSecret)
	QuotsClient = q
}

// var q goquots.IQuots
// q = goquots.InitQuots("http://localhost:8000", "GOQUOTS", "IlFELGMLf^BmJg2MVV")
// var IQuots = q
