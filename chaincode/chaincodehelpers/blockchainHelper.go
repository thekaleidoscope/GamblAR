package chaincodehelpers

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func WriteToBlockchain(stub shim.ChaincodeStubInterface, args []string) (string, error) {

	if len(args) != 2 {
		return "", fmt.Errorf("Incorrect Arguments. Expecting a key and value")
	}

	err := stub.PutState(args[0], []byte(args[1]))

	if err != nil {
		return "", fmt.Errorf("Failed to set account with info: %s", args[0])
	}

	// Return this value in response
	return "", nil
}

func ReadFromBlockchain(stub shim.ChaincodeStubInterface, args []string) (string, error) {

	if len(args) != 1 {
		return "", fmt.Errorf("Incorrect Arguments. Expecting a key")
	}

	value, err := stub.GetState(args[0])

	if err != nil {
		return "", fmt.Errorf("Failed to get account with info: %s with error", args[0], err)
	}

	if value == nil {
		return "", fmt.Errorf("account not found: %s", args[0])
	}

	return string(value), nil
}
