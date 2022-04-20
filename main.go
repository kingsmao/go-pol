package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/shopspring/decimal"
	"github.com/tyler-smith/go-bip39"
	"io/ioutil"
	"math/big"
	"net/http"
	"time"
)

type Result struct {
	Balance string `json:"result"`
}

var weiScale = big.NewInt(100000000000000000)

const MainEth = "m/44'/60'/0'/0/0'"
const InfuraProjectId = "df5d481c6d6d40759f2d6172f747b47b"
const GetBalanceUrl = "https://mainnet.infura.io/v3/" + InfuraProjectId

func main() {
	for {
		luck()
	}
}

func random12Mnemonic() string {
	entropy, _ := bip39.NewEntropy(128)
	mnemonic, _ := bip39.NewMnemonic(entropy)
	//fmt.Printf("生成助记词: %s\n", mnemonic)
	return mnemonic
}

func mnemonicToAddress(mnemonic string) string {
	path := hdwallet.MustParseDerivationPath(MainEth)
	wallet, _ := hdwallet.NewFromMnemonic(mnemonic)
	account, _ := wallet.Derive(path, false)
	address := account.Address.Hex()
	//fmt.Printf("地址: %s\n", address)
	return address
}

func hexToBigInt(hex string) *big.Int {
	n := new(big.Int)
	n, _ = n.SetString(hex[2:], 16)
	return n
}

func getBalance(address string) string {
	client := &http.Client{}
	data := make(map[string]interface{})
	data["jsonrpc"] = "2.0"
	data["method"] = "eth_getBalance"
	data["id"] = 1
	var params = [2]string{address, "latest"}
	data["params"] = params
	bytesData, _ := json.Marshal(data)
	req, _ := http.NewRequest(http.MethodPost, GetBalanceUrl, bytes.NewReader(bytesData))
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	s := string(body)
	//fmt.Println(s)
	var result Result
	json.Unmarshal([]byte(s), &result)
	return result.Balance
}

func luck() {
	mnemonic := random12Mnemonic()
	address := mnemonicToAddress(mnemonic)
	//address = "0xfC9fe0DfAEc95fbb5dF550CD07299CBD7727e09f"
	hex := getBalance(address)
	wei := hexToBigInt(hex)
	decimalD1 := decimal.NewFromBigInt(wei, 32)
	decimalD2 := decimal.NewFromBigInt(weiScale, 32)
	decimalResult := decimalD1.Div(decimalD2)
	balance, _ := decimalResult.Float64()
	if balance > 0 {
		fmt.Printf("助记词：%s\n地址：%s\n余额：%v\n", mnemonic, address, balance)
		recordToLocal(mnemonic)
	} else {
		fmt.Println("地址余额为0 --> 跳过", address)
	}
}

func recordToLocal(mnemonic string) {
	content := []byte(mnemonic)
	filename := "balance_" + time.Now().Format("2006-01-02 15:04:05")
	err := ioutil.WriteFile(filename, content, 0644)
	if err != nil {
		panic(err)
	}
}
