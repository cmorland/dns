package whois

import (
	"encoding/json"
	db "github.com/jawr/dns/database/models/whois"
	"github.com/jawr/dns/rest/util"
	"net/http"
)

type QueryPost struct {
	Email string `json:"email"`
}

func PostQuery(w http.ResponseWriter, r *http.Request) {
	query := QueryPost{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&query)
	if err != nil {
		util.ToJSON(query, err, w)
		return
	}

	list, err := db.GetByHasEmail(query.Email).GetAll()

	util.ToJSON(list, err, w)
}
