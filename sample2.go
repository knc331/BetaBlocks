package main

import (
 	"errors"
	"fmt"
	"strconv"
	"time"
    "encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"math/rand"
	//"github.com/tidwall/gjson"
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

type Transaction struct{
	BuyerName string
	SellerName string
	TransactionValue int
	Units int
	Time time.Time
}

type TransactionManager struct{
	TransactionLedger []Transaction 
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
    } else if function =="tradeManagerFunction" {
    	return t.tradeManagerFunction(stub,args)
    } else if function =="transactionManagerFunction"{
    	return t.transactionManagerFunction(stub,args)
    } else if function =="PerformSettlement" {
    	return t.PerformSettlement(stub,args)
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
	  //var unmarshalUser User                         //read a variable
	    return t.readUser(stub, args)
        //return json.Unmarshal(t.readUser(stub, args), &unmarshalUser)
    }else if function == "readTradeManager" {
    	//unmarshalTradeManager:= new(TradeManager)
	    fmt.Println("inside readTradeManager")
    	return t.readTradeManager(stub,args)

    } else if function == "readTransactionManager"{
    	return t.readTransactionManager(stub,args)
    }
    fmt.Println("query did not find func: " + function)

    return nil, errors.New("Received unknown function query: " + function)
}


// Reading the OrderManager ...called by Query
func (t *SimpleChaincode) readUser(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error
	//var unmarshalUser User
    if len(args) != 1 {
        return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
    }

    key = args[0]
    valAsbytes, err := stub.GetState(key)
    if err != nil {
        jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
        return nil, errors.New(jsonResp)
    } /*else {
    	json.Unmarshal(valAsbytes,&unmarshalUser)
    }*/
	fmt.Println("unmarshalUser " + string(valAsbytes))
    return valAsbytes, nil
}

// Reading the OrderManager ...called by Query
func (t *SimpleChaincode) tradeManagerFunction(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	
	var Name, Ordertype, tradeManagerKey string
	var Price, Units int
	var Time time.Time
	var err error
	//var tradeManager TradeManager

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
	}

	Time = time.Now()
	Ordertype = args[3]
	tradeManagerKey = args[4]


	trade := Trade {Name, Price, Units, Time, Ordertype}
	
	//readTradeManager
	var marshalltradeMgr []byte 
	var unmarshalltradeMgr TradeManager
	var stringArr []string
	stringArr = append(stringArr, "trademanager")
	marshalltradeMgr, err = t.readTradeManager(stub, stringArr)
	json.Unmarshal(marshalltradeMgr,&unmarshalltradeMgr)
	if Ordertype == "buy" {
			unmarshalltradeMgr.BuySide = append(unmarshalltradeMgr.BuySide, trade);
	}else if Ordertype == "sell" {
			unmarshalltradeMgr.SellSide = append(unmarshalltradeMgr.SellSide, trade);
	}
	//fmt.Println("unmarshalTradeManager " + gjson.Parse(json).Get("BuySide.Name"))

	//check Order Type 
	//if buy 
	//append in buy 
	//put state



	//var tradeArray []Trade;
	//tradeArray = append(tradeArray, trade); 
	//tradeManager.BuySide=tradeArray;

	tradeManagerByteArray, err := json.Marshal(unmarshalltradeMgr)

	//tradeManagerByteArray, err := json.Marshal(tradeManager)
	//gjson.Parse(tradeManager).Get("name").Get("last")
	// Write the state to the ledger
	err = stub.PutState(tradeManagerKey, tradeManagerByteArray)
	if err != nil {
		return nil, err
	}

	return nil, nil

}

// Reading the Trades 
func (t *SimpleChaincode) readTradeManager (stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error
	var tradeManager []byte
	//var unmarshalTradeManager TradeManager

    if len(args) != 1 {
        return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
    }

    key = args[0]
    tradeManager, err = stub.GetState(key)
    if err != nil {
        jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
        return nil, errors.New(jsonResp)
    } /*else{
    	json.Unmarshal(tradeManager,&unmarshalTradeManager)
    }
    output,err :=  json.Marshal(unmarshalTradeManager.BuySide)
    if err != nil {
        //jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
        return nil, err
    }
    //fmt.Println("unmarshalTradeManager " + string(output)) */   
    return tradeManager, nil		
}

// Reading the Trades 
func (t *SimpleChaincode) readTransactionManager (stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error
	var transactionManagerByteArray []byte
	//var unmarshalTradeManager TradeManager

    if len(args) != 1 {
        return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
    }

    key = args[0]
    transactionManagerByteArray, err = stub.GetState(key)
    if err != nil {
        jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
        return nil, errors.New(jsonResp)
    } /*else{
    	json.Unmarshal(tradeManager,&unmarshalTradeManager)
    }
    output,err :=  json.Marshal(unmarshalTradeManager.BuySide)
    if err != nil {
        //jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
        return nil, err
    }
    //fmt.Println("unmarshalTradeManager " + string(output)) */   
    return transactionManagerByteArray, nil		
}

// Reading the Trades 
func (t *SimpleChaincode) PerformSettlement (stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	

	var marshalltradeMgr []byte 
	var unmarshalltradeMgr TradeManager
	var stringArr []string
	var err error
	stringArr = append(stringArr, "trademanager")
	marshalltradeMgr, err = t.readTradeManager(stub, stringArr)
	if err!=nil {
		return nil, err
	}

	json.Unmarshal(marshalltradeMgr,&unmarshalltradeMgr)

	var BuySide []Trade
	var SellSide []Trade

	BuySide = unmarshalltradeMgr.BuySide
	SellSide = unmarshalltradeMgr.SellSide
	
	BuySide = qsort(BuySide)
	SellSide = qsort(SellSide)

	//perform settlement
	for i := range BuySide {
		 for j := range BuySide {   
		    if BuySide[i].Price >= SellSide[j].Price {
		      if BuySide[i].Units  <= SellSide[j].Units{
		      	//execute Transaction
		      	var transactionValue int
		      	var stringArr []string
		      	stringArr= append(stringArr,BuySide[i].Name)
		      	stringArr= append(stringArr,SellSide[j].Name)
		      	transactionValue = BuySide[i].Price*BuySide[i].Units
		      	stringArr= append(stringArr,strconv.Itoa(transactionValue))
		      	stringArr = append(stringArr,strconv.Itoa(BuySide[i].Units))
		      	t.transactionManagerFunction(stub,stringArr)
		      }
		    }
		 }
	}

	return nil,nil

	//APrice = gjson.Parse(unmarshallAval).Get("Balance")


	/*var key, jsonResp string
	var err error
	var tradeManager []byte
	//var unmarshalTradeManager TradeManager

    if len(args) != 1 {
        return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
    }

    key = args[0]
    tradeManager, err = stub.GetState(key)
    if err != nil {
        jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
        return nil, errors.New(jsonResp)
    } /*else{
    	json.Unmarshal(tradeManager,&unmarshalTradeManager)
    }
    output,err :=  json.Marshal(unmarshalTradeManager.BuySide)
    if err != nil {
        //jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
        return nil, err
    }*/
    //fmt.Println("unmarshalTradeManager " + string(output)) */   
   //return tradeManager, nil		
}

// Reading the OrderManager ...called by Query
func (t *SimpleChaincode) transactionManagerFunction(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	
	//var buyerName, sellerName, 
	var transactionKey string
	//var Price, Units int
	//var Time time.Time
	var err error
	var transaction Transaction
	//var tradeManager TradeManager

    if len(args) != 5 {
        return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
    }

   

    // Initialize the  seller with key arg[]
	transaction.BuyerName = args[0]
	
	transaction.SellerName = args[1]
	transaction.TransactionValue, err = strconv.Atoi(args[2])
	if err != nil {
		return nil, errors.New("Expecting integer value for Price")
	}
	transaction.Units, err = strconv.Atoi(args[3])
	if err != nil {
		return nil, errors.New("Expecting integer value for Units")
	}

	transaction.Time =time.Now()

	Avalbytes, err := stub.GetState(transaction.BuyerName ) 
	if err != nil { 
		return nil, errors.New("Failed to get state") 
	} 
	if Avalbytes == nil { 
		return nil, errors.New("Entity not found") 
	} 

	 transactionKey=args[4]
 
	Bvalbytes, err := stub.GetState(transaction.SellerName) 
		if err != nil { 
			return nil, errors.New("Failed to get state") 
		} 
		if Bvalbytes == nil { 
			return nil, errors.New("Entity not found") 
		} 
	 

	var unmarshallAval User
	json.Unmarshal(Avalbytes, &unmarshallAval)

	var unmarshallBval User
	json.Unmarshal(Bvalbytes,&unmarshallBval)

	var APrice int 
	//APrice = gjson.Parse(unmarshallAval).Get("Balance")
	APrice=unmarshallAval.Balance
	var BPrice int 
	//BPrice = gjson.Parse(unmarshallAval).Get("Balance")
	BPrice = unmarshallBval.Balance
	var AUnits int 
	//AUnits = gjson.Parse(unmarshallAval).Get("Units")
	AUnits = unmarshallAval.Units
	var BUnits int 
	//BUnits = gjson.Parse(unmarshallAval).Get("Units")
	BUnits = unmarshallBval.Units
	
	BPrice = BPrice + APrice
	APrice = APrice - BPrice

	BUnits = BUnits - AUnits
	AUnits = AUnits + BUnits

	unmarshallAval.Balance = APrice
	unmarshallAval.Units = AUnits

	unmarshallBval.Balance = BPrice
	unmarshallBval.Units = BUnits

	AvalArray, err := json.Marshal(unmarshallAval)
	BvalArray, err := json.Marshal(unmarshallBval)

	// Write the state to the ledger
	err = stub.PutState(transaction.BuyerName, AvalArray) 
	if err != nil { 
 		return nil, err 
 	} 
 
 
 	err = stub.PutState(transaction.SellerName, BvalArray) 
 	if err != nil { 
 		return nil, err 
 	} 

	var marshallTransactionMgr []byte
	var unmarshalTransactionMr []Transaction
 	marshallTransactionMgr, err = stub.GetState(transactionKey) 
 	json.Unmarshal(marshallTransactionMgr,&unmarshalTransactionMr)
 	if err != nil { 
 		return nil, err 
 	} 
 	unmarshalTransactionMr = append(unmarshalTransactionMr, transaction);

 	transactionManagerByteArray, err := json.Marshal(unmarshalTransactionMr)
	err = stub.PutState(transactionKey, transactionManagerByteArray)
	if err != nil {
		return nil, err
	}

	return nil, nil
}



func qsort(BuySide []Trade) []Trade {
		 	 if len(BuySide) < 2 { return BuySide }

		  left, right := 0, len(BuySide) - 1

		  // Pick a pivot
		  pivotIndex := rand.Int() % len(BuySide)

		  // Move the pivot to the right
		  BuySide[pivotIndex], BuySide[right] = BuySide[right], BuySide[pivotIndex]

		  // Pile elements smaller than the pivot on the left
		  for i := range BuySide {
		    if BuySide[i].Price < BuySide[right].Price {
		      BuySide[i], BuySide[left] = BuySide[left], BuySide[i]
		      left++
		    }
		  }

		  // Place the pivot after the last smaller element
		  BuySide[left], BuySide[right] = BuySide[right], BuySide[left]

		  // Go down the rabbit hole
		  qsort(BuySide[:left])
		  qsort(BuySide[left + 1:])


		  return BuySide
	}