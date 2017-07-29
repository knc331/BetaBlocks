import ( time

)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type User struct {
    Name string
	Balance int64 `json:"balance,string"`
    Units int64 `json:"units,string"`
    
}	

type TradeManager struct {
    BuySide []struct {
    Name string
    Price int64
    Units int64
    Time time.Time
    
	}	

	SellSide []struct {
    Name string
    Price int64
    Units int64
    Time time.Time
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
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 4")
	}

	// Initialize the  seller with key arg[]
	Name = args[0]
	Balance, err = strconv.Atoi(args[1])
	if err != nil {
		return nil, errors.New("Expecting integer value for buyer A ")
	}
	Units, err = strconv.Atoi(args[2])
	if err != nil {
		return nil, errors.New("Expecting integer value for buyer A ")
	}


u1 := User{Name, Balance, Units}
userByteArray, err := json.Marshal(u1)
	//display the input values 
	fmt.Printf("Name = %s, Balance = %d , Units = %d\n ", Name, Balance, Units)

	// Write the state to the ledger
	err = stub.PutState(Name, userByteArray)
	if err != nil {
		return nil, err
	}

	return nil, nil
}


/////////////////place order function


/*func (t *SimpleChaincode) placeOrder(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, value string
	var err error
	fmt.Println("running placeOrder()")

	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	}

	user = args[0]    //Order id generate????
	value = args[1]	 //Order value ( includes username, value, type_of_order ie buy or sell)
	err = stub.PutState(key, []byte(value))  //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}*/


/////////INVOKE function


func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

    // Handle different functions
    if function == "init" {
        return t.Init(stub, "init", args)
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
    	unmarshalUser := User{}                         //read a variable
        return json.Unmarshal(t.readUser(stub, args), &unmarshalUser)
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
