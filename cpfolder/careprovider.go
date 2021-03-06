package main

import (
"errors"
"fmt"

"encoding/json"

"github.com/hyperledger/fabric/core/util"
"github.com/hyperledger/fabric/core/chaincode/shim"
)

var EVENT_COUNTER = "event_counter"
type ManageCareProvider struct {
}
var CareProviderIndexStr = "_CareProviderindex"
type Patient struct{             // Attributes of a Patient      
  PatientID string `json:"PatientID"`
  PatientName string `json:"PatientName"`
  Address   string `json:"Address"`         
  Problems string `json:"Problems"`
  Gender string `json:"Gender"`
  PatientMobile string `json:"PatientMobile"`
  Medications string `json:"Medications"`
  Remarks string `json: "Remarks"`
  PatientEmail string `json: "PatientEmail"`
  User string `json: "User"`
  IStatus string `json: "IStatus"`
  }

  func main() {     
  err := shim.Start(new(ManageCareProvider))
  if err != nil {
    fmt.Printf("Error starting ManageCareProvider chaincode: %s", err)
  }
  }
func (t *ManageCareProvider) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
  var msg string
  var err error
  if len(args) != 1 {
    return nil, errors.New("Incorrect number of arguments. Expecting 1")
  }
  // Initialize the chaincode
  msg = args[0]
  fmt.Println("ManageCareProvider chaincode is deployed successfully.");
  
  // Write the state to the ledger
  err = stub.PutState("abc", []byte(msg))       //making a test var "abc", I find it handy to read/write to it right away to test the network
  if err != nil {
    return nil, err
  }
  
  var empty []string
  jsonAsBytes, _ := json.Marshal(empty)               //marshal an emtpy array of strings to clear the index
  err = stub.PutState(CareProviderIndexStr, jsonAsBytes)
  if err != nil {
    return nil, err
  }
  err = stub.PutState(EVENT_COUNTER, []byte("1"))
  if err != nil {
    return nil, err
  }
  return nil, nil
}

func (t *ManageCareProvider) Run(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
    fmt.Println("run is running " + function)
    return t.Invoke(stub, function, args)
  }

  func (t *ManageCareProvider) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
    fmt.Println("invoke is running " + function)

  // Handle different functions
  if function == "init" {                         //initialize the chaincode state, used as reset
    return t.Init(stub, "init", args)
  } else if function == "cupdate_patient" {
    return t.cupdate_patient(stub, args)
  } else if function == "share_patient" {
    return t.share_patient(stub,args)
  }/* else if function == "update_istatus" {
    return t.update_istatus(stub, args)
  }*/

   fmt.Println("invoke did not find func: " + function)          //error
  
  return nil, errors.New("Received unknown function invocation")
}

func (t *ManageCareProvider) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
  fmt.Println("query is running " + function)

  // Handle different functions
  if function == "getPatient_byID" {                         //Read a Patient by transId
    return t.getPatient_byID(stub, args)
  } else if function == "get_byCareProviderID" {
    return t.get_byCareProviderID(stub,args)
  } 
  fmt.Println("query did not find func: " + function)           //error
  return nil, errors.New("Received unknown function query")
}

func (t *ManageCareProvider) getPatient_byID(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3 args")
	}
	PatientChaincode := args[0]
	PatientID := args[1]
	f1 := "getPatient_byID"
	queryArgs1 := util.ToChaincodeArgs(f1, PatientID)
	patientAsBytes, err := stub.QueryChaincode(PatientChaincode, queryArgs1)
	if err != nil {
		errStr := fmt.Sprintf("Failed to query chaincode. Got error: %s", err.Error())
		fmt.Printf(errStr)
		return nil, errors.New(errStr)
	}
	res := Patient{}
	json.Unmarshal(patientAsBytes, &res)
	fmt.Println(res)
	if res.PatientID == PatientID {
		fmt.Println("Patient found with PatientID : " + PatientID)
	} else {
		return nil, errors.New("PatientID not found")
	}
	return nil,nil
}
func (t *ManageCareProvider) get_byCareProviderID(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3 args")
	}
	PatientChaincode := args[0]
	CareProviderID := args[1]
	f1 := "get_byCareProviderID"
	queryArgs1 := util.ToChaincodeArgs(f1, CareProviderID)
	patientAsBytes, err := stub.QueryChaincode(PatientChaincode, queryArgs1)
	if err != nil {
		errStr := fmt.Sprintf("Failed to query chaincode. Got error: %s", err.Error())
		fmt.Printf(errStr)
		return nil, errors.New(errStr)
	}
	return patientAsBytes,nil
}
func (t *ManageCareProvider) cupdate_patient(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 5 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3 args")
	}
	PatientChaincode := args[0]
    PatientID := args[1]
    Medications := args[2]
    Remarks := args[3]
    User := args[4]
	f1 := "cupdate_patient"
	queryArgs1 := util.ToChaincodeArgs(f1, PatientID,Medications,Remarks,User)
	_, err := stub.InvokeChaincode(PatientChaincode, queryArgs1)
	if err != nil {
		errStr := fmt.Sprintf("Failed to query chaincode. Got error: %s", err.Error())
		fmt.Printf(errStr)
		return nil, errors.New(errStr)
	}
	
	return nil,nil
}

func (t *ManageCareProvider) share_patient(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
 fmt.Println("enter share function")
 if len(args) != 3 {
    return nil, errors.New("Incorrect number of arguments. Expecting 1")
  }
  PatientChaincode := args[0]
  PatientID := args[1]
  DoctorID := args[2]
  f1 := "share_patient"
 queryArgs1 := util.ToChaincodeArgs(f1, PatientID,DoctorID)
 // queryArgs1 := util.ToChaincodeArgs(f1, PatientID)
_, err := stub.InvokeChaincode(PatientChaincode, queryArgs1)
  if err != nil {
    errStr := fmt.Sprintf("Failed to query chaincode. Got error: %s", err.Error())
    fmt.Printf(errStr)
    return nil, errors.New(errStr)
  }
	
  return nil, nil
}
