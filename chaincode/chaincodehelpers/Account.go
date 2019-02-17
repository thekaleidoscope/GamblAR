package chaincodehelpers

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func SetAccountInfo(stub shim.ChaincodeStubInterface, args []string) (string, error) {

	if len(args) != 2 {
		return "", fmt.Errorf("Incorrect Arguments. Expecting a key and value")
	}

	err := stub.PutState(args[0], []byte(args[1]))

	if err != nil {
		return "", fmt.Errorf("Failed to set asset: %s", args[0])
	}

	// Notify listeners that an event "eventInvoke" have been executed
	err = stub.SetEvent("eventInvoke", []byte{})
	if err != nil {
		return "", fmt.Errorf("Failed to set event: %s", err)
	}

	// Return this value in response
	return "", nil
}

func GetAccountInfo(stub shim.ChaincodeStubInterface, args []string) (string, error) {

	if len(args) != 1 {
		return "", fmt.Errorf("Incorrect Arguments. Expecting a key")
	}

	value, err := stub.GetState(args[0])

	if err != nil {
		return "", fmt.Errorf("Failed to set asset: %s with error", args[0], err)
	}

	if value == nil {
		return "", fmt.Errorf("Asset not found: %s", args[0])
	}

	return string(value), nil
}
