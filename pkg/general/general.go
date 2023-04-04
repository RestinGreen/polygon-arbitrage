package general

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type General struct {
	IpcEndpoint    string
	InfuraEndpoint string
	Ctx            context.Context

	//database credentials
	Host     string
	Port     string
	User     string
	Password string
	DBName   string

	UniV2Methods map[string]string
}

func NewGeneral() *General {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Failed to load .env file.")
		panic(err)
	}

	var g General

	g.IpcEndpoint = os.Getenv("IPC_ENDPOINT")
	g.InfuraEndpoint = os.Getenv("INFURA_ENDPOINT")
	g.Host = os.Getenv("HOST")
	g.Port = os.Getenv("PORT")
	g.User = os.Getenv("USER")
	g.Password = os.Getenv("PASSWORD")
	g.DBName = os.Getenv("DBNAME")

	g.UniV2Methods = map[string]string{
		"fb3bdb41": "swapETHForExactTokens",
		"7ff36ab5": "swapExactETHForTokens",
		"b6f9de95": "swapExactETHForTokensSupportingFeeOnTransferTokens", //fee on transfer
		"18cbafe5": "swapExactTokensForETH",
		"791ac947": "swapExactTokensForETHSupportingFeeOnTransferTokens", //fee on transfer
		"38ed1739": "swapExactTokensForTokens",
		"5c11d795": "swapExactTokensForTokensSupportingFeeOnTransferTokens", //fee on transfer
		"4a25d94a": "swapTokensForExactETH",
		"8803dbee": "swapTokensForExactTokens",
	}
	return &g
}
