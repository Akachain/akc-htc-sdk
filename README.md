
# Akachain High Throughput SDK for Chaincode

## I. Overview

The Akachain High Throughput Chaincode (AKC HTC) is designed for applications handling hundreds or thousands transaction per second which all read or update the same asset (key) in the ledger.

This document provides the AKC HTC template interface and how to use.

## II. AKC HTC Interface

The AKC HTC sdk is packaged into akc_htc package which provide the following interfaces

##### Insert: The insert function inserts the value into the temporary storage (the state db that may be deleted later) as a single row. The key is unique and created by combining the input and transaction id.

```go
htc.NewHighThroughput(stub, <name>, <key>, <operation>, <ref>, <value>)
```

- Name: The name, object or attribute that applied the high throughput chaincode. Example: merchant, user ...
- Key: The key identify the object. Example: merchant address, user id ...
- Operation: The operation that used for aggregation. Currently support OP_PLUS (+) and OP_MINUS (-)
- Ref: reference for high throughput record.
- Value: The value of key. Currently for aggregation purpose only, so it should be in numeric type

#### Get: Get the value from temporary storage

```go
htc.PruneHighThroughput(stub)
```

## III. How to use

To use AKC HTC, the package akchtc must be imported to chaincode file. Ex:

```go
import (
  "encoding/json"
  "errors"
  "fmt"
  "reflect"
  "github.com/hyperledger/fabric/core/chaincode/shim"
  "github.com/Akachain/akc-htc-sdk/htc"
)

// example code insert using Akachain High throughput
func insertHTC(stub shim.ChaincodeStubInterface, args []string) (string, error) {
  // Init Akachain High Throughput
  err := htc.NewHighThroughput(stub, "test", "test", "+", "abc", 100)
  if err != nil {
  	resErr := ResponseError{ResCode: "500", Msg: ""}
  	return RespondError(resErr)
  }
  if res != nil {
    return fmt.Sprintf("Failure"), res
  }
  return fmt.Sprintf("Success"), nil
}

// Example prune data HTC
// This func response JSON data for "variableName" after prune success.
func pruneHTC(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	resp, err := akc.PruneHighThroughput(stub)

  if err != nil {
    return nil, err
  }

  return resp, nil
}
```