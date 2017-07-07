package main

import (
"errors"
"fmt"
"strconv"
"encoding/json"
"strings"
"github.com/hyperledger/fabric/core/chaincode/shim"
"github.com/hyperledger/fabric/protos/peer"
) 

var EVENT_COUNTER = "event_counter"

type ManagePatient struct {
}

var PatientIndexStr = "_Patientindex" 
var DoctorIndexStr = "_Doctorindex"
var CareProviderIndexStr = "_CareProviderindex"
var InsuranceProviderStr = "_InsuranceProviderindex"      //name for the key/value that will store a list of all known Patients

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
  err := shim.Start(new(ManagePatient))
  if err != nil {
    fmt.Printf("Error starting ManagePatient chaincode: %s", err)
  }
}

func (t *ManagePatient) Init(stub shim.ChaincodeStubInterface) peer.Response {

	args := stub.GetStringArgs()
  var msg string
  var err error
  if len(args) != 1 {
    return shim.Error("Incorrect number of arguments. Expecting 1")
  }
  // Initialize the chaincode
  msg = args[0]
  fmt.Println("ManagePatient chaincode is deployed successfully.");
  
  // Write the state to the ledger
  err = stub.PutState("abc", []byte(msg))       //making a test var "abc", I find it handy to read/write to it right away to test the network
  if err != nil {
    return shim.Error(fmt.Sprintf("Failed to create asset: %s", args[0]))
  }
  
  var empty []string
  jsonAsBytes, _ := json.Marshal(empty)               //marshal an emtpy array of strings to clear the index
  err = stub.PutState(PatientIndexStr, jsonAsBytes)
  if err != nil {
    return shim.Error(fmt.Sprintf("Failed to create asset in patientindex: %s", args[0]))
  }
  err = stub.PutState(EVENT_COUNTER, []byte("1"))
  if err != nil {
    return shim.Error(fmt.Sprintf("Failed to create asset in event counter: %s", args[0]))
  }
  return shim.Success(nil)
}

  func (t *ManagePatient) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
    //fmt.Println("invoke is running " + function)
       fn, args := stub.GetFunctionAndParameters()
  // Handle different functions
       var result string
    var err error
  if fn == "update_istatus" {                         //initialize the chaincode state, used as reset
    result, err = update_istatus(stub, args)
  } else if fn == "create_patient"{
    result, err = create_patient(stub,args)
  } else if fn == "delete" {                     //delete a new Patient
    result, err = delete(stub, args)
  } else if fn == "update_patient" {
    result, err = update_patient(stub,args)
  } else if fn == "share_patient" {
    result, err = share_patient(stub,args)
  } else if fn == "dupdate_patient" {
    result, err = dupdate_patient(stub, args)
  } else if fn == "cupdate_patient" {
    result, err = cupdate_patient(stub, args)
  } else if fn == "getPatient_byID" {                         //Read a Patient by transId
    result, err = getPatient_byID(stub, args)
  } else if fn == "getPatient_byEmail" {
    result, err = getPatient_byEmail(stub,args)
  } else if fn == "get_byDoctorID" {
    result, err = get_byDoctorID(stub,args)
  } else if fn == "get_byCareProviderID" {
    result, err = get_byCareProviderID(stub,args)
  } else if fn == "get_byInsuranceProviderID" {
    result, err = get_byInsuranceProviderID(stub,args)
  }
  if err != nil {
            return shim.Error(err.Error())
    }
   fmt.Println("invoke did not find func: " + fn)          //error
  
  return shim.Success([]byte(result))
}

func getPatient_byID(stub shim.ChaincodeStubInterface, args []string) (string, error) {
  var PatientID, jsonResp string
  var err error
  fmt.Println("start getPatient_byID")
  if len(args) != 1 {
    return "", fmt.Errorf("Incorrect arguments. Expecting a key")
  }
  // set PatientID
  PatientID = args[0]
  valAsbytes, err := stub.GetState(PatientID)                  //get the PatientID from chaincode state
  if err != nil {
       return "", fmt.Errorf("Failed to get asset: %s with error: %s", args[0], err)
  }
  if value == nil {
            return "", fmt.Errorf("Asset not found: %s", args[0])
    }
    return string(value), nil
  //fmt.Print("valAsbytes : ")
  //fmt.Println(valAsbytes)                        //send it onward
}

func getPatient_byEmail(stub shim.ChaincodeStubInterface, args []string) (string, error) {
  var PatientEmail, jsonResp, errResp string
  var err error
  var valIndex Patient
  fmt.Println("start getPatient_byEmail")
  if len(args) != 1 {
    return "", fmt.Errorf("Incorrect number of arguments. Expecting ID of the patient to query")
  }
  // set PatientID
   PatientEmail= args[0]
  PatientAsBytes, err := stub.GetState(PatientIndexStr)                  //get the PatientID from chaincode state
  if err != nil {
    jsonResp = "{\"Error\":\"Failed to get state for " + PatientEmail + "\"}"
    return "",fmt.Errorf(jsonResp)
  }
 
    var PatientIndex []string
  json.Unmarshal(PatientAsBytes, &PatientIndex) 

  jsonResp = "{"
  for i,val := range PatientIndex{

    fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for getPatient_byID")
    valueAsBytes, err := stub.GetState(val)
    if err != nil {
      errResp = "{\"Error\":\"Failed to get state for " + val + "\"}"
      return "", fmt.Errorf(errResp)
    }
  
    var err1 error
    err1 = json.Unmarshal(valueAsBytes, &valIndex)
    if err1 != nil {
      fmt.Println(err1)
  }
      
    if valIndex.PatientEmail == PatientEmail{
      fmt.Println("Patientfound")
      jsonResp = jsonResp + "\""+ val + "\":" + string(valueAsBytes[:])
      if i < len(PatientIndex)-1 {
        jsonResp = jsonResp + ","
      }}}
  //fmt.Print("valAsbytes : ")
  //fmt.Println(valAsbytes)
  jsonResp = jsonResp + "}"
  //fmt.Println("jsonResp : " + jsonResp)
  //fmt.Print("jsonResp in bytes : ")
  //fmt.Println([]byte(jsonResp))
  fmt.Println("end getby patientemail")
  return string(jsonResp), nil 
}

func create_patient(stub shim.ChaincodeStubInterface, args []string) (string, error) {
  var err error

  if len(args) != 11{
    return "", fmt.Errorf("Incorrect number of arguments. Expecting 9")
  }
  fmt.Println("start create_Patient OK")

  PatientID := args[0]
  Address := args[2]
  Problems := args[3]
  PatientName:= args[1]
  Gender := args[4]
  PatientMobile := args[5]
  Medications := args[6]
  Remarks := args[7]
  PatientEmail := args[8]
  User := args[9]
  IStatus := args[10]
  
  //fmt.Println("start create_Patient 1")
  PatientAsBytes, err := stub.GetState(PatientID)
  if err != nil {
    //fmt.Println("start create_Patient 2")
    return "", fmt.Errorf("Failed to get Patient Patient_id")
  }
  var res Patient
  //fmt.Println("start create_Patient 3")
  json.Unmarshal(PatientAsBytes, &res)
  
  fmt.Println(res.PatientID)
  if res.PatientID == PatientID{
  fmt.Println("This patient already exist")
  return "", fmt.Errorf("This Patient arleady exists")       
                                                           //all stop a patient by this name exists
  }
     fmt.Println("start create_Patient 4")
     //build the CreatePatient json string manually
      PatientDetails :=  `{`+
    `"PatientID": "` + PatientID + `" , `+
    `"Address": "` + Address + `" , `+
    `"Problems": "` + Problems + `" , `+
    `"PatientName": "` + PatientName + `" , `+
    `"Gender": "` + Gender + `" , `+ 
    `"PatientMobile": "` + PatientMobile + `" , `+ 
    `"Medications": "` + Medications + `" , `+ 
    `"Remarks": "` + Remarks + `" , `+ 
    `"PatientEmail": "` + PatientEmail + `" , `+
    `"User": "` + User + `" , `+
    `"IStatus": "` + IStatus + `" `+
    `}`



    fmt.Print("Patient details in array: ")
    fmt.Println(PatientDetails)
    err = stub.PutState(PatientID, []byte(PatientDetails))                  //store Patient with PatientID as key
    if err != nil {
    return "", fmt.Errorf(err)
  }
  //get the patient
  PatientIndexAsBytes, err := stub.GetState(PatientIndexStr)
  if err != nil {
    return "", fmt.Errorf("Failed to get Patient index")
  }
  var PatientIndex []string
  //fmt.Print("PatientIndexAsBytes: ")
  //fmt.Println(PatientIndexAsBytes)
  
  json.Unmarshal(PatientIndexAsBytes, &PatientIndex)              //un stringify it aka JSON.parse()
    
  PatientIndex = append(PatientIndex, PatientID)                 //add Patient transID to index list
  //fmt.Println("! Patient index after appending transId: ", poIndex)
  jsonAsBytes, _ := json.Marshal(PatientIndex)
  //fmt.Print("jsonAsBytes: ")
  //fmt.Println(jsonAsBytes)
  err = stub.PutState(PatientIndexStr, jsonAsBytes)            //store name of Patient
  if err != nil {
    return "", fmt.Errorf(err)
  }

  fmt.Println("end create_Patient")
  return string(PatientDetails), nil
}

func delete(stub shim.ChaincodeStubInterface, args []string) (string, error) {
  fmt.Println("enter delete function")
  if len(args) != 1 {
    return "", fmt.Errorf("Incorrect number of arguments. Expecting 1")
  }
  // set PatientID
 PatientID := args[0]
 err := stub.DelState(PatientID)
// fmt.Println(1)
  //get the patient index
 PatientIndexAsBytes, err := stub.GetState(PatientIndexStr)
// fmt.Println(2)
 if err != nil {
 return "", fmt.Errorf("Failed to get Patient index")
  }
 // fmt.Println(3)

  //fmt.Println("poAsBytes in delete po")
  //fmt.Println(poAsBytes);
  var PatientIndex []string
//  fmt.Println(4)
  json.Unmarshal(PatientIndexAsBytes, &PatientIndex)               //un stringify it aka JSON.parse()
  //fmt.Println("poIndex in delete po")
  //fmt.Println(poIndex);
  //remove marble from index
 // fmt.Println(5)
  for i,val := range PatientIndex{
  fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for " + PatientID)
  if val == PatientID{                             //find the correct patient
      fmt.Println("found patient with matching patientID")
      PatientIndex = append(PatientIndex[:i], PatientIndex[i+1:]...)     //remove it
     // fmt.Println(6)
      for x:= range PatientIndex{                      //debug prints...
        fmt.Println(string(x) + " - " + PatientIndex[x])
      }
      break
    }
  }
  //fmt.Println(6)
  jsonAsBytes, _ := json.Marshal(PatientIndex)                 //save new index
  err = stub.PutState(PatientIndexStr, jsonAsBytes)
  return string(PatientIndexStr), nil
}

func update_patient(stub shim.ChaincodeStubInterface, args []string) (string, error) {
  var jsonResp string
  var err error
  fmt.Println("start update_patient")

  PatientID := args[0]
  Address := args[1]
  Problems := args[2]
  PatientName:= args[3]
  Gender := args[4]
  PatientMobile := args[5]
  Medications := args[6]
Remarks := args[7]
PatientEmail := args[8]
User := args[9]

  if len(args) != 10 {
    return "", fmt.Errorf("Incorrect number of arguments. Expecting 9.")
  }
  // set PatientID
  //PatientID := args[0]
  PatientAsBytes, err := stub.GetState(PatientID)                 //get the patient for the specified patientID from chaincode state
  if err != nil {
    jsonResp = "{\"Error\":\"Failed to get state for " + PatientID + "\"}"
    return "", fmt.Errorf(jsonResp)
  }
  //fmt.Print("vesselAsBytes in update vessel")
  //fmt.Println(vesselAsBytes);
   res := Patient{}
  json.Unmarshal( PatientAsBytes, &res)
  if res.PatientID == PatientID{
    fmt.Println("Patient found with PatientID : " + PatientID)
    //fmt.Println(res);
    res.Address = args[1]
    res.Problems  = args[2]
    res.PatientName = args[3]
    res.Gender = args[4]
    res.PatientMobile = args[5]
    res.Medications = args[6]
     res.Remarks = args[7]
    res.User = args[9]
    }
     IStatus := res.IStatus
  
  //build the CreatePatient json string manually
  PatientDetails :=  `{`+
    `"PatientID": "` + PatientID + `" , `+
    `"Address": "` + Address + `" , `+
    `"Problems": "` + Problems + `" , `+
    `"PatientName": "` + PatientName + `" , `+
    `"Gender": "` + Gender + `" , `+ 
    `"PatientMobile": "` + PatientMobile + `" , `+ 
    `"Medications": "` + Medications + `" , `+ 
    `"Remarks": "` + Remarks + `" , `+ 
    `"PatientEmail": "` + PatientEmail + `" , `+
    `"User": "` + User + `" , `+
    `"IStatus": "` + IStatus + `" `+
    `}`
  err = stub.PutState(PatientID, []byte(PatientDetails))                  //store patient with id as key
  if err != nil {
    return "", fmt.Errof(err)
  }
  return string(PatientDetails), nil
}

func share_patient(stub shim.ChaincodeStubInterface, args []string) (string, error) {
 fmt.Println("enter share function")
 if len(args) != 2 {
    return "", fmt.Errorf("Incorrect number of arguments. Expecting 1")
  }
  PatientID := args[0]
  DoctorID := args[1]
  s := strings.HasPrefix(DoctorID,"ip")
  if s == true {
    /*f1 := "update_istatus"
  invokeArgs2 := util.ToChaincodeArgs(f1, PatientID, "Claimed")
  result2, err := stub.InvokeChaincode(Pchain, invokeArgs2)
  if err != nil {
    errStr := fmt.Sprintf("Failed to update Transaction status from 'PatientID' chaincode. Got error: %s", err.Error())
    fmt.Printf(errStr)
    return nil, errors.New(errStr)
  }
  fmt.Print("Transaction hash returned: ")
  fmt.Println(result2)
  fmt.Println("Successfully updated istatus to 'Claimed'")*/
  var jsonResp string
  var err error
  fmt.Println("start update_istatus")
  /*if len(args) != 2 {
    return nil, errors.New("Incorrect number of arguments. Expecting 3.")
  }*/
  // set vesselID
  //PatientID := args[0]
  PatientAsBytes, err := stub.GetState(PatientID)                  //get the Berth for the specified vesselID from chaincode state
  if err != nil {
    jsonResp = "{\"Error\":\"Failed to get state for " + PatientID + "\"}"
    return "", fmt.Errorf(jsonResp)
  }
  //fmt.Print("berthAsBytes in update berth")
  //fmt.Println(berthAsBytes);
  res := Patient{}
  json.Unmarshal(PatientAsBytes, &res)
  if res.PatientID == PatientID{
    fmt.Println("Patient found with PatientID : " + PatientID)
    if res.IStatus == "Claim Initiated"{
      return "", fmt.Errorf("Insurance already shared and claimed")
    } else if res.IStatus == "Approved"{
      return "", fmt.Errorf("Insurance already approved")
    } else if res.IStatus == "Rejected"{
      return "", fmt.Errorf("Insurance already rejected")
    }
    res.IStatus = "Claim Initiated"
  
  }
  Address := res.Address
  Problems := res.Problems
  PatientName:= res.PatientName
  Gender := res.Gender
  PatientMobile := res.PatientMobile
  PatientEmail := res.PatientEmail
  Medications := res.Medications
  Remarks := res.Remarks
  User := res.User
  IStatus := res.IStatus
  //build the Berth json string manually
  PatientDetails :=  `{`+
    `"PatientID": "` + PatientID + `" , `+
    `"Address": "` + Address + `" , `+
    `"Problems": "` + Problems + `" , `+
    `"PatientName": "` + PatientName + `" , `+
    `"Gender": "` + Gender + `" , `+ 
    `"PatientMobile": "` + PatientMobile + `" , `+ 
    `"Medications": "` + Medications + `" , `+ 
    `"Remarks": "` + Remarks + `" , `+ 
    `"PatientEmail": "` + PatientEmail + `" , `+
    `"User": "` + User + `" , `+
    `"IStatus": "` + IStatus + `" `+
    `}`
  err = stub.PutState(PatientID, []byte(PatientDetails))                 //store Berth with id as key
  if err != nil {
    return "", fmt.Errof(err)
  }
  }
  /*PatientDetails :=  `{`+
    `"PatientID": "` + PatientID + `" , `+
    `"DoctorID": "` + DoctorID + `" , `+
    `}`*/
  DoctorAsBytes, err := stub.GetState(DoctorID)
  if err != nil {
    return "", fmt.Errorf("Failed to get Doctor index")
  }
  var DoctorIndex []string
  json.Unmarshal(DoctorAsBytes, &DoctorIndex)
   DoctorIndex = append(DoctorIndex, PatientID)
   jsonAsBytes, _ := json.Marshal(DoctorIndex)
  err = stub.PutState(DoctorID, jsonAsBytes)            //store name of Patient
  if err != nil {
    return "", fmt.Errorf(err)
  }
  return string(DoctorID), nil
}

func get_byDoctorID(stub shim.ChaincodeStubInterface, args []string) (string, error){
  var doctorIndex []string
var DoctorID, jsonResp, errResp string
  var err error
fmt.Println("start get_byDoctorID")
  if len(args) != 1 {
    return "", fmt.Errorf("Incorrect number of arguments. Expecting ID of the doctor to query")
  }
 DoctorID = args[0]
 valAsbytes, err := stub.GetState(DoctorID)
 if err != nil {
    jsonResp = "{\"Error\":\"Failed to get state for " + DoctorID + "\"}"
    return "", fmt.Errorf(jsonResp)
  }
 json.Unmarshal(valAsbytes, &doctorIndex)
 jsonResp = "{"
  for i,val := range doctorIndex{
    fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for get_byDoctorID")
    valueAsBytes, err := stub.GetState(val)
    if err != nil {
      errResp = "{\"Error\":\"Failed to get state for " + val + "\"}"
      return "", fmt.Errorf(errResp)
    }

      jsonResp = jsonResp + "\""+ val + "\":" + string(valueAsBytes[:])
      //fmt.Println("jsonResp inside if")
      //fmt.Println(jsonResp)
      if i < len(doctorIndex)-1 {
        jsonResp = jsonResp + ","
      }
    }
    
  jsonResp = jsonResp + "}"
  fmt.Println("end get_byDoctorID")
  return string(jsonResp), nil
}

func get_byCareProviderID(stub shim.ChaincodeStubInterface, args []string) (string, error){
  var CareProviderIndex []string
var CareProviderID, jsonResp, errResp string
  var err error
fmt.Println("start get_byCareProviderID")
  if len(args) != 1 {
    return "", fmt.Errorf("Incorrect number of arguments. Expecting ID of the CareProvider to query")
  }
 CareProviderID = args[0]
 valAsbytes, err := stub.GetState(CareProviderID)
 if err != nil {
    jsonResp = "{\"Error\":\"Failed to get state for " + CareProviderID + "\"}"
    return "", fmt.Errorf(jsonResp)
  }
 json.Unmarshal(valAsbytes, &CareProviderIndex)
 jsonResp = "{"
  for i,val := range CareProviderIndex{
    fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for get_byCareProviderID")
    valueAsBytes, err := stub.GetState(val)
    if err != nil {
      errResp = "{\"Error\":\"Failed to get state for " + val + "\"}"
      return "", fmt.Errorf(errResp)
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
  return string(jsonResp), nil
}

func get_byInsuranceProviderID(stub shim.ChaincodeStubInterface, args []string) (string, error){
var InsuranceProviderIndex []string
var InsuranceProviderID, jsonResp, errResp string
  var err error
fmt.Println("start get_byInsuranceProviderID")
  if len(args) != 1 {
    return "", fmt.Errorf("Incorrect number of arguments. Expecting ID of the InsuranceProvider to query")
  }
 InsuranceProviderID = args[0]
 valAsbytes, err := stub.GetState(InsuranceProviderID)
 if err != nil {
    jsonResp = "{\"Error\":\"Failed to get state for " + InsuranceProviderID + "\"}"
    return "", fmt.Errorf(jsonResp)
  }
 json.Unmarshal(valAsbytes, &InsuranceProviderIndex)
 jsonResp = "{"
  for i,val := range InsuranceProviderIndex{
    fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for get_byInsuranceProviderID")
    valueAsBytes, err := stub.GetState(val)
    if err != nil {
      errResp = "{\"Error\":\"Failed to get state for " + val + "\"}"
      return "", fmt.Errorf(errResp)
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
  return string(jsonResp), nil
}

func dupdate_patient(stub shim.ChaincodeStubInterface, args []string) (string, error) {
  var jsonResp string
  var err error
  fmt.Println("start dupdate_patient")

  PatientID := args[0]
  Medications := args[1]
  Remarks := args[2]
  User := args[3]

  if len(args) != 4 {
    return "", fmt.Errorf("Incorrect number of arguments. Expecting 3.")
  }
  // set PatientID
  //PatientID := args[0]
  PatientAsBytes, err := stub.GetState(PatientID)                 //get the patient for the specified patientID from chaincode state
  if err != nil {
    jsonResp = "{\"Error\":\"Failed to get state for " + PatientID + "\"}"
    return "", fmt.Errorf(jsonResp)
  }
  //fmt.Print("vesselAsBytes in update vessel")
  //fmt.Println(vesselAsBytes);
   res := Patient{}
  json.Unmarshal( PatientAsBytes, &res)
  if res.PatientID == PatientID{
    fmt.Println("Patient found with PatientID : " + PatientID)
    //fmt.Println(res);
    res.Medications = args[1]
    res.Remarks = args[2]
    res.User = args[3]
    }
  

  Address := res.Address
  Problems := res.Problems
  PatientName:= res.PatientName
  Gender := res.Gender
  PatientMobile := res.PatientMobile
  PatientEmail := res.PatientEmail
  IStatus := res.IStatus
  
  //build the CreatePatient json string manually
  PatientDetails :=  `{`+
    `"PatientID": "` + PatientID + `" , `+
    `"Address": "` + Address + `" , `+
    `"Problems": "` + Problems + `" , `+
    `"PatientName": "` + PatientName + `" , `+
    `"Gender": "` + Gender + `" , `+ 
    `"PatientMobile": "` + PatientMobile + `" , `+ 
    `"Medications": "` + Medications + `" , `+ 
    `"Remarks": "` + Remarks + `" , `+ 
    `"PatientEmail": "` + PatientEmail + `" , `+
    `"User": "` + User + `" , `+
    `"IStatus": "` + IStatus + `" `+
    `}`
  err = stub.PutState(PatientID, []byte(PatientDetails))                  //store patient with id as key
  if err != nil {
    return "", fmt.Errorf(err)
  }
  return string(PatientDetails), nil
}

func cupdate_patient(stub shim.ChaincodeStubInterface, args []string) (string, error) {
  var jsonResp string
  var err error
  fmt.Println("start cupdate_patient")

  PatientID := args[0]
  Medications := args[1]
  Remarks := args[2]
  User := args[3]

  if len(args) != 4 {
    return "", fmt.Errorf("Incorrect number of arguments. Expecting 3.")
  }
  // set PatientID
  //PatientID := args[0]
  PatientAsBytes, err := stub.GetState(PatientID)                 //get the patient for the specified patientID from chaincode state
  if err != nil {
    jsonResp = "{\"Error\":\"Failed to get state for " + PatientID + "\"}"
    return "", fmt.Errorf(jsonResp)
  }
  //fmt.Print("vesselAsBytes in update vessel")
  //fmt.Println(vesselAsBytes);
   res := Patient{}
  json.Unmarshal( PatientAsBytes, &res)
  if res.PatientID == PatientID{
    fmt.Println("Patient found with PatientID : " + PatientID)
    //fmt.Println(res);
    res.Medications = args[1]
    res.Remarks = args[2]
    res.User = args[3]
    }
  

  Address := res.Address
  Problems := res.Problems
  PatientName:= res.PatientName
  Gender := res.Gender
  PatientMobile := res.PatientMobile
  PatientEmail := res.PatientEmail
  IStatus := res.IStatus
  
  //build the CreatePatient json string manually
  PatientDetails :=  `{`+
    `"PatientID": "` + PatientID + `" , `+
    `"Address": "` + Address + `" , `+
    `"Problems": "` + Problems + `" , `+
    `"PatientName": "` + PatientName + `" , `+
    `"Gender": "` + Gender + `" , `+ 
    `"PatientMobile": "` + PatientMobile + `" , `+ 
    `"Medications": "` + Medications + `" , `+ 
    `"Remarks": "` + Remarks + `" , `+ 
    `"PatientEmail": "` + PatientEmail + `" , `+
    `"User": "` + User + `" , `+
    `"IStatus": "` + IStatus + `" `+
    `}`
  err = stub.PutState(PatientID, []byte(PatientDetails))                  //store patient with id as key
  if err != nil {
    return "", fmt.Errorf(err)
  }
  return string(PatientDetails), nil
}

func update_istatus(stub shim.ChaincodeStubInterface, args []string) (string, error) {
  var jsonResp string
  var err error
  fmt.Println("start update_istatus")
  if len(args) != 2 {
    return "", fmt.Errorf("Incorrect number of arguments. Expecting 3.")
  }
  // set vesselID
  PatientID := args[0]
  IStatus := args[1]
  PatientAsBytes, err := stub.GetState(PatientID)                  //get the Berth for the specified vesselID from chaincode state
  if err != nil {
    jsonResp = "{\"Error\":\"Failed to get state for " + PatientID + "\"}"
    return "", fmt.Errorf(jsonResp)
  }
  //fmt.Print("berthAsBytes in update berth")
  //fmt.Println(berthAsBytes);
  res := Patient{}
  json.Unmarshal(PatientAsBytes, &res)
  if res.PatientID == PatientID{
    fmt.Println("Patient found with PatientID : " + PatientID)
    if res.IStatus != "Claim Initiated"{
      return "", fmt.Errorf("claim already approved or rejected once")
    }
    res.IStatus = args[1]
  
  }
  Address := res.Address
  Problems := res.Problems
  PatientName:= res.PatientName
  Gender := res.Gender
  PatientMobile := res.PatientMobile
  PatientEmail := res.PatientEmail
  Medications := res.Medications
  Remarks := res.Remarks
  User := res.User
  
  //build the Berth json string manually
  PatientDetails :=  `{`+
    `"PatientID": "` + PatientID + `" , `+
    `"Address": "` + Address + `" , `+
    `"Problems": "` + Problems + `" , `+
    `"PatientName": "` + PatientName + `" , `+
    `"Gender": "` + Gender + `" , `+ 
    `"PatientMobile": "` + PatientMobile + `" , `+ 
    `"Medications": "` + Medications + `" , `+ 
    `"Remarks": "` + Remarks + `" , `+ 
    `"PatientEmail": "` + PatientEmail + `" , `+
    `"User": "` + User + `" , `+
    `"IStatus": "` + IStatus + `" `+
    `}`
  err = stub.PutState(PatientID, []byte(PatientDetails))                 //store Berth with id as key
  if err != nil {
    return "",fmt.Errorf(err)
  }
  return string(PatientDetails), nil
}
