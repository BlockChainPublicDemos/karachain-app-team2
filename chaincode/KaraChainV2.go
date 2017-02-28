package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"regexp"
	//	"strconv"
	//	"strings"
)

var logger = shim.NewLogger("CLDChaincode")

//==============================================================================================================================
//	 Participant types - Each participant type is mapped to an integer which we use to compare to the value stored in a
//						 user's eCert
//==============================================================================================================================
//CURRENT WORKAROUND USES ROLES CHANGE WHEN OWN USERS CAN BE CREATED SO THAT IT READ 1, 2, 3, 4, 5
const AUTHORITY = "regulator"
const SINGER = "singer"
const AUDIENCE = "AUDIENCE"
const COPYRIGHT_AUTHORITY = "copyright_authority "

//==============================================================================================================================
//	 Status types - Asset lifecycle is broken down into z statuses, this is part of the business logic to determine what can
//					be done to the song at points in it's lifecycle
// Not sure if we need them
//==============================================================================================================================
const STATE_TEMPLATE = 0
const STATE_RECORDED = 1
const STATE_UPDATED = 2
const STATE_VOTED = 3
const STATE_CONTRACT_PROVIDED = 4
const STATE_CONTRACT_REJECTED = 5
const STATE_CONTRACT_ACCEPTED = 6
const STATE_OBSOLETE = 7

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
	Song_ID                    string `json:"Song_ID"`
	Date_created               string `json:"Date_created"`
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
	Obsolete                   bool   `json:"Obsolete"`
	Status                     string `json:"Status"`
}

//==============================================================================================================================
//	Song Holder - Defines the structure that holds all the Song_IDs for songs that have been created.
//				Used as an index when querying all songs
//==============================================================================================================================

type Song_Holder struct {
	Song_IDs []string `json:"Song_IDs"`
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

	bytes, err := json.Marshal(Song_IDs)

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
//	 retrieve_Song_ID - Gets the state of the data at Song_ID in the ledger then converts it from the stored
//					JSON into the Song struct for use in the contract. Returns the song struct.
//					Returns empty v if it errors.
//==============================================================================================================================
func (t *SimpleChaincode) retrieve_Song_ID(stub shim.ChaincodeStubInterface, Song_ID string) (Song, error) {

	var s Song

	bytes, err := stub.GetState(Song_ID)

	if err != nil {
		fmt.Printf("RETRIEVE_Song_ID: Failed to invoke song_code: %s", err)
		return s, errors.New("RETRIEVE_Song_ID: Error retrieving song with Song_ID = " + Song_ID)
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

	if function == "create_song" { // we create a song from scratch
		return t.create_song(stub, caller, caller_affiliation, args[0])
	} else if function == "ping" {
		return t.ping(stub)
	} else { // If the function is not a create then there must be a Song so we need to retrieve the Song.
		argPos := 1

		s, err := t.retrieve_Song_ID(stub, args[argPos])

		if err != nil {
			fmt.Printf("INVOKE: Error retrieving Song: %s", err)
			return nil, errors.New("Error retrieving Song")
		}

		if function == "Set_Contract" { // Only the copyright authority is allowed to set a contract
			return t.set_contract(stub, s, caller, caller_affiliation, args[0])
		} else if function == "Set_Rating" { // Rating can be set by anybody, but only once and not by the singer
			return t.set_rating(stub, s, caller, caller_affiliation, args[0])
		} else if function == "Set_Contract_Response" { // Function my only be called by the singer
			return t.set_contract_response(stub, s, caller, caller_affiliation, args[0])
		} else if function == "update_song" { // Function may only be called by the singer
			return t.update_song(stub, s, caller, caller_affiliation, args[0])
		} else if function == "update_contract" { // Function may only be called by the copyright institution if the existing contract
			return t.update_contract(stub, s, caller, caller_affiliation, args[0])
		} else if function == "update_rating" { // Rating can be set by anybody, but only once and not by the singer. In this case, a previous rating of the user must exist
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

	if function == "Get_Song" { // Allowed by anybody to get the latest song details. Audience should not see contract details
		if len(args) != 1 {
			fmt.Printf("Incorrect number of arguments passed")
			return nil, errors.New("QUERY: Incorrect number of arguments passed")
		}
		s, err := t.retrieve_Song_ID(stub, args[0])
		if err != nil {
			fmt.Printf("QUERY: Error retrieving Song: %s", err)
			return nil, errors.New("QUERY: Error retrieving Song " + err.Error())
		}
		return t.get_song_details(stub, s, caller, caller_affiliation)
	} else if function == "Get_Rating" {
		return t.get_rating(stub, args[0], caller, caller_affiliation) // A user should be able to get his own rating that was made in the past for a particular song
	} else if function == "Get_Contract" { // Only allowed for singer or copyright authority to see the latest contract
		return t.get_contract(stub, args[0], caller, caller_affiliation)
	} else if function == "Get_overall_Rating" { //Anybody should be able to see the overall (average) rating of a song
		return t.get_overall_rating(stub, args[0], caller, caller_affiliation)
	} else if function == "Get_Songs" {
		return t.get_songs(stub, caller, caller_affiliation)
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
//	 Create Song - Creates the initial JSON for the Song and then saves it to the ledger.
//=================================================================================================================================
func (t *SimpleChaincode) create_song(stub shim.ChaincodeStubInterface, caller string, caller_affiliation string, Song_ID_r string) ([]byte, error) {
	var s Song

	Song_ID := "\"Song_ID\":\"" + Song_ID_r + "\", " // Variables to define the JSON
	Date_created := "\"Date_created\":\"UNDEFINED\", "
	SmartContract_Unique_ID := "\"SmartContract_Unique_ID\":0, "
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
	User_Id := "\"User_Id\":\"UNDEFINED\""
	User_role := "\"User_role\":\"UNDEFINED\""
	Obsolete := "\"Obsolete\":\"False\""
	Status := "\"Status\":\"False\""

	Song_json := "{" + Song_ID + Date_created + SmartContract_Unique_ID + Singer_Id + Singer_Name + Video_Id + owner + Video_Link + Video_date_created + Video_QR_code_Id +
		Copyright_Id + Copyright_date_created + Copyright_date_accepted + Copyright_date_rejected + Copyright_Institution_Id + Copyright_Institution_Name + Copyright_State +
		Venue_Id + Venue_Name + Copyright_Institution_Name + User_Id + User_role + Obsolete + Status + "}" // Concatenates the variables to create the total JSON object

	// Do we need a certain criteria for a song ID?
	_, err := regexp.Match("^[A-z][A-z][0-9]{7}", []byte(Song_ID)) // matched = true if the Song ID passed fits format of two letters followed by seven digits

	if err != nil {
		fmt.Printf("CREATE_Song: Invalid Song_ID: %s", err)
		return nil, errors.New("Invalid Song_ID")
	}

	if Song_ID_r == "" {
		fmt.Printf("CREATE_SONG: Invalid Song_ID provided")
		return nil, errors.New("Invalid Song_ID provided")
	}

	err = json.Unmarshal([]byte(Song_json), &s) // Convert the JSON defined above into a Song object for go

	if err != nil {
		return nil, errors.New("Invalid JSON object")
	}

	record, err := stub.GetState(s.Song_ID) // If not an error then a record exists so cant create a new car with this Song_ID as it must be unique

	if record != nil {
		return nil, errors.New("Song already exists")
	}

	if caller_affiliation != SINGER { // Only the singer can create a new Song
		return nil, errors.New(fmt.Sprintf("Permission Denied. create_song. %v === %v", caller_affiliation, SINGER))
	}

	_, err = t.save_changes(stub, s)

	if err != nil {
		fmt.Printf("CREATE_SONG: Error saving changes: %s", err)
		return nil, errors.New("Error saving changes")
	}

	bytes, err := stub.GetState("Song_IDs")

	if err != nil {
		return nil, errors.New("Unable to get Song_ID")
	}

	var Song_IDs Song_Holder // Not sure what this holder means. Need to check

	err = json.Unmarshal(bytes, &Song_IDs)

	if err != nil {
		return nil, errors.New("Corrupt Song_Holder record") // Not sure what this holder means. Need to check
	}

	err = stub.PutState("Song_IDs", bytes)

	if err != nil {
		return nil, errors.New("Unable to put the state")
	}

	return nil, nil

}

//=================================================================================================================================
//	 Transfer Functions
//=================================================================================================================================
//	 singer_to_authority - might be needed to transfer a song from a singer to a copyright authority
//=================================================================================================================================
//func (t *SimpleChaincode) singer_to_authority(stub shim.ChaincodeStubInterface, s Song, caller string, caller_affiliation string, recipient_name string, recipient_affiliation string) ([]byte, error) {
//
//	if  S.Owner == caller &&
//		S.Obsolete == false { // If the roles and users are ok
//
//		S.Owner = COPYRIGHT_AUTHORITY       // then make the owner the new owner
//		S.Status = STATE_CONTRACT_ACCEPTED  // and mark it in the state of manufacture
//
//	} else { // Otherwise if there is an error
//		fmt.Printf("singer_to_authority: Permission Denied")
//		return nil, errors.New(fmt.Sprintf("Permission Denied. singer_to_authority. %v %v %v %v %v ", s, s.Status, v.Owner, caller, caller_affiliation))
//
//	}
//
//	_, err := t.save_changes(stub, s) // Write new state
//
//	if err != nil {
//		fmt.Printf("singer_to_authority: Error saving changes: %s", err)
//		return nil, errors.New("Error saving changes")
//	}
//
//	return nil, nil // We are Done
//
//}
//

////=================================================================================================================================
////	 manufacturer_to_private
////=================================================================================================================================
//func (t *SimpleChaincode) manufacturer_to_private(stub shim.ChaincodeStubInterface, v Vehicle, caller string, caller_affiliation string, recipient_name string, recipient_affiliation string) ([]byte, error) {
//
//	if v.Make == "UNDEFINED" ||
//		v.Model == "UNDEFINED" ||
//		v.Reg == "UNDEFINED" ||
//		v.Colour == "UNDEFINED" ||
//		v.VIN == 0 { //If any part of the car is undefined it has not bene fully manufacturered so cannot be sent
//		fmt.Printf("MANUFACTURER_TO_PRIVATE: Car not fully defined")
//		return nil, errors.New(fmt.Sprintf("Car not fully defined. %v", v))
//	}
//
//	if v.Status == STATE_MANUFACTURE &&
//		v.Owner == caller &&
//		caller_affiliation == MANUFACTURER &&
//		recipient_affiliation == PRIVATE_ENTITY &&
//		v.Scrapped == false {
//
//		v.Owner = recipient_name
//		v.Status = STATE_PRIVATE_OWNERSHIP
//
//	} else {
//		return nil, errors.New(fmt.Sprintf("Permission Denied. manufacturer_to_private. %v %v === %v, %v === %v, %v === %v, %v === %v, %v === %v", v, v.Status, STATE_PRIVATE_OWNERSHIP, v.Owner, caller, caller_affiliation, PRIVATE_ENTITY, recipient_affiliation, SCRAP_MERCHANT, v.Scrapped, false))
//	}
//
//	_, err := t.save_changes(stub, v)
//
//	if err != nil {
//		fmt.Printf("MANUFACTURER_TO_PRIVATE: Error saving changes: %s", err)
//		return nil, errors.New("Error saving changes")
//	}
//
//	return nil, nil
//
//}
//
////=================================================================================================================================
////	 private_to_private
////=================================================================================================================================
//func (t *SimpleChaincode) private_to_private(stub shim.ChaincodeStubInterface, v Vehicle, caller string, caller_affiliation string, recipient_name string, recipient_affiliation string) ([]byte, error) {
//
//	if v.Status == STATE_PRIVATE_OWNERSHIP &&
//		v.Owner == caller &&
//		caller_affiliation == PRIVATE_ENTITY &&
//		recipient_affiliation == PRIVATE_ENTITY &&
//		v.Scrapped == false {
//
//		v.Owner = recipient_name
//
//	} else {
//		return nil, errors.New(fmt.Sprintf("Permission Denied. private_to_private. %v %v === %v, %v === %v, %v === %v, %v === %v, %v === %v", v, v.Status, STATE_PRIVATE_OWNERSHIP, v.Owner, caller, caller_affiliation, PRIVATE_ENTITY, recipient_affiliation, SCRAP_MERCHANT, v.Scrapped, false))
//	}
//
//	_, err := t.save_changes(stub, v)
//
//	if err != nil {
//		fmt.Printf("PRIVATE_TO_PRIVATE: Error saving changes: %s", err)
//		return nil, errors.New("Error saving changes")
//	}
//
//	return nil, nil
//
//}
//
////=================================================================================================================================
////	 private_to_lease_company
////=================================================================================================================================
//func (t *SimpleChaincode) private_to_lease_company(stub shim.ChaincodeStubInterface, v Vehicle, caller string, caller_affiliation string, recipient_name string, recipient_affiliation string) ([]byte, error) {
//
//	if v.Status == STATE_PRIVATE_OWNERSHIP &&
//		v.Owner == caller &&
//		caller_affiliation == PRIVATE_ENTITY &&
//		recipient_affiliation == LEASE_COMPANY &&
//		v.Scrapped == false {
//
//		v.Owner = recipient_name
//
//	} else {
//		return nil, errors.New(fmt.Sprintf("Permission denied. private_to_lease_company. %v === %v, %v === %v, %v === %v, %v === %v, %v === %v", v.Status, STATE_PRIVATE_OWNERSHIP, v.Owner, caller, caller_affiliation, PRIVATE_ENTITY, recipient_affiliation, SCRAP_MERCHANT, v.Scrapped, false))
//
//	}
//
//	_, err := t.save_changes(stub, v)
//	if err != nil {
//		fmt.Printf("PRIVATE_TO_LEASE_COMPANY: Error saving changes: %s", err)
//		return nil, errors.New("Error saving changes")
//	}
//
//	return nil, nil
//
//}
//
////=================================================================================================================================
////	 lease_company_to_private
////=================================================================================================================================
//func (t *SimpleChaincode) lease_company_to_private(stub shim.ChaincodeStubInterface, v Vehicle, caller string, caller_affiliation string, recipient_name string, recipient_affiliation string) ([]byte, error) {
//
//	if v.Status == STATE_PRIVATE_OWNERSHIP &&
//		v.Owner == caller &&
//		caller_affiliation == LEASE_COMPANY &&
//		recipient_affiliation == PRIVATE_ENTITY &&
//		v.Scrapped == false {
//
//		v.Owner = recipient_name
//
//	} else {
//		return nil, errors.New(fmt.Sprintf("Permission Denied. lease_company_to_private. %v %v === %v, %v === %v, %v === %v, %v === %v, %v === %v", v, v.Status, STATE_PRIVATE_OWNERSHIP, v.Owner, caller, caller_affiliation, PRIVATE_ENTITY, recipient_affiliation, SCRAP_MERCHANT, v.Scrapped, false))
//	}
//
//	_, err := t.save_changes(stub, v)
//	if err != nil {
//		fmt.Printf("LEASE_COMPANY_TO_PRIVATE: Error saving changes: %s", err)
//		return nil, errors.New("Error saving changes")
//	}
//
//	return nil, nil
//
//}
//
////=================================================================================================================================
////	 private_to_scrap_merchant
////=================================================================================================================================
//func (t *SimpleChaincode) private_to_scrap_merchant(stub shim.ChaincodeStubInterface, v Vehicle, caller string, caller_affiliation string, recipient_name string, recipient_affiliation string) ([]byte, error) {
//
//	if v.Status == STATE_PRIVATE_OWNERSHIP &&
//		v.Owner == caller &&
//		caller_affiliation == PRIVATE_ENTITY &&
//		recipient_affiliation == SCRAP_MERCHANT &&
//		v.Scrapped == false {
//
//		v.Owner = recipient_name
//		v.Status = STATE_BEING_SCRAPPED
//
//	} else {
//		return nil, errors.New(fmt.Sprintf("Permission Denied. private_to_scrap_merchant. %v %v === %v, %v === %v, %v === %v, %v === %v, %v === %v", v, v.Status, STATE_PRIVATE_OWNERSHIP, v.Owner, caller, caller_affiliation, PRIVATE_ENTITY, recipient_affiliation, SCRAP_MERCHANT, v.Scrapped, false))
//	}
//
//	_, err := t.save_changes(stub, v)
//
//	if err != nil {
//		fmt.Printf("PRIVATE_TO_SCRAP_MERCHANT: Error saving changes: %s", err)
//		return nil, errors.New("Error saving changes")
//	}
//
//	return nil, nil
//
//}

//=================================================================================================================================
//	 Update Functions
//=================================================================================================================================
//	 update_song
//=================================================================================================================================
func (t *SimpleChaincode) update_song(stub shim.ChaincodeStubInterface, s Song, caller string, caller_affiliation string, new_value string) ([]byte, error) {

	var err error
	if s.Obsolete == true {

		return nil, errors.New(fmt.Sprintf("Song is obsolete and cannot be updated. %v", s.Obsolete))

	}

	_, err = t.save_changes(stub, s) // Save the changes in the blockchain

	if err != nil {
		fmt.Printf("UPDATE_SONG: Error saving changes: %s", err)
		return nil, errors.New("Error saving changes")
	}

	return nil, nil

}

//=================================================================================================================================
//	 set_rating - A song is voted by an user. Only 1 vote allowed per user and song.
//=================================================================================================================================
func (t *SimpleChaincode) set_rating(stub shim.ChaincodeStubInterface, s Song, caller string, caller_affiliation string, new_value string) ([]byte, error) {

	// to be implemented
	return nil, nil

}

//=================================================================================================================================
//	 set_contract - The CA may provide a contract to a singer for a particular song. Only 1 contract allowed per copyright authority
//=================================================================================================================================
func (t *SimpleChaincode) set_contract(stub shim.ChaincodeStubInterface, s Song, caller string, caller_affiliation string, new_value string) ([]byte, error) {

	// to be implemented
	return nil, nil

}

//=================================================================================================================================
//	 set_contract_response - The singer can either accept or reject a contract.
//=================================================================================================================================
func (t *SimpleChaincode) set_contract_response(stub shim.ChaincodeStubInterface, s Song, caller string, caller_affiliation string, new_value string) ([]byte, error) {

	// to be implemented
	return nil, nil

}

//=================================================================================================================================
//	 update_contract - A new contract can be provided by the CA. This again has to be approved or rejected by the singer.
//=================================================================================================================================
func (t *SimpleChaincode) update_contract(stub shim.ChaincodeStubInterface, s Song, caller string, caller_affiliation string, new_value string) ([]byte, error) {

	// to be implemented
	return nil, nil

}

//=================================================================================================================================
//	 update_rating - A user can update its own rating that was done in the past.
//=================================================================================================================================
func (t *SimpleChaincode) update_rating(stub shim.ChaincodeStubInterface, s Song, caller string, caller_affiliation string, new_value string) ([]byte, error) {

	// to be implemented
	return nil, nil

}

//=================================================================================================================================
//	 update_registration
//=================================================================================================================================
//func (t *SimpleChaincode) update_registration(stub shim.ChaincodeStubInterface, v Vehicle, caller string, caller_affiliation string, new_value string) ([]byte, error) {
//
//	if v.Owner == caller &&
//		caller_affiliation != SCRAP_MERCHANT &&
//		v.Scrapped == false {
//
//		v.Reg = new_value
//
//	} else {
//		return nil, errors.New(fmt.Sprint("Permission denied. update_registration"))
//	}
//
//	_, err := t.save_changes(stub, v)
//
//	if err != nil {
//		fmt.Printf("UPDATE_REGISTRATION: Error saving changes: %s", err)
//		return nil, errors.New("Error saving changes")
//	}
//
//	return nil, nil
//
//}

//=================================================================================================================================
//	 update_colour
////=================================================================================================================================
//func (t *SimpleChaincode) update_colour(stub shim.ChaincodeStubInterface, v Vehicle, caller string, caller_affiliation string, new_value string) ([]byte, error) {
//
//	if v.Owner == caller &&
//		caller_affiliation == MANUFACTURER && /*((v.Owner				== caller			&&
//		caller_affiliation	== MANUFACTURER)		||
//		caller_affiliation	== AUTHORITY)			&&*/
//		v.Scrapped == false {
//
//		v.Colour = new_value
//	} else {
//
//		return nil, errors.New(fmt.Sprint("Permission denied. update_colour %t %t %t"+v.Owner == caller, caller_affiliation == MANUFACTURER, v.Scrapped))
//	}
//
//	_, err := t.save_changes(stub, v)
//
//	if err != nil {
//		fmt.Printf("UPDATE_COLOUR: Error saving changes: %s", err)
//		return nil, errors.New("Error saving changes")
//	}
//
//	return nil, nil
//
//}

//=================================================================================================================================
//	 update_make
////=================================================================================================================================
//func (t *SimpleChaincode) update_make(stub shim.ChaincodeStubInterface, v Vehicle, caller string, caller_affiliation string, new_value string) ([]byte, error) {
//
//	if v.Status == STATE_MANUFACTURE &&
//		v.Owner == caller &&
//		caller_affiliation == MANUFACTURER &&
//		v.Scrapped == false {
//
//		v.Make = new_value
//	} else {
//
//		return nil, errors.New(fmt.Sprint("Permission denied. update_make %t %t %t"+v.Owner == caller, caller_affiliation == MANUFACTURER, v.Scrapped))
//
//	}
//
//	_, err := t.save_changes(stub, v)
//
//	if err != nil {
//		fmt.Printf("UPDATE_MAKE: Error saving changes: %s", err)
//		return nil, errors.New("Error saving changes")
//	}
//
//	return nil, nil
//
//}

//=================================================================================================================================
////	 update_model
////=================================================================================================================================
//func (t *SimpleChaincode) update_model(stub shim.ChaincodeStubInterface, v Vehicle, caller string, caller_affiliation string, new_value string) ([]byte, error) {
//
//	if v.Status == STATE_MANUFACTURE &&
//		v.Owner == caller &&
//		caller_affiliation == MANUFACTURER &&
//		v.Scrapped == false {
//
//		v.Model = new_value
//
//	} else {
//		return nil, errors.New(fmt.Sprint("Permission denied. update_model %t %t %t"+v.Owner == caller, caller_affiliation == MANUFACTURER, v.Scrapped))
//
//	}
//
//	_, err := t.save_changes(stub, v)
//
//	if err != nil {
//		fmt.Printf("UPDATE_MODEL: Error saving changes: %s", err)
//		return nil, errors.New("Error saving changes")
//	}
//
//	return nil, nil
//
//}

//=================================================================================================================================
//	 Song Obsolete - not sure if we need this function. I have just implemented it if we want to make songs obsolete
//=================================================================================================================================
func (t *SimpleChaincode) song_obsolete(stub shim.ChaincodeStubInterface, s Song, caller string, caller_affiliation string) ([]byte, error) {

	if s.Obsolete != true {

		return nil, errors.New("Cannot make song obsolete")
	} else {
		_, err := t.save_changes(stub, s)

		if err != nil {
			fmt.Printf("SONG_OBSOLETE: Error saving changes: %s", err)
			return nil, errors.New("SONG_OBSOLETE: Error saving changes")
		} else {
			return nil, nil

		}
	}
}

//=================================================================================================================================
//	 Read Functions
//=================================================================================================================================
//	 get_song_details - Returns details of a song
//=================================================================================================================================
func (t *SimpleChaincode) get_song_details(stub shim.ChaincodeStubInterface, s Song, caller string, caller_affiliation string) ([]byte, error) {

	bytes, err := json.Marshal(s)

	if err != nil {
		return nil, errors.New("GET_SONG_DETAILS: Invalid song object")
	} else {
		return bytes, nil
	}

}

//=================================================================================================================================
//	 get_songs -  Returns all songs and details
//=================================================================================================================================

func (t *SimpleChaincode) get_songs(stub shim.ChaincodeStubInterface, caller string, caller_affiliation string) ([]byte, error) {
	bytes, err := stub.GetState("Song_IDs")

	if err != nil {
		return nil, errors.New("Unable to get Song_IDs")
	}

	var Song_IDs Song_Holder

	err = json.Unmarshal(bytes, &Song_IDs)

	if err != nil {
		return nil, errors.New("Corrupt Song")
	}

	result := "["

	var temp []byte
	var s Song

	for _, Song_ID := range Song_IDs.Song_IDs {

		s, err = t.retrieve_Song_ID(stub, Song_ID)

		if err != nil {
			return nil, errors.New("Failed to retrieve Song")
		}

		temp, err = t.get_song_details(stub, s, caller, caller_affiliation)

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
//	 check_unique_Song_ID - Song ID must be unique in the BlockChain when a new song is created
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
//	 Get_rating - Read the rating that was done by a user for a certain song
//=================================================================================================================================
func (t *SimpleChaincode) get_rating(stub shim.ChaincodeStubInterface, Song_ID string, caller string, caller_affiliation string) ([]byte, error) {
	// to be implemented
	return nil, nil
}

//=================================================================================================================================
//	 Get_overall_rating - Calculate the average rating of a song
//=================================================================================================================================
func (t *SimpleChaincode) get_overall_rating(stub shim.ChaincodeStubInterface, Song_ID string, caller string, caller_affiliation string) ([]byte, error) {
	// to be implemented
	return nil, nil
}

//=================================================================================================================================
//	 Get_contract - If a contract was provided, it can be shown to the singer or CA
//=================================================================================================================================
func (t *SimpleChaincode) get_contract(stub shim.ChaincodeStubInterface, Song_ID string, caller string, caller_affiliation string) ([]byte, error) {
	// to be implemented
	return nil, nil
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
