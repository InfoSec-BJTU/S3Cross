/*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
* http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

package main

/* Imports
 * 4 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
	"golang.org/x/crypto/sha3"
)

// Define the Smart Contract structure
type SmartContract struct {
}

type Psd struct {
	// PKP   string `json:"PKP"`
	Ct    string `json:"Ct"`
	H_pow string `json:"H_pow"`
	Acc   string `json:"Acc"`
	Time  string `json:"Time"`
}

/*
 * The Init method is called when the Smart Contract "fabcar" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract "fabcar"
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "queryPsd" {
		return s.queryPsd(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "createPsd" {
		return s.createPsd(APIstub, args)
	} else if function == "queryAllPsds" {
		return s.queryAllPsds(APIstub)
	} else if function == "deleteAllPsds" {
		return s.queryAllPsds(APIstub)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) queryPsd(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	psdAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(psdAsBytes)
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	time_t := strconv.FormatInt(time.Now().Unix(), 10)

	psds := []Psd{
		{Ct: "c1+c2", H_pow: "Hash-of-PKP-Ct-and-the-proofs", Acc: "Doamin-Accumulator", Time: time_t},
	}

	// PKP: "Public-key-of-the-pseudonym"
	for _, psd := range psds {
		psdAsBytes, _ := json.Marshal(psd)
		APIstub.PutState("Public-key-of-the-pseudonym", psdAsBytes)
		fmt.Println("Added", psd)
	}

	return shim.Success(nil)
}

// Demo input:
// pkp: 02d37c24ddd0aaec7c3f1efe95d9cfbbec6b5c9c90291aac1ab1556a36d8cac95a
// ct: c1: 02d37c24ddd0aaec7c3f1efe95d9cfbbec6b5c9c90291aac1ab1556a36d8cac95a, c2: 03069eccec991367ac7f65da48e8b72996fb85ee79f3c24a04c6c665a2b09370f6
// hpow: bcb607133f8f8b7622eeb54894fcef53c2ba51cf57138d3993b48dcaf39c1296
// acc: e3dc3fc8df6b2af80428d69773bf37eec17b3b00017232df9f95578cb1c1c49d398755e87ff4a009db2c561d4bec55a1b31f21e88866832babff0b0c5b673d78e4713e613424624852cfd9558a3e2f4dbcb50f720a469e75a51e0d51d55e8cfe5a73e7b8fe5f82a1fb2b86df6317ad006166f21172b207b183fcbab6bad3e31a0b8f20e7aff6b3f7392856126d5e201c03 d4c9eb1be8091c9e4a42cb8d30a8047dc1b3ef8067a9ea08accd99df293aaf5ecd02c3de0f076eadafdaf77290a79b0f3ad0439949b490172f13ea9d9df2d928177ed58c5f8fbda9728f2326d37b9e9b349c88a974d1d972af48e8c1cdf524243edc3ffa3
// nonce: 1013edaec5b2b987d6a470b67a97e993e1fe67de18ac05e6acc7e6e321d9d1c9
func (s *SmartContract) createPsd(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	// omit the result of pow (<2^{\iota}) for simpilicity
	var comb []byte
	comb = append([]byte(args[2]), []byte(args[4])...)
	_ = sha3.Sum256(comb)

	time_t := strconv.FormatInt(time.Now().Unix(), 10)

	var psd = Psd{Ct: args[1], H_pow: args[2], Acc: args[3], Time: time_t}

	psdAsBytes, _ := json.Marshal(psd)
	APIstub.PutState(args[0], psdAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) queryAllPsds(APIstub shim.ChaincodeStubInterface) sc.Response {

	resultsIterator, err := APIstub.GetStateByRange("", "")
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllPsds:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) deleteAllPsds(APIstub shim.ChaincodeStubInterface) sc.Response {

	resultsIterator, err := APIstub.GetStateByRange("", "")
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		APIstub.DelState(queryResponse.Key)
	}

	fmt.Printf("- deleteAllPsds\n")

	return shim.Success(nil)
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
