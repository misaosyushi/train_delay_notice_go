package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/line/line-bot-sdk-go/linebot"
)

func main() {
	lambda.Start(postLineMessage)
}

func postLineMessage() {
	bot, err := linebot.New(os.Getenv("CHANNEL_SECRET"), os.Getenv("CHANNEL_TOKEN"))
	if err != nil {
		fmt.Println(err)
	}
	if _, err := bot.PushMessage(os.Getenv("USER_ID"), linebot.NewTextMessage(createTrainDelayInfo())).Do(); err != nil {
		fmt.Println(err)
	}
}

func createTrainDelayInfo() string {
	targetCompany := "JR東日本"
	targetNames := []string{"埼京線", "湘南新宿ライン", "京浜東北線", "高崎線", "宇都宮線"}
	var notifyNames []string

	for _, train := range getTrainDelayInfo() {
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
	}
	return strings.Join(notifyNames, ", ") + "が遅延しています。詳細をご確認ください。\nhttps://transit.yahoo.co.jp/traininfo/area/4/"
}

type TrainDelayInfo struct {
	Name    string `json:"name"`
	Company string `json:"company"`
}

func getTrainDelayInfo() []TrainDelayInfo {
	url := "https://tetsudo.rti-giken.jp/free/delay.json"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
	}

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	byteArray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	var trainDelayInfos []TrainDelayInfo
	if err := json.Unmarshal([]byte(byteArray), &trainDelayInfos); err != nil {
		fmt.Println("JSON Unmarshal error:", err)
	}
	return trainDelayInfos
}
