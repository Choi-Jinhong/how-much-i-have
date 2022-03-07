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
type CoinGekco struct {
	ID                  string              `json:"id"`
	Symbol              string              `json:"symbol"`
	Name                string              `json:"name"`
	Localization        Localization        `json:"localization"`
	Image               Image               `json:"image"`
	MarketData          MarketData          `json:"market_data"`
	CommunityData       CommunityData       `json:"community_data"`
	DeveloperData       DeveloperData       `json:"developer_data"`
	PublicInterestStats PublicInterestStats `json:"public_interest_stats"`
}
type Localization struct {
	En   string `json:"en"`
	De   string `json:"de"`
	Es   string `json:"es"`
	Fr   string `json:"fr"`
	It   string `json:"it"`
	Pl   string `json:"pl"`
	Ro   string `json:"ro"`
	Hu   string `json:"hu"`
	Nl   string `json:"nl"`
	Pt   string `json:"pt"`
	Sv   string `json:"sv"`
	Vi   string `json:"vi"`
	Tr   string `json:"tr"`
	Ru   string `json:"ru"`
	Ja   string `json:"ja"`
	Zh   string `json:"zh"`
	ZhTw string `json:"zh-tw"`
	Ko   string `json:"ko"`
	Ar   string `json:"ar"`
	Th   string `json:"th"`
	ID   string `json:"id"`
}
type Image struct {
	Thumb string `json:"thumb"`
	Small string `json:"small"`
}
type CurrentPrice struct {
	Aed  float64 `json:"aed"`
	Ars  float64 `json:"ars"`
	Aud  float64 `json:"aud"`
	Bch  float64 `json:"bch"`
	Bdt  float64 `json:"bdt"`
	Bhd  float64 `json:"bhd"`
	Bmd  float64 `json:"bmd"`
	Bnb  float64 `json:"bnb"`
	Brl  float64 `json:"brl"`
	Btc  float64 `json:"btc"`
	Cad  float64 `json:"cad"`
	Chf  float64 `json:"chf"`
	Clp  float64 `json:"clp"`
	Cny  float64 `json:"cny"`
	Czk  float64 `json:"czk"`
	Dkk  float64 `json:"dkk"`
	Dot  float64 `json:"dot"`
	Eos  float64 `json:"eos"`
	Eth  float64 `json:"eth"`
	Eur  float64 `json:"eur"`
	Gbp  float64 `json:"gbp"`
	Hkd  float64 `json:"hkd"`
	Huf  float64 `json:"huf"`
	Idr  float64 `json:"idr"`
	Ils  float64 `json:"ils"`
	Inr  float64 `json:"inr"`
	Jpy  float64 `json:"jpy"`
	Krw  float64 `json:"krw"`
	Kwd  float64 `json:"kwd"`
	Lkr  float64 `json:"lkr"`
	Ltc  float64 `json:"ltc"`
	Mmk  float64 `json:"mmk"`
	Mxn  float64 `json:"mxn"`
	Myr  float64 `json:"myr"`
	Ngn  float64 `json:"ngn"`
	Nok  float64 `json:"nok"`
	Nzd  float64 `json:"nzd"`
	Php  float64 `json:"php"`
	Pkr  float64 `json:"pkr"`
	Pln  float64 `json:"pln"`
	Rub  float64 `json:"rub"`
	Sar  float64 `json:"sar"`
	Sek  float64 `json:"sek"`
	Sgd  float64 `json:"sgd"`
	Thb  float64 `json:"thb"`
	Try  float64 `json:"try"`
	Twd  float64 `json:"twd"`
	Uah  float64 `json:"uah"`
	Usd  float64 `json:"usd"`
	Vef  float64 `json:"vef"`
	Vnd  float64 `json:"vnd"`
	Xag  float64 `json:"xag"`
	Xau  float64 `json:"xau"`
	Xdr  float64 `json:"xdr"`
	Xlm  float64 `json:"xlm"`
	Xrp  float64 `json:"xrp"`
	Yfi  float64 `json:"yfi"`
	Zar  float64 `json:"zar"`
	Bits float64 `json:"bits"`
	Link float64 `json:"link"`
	Sats float64 `json:"sats"`
}
type MarketCap struct {
	Aed  float64 `json:"aed"`
	Ars  float64 `json:"ars"`
	Aud  float64 `json:"aud"`
	Bch  float64 `json:"bch"`
	Bdt  float64 `json:"bdt"`
	Bhd  float64 `json:"bhd"`
	Bmd  float64 `json:"bmd"`
	Bnb  float64 `json:"bnb"`
	Brl  float64 `json:"brl"`
	Btc  float64 `json:"btc"`
	Cad  float64 `json:"cad"`
	Chf  float64 `json:"chf"`
	Clp  float64 `json:"clp"`
	Cny  float64 `json:"cny"`
	Czk  float64 `json:"czk"`
	Dkk  float64 `json:"dkk"`
	Dot  float64 `json:"dot"`
	Eos  float64 `json:"eos"`
	Eth  float64 `json:"eth"`
	Eur  float64 `json:"eur"`
	Gbp  float64 `json:"gbp"`
	Hkd  float64 `json:"hkd"`
	Huf  float64 `json:"huf"`
	Idr  float64 `json:"idr"`
	Ils  float64 `json:"ils"`
	Inr  float64 `json:"inr"`
	Jpy  float64 `json:"jpy"`
	Krw  float64 `json:"krw"`
	Kwd  float64 `json:"kwd"`
	Lkr  float64 `json:"lkr"`
	Ltc  float64 `json:"ltc"`
	Mmk  float64 `json:"mmk"`
	Mxn  float64 `json:"mxn"`
	Myr  float64 `json:"myr"`
	Ngn  float64 `json:"ngn"`
	Nok  float64 `json:"nok"`
	Nzd  float64 `json:"nzd"`
	Php  float64 `json:"php"`
	Pkr  float64 `json:"pkr"`
	Pln  float64 `json:"pln"`
	Rub  float64 `json:"rub"`
	Sar  float64 `json:"sar"`
	Sek  float64 `json:"sek"`
	Sgd  float64 `json:"sgd"`
	Thb  float64 `json:"thb"`
	Try  float64 `json:"try"`
	Twd  float64 `json:"twd"`
	Uah  float64 `json:"uah"`
	Usd  float64 `json:"usd"`
	Vef  float64 `json:"vef"`
	Vnd  float64 `json:"vnd"`
	Xag  float64 `json:"xag"`
	Xau  float64 `json:"xau"`
	Xdr  float64 `json:"xdr"`
	Xlm  float64 `json:"xlm"`
	Xrp  float64 `json:"xrp"`
	Yfi  float64 `json:"yfi"`
	Zar  float64 `json:"zar"`
	Bits float64 `json:"bits"`
	Link float64 `json:"link"`
	Sats float64 `json:"sats"`
}
type TotalVolume struct {
	Aed  float64 `json:"aed"`
	Ars  float64 `json:"ars"`
	Aud  float64 `json:"aud"`
	Bch  float64 `json:"bch"`
	Bdt  float64 `json:"bdt"`
	Bhd  float64 `json:"bhd"`
	Bmd  float64 `json:"bmd"`
	Bnb  float64 `json:"bnb"`
	Brl  float64 `json:"brl"`
	Btc  float64 `json:"btc"`
	Cad  float64 `json:"cad"`
	Chf  float64 `json:"chf"`
	Clp  float64 `json:"clp"`
	Cny  float64 `json:"cny"`
	Czk  float64 `json:"czk"`
	Dkk  float64 `json:"dkk"`
	Dot  float64 `json:"dot"`
	Eos  float64 `json:"eos"`
	Eth  float64 `json:"eth"`
	Eur  float64 `json:"eur"`
	Gbp  float64 `json:"gbp"`
	Hkd  float64 `json:"hkd"`
	Huf  float64 `json:"huf"`
	Idr  float64 `json:"idr"`
	Ils  float64 `json:"ils"`
	Inr  float64 `json:"inr"`
	Jpy  float64 `json:"jpy"`
	Krw  float64 `json:"krw"`
	Kwd  float64 `json:"kwd"`
	Lkr  float64 `json:"lkr"`
	Ltc  float64 `json:"ltc"`
	Mmk  float64 `json:"mmk"`
	Mxn  float64 `json:"mxn"`
	Myr  float64 `json:"myr"`
	Ngn  float64 `json:"ngn"`
	Nok  float64 `json:"nok"`
	Nzd  float64 `json:"nzd"`
	Php  float64 `json:"php"`
	Pkr  float64 `json:"pkr"`
	Pln  float64 `json:"pln"`
	Rub  float64 `json:"rub"`
	Sar  float64 `json:"sar"`
	Sek  float64 `json:"sek"`
	Sgd  float64 `json:"sgd"`
	Thb  float64 `json:"thb"`
	Try  float64 `json:"try"`
	Twd  float64 `json:"twd"`
	Uah  float64 `json:"uah"`
	Usd  float64 `json:"usd"`
	Vef  float64 `json:"vef"`
	Vnd  float64 `json:"vnd"`
	Xag  float64 `json:"xag"`
	Xau  float64 `json:"xau"`
	Xdr  float64 `json:"xdr"`
	Xlm  float64 `json:"xlm"`
	Xrp  float64 `json:"xrp"`
	Yfi  float64 `json:"yfi"`
	Zar  float64 `json:"zar"`
	Bits float64 `json:"bits"`
	Link float64 `json:"link"`
	Sats float64 `json:"sats"`
}
type MarketData struct {
	CurrentPrice CurrentPrice `json:"current_price"`
	MarketCap    MarketCap    `json:"market_cap"`
	TotalVolume  TotalVolume  `json:"total_volume"`
}
type CommunityData struct {
	FacebookLikes            interface{} `json:"facebook_likes"`
	TwitterFollowers         interface{} `json:"twitter_followers"`
	RedditAveragePosts48H    float64     `json:"reddit_average_posts_48h"`
	RedditAverageComments48H float64     `json:"reddit_average_comments_48h"`
	RedditSubscribers        interface{} `json:"reddit_subscribers"`
	RedditAccountsActive48H  interface{} `json:"reddit_accounts_active_48h"`
}
type CodeAdditionsDeletions4Weeks struct {
	Additions int `json:"additions"`
	Deletions int `json:"deletions"`
}
type DeveloperData struct {
	Forks                        int                          `json:"forks"`
	Stars                        int                          `json:"stars"`
	Subscribers                  int                          `json:"subscribers"`
	TotalIssues                  int                          `json:"total_issues"`
	ClosedIssues                 int                          `json:"closed_issues"`
	PullRequestsMerged           int                          `json:"pull_requests_merged"`
	PullRequestContributors      int                          `json:"pull_request_contributors"`
	CodeAdditionsDeletions4Weeks CodeAdditionsDeletions4Weeks `json:"code_additions_deletions_4_weeks"`
	CommitCount4Weeks            int                          `json:"commit_count_4_weeks"`
}
type PublicInterestStats struct {
	AlexaRank   int         `json:"alexa_rank"`
	BingMatches interface{} `json:"bing_matches"`
}

var (
	Session       *discordgo.Session
	body          Body
	stakingBody   StakingBody
	coinGekco     CoinGekco
	OsmosisApiKey string
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

func messageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {
	if session.State.User.ID == message.Author.ID {
		return
	}

	if message.Content == "!Osmosis" {
		balance := checkBalance()
		krw := checkCoingecko()
		fmt.Printf("BALANCE: %f \n", balance)
		fmt.Printf("KRW: %f \n", krw)
		fmt.Printf("Osmosis price: %f \n", balance*krw)
		result := strconv.FormatFloat(balance*krw, 'f', -1, 32) // return string
		session.ChannelMessageSend(message.ChannelID, result)
	}

	if message.Content == "!pong" {
		session.ChannelMessageSend(message.ChannelID, "ping!")
	}
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
	req, err := http.NewRequest("GET", "https://api.coingecko.com/api/v3/coins/osmosis/history?date=27-02-2022", nil)
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

	json.Unmarshal([]byte(resp), &coinGekco)
	krw := coinGekco.MarketData.CurrentPrice.Krw
	return krw
}
