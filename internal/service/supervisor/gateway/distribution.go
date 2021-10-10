package gateway

import (
	"bytes"
	"encoding/json"
	"github.com/foxfurry/go_kitchen/internal/domain/dto"
	"github.com/spf13/viper"
	"log"
	"net/http"
)

func Distribute(order *dto.Order) (*http.Response, error){
	resp := dto.Distribution{}
	resp.Order = *order

	jsonBody, err := json.Marshal(resp)
	if err != nil {
		log.Panic(err)
	}
	contentType := "application/json"

	return http.Post(viper.GetString("dining_host")+"/Distribute", contentType, bytes.NewReader(jsonBody))
}
