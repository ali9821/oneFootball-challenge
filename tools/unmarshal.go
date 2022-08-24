package tools

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"project/pkg/model"
)

func Unmarshal(response *http.Response, team *model.Team) {
	data, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(data, team)
}
