package atomic_elements

import (
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/pkg/errors"
)

var (
	game = make(Game)
)

func (handler *Handeler) SetAsset(args []string) (string, error) {

	// Prepare arguments

	args = append(args, "setAccountInfo")

	// Add data that will be visible in the proposal, like a description of the invoke request
	transientDataMap := make(map[string][]byte)
	transientDataMap["result"] = []byte("Transient data in SetAsset")

	clientContext := handler.sdk.ChannelContext(handler.ChannelID, fabsdk.WithUser("Admin"))
	client, err := channel.New(clientContext)
	if err != nil {
		return "", errors.WithMessage(err, "failed to create new channel client ")
	}

	// Create a request (proposal) and send it
	response, err := client.Execute(channel.Request{ChaincodeID: handler.ChainCodeID, Fcn: args[len(args)-1], Args: [][]byte{[]byte(args[0]), []byte(args[1])}, TransientMap: transientDataMap})
	if err != nil {
		return "", fmt.Errorf("failed to edit asset value: %v", err)
	}

	return string(response.TransactionID), nil

}

func (handler *Handeler) QueryAsset(asset string) (string, error) {

	// Prepare arguments
	var args []string
	args = append(args, "getAccountInfo")
	args = append(args, asset)

	clientContext := handler.sdk.ChannelContext(handler.ChannelID, fabsdk.WithUser("Admin"))
	client, err := channel.New(clientContext)
	if err != nil {
		return "", errors.WithMessage(err, "failed to create new channel client ")
	}
	response, err := client.Query(channel.Request{ChaincodeID: handler.ChainCodeID, Fcn: args[0], Args: [][]byte{[]byte(args[1])}})
	if err != nil {
		return "", fmt.Errorf("failed to query: %v", err)
	}

	return string(response.Payload), nil
}
