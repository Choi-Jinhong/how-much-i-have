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

type Delegation struct {
	Delegator_address string
	Validator_address string
	Shares            string
}

type Staking struct {
	Delegation Delegation
	Balance    Balance
}

type Body struct {
	Height string
	Result []Balance
}

type StakingBody struct {
	Height string
	Result []Staking
}

var (
	Token       = ""
	Session     *discordgo.Session
	body        Body
	stakingBody StakingBody
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
		result := checkBalance()
		session.ChannelMessageSend(message.ChannelID, result)
	}

	if message.Content == "!pong" {
		session.ChannelMessageSend(message.ChannelID, "ping!")
	}
}

func curlCosmos(url string, apiKey string, types string) int {
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("x-allthatnode-api-key", apiKey)
	if err != nil {
		// handle err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
	}

	defer res.Body.Close()

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		// handle err
	}

	resp := string(bytes)
	var tokens int

	if types == "rest" {
		json.Unmarshal([]byte(resp), &body)
		tokens, err = strconv.Atoi(body.Result[len(body.Result)-1].Amount)
	} else if types == "staking" {
		json.Unmarshal([]byte(resp), &stakingBody)
		tokens, err = strconv.Atoi(stakingBody.Result[len(stakingBody.Result)-1].Balance.Amount)
	}
	if err != nil {
		// handle err
	}
	return tokens
}

func checkBalance() string {
	// Request how many I have tokens.
	restTokens := curlCosmos("https://osmosis-mainnet-rpc.allthatnode.com:1317/bank/balances/ADDRESS", "APIKEY", "rest")

	// Request how many I staking in this chain.
	stakingTokens := curlCosmos("https://osmosis-mainnet-rpc.allthatnode.com:1317/staking/delegators/ADDRESS/delegations", "APIKEY", "staking")

	totalBalance := float64(restTokens+stakingTokens) / 1000000
	return strconv.FormatFloat(totalBalance, 'f', -1, 32)
}

//func checkCoingecko() int64 {
//
//}
