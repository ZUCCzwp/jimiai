package gaode

import (
	"fmt"
	"jiyu/config"
	"jiyu/model/positionModel"
	"jiyu/util/http"
	"log"
)

type reGeoResponse struct {
	Status    string `json:"status"`
	Info      string `json:"info"`
	InfoCode  string `json:"infocode"`
	ReGeoCode struct {
		FormattedAddress string `json:"formatted_address"`
		AddressComponent struct {
			Country  string `json:"country"`
			Province string `json:"province"`
			City     string `json:"city"`
		} `json:"addressComponent"`
	} `json:"regeocode"`
}

// ReGeo 地图逆编码 返回城市名称
func ReGeo(position positionModel.Position) string {
	url := fmt.Sprintf("https://restapi.amap.com/v3/geocode/regeo?key=%s&location=%v,%v&poitype=&radius=0&extensions=all&batch=false&roadlevel=0", config.GaodeConfig.Key, position.Longitude, position.Latitude)

	var response reGeoResponse
	_, body, err := http.Get(url, "", &response)
	if err != nil {
		log.Println("util.gaode.ReGeo http.Get error:", err, body)
		return ""
	}
	r := response.ReGeoCode.AddressComponent

	if r.City != "" {
		return r.City
	} else if r.Province != "" {
		return r.Province
	} else {
		return "未知"
	}
}
