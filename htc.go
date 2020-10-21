package main

import (
	"fmt"
	. "github.com/Akachain/akc-go-sdk/common"
	"github.com/Akachain/akc-htc-sdk/pkg/htc"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// Chaincode implementation
type Chaincode struct {
}

func (s *Chaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func createData(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	err := htc.NewHighThroughput(stub, "test", "test", "+", "abc", 100)
	if err != nil {
		resErr := ResponseError{ResCode: "500", Msg: ""}
		return RespondError(resErr)
	}
	resSuc := ResponseSuccess{ResCode: SUCCESS, Msg: ResCodeDict[SUCCESS], Payload: ""}
	return RespondSuccess(resSuc)
}

func (s *Chaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	// Retrieve the requested Smart Contract function and arguments
	function, args := stub.GetFunctionAndParameters()
	switch function {
	//CreateAdmin
	case "CreateNew":
		return createData(stub, args)
	}
	return shim.Error(fmt.Sprintf("Invoke cannot find function " + function))
}

func main() {
	// Create a new Chain code
	err := shim.Start(new(Chaincode))
	if err != nil {
		fmt.Printf("Error creating new Chain code: %s", err)
	}
}
