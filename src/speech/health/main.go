package main

import (
	"encoding/json"
	"net/http"
	"fmt"
)


type ErrorModel struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}


type HealthModel struct{
	Status  string `json:"status"`
	Message string `json:"message"`
	Version string `json:"version"`
}

func main() {
	baseURL:= "https://api.amerandish.com/v1";
	actionURL:= "/speech/healthcheck";
	authKey:= "<YOUR_API_KEY>";

	url := baseURL + actionURL;

	req, err := http.NewRequest("GET", url, nil);
	if(err != nil){
		fmt.Printf("request error: %v",err)
		return
	}
	req.Header.Add("Authorization", fmt.Sprintf("bearer %v", authKey));

	client := &http.Client{}
	res, err:= client.Do(req)

	if(err != nil){
		fmt.Printf("parse error: %v",err)
		return
	}
	defer res.Body.Close()

	if(res.StatusCode == 200){
		// success
		jsonBody := new(HealthModel)
		fmt.Println(res.StatusCode)
		json.NewDecoder(res.Body).Decode(jsonBody)
		fmt.Println(jsonBody)
	}else{
		// error
		jsonBody := new(ErrorModel)
		fmt.Println(res.StatusCode)
		json.NewDecoder(res.Body).Decode(jsonBody)
		fmt.Println(res.Body)
		fmt.Println(jsonBody.Message)
	}
}