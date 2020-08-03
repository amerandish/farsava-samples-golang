package main


import (
	"fmt"
	"os"
	"bufio"
	"io/ioutil"
	"encoding/base64"
	"github.com/gorilla/websocket"
	"time"
	"os/signal"
)


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
	const BUFFER_SIZE=16000
	baseURL:= "wss://api.amerandish.com/v1";
	actionURL:= "/speech/asrlive";
	authKey:= "<YOUR_API_KEY>";
	interrupt := make(chan os.Signal, 1)
    signal.Notify(interrupt, os.Interrupt)
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

	url := baseURL + actionURL + "?jwt=" + authKey ;

	connection, _ , err := websocket.DefaultDialer.Dial(url, nil)
	if(err != nil){
		fmt.Printf("socket connection error: %v",err)
		return
	}
	fmt.Println("connected")
	defer connection.Close()
	done := make(chan struct{})
	go func() {
        defer close(done)
        for {
            connection.SetReadDeadline(time.Now().Add(2 * time.Minute))
			jsonBody := new(ASRModel)
			 err := connection.ReadJSON(jsonBody)
            if err != nil {
				fmt.Println("read:", err)
				defer connection.Close()
                return
            }
            fmt.Printf("recv: %v\n", jsonBody)
        }
    }()


	ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()

	for index := 0; index < len(encoded); index+=BUFFER_SIZE {


	}
	index:=0
	for {
		select {
		case <-done:
			err := connection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				fmt.Println("write close:", err)
				return
			}
			return
		case _ = <-ticker.C:
			if(index<len(encoded)){
				if(index+BUFFER_SIZE<len(encoded)){
					fmt.Printf("send: %v->%v\n", index,index+BUFFER_SIZE)
					connection.WriteMessage(websocket.BinaryMessage, []byte(encoded[index:index+BUFFER_SIZE]))
				}else{
					fmt.Printf("send: %v->%v\n", index,len(encoded))
					connection.WriteMessage(websocket.BinaryMessage, []byte(encoded[index:]))
				}
			}
			if err != nil {
				fmt.Println("write:", err)
				return
			}
			index+=BUFFER_SIZE
		case <-interrupt:
			fmt.Println("interrupt")

			err := connection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				fmt.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}