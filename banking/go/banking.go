package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
	//"github.com/hyperledger/fabric-chaincode-go/shim"
	//sc "github.com/hyperledger/fabric-protos-go/peer"
	"strconv"
)

type SmartContract struct {
}

type Account struct {
	Owner string `json:"owner"`
	//Token   string `json:"token"`
	Balance int `json:"balance"`
}

func (t *SmartContract) Init() sc.Response {

	//fmt.Println("Init")
	//var err error
	//// Initialize the chaincode
	//fmt.Printf("Aval = %d, Bval = %d\n", Aval, Bval)
	//err = stub.PutState(A, []byte(strconv.Itoa(Aval)))
	//if err != nil {
	//	return shim.Error("")
	//}
	//
	//err = stub.PutState(B, []byte(strconv.Itoa(Bval)))
	//if err != nil {
	//	return shim.Error("")
	//}
	//
	//return shim.Success(nil)
	//fmt.Println("Init")
	//_, args := stub.GetFunctionAndParameters()
	//var A, B string           // Entities
	//var Atoken, Btoken string // Token name holdings
	//var Aval, Bval int        // Asset holdings
	//var err error
	//
	//if len(args) != 4 {
	//	return shim.Error("Incorrect number of arguments. Expecting 4")
	//}
	//
	//// Initialize the chaincode
	//A = args[0]
	//Atoken = args[1]
	//Aval, err = strconv.Atoi(args[2])
	//if err != nil {
	//	return shim.Error("Expecting integer value for asset holding")
	//}
	//B = args[3]
	//Btoken = args[4]
	//Bval, err = strconv.Atoi(args[5])
	//if err != nil {
	//	return shim.Error("Expecting integer value for asset holding")
	//}
	//fmt.Printf("%s has %s !! val = %d\n, %s has %s !! val = %d", A, Atoken, Aval, B, Btoken, Bval)
	//
	//// Write the state to the ledger
	//err = stub.PutState(A, []byte(strconv.Itoa(Aval)))
	//if err != nil {
	//	return shim.Error(err.Error())
	//}
	//
	//err = stub.PutState(B, []byte(strconv.Itoa(Bval)))
	//if err != nil {
	//	return shim.Error(err.Error())
	//}

	return shim.Success(nil)
}

// Transaction makes payment of X units from A to B
func (t *SmartContract) Invoke(stub shim.ChaincodeStubInterface) sc.Response {
	function, args := stub.GetFunctionAndParameters()
	if function == "transfer" {
		// Make payment of X units from A to B
		return t.transfer(stub, args)
	} else if function == "deleteAccount" {
		// Deletes an entity from its state
		return t.deleteAccount(stub, args)
	} else if function == "query" {
		return t.query(stub, args)
	} else if function == "createAccount" {
		return t.createAccount(stub, args)
	}

	return shim.Error("Invalid invoke function name. Expecting \"invoke\" \"delete\" \"query\"")
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	accounts := []Account{
		Account{Owner: "Tomoko", Balance: 200},
		Account{Owner: "Brad", Balance: 100},
		Account{Owner: "Ji Soo", Balance: 500},
		Account{Owner: "Max", Balance: 1000},
		Account{Owner: "Amy", Balance: 300},
		Account{Owner: "Michel", Balance: 400},
		Account{Owner: "Ann", Balance: 600},
		Account{Owner: "Pari", Balance: 150},
		Account{Owner: "Valeria", Balance: 700},
		Account{Owner: "Shotaro", Balance: 800},
	}

	i := 0
	for i < len(accounts) {
		fmt.Println("i is ", i)
		accountAsBytes, _ := json.Marshal(accounts[i])
		owner := accounts[i].Owner
		err := APIstub.PutState(owner, accountAsBytes)
		if err != nil {
			return sc.Response{}
		}
		fmt.Println("Added", accounts[i])
		i = i + 1
	}

	return shim.Success(nil)
}

func (s *SmartContract) createAccount(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}
	strBalance, err := strconv.Atoi(args[2])
	if err != nil {
		return shim.Error(err.Error())
	}
	var account = Account{Owner: args[1], Balance: strBalance}

	accountAsBytes, _ := json.Marshal(account)
	err = APIstub.PutState(args[0], accountAsBytes)
	if err != nil {
		return sc.Response{}
	}

	return shim.Success(nil)
}

func (t *SmartContract) transfer(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	//(a,b,10)
	var err error
	var sender, receiver string // Entities
	var remittance int
	var Sval, Rval int

	sender = args[0]
	receiver = args[1]

	senderAccountAsBytes, err := stub.GetState(sender)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if senderAccountAsBytes == nil {
		return shim.Error("Entity not found")
	}
	Sval, _ = strconv.Atoi(string(senderAccountAsBytes))

	receiverAccountAsBytes, err := stub.GetState(receiver)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if receiverAccountAsBytes == nil {
		return shim.Error("Entity not found")
	}
	Rval, _ = strconv.Atoi(string(receiverAccountAsBytes))

	//	token = args[2]
	remittance, err = strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("Invalid transaction amount, expecting a integer value")
	}

	Sval = Sval - remittance
	Rval = Rval + remittance
	fmt.Printf("Sval = %d, Rval = %d\n", Sval, Rval)

	// Write the state back to the ledger
	err = stub.PutState(sender, []byte(strconv.Itoa(Sval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(receiver, []byte(strconv.Itoa(Rval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// Delete  an entity from state
func (t *SmartContract) deleteAccount(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	A := args[0]
	// Delete the key from the state in ledger
	err := stub.DelState(A)
	if err != nil {
		return shim.Error("Failed to delete state")
	}

	return shim.Success(nil)
}

// Query callback representing the query of a chaincode
func (t *SmartContract) query(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	var err error
	A := args[0]
	// Get the state from the ledger
	Avalbytes, err := stub.GetState(A)
	if err != nil {

		jsonResp := "{\"Error\":\"Failed to get state for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	if Avalbytes == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	jsonResp := "{\"Name\":\"" + A + "\",\"Amount\":\"" + string(Avalbytes) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return shim.Success(nil)
}

func main() {
	err := shim.Start(new(SmartContract))
	if err != nil {
		panic(err.Error())
	}
	//if err := cc.Start(); err != nil {
	//	fmt.Printf("Error starting ABstore chaincode: %s", err)
	//}
}
