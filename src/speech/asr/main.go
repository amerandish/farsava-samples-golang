package main


import (
	"bytes"
	"encoding/json"
	"net/http"
	"fmt"
	"os"
	"bufio"
	"io/ioutil"
	"encoding/base64"
)


type ErrorModel struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}


// request models
type ASRRequestModel struct {
	Config Config `json:"config"`
	Audio  Audio  `json:"audio"`
}

type Audio struct {
	Data string `json:"data"`
}
// response models
type Config struct {
	AudioEncoding   string `json:"audioEncoding"`
	SampleRateHertz int64  `json:"sampleRateHertz"`
	LanguageCode    string `json:"languageCode"`
	MaxAlternatives int64  `json:"maxAlternatives"`
	ProfanityFilter bool   `json:"profanityFilter"`
	ASRModel        string `json:"asrModel"`
	LanguageModel   string `json:"languageModel"`
}


type ASRModel struct {
	TranscriptionID string   `json:"transcriptionId"`
	Duration        int64    `json:"duration"`
	InferenceTime   int64    `json:"inferenceTime"`
	Status          string   `json:"status"`
	Results         []Result `json:"results"`
}

type Result struct {
	Transcript string  `json:"transcript"`
	Confidence float64 `json:"confidence"`
	Words      []Word  `json:"words"`
}

type Word struct {
	StartTime  float64 `json:"startTime"`
	EndTime    float64 `json:"endTime"`
	Word       string  `json:"word"`
	Confidence float64 `json:"confidence"`
}

func main() {
	baseURL:= "https://api.amerandish.com/v1";
	actionURL:= "/speech/asr";
	authKey:= "<YOUR_API_KEY>";
	filePath := "<YOUR_WAV_FILE_PATH>";

	file, err := os.Open(filePath)
	if(err != nil){
		fmt.Printf("file open error: %v",err)
		return
	}

	reader := bufio.NewReader(file)
	fileData, err := ioutil.ReadAll(reader)
	if(err != nil){
		fmt.Printf("file reader error: %v",err)
		return
	}

	encoded:=base64.StdEncoding.EncodeToString(fileData)
	payload:=ASRRequestModel{
		Audio: Audio{
			Data: encoded,
		},
		Config: Config{
			AudioEncoding: "LINEAR16",
        	SampleRateHertz: 16000,
        	LanguageCode: "fa",
        	MaxAlternatives: 1,
        	ProfanityFilter: true,
        	ASRModel: "default",
        	LanguageModel: "general",
		},
	}

	url := baseURL + actionURL;
	payloadJSON, err := json.Marshal(payload)
	if(err != nil){
		fmt.Printf("payload json marshal error: %v",err)
		return
	}
	payloadBuffer := bytes.NewBuffer(payloadJSON)
	req, err := http.NewRequest("POST", url, payloadBuffer);

	if(err != nil){
		fmt.Printf("request error: %v",err)
		return
	}
	req.Header.Add("Authorization", fmt.Sprintf("bearer %v", authKey));
	req.Header.Add("Content-Type", "application/json");

	client := &http.Client{}
	res, err:= client.Do(req)

	if(err != nil){
		fmt.Printf("parse error: %v",err)
		return
	}
	defer res.Body.Close()

	if(res.StatusCode == 200){
		// success
		jsonBody := new(ASRModel)
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