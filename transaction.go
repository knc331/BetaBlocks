
package main

import (
 	"errors"
	"fmt"
	"strconv"
	"time"
    "encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/tidwall/gjson"
)

// Reading the OrderManager ...called by Query
func (t *SimpleChaincode) transactionManagerFunction(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	
	var buyerName, sellerName string
	var Price, Units int
	var Time time.Time
	var err error
	var trasaction Trasaction
	//var tradeManager TradeManager

    if len(args) != 4 {
        return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
    }

   

    // Initialize the  seller with key arg[]
	trasaction.BuyerName = args[0]
	
	transaction.SellerName = args[1]
	trasaction.TransactionValue, err = strconv.Atoi(args[2])
	if err != nil {
		return nil, errors.New("Expecting integer value for Price")
	}
	Units, err = strconv.Atoi(args[3])
	if err != nil {
		return nil, errors.New("Expecting integer value for Units")
	}

	transaction.Time =time.Now()

	Avalbytes, err := stub.GetState(trasaction.BuyerName ) 
	if err != nil { 
		return nil, errors.New("Failed to get state") 
	} 
	if Avalbytes == nil { 
		return nil, errors.New("Entity not found") 
	} 


 
	Bvalbytes, err := stub.GetState(transaction.SellerName) 
		if err != nil { 
			return nil, errors.New("Failed to get state") 
		} 
		if Bvalbytes == nil { 
			return nil, errors.New("Entity not found") 
		} 
	 

	var unmarshallAval byte[]
	json.Unmarshal(Avalbytes,&unmarshallAval)

	var unmarshallBval byte[]
	json.Unmarshal(Bvalbytes,&unmarshallBval)

	var APrice int 
	APrice = gjson.Parse(unmarshallAval).Get("Balance")
	var BPrice int 
	BPrice = gjson.Parse(unmarshallAval).Get("Balance")
	var AUnits int 
	AUnits = gjson.Parse(unmarshallAval).Get("Units")
	var BUnits int 
	BUnits = gjson.Parse(unmarshallAval).Get("Units")

	
	BPrice = BPrice + APrice
	Aprice = Aprice - BPrice

	BUnits = BUnits - AUnits
	AUnits = AUnits + BUnits

	unmarshallAval["Balance"] = Aprice
	unmarshallAval["Units"] = AUnits

	unmarshallBval["Balance"] = Bprice
	unmarshallBval["Units"] = BUnits

	AvalArray, err := json.Marshal(unmarshallAval)
	BvalArray, err := json.Marshal(unmarshallBval)

	// Write the state to the ledger
	err = stub.PutState(trasaction.BuyerName, AvalArray)) 
	if err != nil { 
 		return nil, err 
 	} 
 
 
 	err = stub.PutState(trasaction.SellerName, BvalArray)) 
 	if err != nil { 
 		return nil, err 
 	} 


 	err = stub.GetState(transactionKey) 
 	if err != nil { 
 		return nil, err 
 	} 

 	
 	err = stub.PutState(transactionKey,trasaction) 
 	if err != nil { 
 		return nil, err 
 	} 


}
}