package main

import (
"errors"
"fmt"
"strconv"
"encoding/json"

"github.com/hyperledger/fabric/core/chaincode/shim"
)

var EVENT_COUNTER = "event_counter"
type ManageCareProvider struct {
}
type CareProvider struct{             // Attributes of a Patient      
  CareProviderID string `json:"CareProviderID"`
}

var CareProviderIndexStr = "_CareProviderindex"

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
  err = stub.PutState("abc", []byte(msg))     
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
  } else if function == "dupdate_patient" {
    return t.dupdate_patient(stub,args)
  }    
   fmt.Println("invoke did not find func: " + function)          //error
  
  return nil, errors.New("Received unknown function invocation")
}

func (t *ManageCareProvider) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
  fmt.Println("query is running " + function)

  // Handle different functions
  if function == "get_byPatientID" {                         //Read a Patient by transId
    return t.get_byPatientID(stub, args)
  } else if function == "get_byCareProviderID" {
    return t.get_byCareProviderID(stub,args)
  } 
  fmt.Println("query did not find func: " + function)           //error
  return nil, errors.New("Received unknown function query")
}

func (t *ManageCareProvider) dupdate_patient(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
  var jsonResp string
  var err error
  fmt.Println("start dupdate_patient")

  PatientID := args[0]
  Remarks := args[1]


  if len(args) != 2 {
    return nil, errors.New("Incorrect number of arguments. Expecting 6.")
  }
  // set PatientID
  //PatientID := args[0]
  PatientAsBytes, err := stub.GetState(PatientID)                 //get the patient for the specified patientID from chaincode state
  if err != nil {
    jsonResp = "{\"Error\":\"Failed to get state for " + PatientID + "\"}"
    return nil, errors.New(jsonResp)
  }
  //fmt.Print("vesselAsBytes in update vessel")
  //fmt.Println(vesselAsBytes);
   res := Patient{}
  json.Unmarshal( PatientAsBytes, &res)
  if res.PatientID == PatientID{
    fmt.Println("Patient found with PatientID : " + PatientID)
    //fmt.Println(res);
    
    res.Remarks = args[1]
    }
	Address := res.Address
	Problems := res.Problems
	PatientName := res.PatientName
	Gender := res.Gender
	PatientMobile := res.PatientMobile
  
  //build the CreatePatient json string manually
  PatientDetails :=  `{`+
    `"PatientID": "` + PatientID + `" , `+
    `"Address": "` + Address + `" , `+
    `"Problems": "` + Problems + `" , `+
    `"PatientName": "` + PatientName + `" , `+
    `"Gender": "` + Gender + `" , `+ 
    `"PatientMobile": "` + PatientMobile + `" , `+ 
    `"Remarks": "` + Remarks + `" `+
    `}`
  err = stub.PutState(PatientID, []byte(PatientDetails))                  //store patient with id as key
  if err != nil {
    return nil, err
  }
  return nil, nil
}

func (t *ManageCareProvider) get_byPatientID(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){
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

func (t *ManageCareProvider) get_byCareProviderID(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){
var CareProviderIndex []string
var CareProviderID, jsonResp, errResp string
  var err error
fmt.Println("start get_byCareProviderID")
  if len(args) != 1 {
    return nil, errors.New("Incorrect number of arguments. Expecting ID of the CareProvider to query")
  }
 CareProviderID = args[0]
 valAsbytes, err := stub.GetState(CareProviderID)
 if err != nil {
    jsonResp = "{\"Error\":\"Failed to get state for " + CareProviderID + "\"}"
    return nil, errors.New(jsonResp)
  }
 json.Unmarshal(valAsbytes, &CareProviderIndex)
 jsonResp = "{"
	for i,val := range CareProviderIndex{
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for get_byCareProviderID")
		valueAsBytes, err := stub.GetState(val)
		if err != nil {
			errResp = "{\"Error\":\"Failed to get state for " + val + "\"}"
			return nil, errors.New(errResp)
		}

			jsonResp = jsonResp + "\""+ val + "\":" + string(valueAsBytes[:])
			//fmt.Println("jsonResp inside if")
			//fmt.Println(jsonResp)
			if i < len(CareProviderIndex)-1 {
				jsonResp = jsonResp + ","
			}
		}
		
	jsonResp = jsonResp + "}"
	fmt.Println("end get_byCareProviderID")
	return []byte(jsonResp), nil
}
