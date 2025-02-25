package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"flag"

	"github.com/dgraph-io/badger/v2"
	"github.com/gojp/goreportcard/check"
)

var domain = flag.String("domain", "goreportcard.com", "Domain used for your goreportcard installation")
var googleAnalyticsKey = flag.String("google_analytics_key", "UA-58936835-1", "Google Analytics Account Id")

// ReportHandler handles the report page
func (gh *GRCHandler) ReportHandler(w http.ResponseWriter, r *http.Request, db *badger.DB, repo string) {
	branch := check.GetBranchNameFromQuery(repo, r.URL.Query().Get("branch"))

	log.Printf("Displaying report: %q branch: %q", repo, branch)

	t, err := gh.loadTemplate("/templates/report.html")
	if err != nil {
		log.Println("ERROR: could not get report template: ", err)
		http.Error(w, err.Error(), 500)
		return
	}

	resp, err := getFromCache(db, repo, branch)
	needToLoad := false
	if err != nil {
		switch err.(type) {
		case notFoundError:
			// don't bother logging - we already log in getFromCache. continue
		default:
			log.Println("ERROR ReportHandler:", err) // log error, but continue
		}
		needToLoad = true
		resp.Branch = branch
	}

	respBytes, err := json.Marshal(resp.CalculateFileURLForFileSummaries())
	if err != nil {
		log.Println("ERROR ReportHandler: could not marshal JSON: ", err)
		http.Error(w, "Failed to load cache object", 500)
		return
	}

	t.Execute(w, map[string]interface{}{
		"repo":                 repo,
		"branch":               branch,
		"response":             string(respBytes),
		"loading":              needToLoad,
		"domain":               domain,
		"google_analytics_key": googleAnalyticsKey,
	})
}
