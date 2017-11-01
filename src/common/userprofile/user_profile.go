package userprofile

import (
	"encoding/json"
)

type ProfileGenre struct {
	Name 	string  `json:"name"`
	Weight  float64 `json:"weight"`
	Play  int     `json:"play"`
	Views	int		`json:"views"`
	Ratio  float64  `json:"ratio"`
}

type ProfileTag struct {
	Name 	string  `json:"name"`
	Weight  float64 `json:"weight"`
	Ratio   float64 `json:"ratio"`
}

type ProfileAuthor struct {
	Name 	string  `json:"name"`
	Weight  float64 `json:"weight"`
	Ratio   float64 `json:"ratio"`
}

type UserProfile struct {
	Uid 		string `json:"uid"`
	Utdid 		string `json:"utdid"`
	Imei		string `json:"imei"`
	//Age 		int	   `json:"age"`
	//Gender 		int    `json:"gender"`
	Country  	string `json:"country"`
	Location 	string `json:"loc"`
	//CreateTs    int64  `json:"create_ts"`
	//UpdateTs 	int64  `json:"update_ts"`
	Active30    int    `json:"active_30"`
	Active7     int    `json:"active_7"`
	//RealActive  int    `json:"real_active"`
	Appver      string `json:"app_ver"`
	//Appid       string `json:"appid"`
	//Views7      int    `json:"views_7"`
	//Play7       int    `json:"play_7"`
	//AvgCtr7     float64 `json:"avg_ctr_7"`
	AcitveBitmap int64  `json:"acitve_bitmap"`
	//Views30     int    `json:"views_30"`
	//Play30	    int    `json:"play_30"`
	//AvgCtr30    float64 `json:"avg_ctr_30"`
	//AdultRatio  float64 `json:"adult_ratio"`
	//FoucsRatio  float64 `json:"foucs_ratio"`
	Genre       []ProfileGenre `json:"genre"`
	Tags        []ProfileTag `json:"tag"`
	Authors        []ProfileAuthor `json:"author"`
}

func UserProfileLoads(strJson string, userProfile *UserProfile) error {
	err := json.Unmarshal([]byte(strJson), userProfile)
	return err
}

func UserProfileDumps(profile UserProfile) (string, error) {
	data, err := json.Marshal(profile)
	return string(data), err
}