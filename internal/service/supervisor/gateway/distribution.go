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

	return http.Post("http://" + viper.GetString("dining_host")+"/distribution", contentType, bytes.NewReader(jsonBody))
}

func DistributeV2(order *dto.Order) (*http.Response, error) {
	jsonBody, err := json.Marshal(order)
	if err != nil {
		log.Panic(err)
	}
	contentType := "application/json"

	return http.Post("http://" + viper.GetString("delivery_host")+"/distribution", contentType, bytes.NewReader(jsonBody))

}
