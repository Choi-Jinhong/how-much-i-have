package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

type Balance struct {
	Denom  string
	Amount string
}

type Body struct {
	Height string
	Result []Balance
}

var (
	Token    = ""
	Session  *discordgo.Session
	body     Body
	balanace []Balance
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

	if message.Content == "!Osmosis" {
		result := checkOsmosis()
		session.ChannelMessageSend(message.ChannelID, result)
	}

	if message.Content == "!pong" {
		session.ChannelMessageSend(message.ChannelID, "ping!")
	}
}

func checkOsmosis() string {
	req, err := http.NewRequest("GET", "https://osmosis-mainnet-rpc.allthatnode.com:1317/bank/balances/osmo13fla7v859d3sqrff2afx84mnc7grumtsa3hllc", nil)
	if err != nil {
		// handle err
	}
	req.Header.Add("x-allthatnode-api-key", "chjUxhz3pahoppe9DF06MLCebipgi2b7")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	result := string(bytes)

	json.Unmarshal([]byte(result), &body)
	balance, err := strconv.Atoi(body.Result[4].Amount)
	if err != nil {
		// handle err
	}

	return strconv.Itoa(balance)
}
