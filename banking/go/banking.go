package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	//"github.com/hyperledger/fabric-chaincode-go/shim"
	//"github.com/hyperledger/fabric-protos-go/peer"
	"strconv"
)

type Chaincode struct {
}

type Account struct {
	OwnerID string `json:"owner_id"`
	Balance int64  `json:"balance"`
}

func (c *Chaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

func (c *Chaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	function, args := stub.GetFunctionAndParameters()
	if function == "createAccount" {
		return chaincode.createAccount(stub, args)
	} else if function == "deleteAccount" {
		return chaincode.deleteAccount(stub, args)
	} else if function == "deposit" {
		return chaincode.deposit(stub, args)
	} else if function == "getAccount" {
		return chaincode.getAccount(stub, args)
	} else if function == "transfer" {
		return chaincode.transfer(stub, args)
	} else if function == "withdrawal" {
		return chaincode.withdrawal(stub, args)
	}
	return shim.Error("Invalid function name")
}

// params:
//   - OwnerID: a name of the account
//   - Balance : initial balance
func (c *Chaincode) createAccount(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	var amount int64

	accountAsBytes, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("Failed to get OwnerID: " + err.Error())
	} else if accountAsBytes != nil {
		fmt.Println("This ownerID already exists: " + args[0])
		return shim.Error("This account already exists: " + args[0])
	}
	amount, err = strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return shim.Error(err.Error())
	}
	var account = Account{OwnerID: args[0], Balance: amount}
	accountAsBytes, err = json.Marshal(account)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(args[0], accountAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(accountAsBytes)
}

// params:
//   - OwnerID: a name of the account
func (c *Chaincode) deleteAccount(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	err := stub.DelState(args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func (c *Chaincode) deposit(stub shim.ChaincodeStubInterface, args []string) peer.Response { //입금
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	var amount int64
	var targetAccount Account

	accountAsBytes, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error(err.Error())
	} else if accountAsBytes == nil {
		return shim.Error("This account dose not exists: " + args[0])
	}
	err = json.Unmarshal(accountAsBytes, &targetAccount)
	if err != nil {
		return shim.Error(err.Error())
	}
	amount, err = strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return shim.Error(err.Error())
	}
	targetAccount.Balance = targetAccount.Balance + amount
	accountAsBytes, err = json.Marshal(targetAccount)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(args[0], accountAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(accountAsBytes)
}

func (c *Chaincode) getAccount(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	//(a)
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	accountAsBytes, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error(err.Error())
	} else if accountAsBytes == nil {
		return shim.Error("This account dose not exists: " + args[0])
	}
	return shim.Success(accountAsBytes)
}

func (c *Chaincode) transfer(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	//(a,b,10)
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}
	var sender, receiver string
	var remittance int64
	var senderAccount, receiverAccount Account

	sender = args[0]
	receiver = args[1]

	if sender == receiver {
		return shim.Error("receiver and sender cannot be the same")
	}
	senderAccountAsBytes, err := stub.GetState(sender)
	if err != nil {
		return shim.Error(err.Error())
	} else if senderAccountAsBytes == nil {
		return shim.Error("SenderAccount not found")
	}
	err = json.Unmarshal(senderAccountAsBytes, &senderAccount)
	if err != nil {
		return shim.Error(err.Error())
	}

	receiverAccountAsBytes, err := stub.GetState(receiver)
	if err != nil {
		return shim.Error(err.Error())
	} else if receiverAccountAsBytes == nil {
		return shim.Error("ReceieverAccount not found")
	}
	err = json.Unmarshal(receiverAccountAsBytes, &receiverAccount)
	if err != nil {
		return shim.Error(err.Error())
	}

	remittance, err = strconv.ParseInt(args[2], 10, 64)
	if err != nil {
		return shim.Error(err.Error())
	}
	if remittance < 0 {
		return shim.Error("remittance cannot be negative")
	} else if remittance > senderAccount.Balance {
		return shim.Error("remittance must be less than the balance")
	}
	senderAccount.Balance = senderAccount.Balance - remittance
	receiverAccount.Balance = receiverAccount.Balance + remittance

	senderAccountAsBytes, err = json.Marshal(senderAccount)
	if err != nil {
		return shim.Error(err.Error())
	}
	receiverAccountAsBytes, err = json.Marshal(receiverAccount)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(sender, senderAccountAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(receiver, receiverAccountAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(senderAccountAsBytes)
}

func (c *Chaincode) withdrawal(stub shim.ChaincodeStubInterface, args []string) peer.Response { //출금
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	var amount int64
	var targetAccount Account
	accountAsBytes, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error(err.Error())
	} else if accountAsBytes == nil {
		return shim.Error("This account dose not exists: " + args[0])
	}
	err = json.Unmarshal(accountAsBytes, &targetAccount)
	if err != nil {
		return shim.Error(err.Error())
	}
	amount, err = strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return shim.Error(err.Error())
	} else if amount > targetAccount.Balance {
		return shim.Error("amount must be less than the balance")
	}
	targetAccount.Balance = targetAccount.Balance - amount
	accountAsBytes, err = json.Marshal(targetAccount)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(args[0], accountAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(accountAsBytes)
}

//func strToInt(strBal string) (iBal int64, err error) {
//	iBal, err = strconv.ParseInt(strBal, 10, 64)
//	if err != nil {
//		return 0, err
//	}
//	return iBal, err
//}

func main() {
	err := shim.Start(new(Chaincode))
	if err != nil {
		panic(err.Error())
	}
}
