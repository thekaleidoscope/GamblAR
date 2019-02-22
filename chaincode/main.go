package main

import (
	"fmt"

	"github.com/GamblAR/chaincode/chaincodehelpers"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type Account struct {
}

func (t *Account) Init(stub shim.ChaincodeStubInterface) peer.Response {

	fn, _ := stub.GetFunctionAndParameters()

	if fn != "init" {
		return shim.Error("Incorrect Function")
	}

	return shim.Success(nil)
}

func (t *Account) Invoke(stub shim.ChaincodeStubInterface) peer.Response {

	fn, args := stub.GetFunctionAndParameters()

	var result string
	var err error

	if fn == "setAccountInfo" {
		result, err = chaincodehelpers.SetAccountInfo(stub, args)
	} else if fn == "getAccountInfo" {
		result, err = chaincodehelpers.GetAccountInfo(stub, args)
	}

	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte(result))
}

func main() {
	if err := shim.Start(new(Account)); err != nil {
		fmt.Printf("Error starting Account chaincode %s", err)
	}
}
