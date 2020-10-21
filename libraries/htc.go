package libraries

import (
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// TimestampISO convert time unix to time ISO 8601
func TimestampISO(timeUnix int64) string {
	return time.Unix(timeUnix, 0).Format(time.RFC3339)
}

// TimestampUnix convert string time to unix time
func TimestampUnix(isoTime string) int64 {
	t, err := time.Parse(time.RFC3339, isoTime)
	if err != nil {
		panic(err)
	}

	ut := t.UnixNano() / int64(time.Millisecond)
	return ut
}

func TimeNow(stub shim.ChaincodeStubInterface) string {
	txTime, _ := stub.GetTxTimestamp()

	return TimestampISO(txTime.Seconds)
}
