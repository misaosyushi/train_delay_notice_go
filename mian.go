package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/line/line-bot-sdk-go/linebot"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func main() {
	lambda.Start(postLineMessage)
}

func postLineMessage()  {
	bot, err := linebot.New(os.Getenv("CHANNEL_SECRET"), os.Getenv("CHANNEL_TOKEN"))
	if err != nil {
		fmt.Println(err)
	}
	if _, err := bot.PushMessage(os.Getenv("USER_ID"), linebot.NewTextMessage(createTrainDelayInfo())).Do(); err != nil {
		fmt.Println(err)
	}
}

func createTrainDelayInfo() string {
	url := "https://tetsudo.rti-giken.jp/free/delay.json"

	req, _ := http.NewRequest("GET", url, nil)

	client := new(http.Client)
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)

	type TrainDelayInfo struct {
		Name string `json:"name"`
		Company string `json:"company"`
	}

	var trainDelayInfos []TrainDelayInfo
	if err := json.Unmarshal([]byte(byteArray), &trainDelayInfos); err != nil {
		fmt.Println("JSON Unmarshal error:", err)
		return ""
	}

	targetCompany := "JR東日本"
	targetNames := []string{"埼京線", "湘南新宿ライン", "京浜東北線", "高崎線", "宇都宮線"}
	var notifyNames []string

	for _, train := range trainDelayInfos {
		company := train.Company
		name := train.Name

		for _, targetName := range targetNames {
			if company == targetCompany && name == targetName {
				notifyNames = append(notifyNames, name)
			}
		}
	}

	if len(notifyNames) == 0 {
		return "遅延情報はありませんでした！\n良い一日を〜♪"
	} else {
		return strings.Join(notifyNames, ", ") + "が遅延しています。詳細をご確認ください。\nhttps://transit.yahoo.co.jp/traininfo/area/4/"
	}
}
