package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"reflect"

	"github.com/bwmarrin/discordgo"
	"github.com/go-resty/resty/v2"
)

var (
	Token   = ""
	Session *discordgo.Session
)

func init() {
	var err error
	Session, err = discordgo.New("Bot " + Token)
	if err != nil {
		log.Fatalf("클라이언트 생성 오류: %v", err)
	}

	err = Session.Open()
	if err != nil {
		log.Fatalf("세션 오픈 오류: %v", err)
	}

	log.Printf("%s (%s)에 로그인 됨", Session.State.User.String(), Session.State.User.ID)
}

func main() {
	Session.AddHandler(messageCreate)

	defer Session.Close()

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("봇 종료됨")
}

func messageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {
	if session.State.User.ID == message.Author.ID {
		return
	}

	if message.Content == "!HowMuchIHave" {
		resp := checkBalance()

		fmt.Printf("\nResponse Body: %v", resp)
		session.ChannelMessageSend(message.ChannelID, resp)
	}

	if message.Content == "!pong" {
		session.ChannelMessageSend(message.ChannelID, "ping!")
	}
}

func checkBalance() string {
	client := resty.New()
	resp, err := client.R().
		SetHeader("x-allthatnode-api-key", "chjUxhz3pahoppe9DF06MLCebipgi2b7").
		Get("https://osmosis-mainnet-rpc.allthatnode.com:1317/bank/balances/osmo13fla7v859d3sqrff2afx84mnc7grumtsa3hllc")
	fmt.Printf("\nERROR: %v", err)

	body := resp.String()
	json.Unmarshal([]byte(str), &body)
	fmt.Printf("%v", reflect.TypeOf(body))
	return body
}
