package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"

	"github.com/Choi-Jinhong/how-much-i-have/configuration"
	"github.com/bwmarrin/discordgo"
	"github.com/carlescere/scheduler"
	"github.com/spf13/viper"
)

// check Omsosis balances
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

// Check coingecko struct

type Coingecko struct {
	Osmosis Osmosis
}

type Osmosis struct {
	Krw float64
}

var (
	Session       *discordgo.Session
	body          Body
	stakingBody   StakingBody
	OsmosisApiKey string
	coingecko     Coingecko
)

func init() {
	setRuntimeConfig()
	var err error
	OsmosisApiKey = configuration.RuntimeConf.Api.OsmosisApiKey
	Session, err = discordgo.New("Bot " + configuration.RuntimeConf.Discord.BotToken)
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

func setRuntimeConfig() {
	viper.AddConfigPath(".")
	viper.SetConfigName("local")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&configuration.RuntimeConf)
	if err != nil {
		panic(err)
	}
}

func messageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {
	job := func() {
		fmt.Printf("Time up !")
		session.ChannelMessageSend(message.ChannelID, calculateOsmosis())
	}
	if session.State.User.ID == message.Author.ID {
		return
	}

	if message.Content == "!Osmosis" {
		fmt.Printf("Check Osmosis !!")
		session.ChannelMessageSend(message.ChannelID, calculateOsmosis())
	}

	if message.Content == "!Daily" {
		scheduler.Every().At("01:00").Run(job)
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

func checkBalance() float64 {
	//OsmosisUrl := configuration.RuntimeConf.Api.OsmosisUrl

	// Request how many I have tokens.
	restTokens := curlCosmos("https://osmosis-mainnet-rpc.allthatnode.com:1317/bank/balances/osmo13fla7v859d3sqrff2afx84mnc7grumtsa3hllc", OsmosisApiKey, "rest")

	// Request how many I staking in this chain.
	stakingTokens := curlCosmos("https://osmosis-mainnet-rpc.allthatnode.com:1317/staking/delegators/osmo13fla7v859d3sqrff2afx84mnc7grumtsa3hllc/delegations", OsmosisApiKey, "staking")

	totalBalance := float64(restTokens+stakingTokens) / 1000000
	//return strconv.FormatFloat(totalBalance, 'f', -1, 32) // return string
	return totalBalance
}

func checkCoingecko() float64 {
	req, err := http.NewRequest("GET", "https://api.coingecko.com/api/v3/simple/price?ids=osmosis&vs_currencies=krw", nil)
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

	json.Unmarshal([]byte(resp), &coingecko)
	krw := coingecko.Osmosis.Krw
	return krw
}

func calculateOsmosis() string {
	balance := checkBalance()
	krw := checkCoingecko()
	result := strconv.FormatFloat(balance*krw, 'f', -1, 32) // return string
	return result
}
