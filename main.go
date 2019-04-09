package main

import (
	"fmt"
	"os"
	"time"

	"github.com/GamblAR/atomic_elements"
	"github.com/pkg/errors"
)

func main() {
	game_handler := atomic_elements.Handeler{
		ChannelID: "townchannel",

		//OrgInfo
		OrgAdmins: map[string]string{
			"Org1": "Admin",
			"Org2": "Admin",
		},
		Organizations: []string{"Org1", "Org2"},
		ConfigFile:    "config.yaml",

		ChainCodeID:     "account",
		ChaincodeGoPath: os.Getenv("GOPATH"),
		ChaincodePath:   "github.com/GamblAR/chaincode/",
	}

	fmt.Println("-------Initializing the game handler.--------")
	err := game_handler.Initializer()
	if err != nil {
		fmt.Println(errors.Wrap(err, "initialization of game handler failed"))

	} else {
		fmt.Println("-------Successful initializing the game handler.--------")
	}

	fmt.Println("------ Creating And Join Channel For each Organizations ----------")

	err = game_handler.CreateAndJoinChannel()
	if err != nil {
		fmt.Println(errors.Wrap(err, "CreateAndJoinChannel game handler failed"))

	} else {
		fmt.Println("------- Completed Chaincode Installation and Instantiation --------")

	}
	time.Sleep(30 * time.Second)

	fmt.Println("-------  Invoking the chaincode to Create Asset and Set Value --------")

	txID, err := game_handler.SetAsset([]string{"Acc1", "1000"})
	if err != nil {
		fmt.Println(errors.Wrap(err, "Account SetAsset game handler failed"))

	} else {
		fmt.Println("Setting Value TxID: ", txID)
		fmt.Println("-------Successful  SetAsset Account--------")
	}

	fmt.Println("-------  Query the chaincode to Get Value of Asset --------")

	txID, err = game_handler.QueryAsset("Acc1")
	if err != nil {
		fmt.Println(errors.Wrap(err, "Account QueryAsset game handler failed"))

	} else {
		fmt.Println("Query Value : ", txID)
		fmt.Println("-------Successful  QueryAsset Account--------")
	}

	game_handler.MakeBet("Game1", "Acc1", "2000")
	game_handler.MakeBet("Game1", "Acc2", "2000")

	err = game_handler.WriteBetsToLedger("Game1")
	if err != nil {
		fmt.Println(errors.Wrap(err, "Writing bets to blockchain failed"))

	} else {

		fmt.Println("-------Successfully written bets to blockchain--------")
	}

}
