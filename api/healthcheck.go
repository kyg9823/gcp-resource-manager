package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/kyg9823/gcp-resource-manager/types"
)

func Healthcheck(w http.ResponseWriter, r *http.Request) {

	result := &types.Result{
		StatusCode: 200,
		Message:    "OK",
	}
	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Printf("error encoding response: %v", err)
		http.Error(w, "Could not marshal JSON output", 500)
		return
	}
	fmt.Fprint(w)
}
