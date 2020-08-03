
# Farsava - ASR live Api (WebSocket)


## install dependencies

```bash
go get github.com/gorilla/websocket
```

## configs
```go
baseURL:= "wss://api.amerandish.com/v1";
actionURL:= "/speech/asrlive";
authKey:= "<YOUR_API_KEY>";
filePath := "<YOUR_WAV_FILE_PATH>"
```

## run

```bash
go run main.go
```

