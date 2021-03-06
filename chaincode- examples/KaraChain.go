package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"regexp"
	"strconv"
	"strings"
)

// CHANGE! CLDChaincode to KCDChaincode
var logger = shim.NewLogger("KCDChaincode")

//==============================================================================================================================
//	 Participant types - Each participant type is mapped to an integer which we use to compare to the value stored in a
//						 user's eCert
//==============================================================================================================================
//CURRENT WORKAROUND USES ROLES CHANGE WHEN OWN USERS CAN BE CREATED SO THAT IT READ 1, 2, 3, 4, 5
const AUTHORITY = "regulator"
const SINGER = "singer"
const PRIVATE_ENTITY = "private"
const COPYRIGHT_AUTHORITY = "copyright_authority"

//==============================================================================================================================
//	 Status types - Asset lifecycle is broken down into z statuses, this is part of the business logic to determine what can
//					be done to the vehicle at points in it's lifecycle
//==============================================================================================================================
const STATE_TEMPLATE = 0
const STATE_RECORDED = 1
const STATE_UPDATED = 2
const STATE_VOTED = 3
const STATE_CONTRACT_PROVIDED = 4
const STATE_CONTRACT_REJECTED = 5
const STATE_CONTRACT_ACCEPTED = 6

//==============================================================================================================================
//	 Structure Definitions
//==============================================================================================================================
//	Chaincode - A blank struct for use with Shim (A HyperLedger included go file used for get/put state
//				and other HyperLedger functions)
//==============================================================================================================================
type SimpleChaincode struct {
}

//==============================================================================================================================
//	Song - Defines the structure for a Song object. JSON on right tells it what JSON fields to map to
//			  that element when reading a JSON object into the struct e.g. JSON make -> Struct Make.
//==============================================================================================================================
type Song struct {
	Date_Created               string `json:"Date_Created"`
	SmartContract_Unique_ID    string `json:"SmartContract_Unique_ID"`
	Singer_Id                  string `json:"Singer_Id"`
	Singer_Name                string `json:"Singer_Name"`
	Video_Id                   string `json:"Video_Id"`
	Owner                      string `json:"Owner"`
	Video_Link                 string `json:"Video_Link"`
	Video_date_created         string `json:"Video_date_created"`
	Video_QR_code_Id           string `json:"Video_QR_code_Id"`
	Copyright_Id               string `json:"Copyright_Id"`
	Copyright_date_created     string `json:"Copyright_date_created"`
	Copyright_date_accepted    string `json:"Copyright_date_accepted"`
	Copyright_date_rejected    string `json:"Copyright_date_rejected"`
	Copyright_Institution_Id   string `json:"Copyright_Institution_Id"`
	Copyright_Institution_Name string `json:"Copyright_Institution_Name"`
	Copyright_State            string `json:"Copyright_State"`
	Venue_Id                   string `json:"Venue_Id"`
	Venue_Name                 string `json:"Venue_Name"`
	User_Id                    string `json:"User_Id"`
	User_role                  string `json:"User_role"`
	User_rating                string `json:"User_role"`
}

//==============================================================================================================================
//	V5C Holder - Defines the structure that holds all the Song_IDs for songs that have been created.
//				Used as an index when querying all songs
//==============================================================================================================================

type Song_Holder struct {
	Song []string `json:"Song"`
}

//==============================================================================================================================
//	User_and_eCert - Struct for storing the JSON of a user and their ecert
//==============================================================================================================================

type User_and_eCert struct {
	Identity string `json:"identity"`
	eCert    string `json:"ecert"`
}

//==============================================================================================================================
//	Init Function - Called when the user deploys the chaincode
//==============================================================================================================================
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	//Args
	//				0
	//			peer_address

	var Song_IDs Song_Holder

	bytes, err := json.Marshal(Song_ID)

	if err != nil {
		return nil, errors.New("Error creating Song record")
	}

	err = stub.PutState("Song_IDs", bytes)

	for i := 0; i < len(args); i = i + 2 {
		t.add_ecert(stub, args[i], args[i+1])
	}

	return nil, nil
}

//==============================================================================================================================
//	 General Functions
//==============================================================================================================================
//	 get_ecert - Takes the name passed and calls out to the REST API for HyperLedger to retrieve the ecert
//				 for that user. Returns the ecert as retrived including html encoding.
//==============================================================================================================================
func (t *SimpleChaincode) get_ecert(stub shim.ChaincodeStubInterface, name string) ([]byte, error) {

	ecert, err := stub.GetState(name)

	if err != nil {
		return nil, errors.New("Couldn't retrieve ecert for user " + name)
	}

	return ecert, nil
}

//==============================================================================================================================
//	 add_ecert - Adds a new ecert and user pair to the table of ecerts
//==============================================================================================================================

func (t *SimpleChaincode) add_ecert(stub shim.ChaincodeStubInterface, name string, ecert string) ([]byte, error) {

	err := stub.PutState(name, []byte(ecert))

	if err == nil {
		return nil, errors.New("Error storing eCert for user " + name + " identity: " + ecert)
	}

	return nil, nil

}

//==============================================================================================================================
//	 get_caller - Retrieves the username of the user who invoked the chaincode.
//				  Returns the username as a string.
//==============================================================================================================================

func (t *SimpleChaincode) get_username(stub shim.ChaincodeStubInterface) (string, error) {

	username, err := stub.ReadCertAttribute("username")
	if err != nil {
		return "", errors.New("Couldn't get attribute 'username'. Error: " + err.Error())
	}
	return string(username), nil
}

//==============================================================================================================================
//	 check_affiliation - Takes an ecert as a string, decodes it to remove html encoding then parses it and checks the
// 				  		certificates common name. The affiliation is stored as part of the common name.
//==============================================================================================================================

func (t *SimpleChaincode) check_affiliation(stub shim.ChaincodeStubInterface) (string, error) {
	affiliation, err := stub.ReadCertAttribute("role")
	if err != nil {
		return "", errors.New("Couldn't get attribute 'role'. Error: " + err.Error())
	}
	return string(affiliation), nil

}

//==============================================================================================================================
//	 get_caller_data - Calls the get_ecert and check_role functions and returns the ecert and role for the
//					 name passed.
//==============================================================================================================================

func (t *SimpleChaincode) get_caller_data(stub shim.ChaincodeStubInterface) (string, string, error) {

	user, err := t.get_username(stub)

	// if err != nil { return "", "", err }

	// ecert, err := t.get_ecert(stub, user);

	// if err != nil { return "", "", err }

	affiliation, err := t.check_affiliation(stub)

	if err != nil {
		return "", "", err
	}

	return user, affiliation, nil
}

//==============================================================================================================================
//	 retrieve_v5c - Gets the state of the data at Song_ID in the ledger then converts it from the stored
//					JSON into the Song struct for use in the contract. Returns the song struct.
//					Returns empty v if it errors.
//==============================================================================================================================
func (t *SimpleChaincode) retrieve_Song_ID(stub shim.ChaincodeStubInterface, Song_ID string) (Song, error) {

	var s Song

	bytes, err := stub.GetState(Song_ID)

	if err != nil {
		fmt.Printf("RETRIEVE_Song_ID: Failed to invoke song_code: %s", err)
		return v, errors.New("RETRIEVE_Song_ID: Error retrieving song with Song_ID = " + Song_ID)
	}

	err = json.Unmarshal(bytes, &s)

	if err != nil {
		fmt.Printf("RETRIEVE_Song_ID: Corrupt song record "+string(bytes)+": %s", err)
		return s, errors.New("RETRIEVE_Song_ID: Corrupt song record" + string(bytes))
	}

	return s, nil
}

//==============================================================================================================================
// save_changes - Writes to the ledger the song struct passed in a JSON format. Uses the shim file's
//				  method 'PutState'.
//==============================================================================================================================
func (t *SimpleChaincode) save_changes(stub shim.ChaincodeStubInterface, s Song) (bool, error) {

	bytes, err := json.Marshal(s)

	if err != nil {
		fmt.Printf("SAVE_CHANGES: Error converting song record: %s", err)
		return false, errors.New("Error converting song record")
	}

	err = stub.PutState(s.Song_ID, bytes)

	if err != nil {
		fmt.Printf("SAVE_CHANGES: Error storing song record: %s", err)
		return false, errors.New("Error storing song record")
	}

	return true, nil
}

//==============================================================================================================================
//	 Router Functions
//==============================================================================================================================
//	Invoke - Called on chaincode invoke. Takes a function name passed and calls that function. Converts some
//		  initial arguments passed to other things for use in the called function e.g. name -> ecert
//==============================================================================================================================
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	caller, caller_affiliation, err := t.get_caller_data(stub)

	if err != nil {
		return nil, errors.New("Error retrieving caller information")
	}

	if function == "Set_Song" {
		return t.Set_Song(stub, caller, caller_affiliation, args[0])
	} else if function == "ping" {
		return t.ping(stub)
	} else { // If the function is not a create then there must be a Song so we need to retrieve the Song.
		argPos := 1

		s, err := t.retrieve_Song_ID(stub, args[argPos])

		if err == nil {
			fmt.Printf("INVOKE: Song ID already exists in BlockChain")
			return nil, errors.New("Song ID already exists in BlockChain")
		}

		if function == "Set_Contract" {
			return t.Set_Contract(stub, s, caller, caller_affiliation, args[0], "manufacturer")
		} else if function == "Set_Rating" {
			return t.Set_Rating(stub, s, caller, caller_affiliation, args[0], "private")
		} else if function == "Set_Song" {
			return t.Set_Song(stub, s, caller, caller_affiliation, args[0], "private")
		} else if function == "Set_Contract_Response" {
			return t.Set_Contract_Response(stub, s, caller, caller_affiliation, args[0], "lease_company")
		} else if function == "update_song" {
			return t.update_song(stub, s, caller, caller_affiliation, args[0])
		} else if function == "update_contract" {
			return t.update_contract(stub, s, caller, caller_affiliation, args[0])
		} else if function == "update_rating" {
			return t.update_rating(stub, s, caller, caller_affiliation, args[0])
		}

		return nil, errors.New("Function of the name " + function + " doesn't exist.")

	}
}

//=================================================================================================================================
//	Query - Called on chaincode query. Takes a function name passed and calls that function. Passes the
//  		initial arguments passed are passed on to the called function.
//=================================================================================================================================
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	caller, caller_affiliation, err := t.get_caller_data(stub)
	if err != nil {
		fmt.Printf("QUERY: Error retrieving caller details", err)
		return nil, errors.New("QUERY: Error retrieving caller details: " + err.Error())
	}

	logger.Debug("function: ", function)
	logger.Debug("caller: ", caller)
	logger.Debug("affiliation: ", caller_affiliation)

	//	if function == "get_vehicle_details" {
	//		if len(args) != 1 { fmt.Printf("Incorrect number of arguments passed"); return nil, errors.New("QUERY: Incorrect number of arguments passed") }
	//		v, err := t.retrieve_v5c(stub, args[0])
	//		if err != nil { fmt.Printf("QUERY: Error retrieving v5c: %s", err); return nil, errors.New("QUERY: Error retrieving v5c "+err.Error()) }
	//		return t.get_vehicle_details(stub, v, caller, caller_affiliation)
	//	} else if function == "check_unique_v5c" {
	//		return t.check_unique_v5c(stub, args[0], caller, caller_affiliation)
	//	} else if function == "get_vehicles" {
	//		return t.get_vehicles(stub, caller, caller_affiliation)
	//	} else if function == "get_ecert" {
	//		return t.get_ecert(stub, args[0])
	//	} else if function == "ping" {
	//		return t.ping(stub)
	//	}

	if function == "Get_overall_Rating" {
		if len(args) != 1 {
			fmt.Printf("Incorrect number of arguments passed")
			return nil, errors.New("QUERY: Incorrect number of arguments passed")
		}
		s, err := t.retrieve_Song_ID(stub, args[0])
		if err != nil {
			fmt.Printf("QUERY: Error retrieving Song: %s", err)
			return nil, errors.New("QUERY: Error retrieving Song " + err.Error())
		}
		return t.Get_overall_Rating(stub, s, caller, caller_affiliation)
	} else if function == "Get_Rating" {
		return t.Get_Rating(stub, args[0], caller, caller_affiliation)
	} else if function == "Get_Contract" {
		return t.Get_Contract(stub, caller, caller_affiliation)
	} else if function == "get_ecert" {
		return t.get_ecert(stub, args[0])
	} else if function == "ping" {
		return t.ping(stub)
	}

	return nil, errors.New("Received unknown function invocation " + function)

}

//=================================================================================================================================
//	 Ping Function
//=================================================================================================================================
//	 Pings the peer to keep the connection alive
//=================================================================================================================================
func (t *SimpleChaincode) ping(stub shim.ChaincodeStubInterface) ([]byte, error) {
	return []byte("Hello, world!"), nil
}

//=================================================================================================================================
//	 Create Function
//=================================================================================================================================
//	 Create Vehicle - Creates the initial JSON for the Song and then saves it to the ledger.
//=================================================================================================================================
func (t *SimpleChaincode) create_song(stub shim.ChaincodeStubInterface, caller string, caller_affiliation string, Song_ID string) ([]byte, error) {
	var s Song

	Song_ID := "\"v5cID\":\"" + Song_ID + "\", " // Variables to define the JSON
	SmartContract_Unique_ID := "\"VIN\":0, "
	Singer_Id := "\"Singer_Id\":\"UNDEFINED\", "
	Singer_Name := "\"Singer_Name\":\"UNDEFINED\", "
	Video_Id := "\"Video_Id \":\"UNDEFINED\", "
	owner := "\"Owner\":\"" + caller + "\", "
	Video_Link := "\"Video_Link\":\"UNDEFINED\", "
	Video_date_created := "\"Video_date_created\":\"UNDEFINED\", "
	Video_QR_code_Id := "\"Video_QR_code_Id\":\"UNDEFINED\", "
	Copyright_Id := "\"Copyright_Id\":\"UNDEFINED\""
	Copyright_date_created := "\"Copyright_date_created\":\"UNDEFINED\""
	Copyright_date_accepted := "\"Copyright_date_accepted\":\"UNDEFINED\""
	Copyright_date_rejected := "\"Copyright_date_rejected\":\"UNDEFINED\""
	Copyright_Institution_Id := "\"Copyright_Institution_Id\":\"UNDEFINED\""
	Copyright_Institution_Name := "\"Copyright_Institution_Name\":\"UNDEFINED\""
	Copyright_State := "\"Copyright_State\":\"UNDEFINED\""
	Venue_Id := "\"Venue_Id\":\"UNDEFINED\""
	Venue_Name := "\"Venue_Name\":\"UNDEFINED\""
	Copyright_Institution_Name := "\"Copyright_Institution_Name\":\"UNDEFINED\""
	User_Id := "\"User_Id\":\"UNDEFINED\""
	User_role := "\"User_role\":\"UNDEFINED\""

	Song_json := "{" + Song_ID + SmartContract_Unique_ID + Singer_Id + Singer_Name + Video_Id + owner + Video_Link + Video_date_created + Video_QR_code_Id +
		Copyright_Id + Copyright_date_created + Copyright_date_accepted + Copyright_date_rejected + Copyright_Institution_Id + Copyright_Institution_Name + Copyright_State +
		Venue_Id + Venue_Name + Copyright_Institution_Name + User_Id + User_role + "}" // Concatenates the variables to create the total JSON object

	matched, err := regexp.Match("^[A-z][A-z][0-9]{7}", []byte(Song_ID)) // matched = true if the Song passed fits format of two letters followed by seven digits

	if err != nil {
		fmt.Printf("CREATE_Song: Invalid Song_ID: %s", err)
		return nil, errors.New("Invalid Song_ID")
	}

	if Song_ID == "" ||
		matched == false {
		fmt.Printf("CREATE_SONG: Invalid Song_ID provided")
		return nil, errors.New("Invalid Song_ID provided")
	}

	err = json.Unmarshal([]byte(Song_json), &v) // Convert the JSON defined above into a Song object for go

	if err != nil {
		return nil, errors.New("Invalid JSON object")
	}

	record, err := stub.GetState(s.Song_ID) // If not an error then a record exists so cant create a new car with this Song_ID as it must be unique

	if record != nil {
		return nil, errors.New("Song already exists")
	}

	if caller_affiliation != AUTHORITY { // Only the regulator can create a new Song

		return nil, errors.New(fmt.Sprintf("Permission Denied. create_song. %v === %v", caller_affiliation, AUTHORITY))

	}

	_, err = t.save_changes(stub, s)

	if err != nil {
		fmt.Printf("CREATE_SONG: Error saving changes: %s", err)
		return nil, errors.New("Error saving changes")
	}

	bytes, err := stub.GetState("Song_ID")

	if err != nil {
		return nil, errors.New("Unable to get Song_ID")
	}

	var Song_ID V5C_Holder

	err = json.Unmarshal(bytes, &Song_ID)

	if err != nil {
		return nil, errors.New("Corrupt V5C_Holder record")
	}

	v5cIDs.V5Cs = append(v5cIDs.V5Cs, v5cID)

	bytes, err = json.Marshal(v5cIDs)

	if err != nil {
		fmt.Print("Error creating V5C_Holder record")
	}

	err = stub.PutState("v5cIDs", bytes)

	if err != nil {
		return nil, errors.New("Unable to put the state")
	}

	return nil, nil

}

//=================================================================================================================================
//	 Transfer Functions
//=================================================================================================================================
//	 authority_to_manufacturer
//=================================================================================================================================
func (t *SimpleChaincode) authority_to_manufacturer(stub shim.ChaincodeStubInterface, v Vehicle, caller string, caller_affiliation string, recipient_name string, recipient_affiliation string) ([]byte, error) {

	if v.Status == STATE_TEMPLATE &&
		v.Owner == caller &&
		caller_affiliation == AUTHORITY &&
		recipient_affiliation == MANUFACTURER &&
		v.Scrapped == false { // If the roles and users are ok

		v.Owner = recipient_name     // then make the owner the new owner
		v.Status = STATE_MANUFACTURE // and mark it in the state of manufacture

	} else { // Otherwise if there is an error
		fmt.Printf("AUTHORITY_TO_MANUFACTURER: Permission Denied")
		return nil, errors.New(fmt.Sprintf("Permission Denied. authority_to_manufacturer. %v %v === %v, %v === %v, %v === %v, %v === %v, %v === %v", v, v.Status, STATE_PRIVATE_OWNERSHIP, v.Owner, caller, caller_affiliation, PRIVATE_ENTITY, recipient_affiliation, SCRAP_MERCHANT, v.Scrapped, false))

	}

	_, err := t.save_changes(stub, v) // Write new state

	if err != nil {
		fmt.Printf("AUTHORITY_TO_MANUFACTURER: Error saving changes: %s", err)
		return nil, errors.New("Error saving changes")
	}

	return nil, nil // We are Done

}

//=================================================================================================================================
//	 manufacturer_to_private
//=================================================================================================================================
func (t *SimpleChaincode) manufacturer_to_private(stub shim.ChaincodeStubInterface, v Vehicle, caller string, caller_affiliation string, recipient_name string, recipient_affiliation string) ([]byte, error) {

	if v.Make == "UNDEFINED" ||
		v.Model == "UNDEFINED" ||
		v.Reg == "UNDEFINED" ||
		v.Colour == "UNDEFINED" ||
		v.VIN == 0 { //If any part of the car is undefined it has not bene fully manufacturered so cannot be sent
		fmt.Printf("MANUFACTURER_TO_PRIVATE: Car not fully defined")
		return nil, errors.New(fmt.Sprintf("Car not fully defined. %v", v))
	}

	if v.Status == STATE_MANUFACTURE &&
		v.Owner == caller &&
		caller_affiliation == MANUFACTURER &&
		recipient_affiliation == PRIVATE_ENTITY &&
		v.Scrapped == false {

		v.Owner = recipient_name
		v.Status = STATE_PRIVATE_OWNERSHIP

	} else {
		return nil, errors.New(fmt.Sprintf("Permission Denied. manufacturer_to_private. %v %v === %v, %v === %v, %v === %v, %v === %v, %v === %v", v, v.Status, STATE_PRIVATE_OWNERSHIP, v.Owner, caller, caller_affiliation, PRIVATE_ENTITY, recipient_affiliation, SCRAP_MERCHANT, v.Scrapped, false))
	}

	_, err := t.save_changes(stub, v)

	if err != nil {
		fmt.Printf("MANUFACTURER_TO_PRIVATE: Error saving changes: %s", err)
		return nil, errors.New("Error saving changes")
	}

	return nil, nil

}

//=================================================================================================================================
//	 private_to_private
//=================================================================================================================================
func (t *SimpleChaincode) private_to_private(stub shim.ChaincodeStubInterface, v Vehicle, caller string, caller_affiliation string, recipient_name string, recipient_affiliation string) ([]byte, error) {

	if v.Status == STATE_PRIVATE_OWNERSHIP &&
		v.Owner == caller &&
		caller_affiliation == PRIVATE_ENTITY &&
		recipient_affiliation == PRIVATE_ENTITY &&
		v.Scrapped == false {

		v.Owner = recipient_name

	} else {
		return nil, errors.New(fmt.Sprintf("Permission Denied. private_to_private. %v %v === %v, %v === %v, %v === %v, %v === %v, %v === %v", v, v.Status, STATE_PRIVATE_OWNERSHIP, v.Owner, caller, caller_affiliation, PRIVATE_ENTITY, recipient_affiliation, SCRAP_MERCHANT, v.Scrapped, false))
	}

	_, err := t.save_changes(stub, v)

	if err != nil {
		fmt.Printf("PRIVATE_TO_PRIVATE: Error saving changes: %s", err)
		return nil, errors.New("Error saving changes")
	}

	return nil, nil

}

//=================================================================================================================================
//	 private_to_lease_company
//=================================================================================================================================
func (t *SimpleChaincode) private_to_lease_company(stub shim.ChaincodeStubInterface, v Vehicle, caller string, caller_affiliation string, recipient_name string, recipient_affiliation string) ([]byte, error) {

	if v.Status == STATE_PRIVATE_OWNERSHIP &&
		v.Owner == caller &&
		caller_affiliation == PRIVATE_ENTITY &&
		recipient_affiliation == LEASE_COMPANY &&
		v.Scrapped == false {

		v.Owner = recipient_name

	} else {
		return nil, errors.New(fmt.Sprintf("Permission denied. private_to_lease_company. %v === %v, %v === %v, %v === %v, %v === %v, %v === %v", v.Status, STATE_PRIVATE_OWNERSHIP, v.Owner, caller, caller_affiliation, PRIVATE_ENTITY, recipient_affiliation, SCRAP_MERCHANT, v.Scrapped, false))

	}

	_, err := t.save_changes(stub, v)
	if err != nil {
		fmt.Printf("PRIVATE_TO_LEASE_COMPANY: Error saving changes: %s", err)
		return nil, errors.New("Error saving changes")
	}

	return nil, nil

}

//=================================================================================================================================
//	 lease_company_to_private
//=================================================================================================================================
func (t *SimpleChaincode) lease_company_to_private(stub shim.ChaincodeStubInterface, v Vehicle, caller string, caller_affiliation string, recipient_name string, recipient_affiliation string) ([]byte, error) {

	if v.Status == STATE_PRIVATE_OWNERSHIP &&
		v.Owner == caller &&
		caller_affiliation == LEASE_COMPANY &&
		recipient_affiliation == PRIVATE_ENTITY &&
		v.Scrapped == false {

		v.Owner = recipient_name

	} else {
		return nil, errors.New(fmt.Sprintf("Permission Denied. lease_company_to_private. %v %v === %v, %v === %v, %v === %v, %v === %v, %v === %v", v, v.Status, STATE_PRIVATE_OWNERSHIP, v.Owner, caller, caller_affiliation, PRIVATE_ENTITY, recipient_affiliation, SCRAP_MERCHANT, v.Scrapped, false))
	}

	_, err := t.save_changes(stub, v)
	if err != nil {
		fmt.Printf("LEASE_COMPANY_TO_PRIVATE: Error saving changes: %s", err)
		return nil, errors.New("Error saving changes")
	}

	return nil, nil

}

//=================================================================================================================================
//	 private_to_scrap_merchant
//=================================================================================================================================
func (t *SimpleChaincode) private_to_scrap_merchant(stub shim.ChaincodeStubInterface, v Vehicle, caller string, caller_affiliation string, recipient_name string, recipient_affiliation string) ([]byte, error) {

	if v.Status == STATE_PRIVATE_OWNERSHIP &&
		v.Owner == caller &&
		caller_affiliation == PRIVATE_ENTITY &&
		recipient_affiliation == SCRAP_MERCHANT &&
		v.Scrapped == false {

		v.Owner = recipient_name
		v.Status = STATE_BEING_SCRAPPED

	} else {
		return nil, errors.New(fmt.Sprintf("Permission Denied. private_to_scrap_merchant. %v %v === %v, %v === %v, %v === %v, %v === %v, %v === %v", v, v.Status, STATE_PRIVATE_OWNERSHIP, v.Owner, caller, caller_affiliation, PRIVATE_ENTITY, recipient_affiliation, SCRAP_MERCHANT, v.Scrapped, false))
	}

	_, err := t.save_changes(stub, v)

	if err != nil {
		fmt.Printf("PRIVATE_TO_SCRAP_MERCHANT: Error saving changes: %s", err)
		return nil, errors.New("Error saving changes")
	}

	return nil, nil

}

//=================================================================================================================================
//	 Update Functions
//=================================================================================================================================
//	 update_vin
//=================================================================================================================================
func (t *SimpleChaincode) update_vin(stub shim.ChaincodeStubInterface, v Vehicle, caller string, caller_affiliation string, new_value string) ([]byte, error) {

	new_vin, err := strconv.Atoi(string(new_value)) // will return an error if the new vin contains non numerical chars

	if err != nil || len(string(new_value)) != 15 {
		return nil, errors.New("Invalid value passed for new VIN")
	}

	if v.Status == STATE_MANUFACTURE &&
		v.Owner == caller &&
		caller_affiliation == MANUFACTURER &&
		v.VIN == 0 && // Can't change the VIN after its initial assignment
		v.Scrapped == false {

		v.VIN = new_vin // Update to the new value
	} else {

		return nil, errors.New(fmt.Sprintf("Permission denied. update_vin %v %v %v %v %v", v.Status, STATE_MANUFACTURE, v.Owner, caller, v.VIN, v.Scrapped))

	}

	_, err = t.save_changes(stub, v) // Save the changes in the blockchain

	if err != nil {
		fmt.Printf("UPDATE_VIN: Error saving changes: %s", err)
		return nil, errors.New("Error saving changes")
	}

	return nil, nil

}

//=================================================================================================================================
//	 update_registration
//=================================================================================================================================
func (t *SimpleChaincode) update_registration(stub shim.ChaincodeStubInterface, v Vehicle, caller string, caller_affiliation string, new_value string) ([]byte, error) {

	if v.Owner == caller &&
		caller_affiliation != SCRAP_MERCHANT &&
		v.Scrapped == false {

		v.Reg = new_value

	} else {
		return nil, errors.New(fmt.Sprint("Permission denied. update_registration"))
	}

	_, err := t.save_changes(stub, v)

	if err != nil {
		fmt.Printf("UPDATE_REGISTRATION: Error saving changes: %s", err)
		return nil, errors.New("Error saving changes")
	}

	return nil, nil

}

//=================================================================================================================================
//	 update_colour
//=================================================================================================================================
func (t *SimpleChaincode) update_colour(stub shim.ChaincodeStubInterface, v Vehicle, caller string, caller_affiliation string, new_value string) ([]byte, error) {

	if v.Owner == caller &&
		caller_affiliation == MANUFACTURER && /*((v.Owner				== caller			&&
		caller_affiliation	== MANUFACTURER)		||
		caller_affiliation	== AUTHORITY)			&&*/
		v.Scrapped == false {

		v.Colour = new_value
	} else {

		return nil, errors.New(fmt.Sprint("Permission denied. update_colour %t %t %t"+v.Owner == caller, caller_affiliation == MANUFACTURER, v.Scrapped))
	}

	_, err := t.save_changes(stub, v)

	if err != nil {
		fmt.Printf("UPDATE_COLOUR: Error saving changes: %s", err)
		return nil, errors.New("Error saving changes")
	}

	return nil, nil

}

//=================================================================================================================================
//	 update_make
//=================================================================================================================================
func (t *SimpleChaincode) update_make(stub shim.ChaincodeStubInterface, v Vehicle, caller string, caller_affiliation string, new_value string) ([]byte, error) {

	if v.Status == STATE_MANUFACTURE &&
		v.Owner == caller &&
		caller_affiliation == MANUFACTURER &&
		v.Scrapped == false {

		v.Make = new_value
	} else {

		return nil, errors.New(fmt.Sprint("Permission denied. update_make %t %t %t"+v.Owner == caller, caller_affiliation == MANUFACTURER, v.Scrapped))

	}

	_, err := t.save_changes(stub, v)

	if err != nil {
		fmt.Printf("UPDATE_MAKE: Error saving changes: %s", err)
		return nil, errors.New("Error saving changes")
	}

	return nil, nil

}

//=================================================================================================================================
//	 update_model
//=================================================================================================================================
func (t *SimpleChaincode) update_model(stub shim.ChaincodeStubInterface, v Vehicle, caller string, caller_affiliation string, new_value string) ([]byte, error) {

	if v.Status == STATE_MANUFACTURE &&
		v.Owner == caller &&
		caller_affiliation == MANUFACTURER &&
		v.Scrapped == false {

		v.Model = new_value

	} else {
		return nil, errors.New(fmt.Sprint("Permission denied. update_model %t %t %t"+v.Owner == caller, caller_affiliation == MANUFACTURER, v.Scrapped))

	}

	_, err := t.save_changes(stub, v)

	if err != nil {
		fmt.Printf("UPDATE_MODEL: Error saving changes: %s", err)
		return nil, errors.New("Error saving changes")
	}

	return nil, nil

}

//=================================================================================================================================
//	 scrap_vehicle
//=================================================================================================================================
func (t *SimpleChaincode) scrap_vehicle(stub shim.ChaincodeStubInterface, v Vehicle, caller string, caller_affiliation string) ([]byte, error) {

	if v.Status == STATE_BEING_SCRAPPED &&
		v.Owner == caller &&
		caller_affiliation == SCRAP_MERCHANT &&
		v.Scrapped == false {

		v.Scrapped = true

	} else {
		return nil, errors.New("Permission denied. scrap_vehicle")
	}

	_, err := t.save_changes(stub, v)

	if err != nil {
		fmt.Printf("SCRAP_VEHICLE: Error saving changes: %s", err)
		return nil, errors.New("SCRAP_VEHICLError saving changes")
	}

	return nil, nil

}

//=================================================================================================================================
//	 Read Functions
//=================================================================================================================================
//	 get_vehicle_details
//=================================================================================================================================
func (t *SimpleChaincode) get_vehicle_details(stub shim.ChaincodeStubInterface, v Vehicle, caller string, caller_affiliation string) ([]byte, error) {

	bytes, err := json.Marshal(v)

	if err != nil {
		return nil, errors.New("GET_VEHICLE_DETAILS: Invalid vehicle object")
	}

	if v.Owner == caller ||
		caller_affiliation == AUTHORITY {

		return bytes, nil
	} else {
		return nil, errors.New("Permission Denied. get_vehicle_details")
	}

}

//=================================================================================================================================
//	 get_vehicles
//=================================================================================================================================

func (t *SimpleChaincode) get_vehicles(stub shim.ChaincodeStubInterface, caller string, caller_affiliation string) ([]byte, error) {
	bytes, err := stub.GetState("v5cIDs")

	if err != nil {
		return nil, errors.New("Unable to get v5cIDs")
	}

	var v5cIDs V5C_Holder

	err = json.Unmarshal(bytes, &v5cIDs)

	if err != nil {
		return nil, errors.New("Corrupt V5C_Holder")
	}

	result := "["

	var temp []byte
	var v Vehicle

	for _, v5c := range v5cIDs.V5Cs {

		v, err = t.retrieve_v5c(stub, v5c)

		if err != nil {
			return nil, errors.New("Failed to retrieve V5C")
		}

		temp, err = t.get_vehicle_details(stub, v, caller, caller_affiliation)

		if err == nil {
			result += string(temp) + ","
		}
	}

	if len(result) == 1 {
		result = "[]"
	} else {
		result = result[:len(result)-1] + "]"
	}

	return []byte(result), nil
}

//=================================================================================================================================
//	 check_unique_Song_ID
//=================================================================================================================================
func (t *SimpleChaincode) check_unique_Song_ID(stub shim.ChaincodeStubInterface, Song_ID string, caller string, caller_affiliation string) ([]byte, error) {
	_, err := t.retrieve_Song_ID(stub, Song_ID)
	if err == nil {
		return []byte("false"), errors.New("Song_ID is not unique")
	} else {
		return []byte("true"), nil
	}
}

//=================================================================================================================================
//	 Main - main - Starts up the chaincode
//=================================================================================================================================
func main() {

	err := shim.Start(new(SimpleChaincode))

	if err != nil {
		fmt.Printf("Error starting Chaincode: %s", err)
	}
}
