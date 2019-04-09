package atomic_elements

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/pkg/errors"
)

type gameMeta struct {
	game      Game
	completed bool
}

func (game Game) AddBet(name string, acc string, bidAmt string) {
	//Add the bet details
	game[acc] = Bid{name, acc, bidAmt}
}

func (handler *Handeler) MakeBet(name string, acc string, bidAmt string) {

	game.AddBet(name, acc, bidAmt)

}

func (handler *Handeler) WriteBetsToLedger(name string) error {

	gameData, err := json.Marshal(game)
	if err != nil {
		return fmt.Errorf("Failed to json marshall message %v", err)
	}
	args := []string{name}
	args = append(args, "setAccountInfo")
	clientContext := handler.sdk.ChannelContext(handler.ChannelID, fabsdk.WithUser("Admin"))
	client, err := channel.New(clientContext)
	if err != nil {
		return errors.WithMessage(err, "failed to create new channel client ")
	}

	// Create a request (proposal) and send it
	_, err = client.Execute(channel.Request{ChaincodeID: handler.ChainCodeID, Fcn: args[len(args)-1], Args: [][]byte{[]byte(args[0]), gameData}})
	if err != nil {
		return fmt.Errorf("failed to edit asset value: %v", err)
	}

	return nil
}
