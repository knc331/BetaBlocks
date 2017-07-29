package main

import (
 	"errors"
	"fmt"
	"strconv"
	"time"
    "encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type User struct {
    Name string
	Balance int `json:"balance,string"`
    Units int `json:"units,string"`

}

type  Trade struct {
    Name string
    Price int
    Units int
    Time time.Time
    Ordertype string

	}
type TradeManager struct {
  	BuySide []Trade
	SellSide []Trade 
	}

}


func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode - %s", err)
	}
}



func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Printf("Init called, initializing chaincode")

	//var user_A, buy_user_B, user_C, user_D string    // Entities
	//var Aval, Bval, Cval, Dval int // Asset holdings
	var Name string
	var Balance, Units int
	var err error

	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 4")
	}

	// Initialize with key arg[]
	Name = args[0]
	Balance, err = strconv.Atoi(args[1])
	if err != nil {
		return nil, errors.New("Expecting integer value for balance")
	}
	Units, err = strconv.Atoi(args[2])
	if err != nil {
		return nil, errors.New("Expecting integer value for units")



}
//writing the user to the blockchain

u1 := User{Name, Balance, Units}
userByteArray, err := json.Marshal(u1)
	//display the input values
	fmt.Printf("Name = %s, Balance = %d , Units = %d\n ", Name, Balance, Units)

	// Write the state to the ledger
	err = stub.PutState(Name, userByteArray)
	if err != nil {
		return nil, err
	}
//initializing the trade manager


	return nil, nil
}




func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

    // Handle different functions
    if function == "init" {
        return t.Init(stub, "init", args)
    } else if function =="tradeManager" {
    	return t.tradeManager(stub,"tradeManager",args)
    }

   /*else if function == "placeOrder" {
        return t.placeOrder(stub, args)
    } else if function == "updateOrderManager"{
    	  return t.updateOrderManager(stub, args)
    }*/
    fmt.Println("invoke did not find func: " + function)

    return nil, errors.New("Received unknown function invocation: " + function)
}


// Queries do not result in blocks being added to the chain, and you cannot use functions like PutState inside of Query or any helper functions it calls.


func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

    // Handle different functions
    if function == "readUser" {
	   unmarshalUser := new(User)                         //read a variable
	    //return t.readUser(stub, args)
        return json.Unmarshal(t.readUser(stub, args), &unmarshalUser)
    } else if function == "readTradeManager" {
    	unmarshalTradeManager:= new(TradeManager)
    	return json.Unmarshal(t.tradeManager(stub,"readTradeManger",args), &unmarshalTradeManager)

    }
    fmt.Println("query did not find func: " + function)

    return nil, errors.New("Received unknown function query: " + function)
}


// Reading the OrderManager ...called by Query
func (t *SimpleChaincode) readUser(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error

    if len(args) != 1 {
        return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
    }

    key = args[0]
    valAsbytes, err := stub.GetState(key)
    if err != nil {
        jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
        return nil, errors.New(jsonResp)
    }

    return valAsbytes, nil
}

// Reading the OrderManager ...called by Query
func (t *SimpleChaincode) tradeManager(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	
	var Name, Ordertype string
	var Price, Units int
	var Time time.Time
	var err error
	var tradeManager TradeManager

    if len(args) != 5 {
        return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
    }

    // Initialize the  seller with key arg[]
	Name = args[0]
	Price, err = strconv.Atoi(args[1])
	if err != nil {
		return nil, errors.New("Expecting integer value for balance")
	}
	Units, err = strconv.Atoi(args[2])
	if err != nil {
		return nil, errors.New("Expecting integer value for units")
	Time = t := time.Now()
	Ordertype = args[3]
	tradeManagerKey = args[4]


	trade := Trade{Name, Price, Units, Time, Ordertype}
	tradeManager.BuySide=append(trade)

	tradeManagerByteArray, err := json.Marshal(tradeManager)
	// Write the state to the ledger
	err = stub.PutState(Name, tradeManagerByteArray)
	if err != nil {
		return nil, err
	}
	return nil, nil

/*	// Read Trade Manager

	TradeManager, err := stub.GetState(TradeManagerKey)


	 if err != nil {
        jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
        return nil, errors.New(jsonResp)
    }
	//display the input values
	fmt.Printf("Name = %s, Balance = %d , Units = %d\n ", Name, Balance, Units)
	if Ordertype=="buy" {
		// Write the state to the ledger
		err = stub.PutState(TradeManagerKey, userByteArray)
		if err != nil {
		return nil, err
		}

	} else if function =="sell"{
		// Write the state to the ledger
		err = stub.PutState(TradeManagerKey, userByteArray)
		if err != nil {
		return nil, err
	}
	


    return valAsbytes, nil*/
}

// Reading the Trades 
func (t *SimpleChaincode) readTradeManager (stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error
	var tradeManager TradeManager

    if len(args) != 1 {
        return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
    }

    key = args[0]
    tradeManager, err := stub.GetState(key)
    if err != nil {
        jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
        return nil, errors.New(jsonResp)
    }

    return tradeManager, nil
}


