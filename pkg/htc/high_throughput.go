package htc

import (
	"encoding/json"
	"errors"

	"github.com/Akachain/akc-go-sdk/util"
	cnst "github.com/Akachain/akc-htc-sdk/constants"
	lib "github.com/Akachain/akc-htc-sdk/libraries"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type HighThroughput struct {
	Name          string  `json:"Name"`
	WalletAddress string  `json:"WalletAddress"`
	Amount        float64 `json:"Amount"`
	Operation     string  `json:"Operation"`
	Status        string  `json:"Status"`
	Reference     string  `json:"Reference"`
	TxID          string  `json:"TxID"`
	CreatedAt     string  `json:"CreatedAt"`
	ProcessedAt   string  `json:"ProcessedAt"`
}

func NewHighThroughput(stub shim.ChaincodeStubInterface, name, wallet, op, ref string, amount float64) error {
	if wallet == "" || amount < 0 || (op != cnst.HTC_OP_MINUS && op != cnst.HTC_OP_PLUS) {
		return errors.New("Your variable name is unrecognized")
	}

	txID := stub.GetTxID()
	h := HighThroughput{
		Name:          name,
		WalletAddress: wallet,
		Amount:        amount,
		Operation:     op,
		Reference:     ref,
		Status:        cnst.Waiting,
		TxID:          txID,
		CreatedAt:     lib.TimeNow(stub),
	}

	err := util.Createdata(stub, cnst.HTC_PREFIX, []string{h.Name, h.WalletAddress, txID}, &h)
	return err
}

// func GetHighThroughputData(stub shim.ChaincodeStubInterface, args []string) pb.Response {
// 	if len(args) != 1 {
// 		resErr := ResponseError{
// 			ResCode: lib.ERR2,
// 			Msg:     fmt.Sprintf(lib.ResponseFomart, lib.ResCodeDict[lib.ERR2], "", GetLine())}
// 		return RespondError(resErr)
// 	}

// 	var pageSize int32
// 	errMarshal := json.Unmarshal([]byte(args[0]), &pageSize)
// 	if errMarshal != nil {
// 		resErr := ResponseError{
// 			ResCode: lib.ERR4,
// 			Msg:     fmt.Sprintf(lib.ResponseFomart, lib.ResCodeDict[lib.ERR4], errMarshal.Error(), GetLine())}
// 		return RespondError(resErr)
// 	}

// 	dataQuery, err := getHighThroughputData(stub, pageSize)
// 	if err != nil {
// 		resErr := ResponseError{
// 			ResCode: lib.ERR3,
// 			Msg:     fmt.Sprintf(lib.ResponseFomart, lib.ResCodeDict[lib.ERR3], err.Error(), GetLine())}
// 		return RespondError(resErr)
// 	}

// 	result, _ := json.Marshal(dataQuery)

// 	resSuc := ResponseSuccess{
// 		ResCode: lib.SUCCESS,
// 		Msg:     lib.ResCodeDict[lib.SUCCESS],
// 		Payload: string(result)}
// 	return RespondSuccess(resSuc)
// }

func PruneHighThroughput(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	var err error

	dataQuery, err := getHighThroughputData(stub, 100)

	// if len(args) != 1 {
	// 	resErr := ResponseError{
	// 		ResCode: lib.ERR2,
	// 		Msg:     fmt.Sprintf(lib.ResponseFomart, lib.ResCodeDict[lib.ERR2], "", GetLine())}
	// 	return RespondError(resErr)
	// }

	// var dataQuery []HighThroughput
	// errMarshal := json.Unmarshal([]byte(args[0]), &dataQuery)
	// if errMarshal != nil {
	// 	resErr := ResponseError{
	// 		ResCode: lib.ERR4,
	// 		Msg:     fmt.Sprintf(lib.ResponseFomart, lib.ResCodeDict[lib.ERR4], errMarshal.Error(), GetLine())}
	// 	return RespondError(resErr)
	// }

	responseMap := make(map[string][]HighThroughput)
	dataCalculate := make(map[string]map[string]float64)

	type HighThroughputReturn struct {
		Success []string
		Errors  []string
	}

	response := new(HighThroughputReturn)

	for _, value := range dataQuery {
		key := value.WalletAddress

		if value.Status == cnst.Processed {
			response.Success = append(response.Success, key)
			continue
		}

		if len(responseMap[key]) > 0 {
			responseMap[key] = append(responseMap[key], value)
		} else {
			var dataMap []HighThroughput
			dataMap = append(dataMap, value)
			responseMap[key] = dataMap

			calculateMap := make(map[string]float64)
			dataCalculate[key] = calculateMap

			response.Success = append(response.Success, key)
		}

		// Update HTC record to Processed
		value.Status = cnst.Processed
		value.ProcessedAt = lib.TimeNow(stub)
		err = util.UpdateExistingData(stub, cnst.HTC_PREFIX, []string{value.Name, value.WalletAddress, value.TxID}, &value)
		if err != nil {
			response.Errors = append(response.Errors, key)
			continue
		}

		if value.Operation == cnst.HTC_OP_PLUS {
			dataCalculate[key][value.Name] += value.Amount
		} else {
			dataCalculate[key][value.Name] -= value.Amount
		}
	}

	result, _ := json.Marshal(response)
	return string(result), nil
}

func getHighThroughputData(stub shim.ChaincodeStubInterface, pageSize int32) ([]HighThroughput, error) {
	var result = new(HighThroughput)
	var list = []HighThroughput{}

	var queryString = `
		{ "selector": 
			{ 	
				"Status": 
					{ "$eq": "Waiting" },
				"_id": 
					{"$gt": "\u0000__AKC_High_Throughput_",
					"$lt": "\u0000__AKC_High_Throughput_\uFFFF"}			
			}
		}`

	resultsIterator, _, err := stub.GetQueryResultWithPagination(queryString, pageSize, "")

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	// Check data response after query in database
	if !resultsIterator.HasNext() {
		return list, nil
	}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(queryResponse.Value, result)
		if err != nil {
			continue
		}
		list = append(list, *result)
	}
	return list, nil
}
