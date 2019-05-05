package atomic_elements

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/pkg/errors"
)

func (gameDB GameDB) AddBet(name string, acc string,
	bidAmt string, selection string) error {
	if !gameDB[name].completed {
		gameDB[name].game[name] =
			Bid{name, acc, bidAmt, selection}
	} else {
		return errors.Errorf("Game already ended.")
	}
	return nil
}

//Begin game and end betting period
func (game GameDB) EndBetting(name string) {
	game[name] = GameMeta{game[name].game, true}

}

func (handler *Handeler) WriteBetsToLedger(game GameDB, name string) error {

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
