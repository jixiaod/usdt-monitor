package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/shopspring/decimal"
)

type Address struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

type Config struct {
	BotToken       string    `json:"botToken"`
	ChatID         int64     `json:"chatID"`
	USDTAPIBaseURL string    `json:"usdtAPIBaseURL"`
	Addresses      []Address `json:"addresses"`
	Threshold      float64   `json:"threshold"`
}

func startMonitoring() {

	log.Print("Monitor started")
	// Load configuration
	configFile, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	var config Config
	if err := json.Unmarshal(configFile, &config); err != nil {
		log.Fatalf("Error parsing config file: %v", err)
	}

	bot, err := tgbotapi.NewBotAPI(config.BotToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	lastBalances := make(map[string]decimal.Decimal)
	threshold := decimal.NewFromFloat(config.Threshold)
	for {
		for _, addr := range config.Addresses {
			apiURL := config.USDTAPIBaseURL + "/api/accountv2?address=" + addr.Address
			currentBalance, err := getUSDTBalance(apiURL)
			if err != nil {
				log.Println("Error fetching USDT balance for", addr.Address, ":", err)
				time.Sleep(3 * time.Second)
				continue
			}
			lastBalance, exists := lastBalances[addr.Address]
			currentBalanceDecimal := decimal.RequireFromString(currentBalance)
			if exists && !lastBalance.IsZero() && currentBalanceDecimal.Sub(lastBalance).Abs().GreaterThan(threshold) {
				msg := tgbotapi.NewMessage(config.ChatID, fmt.Sprintf("USDT balance changed for: \n name: %s \n address: %s \n balance: %s", addr.Name, addr.Address, currentBalanceDecimal.String()))
				bot.Send(msg)
			}

			lastBalances[addr.Address] = decimal.RequireFromString(currentBalance)
		}
		time.Sleep(3 * time.Second)
	}
}

func getUSDTBalance(apiURL string) (string, error) {
	resp, err := http.Get(apiURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response struct {
		WithPriceTokens []struct {
			TokenAbbr string `json:"tokenAbbr"`
			Balance   string `json:"balance"`
		} `json:"withPriceTokens"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}

	for _, token := range response.WithPriceTokens {
		if token.TokenAbbr == "USDT" {
			return token.Balance, nil
		}
	}

	return "", fmt.Errorf("USDT balance not found")
}
