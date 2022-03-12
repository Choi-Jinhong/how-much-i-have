package token

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
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
	body        Body
	stakingBody StakingBody
)

func NumberOfToken(url string, apiKey string, types string) int {
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
