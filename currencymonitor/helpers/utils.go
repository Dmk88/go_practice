package helpers

import (
	"encoding/json"

	"github.com/Dmk88/go_practice/currencymonitor/models"
)

func ParseAPIResponse(apiresponse string) (models.APIResponse, error) {
	var resp models.APIResponse
	err := json.Unmarshal([]byte(apiresponse), &resp)
	return resp, err
}
