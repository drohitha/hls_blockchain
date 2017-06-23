package main

import (
"errors"
"fmt"
"strconv"
"encoding/json"

"github.com/hyperledger/fabric/core/chaincode/shim"
)

var EVENT_COUNTER = "event_counter"
type ManageInsuranceProvider struct {
}
type InsuranceProvider struct{             // Attributes of a Patient      
  InsuranceProviderID string `json:"InsuranceProviderID"`
}

var InsuranceProviderIndexStr = "_InsuranceProviderindex"

type Patient struct{             // Attributes of a Patient      
  PatientID string `json:"PatientID"`
  PatientName string `json:"PatientName"`
  Address   string `json:"Address"`         
  Problems string `json:"Problems"`
  Gender string `json:"Gender"`
  PatientMobile string `json:"PatientMobile"`
  Remarks string `json: "Remarks"`
}

func main() {     
  err := shim.Start(new(ManageInsuranceProvider))
  if err != nil {
    fmt.Printf("Error starting ManageInsuranceProvider chaincode: %s", err)
  }
}

func (t *ManageInsuranceProvider) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
  var msg string
  var err error
  if len(args) != 1 {
    return nil, errors.New("Incorrect number of arguments. Expecting 1")
  }
  // Initialize the chaincode
  msg = args[0]
  fmt.Println("ManageInsuranceProvider chaincode is deployed successfully.");
  
  // Write the state to the ledger
  err = stub.PutState("abc", []byte(msg))     
  if err != nil {
    return nil, err
  }
  
  var empty []string
  jsonAsBytes, _ := json.Marshal(empty)               //marshal an emtpy array of strings to clear the index
  err = stub.PutState(InsuranceProviderIndexStr, jsonAsBytes)
  if err != nil {
    return nil, err
  }
  err = stub.PutState(EVENT_COUNTER, []byte("1"))
  if err != nil {
    return nil, err
  }
  return nil, nil
}

func (t *ManageInsuranceProvider) Run(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
    fmt.Println("run is running " + function)
    return t.Invoke(stub, function, args)
}

func (t *ManageInsuranceProvider) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
    fmt.Println("invoke is running " + function)

  // Handle different functions
  if function == "init" {                         //initialize the chaincode state, used as reset
    return t.Init(stub, "init", args)
  }   
   fmt.Println("invoke did not find func: " + function)          //error
  
  return nil, errors.New("Received unknown function invocation")
}

func (t *ManageInsuranceProvider) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
  fmt.Println("query is running " + function)

  // Handle different functions
  if function == "get_byPatientID" {                         //Read a Patient by transId
    return t.get_byPatientID(stub, args)
  } else if function == "get_byInsuranceProviderID" {
    return t.get_byInsuranceProviderID(stub,args)
  } 
  fmt.Println("query did not find func: " + function)           //error
  return nil, errors.New("Received unknown function query")
}



func (t *ManageInsuranceProvider) get_byPatientID(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){
	 var PatientID, jsonResp string
  var err error
  fmt.Println("start get_byPatientID")
  if len(args) != 1 {
    return nil, errors.New("Incorrect number of arguments. Expecting ID of the patient to query")
  }
  // set PatientID
  PatientID = args[0]
  valAsbytes, err := stub.GetState(PatientID)                  //get the PatientID from chaincode state
  if err != nil {
    jsonResp = "{\"Error\":\"Failed to get state for " + PatientID + "\"}"
    return nil, errors.New(jsonResp)
  }
  //fmt.Print("valAsbytes : ")
  //fmt.Println(valAsbytes)
  fmt.Println("end get_byPatientID")
  return valAsbytes, nil  
}

func (t *ManageInsuranceProvider) get_byInsuranceProviderID(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){
var InsuranceProviderIndex []string
var InsuranceProviderID, jsonResp, errResp string
  var err error
fmt.Println("start get_byInsuranceProviderID")
  if len(args) != 1 {
    return nil, errors.New("Incorrect number of arguments. Expecting ID of the InsuranceProvider to query")
  }
 InsuranceProviderID = args[0]
 valAsbytes, err := stub.GetState(InsuranceProviderID)
 if err != nil {
    jsonResp = "{\"Error\":\"Failed to get state for " + InsuranceProviderID + "\"}"
    return nil, errors.New(jsonResp)
  }
 json.Unmarshal(valAsbytes, &InsuranceProviderIndex)
 jsonResp = "{"
	for i,val := range InsuranceProviderIndex{
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for get_byInsuranceProviderID")
		valueAsBytes, err := stub.GetState(val)
		if err != nil {
			errResp = "{\"Error\":\"Failed to get state for " + val + "\"}"
			return nil, errors.New(errResp)
		}

			jsonResp = jsonResp + "\""+ val + "\":" + string(valueAsBytes[:])
			//fmt.Println("jsonResp inside if")
			//fmt.Println(jsonResp)
			if i < len(InsuranceProviderIndex)-1 {
				jsonResp = jsonResp + ","
			}
		}
		
	jsonResp = jsonResp + "}"
	fmt.Println("end get_byInsuranceProviderID")
	return []byte(jsonResp), nil
}
