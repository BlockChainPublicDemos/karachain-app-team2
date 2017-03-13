package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	//"regexp"
	"strconv"
	//"strings"
	//"time"
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
const COPYRIGHT_AUTHORITY = "copyright_authority"

//==============================================================================================================================
//	 Status types - Asset lifecycle is broken down into z statuses, this is part of the business logic to determine what can
//					be done to the song at points in it's lifecycle
// Not sure if we need them
//==============================================================================================================================
const STATE_TEMPLATE = "0"
const STATE_RECORDED = "1"
const STATE_UPDATED = "2"
const STATE_VOTED = "3"
const STATE_CONTRACT_PROVIDED = "4"
const STATE_CONTRACT_REJECTED = "5"
const STATE_CONTRACT_ACCEPTED = "6"
const STATE_OBSOLETE = "7"

//==============================================================================================================================
//	 Structure Definitions
//==============================================================================================================================
//	Chaincode - A blank struct for use with Shim (A HyperLedger included go file used for get/put state
//				and other HyperLedger functions)
//==============================================================================================================================
type SimpleChaincode struct {
}

//Global keys
var SongKey = "_allsongsindex"

//var SingerKey = "_allsingerindex"

//==============================================================================================================================
//	Song - Defines the structure for a Song object. JSON on right tells it what JSON fields to map to
//			  that element when reading a JSON object into the struct e.g. JSON make -> Struct Make.
// 			Key for that structure in the BlockChain is the Song_ID
//==============================================================================================================================
type Song struct {
	Song_ID      string `json:"Song_ID"`
	Date_created string `json:"Date_created"`
	//	SmartContract_Unique_ID    string   `json:"SmartContract_Unique_ID"`
	Singer_Id          string `json:"Singer_Id"`
	Singer_Name        string `json:"Singer_Name"`
	Video_Id           string `json:"Video_Id"`
	Owner              string `json:"Owner"`
	Video_Link         string `json:"Video_Link"`
	Video_date_created string `json:"Video_date_created"`
	Video_QR_code_Id   string `json:"Video_QR_code_Id"`
	//	Copyright_Id               string   `json:"Copyright_Id"`
	//	Copyright_date_created     string   `json:"Copyright_date_created"`
	//	Copyright_date_accepted    string   `json:"Copyright_date_accepted"`
	//	Copyright_date_rejected    string   `json:"Copyright_date_rejected"`
	//	Copyright_Institution_Id   string   `json:"Copyright_Institution_Id"`
	//	Copyright_Institution_Name string   `json:"Copyright_Institution_Name"`
	//	Copyright_State            string   `json:"Copyright_State"`
	Venue_Id   string `json:"Venue_Id"`
	Venue_Name string `json:"Venue_Name"`
	//User_Id     []string `json:"User_Id"`
	//User_role   []string `json:"User_role"`
	//User_rating []string `json:"User_rating"`
	User_rating map[string]int
	Obsolete    bool    `json:"Obsolete"`
	Status      string  `json:"Status"`
	Song_Name   string  `json:"Song_Name"`
	AVG_Rating  float32 `json:"AVG_Rating"`
	//	Contract_date_from         string   `json:"Contract_date_from"`
	//	Contract_date_to           string   `json:"Contract_date_to"`
}

//==============================================================================================================================
//	Song Holder - Defines the structure that holds all the Song_IDs for songs that have been created.
//				Used as an index when querying all songs. Key in the BlockChain for that structure is the Singer ID
//==============================================================================================================================

type Song_Holder struct {
	Songs []Song_basics
	//	Song_IDs        []string `json:"Song_IDs"`
	//	Song_AVG_Rating []string `json:"Song_AVG_Rating"`
	//	Song_Name       []string `json:"Song_Name"`
	//	Singer_Name     []string `json:"Singer_Name"`
	// Songs []Song
}

//==============================================================================================================================
//	Song basic info - Contains basic song information and can be stored in the song holder
//==============================================================================================================================

type Song_basics struct {
	Song_ID         string `json:"Song_ID"`
	Song_AVG_Rating string `json:"Song_AVG_Rating"`
	Song_Name       string `json:"Song_Name"`
	Singer_Name     string `json:"Singer_Name"`
	// Songs []Song
}

//==============================================================================================================================
//	Contract Holder - Defines the structure that holds all the Contracts for a singer that have been created.
//				Used as an index when querying all contracts per Singer
//==============================================================================================================================

type Contract_holder struct {
	Contracts []Contract
}

//==============================================================================================================================
//	Contract  - Defines the structure for a single contract. Can be added to the contract holder
//==============================================================================================================================

type Contract struct {
	Copyright_Id               string `json:"Copyright_Ids"`
	Copyright_date_created     string `json:"Copyright_date_created"`
	Copyright_date_accepted    string `json:"Copyright_date_accepted"`
	Copyright_date_rejected    string `json:"Copyright_date_rejected"`
	Copyright_Institution_Id   string `json:"Copyright_Institution_Id"`
	Copyright_Institution_Name string `json:"Copyright_Institution_Name"`
	Copyright_State            string `json:"Copyright_State"`
	Contract_date_from         string `json:"Contract_date_from"`
	Contract_date_to           string `json:"Contract_date_to"`
	SmartContract_ID           string `json:"SmartContract_ID"`
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

	var err error
	fmt.Printf("INIT: Karachain function: %s ", function)
	//test chaincode
	// Initialize the chaincode
	//	Aval, err = strconv.Atoi(args[0])
	//	if err != nil {
	//		return nil, errors.New("Expecting integer value for asset holding")
	//	}

	// Write the state to the ledger
	err = stub.PutState("karachain", []byte(args[0])) //making a test var "karachain", I find it handy to read/write to it right away to test the network
	if err != nil {
		return nil, err
	}

	//

	var Song_IDs Song_Holder

	bytes, err := json.Marshal(Song_IDs)

	if err != nil {
		return nil, errors.New("Error creating initial song placeholders")
	}

	err = stub.PutState(SongKey, bytes)
	if err != nil {
		return nil, err
	}
	return nil, nil

	//	var Contracts Contract_Holder
	//
	//	bytes_c, err := json.Marshal(Contracts)
	//
	//	if err != nil {
	//		return nil, errors.New("Error creating initial Contract placeholders")
	//	}
	//
	//	err = stub.PutState(SingerKey, bytes_c)
	//	if err != nil {
	//		return nil, err
	//	}
	//	return nil, nil

	//
	//	for i := 0; i < len(args); i = i + 2 {
	//		t.add_ecert(stub, args[i], args[i+1])
	//	}
	//
	//	fmt.Printf("INII: Karachain exit")
	//	return nil, nil
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
//	 song_constructor - Creates a song structure and assigns values to it.
//==============================================================================================================================
func (t *SimpleChaincode) song_constructur(stub shim.ChaincodeStubInterface, caller string, caller_affiliation string, args []string) (Song, error) {

	var s Song
	s.User_rating = make(map[string]int)
	if args[0] == "" {
		fmt.Printf("CREATE_SONG: Invalid Song_ID provided")
		return s, errors.New("Invalid Song_ID provided")
	}

	if len(args) == 11 {
		s.Song_ID = args[0]
		s.Date_created = args[1]
		s.Singer_Id = args[8]
		s.Video_Id = args[2]
		s.Video_Link = args[3]
		s.Video_date_created = args[4]
		s.Video_QR_code_Id = args[5]
		s.Singer_Name = args[9]
		s.Venue_Id = args[6]
		s.Venue_Name = args[7]
		//		s.User_Id = "UNDEFINED"
		//		s.User_role = "UNDEFINED"
		s.Obsolete = false
		s.Status = "UNDEFINED"
		s.Song_Name = args[10]
		s.AVG_Rating = 0.0
	} else {
		return s, errors.New("Not enough attributes to create a song")
	}

	return s, nil
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
		fmt.Printf(" SAVE_CHANGES: Error storing song record: %s", err)
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

	//	for i := 0; i < len(args); i = i + 2 {
	//		t.add_ecert(stub, args[i], args[i+1])
	//	}
	fmt.Printf("INVOKE: Karachain function: %s ", function)

	if err != nil {
		fmt.Printf("QUERY: Error retrieving caller details", err)
		//	return nil, errors.New("Error retrieving caller information")
	}

	if function == "create_song" { // we create a song from scratch
		return t.create_song(stub, caller, caller_affiliation, args)
	} else if function == "ping" {
		return t.ping(stub)
	} else if function == "Set_Contract" { // Only the copyright authority is allowed to set a contract
		return t.set_contract(stub, caller, caller_affiliation, args)
	} else if function == "Set_Contract_Response" { // Function my only be called by the singer
		return t.set_contract_response(stub, caller, caller_affiliation, args)
	} else { // If the function is not a create then there must be a Song so we need to retrieve the Song.

		Song_ID := args[0]

		s, err := t.retrieve_Song_ID(stub, Song_ID)

		if err != nil {
			fmt.Printf("INVOKE: Error retrieving Song: %s", err)
			return nil, errors.New("Error retrieving Song")
		}

		if function == "Set_Rating" { // Rating can be set by anybody, but only once and not by the singer
			return t.set_rating(stub, s, caller, caller_affiliation, args)
		} else if function == "update_song" { // Function may only be called by the singer
			return t.update_song(stub, s, caller, caller_affiliation, args)
		} else if function == "set_obsolete" { // Function may only be called by the singer
			return t.set_obsolete(stub, s, caller, caller_affiliation)
		}
		//		else if function == "set_obsolete" { // Function may only be called by the singer
		//			return t.set_obsolete(stub, s, caller, caller_affiliation, args[0])
		//		}

		return nil, errors.New("Function of the name " + function + " doesn't exist.")

	}
}

//=================================================================================================================================
//	Query - Called on chaincode query. Takes a function name passed and calls that function. Passes the
//  		initial arguments passed are passed on to the called function.
//=================================================================================================================================
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	//TODO add in authentication and certificate management
	caller, caller_affiliation, err := t.get_caller_data(stub)
	if err != nil {
		fmt.Printf("QUERY: Error retrieving caller details", err)
		//			return nil, errors.New("QUERY: Error retrieving caller details: " + err.Error())
	}

	fmt.Println("query is running " + function)

	logger.Debug("function: ", function)
	logger.Debug("caller: ", caller)
	logger.Debug("affiliation: ", caller_affiliation)

	/**TODO  leave out for now */
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
		return t.get_rating(stub, caller, caller_affiliation, args) // A user should be able to get his own rating that was made in the past for a particular song
	} else if function == "Get_Contract" { // Only allowed for singer or copyright authority to see the latest contract
		return t.get_contract(stub, caller, caller_affiliation, args)
	} else if function == "Get_Contracts" { // Only allowed for singer or copyright authority to see the latest contract
		return t.get_contracts(stub, args[0], caller, caller_affiliation)
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
func (t *SimpleChaincode) create_song(stub shim.ChaincodeStubInterface, caller string, caller_affiliation string, args []string) ([]byte, error) {
	var s Song
	//Song_ID := "\"Song_ID\":\"UNDEFINED\"" // Variables to define the JSON
	//	Song_ID := "\"Song_ID\":\"" + Song_ID_r + "\", " // Variables to define the JSON
	//	Date_created := "\"Date_created\":\"UNDEFINED\""
	//	Song_ID := "\"Song_ID\":\"" + args[0] + "\", " // Variables to define the JSON
	//	Date_created := "\"Date_created\":\"" + args[1] + "\", "
	//	SmartContract_Unique_ID := "\"SmartContract_Unique_ID\":\"UNDEFINED\", "
	//	Singer_Id := "\"Singer_Id\":\"" + args[8] + "\", "
	//	Singer_Name := "\"Singer_Name\":\"UNDEFINED\", "
	//	Video_Id := "\"Video_Id\":\"" + args[2] + "\", "
	//	Owner := "\"Owner\":\"UNDEFINED\", "
	//	Video_Link := "\"Video_Link\":\"" + args[3] + "\", "
	//	Video_date_created := "\"Video_date_created\":\"" + args[4] + "\", "
	//	Video_QR_code_Id := "\"Video_QR_code_Id\":\"" + args[5] + "\", "
	//	//	Copyright_Id := "\"Copyright_Id\":\"UNDEFINED\", "
	//	//	Copyright_date_created := "\"Copyright_date_created\":\"UNDEFINED\", "
	//	//	Copyright_date_accepted := "\"Copyright_date_accepted\":\"UNDEFINED\", "
	//	//	Copyright_date_rejected := "\"Copyright_date_rejected\":\"UNDEFINED\", "
	//	//	Copyright_Institution_Id := "\"Copyright_Institution_Id\":\"UNDEFINED\", "
	//	//	Copyright_Institution_Name := "\"Copyright_Institution_Name\":\"UNDEFINED\", "
	//	//	Copyright_State := "\"Copyright_State\":\"UNDEFINED\", "
	//	Venue_Id := "\"Venue_Id\":\"" + args[6] + "\", "
	//	Venue_Name := "\"Venue_Name\":\"" + args[7] + "\", "
	//	User_Id := "\"User_Id\":\"UNDEFINED\", "
	//	User_role := "\"User_role\":\"UNDEFINED\", "
	//	User_rating := "\"User_rating\":\"UNDEFINED\", "
	//	Obsolete := "\"Obsolete\":false, "
	//	Status := "\"Status\":\"UNDEFINED\""
	//	Song_Name := "\"Song_Name\":\"" + args[8] + "\", "
	//	AVG_Rating := "\"AVG_Rating\":\"UNDEFINED\""

	//fmt.Printf("Song ID is %s", args[0])
	//Song_json := "{" + Song_ID + Date_created + SmartContract_Unique_ID + Singer_Id + Singer_Name + Video_Id + Owner + Video_Link + Video_date_created + Video_QR_code_Id + Copyright_Id + Copyright_date_created + Copyright_date_accepted + Copyright_date_rejected + Copyright_Institution_Id + Copyright_Institution_Name + Copyright_State + Venue_Id + Venue_Name + User_Id + User_role + User_rating + Obsolete + Status + "}" // Concatenates the variables to create the total JSON object
	//Song_json := "{" + Song_ID + Date_created + SmartContract_Unique_ID + Singer_Id + Singer_Name + Video_Id + Owner + Video_Link + Video_date_created + Video_QR_code_Id + Venue_Id + Venue_Name + User_Id + User_role + User_rating + Obsolete + Status + "}" // Concatenates the variables to create the total JSON object
	// No need for a song criteria so far
	//	_, err := regexp.Match("^[A-z][A-z][0-9]{7}", []byte(args[0])) // matched = true if the Song ID passed fits format of two letters followed by seven digits
	//
	//	if err != nil {
	//		fmt.Printf("CREATE_Song: Invalid Song_ID: %s", err)
	//		return nil, errors.New("Invalid Song_ID")
	//	}

	//	err = json.Unmarshal([]byte(Song_json), &s) // Convert the JSON defined above into a Song object for go
	//
	//
	//
	//	if err != nil {
	//		error_s := err.Error()
	//		fmt.Printf("Invalid JSON object. Error is : %s", err)
	//		return nil, errors.New(error_s)
	//	}

	s, err := t.song_constructur(stub, caller, caller_affiliation, args)
	record, err := stub.GetState(s.Song_ID) // If not an error then a record exists so cant create a new car with this Song_ID as it must be unique

	if record != nil {
		return nil, errors.New("Song already exists")
	}

	if caller_affiliation != SINGER { // Only the singer can create a new Song
		fmt.Printf("CREATE_SONG: Invalid Song_ID provided")
		//		return nil, errors.New(fmt.Sprintf("Permission Denied. create_song. %v === %v", caller_affiliation, SINGER))
	}

	//saves the song as a unique object in the ledger identified by song id (putstate)
	_, err = t.save_changes(stub, s)

	if err != nil {
		fmt.Printf("CREATE_SONG: Error saving changes: %s", err)
		return nil, errors.New("Error saving changes")
	}

	//gets the structure that contains an array of songs  ..
	bytes, err := stub.GetState(SongKey)

	if err != nil {
		return nil, errors.New("Unable to get Song_ID")
	}

	var Songs Song_Holder // Hold an array of song IDs
	var Song Song_basics

	err = json.Unmarshal(bytes, &Songs)

	if err != nil {
		return nil, errors.New("Corrupt Song_Holder record")
	}
	//need to add the new song to this container of songs (which is either a list of all songs or for a specific singer)
	//	Song_IDs.Song_IDs = append(Song_IDs.Song_IDs, s.Song_ID)
	//	Song_IDs.Singer_Name = append(Song_IDs.Singer_Name, s.Singer_Name)
	//	Song_IDs.Song_AVG_Rating = append(Song_IDs.Song_AVG_Rating, "")
	//	Song_IDs.Song_Name = append(Song_IDs.Song_Name, s.Song_Name)
	Song.Singer_Name = s.Singer_Name
	Song.Song_AVG_Rating = ""
	Song.Song_Name = s.Song_Name
	Song.Song_ID = s.Song_ID
	Songs.Songs = append(Songs.Songs, Song)

	bytes, err = json.Marshal(Songs)

	if err != nil {
		fmt.Print("Error creating Song_Holder record")
	}

	err = stub.PutState(SongKey, bytes)

	if err != nil {
		return nil, errors.New("Unable to put the state of Song_Holder")
	}

	//gets the structure that contains an array of contracts

	// Here we create our composite keys

	//	qr_indexName := "QR~name"
	//
	//	qrIndexKey, err := stub.CreateCompositeKey(qr_indexName, []string{s.Video_QR_code_Id, s.Song_ID})
	//	if err != nil {
	//		return shim.Error(err.Error())
	//	}
	//	//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the marble.
	//	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	//	value := []byte{0x00}
	//	stub.PutState(qrIndexKey, value)
	//
	//	singer_indexName := "Singer~name"
	//
	//	singerIndexKey, err := stub.CreateCompositeKey(singer_indexName, []string{s.Singer_Id, s.Song_ID})
	//	if err != nil {
	//		return shim.Error(err.Error())
	//	}
	//	//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the marble.
	//	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	//	// value := []byte{0x00}
	//	stub.PutState(singerIndexKey, value)
	//
	//	venue_indexName := "Venue~name"
	//
	//	venueIndexKey, err := stub.CreateCompositeKey(venue_indexName, []string{s.Venue_Id, s.Song_ID})
	//	if err != nil {
	//		return shim.Error(err.Error())
	//	}
	//	//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the marble.
	//	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	//	// value := []byte{0x00}
	//	stub.PutState(venueIndexKey, value)

	return nil, nil

}

//=================================================================================================================================
//	 Update Functions
//=================================================================================================================================
//	 update_song
//=================================================================================================================================
func (t *SimpleChaincode) update_song(stub shim.ChaincodeStubInterface, s Song, caller string, caller_affiliation string, args []string) ([]byte, error) {

	var err error
	if s.Obsolete == true {

		return nil, errors.New(fmt.Sprintf("Song is obsolete and cannot be updated."))

	}

	if len(args) == 11 {
		//s.Song_ID = args[0]
		s.Date_created = args[1]
		s.Singer_Id = args[8]
		s.Video_Id = args[2]
		s.Video_Link = args[3]
		s.Video_date_created = args[4]
		s.Video_QR_code_Id = args[5]
		s.Singer_Name = args[9]
		s.Venue_Id = args[6]
		s.Venue_Name = args[7]
		//		s.User_Id = "UNDEFINED"
		//		s.User_role = "UNDEFINED"
		s.Obsolete = false
		//s.Status = "UNDEFINED"
		s.Song_Name = args[10]
		//s.AVG_Rating = "UNDEFINED"
	} else {
		return nil, errors.New("Not enough attributes to create a song")
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
func (t *SimpleChaincode) set_rating(stub shim.ChaincodeStubInterface, s Song, caller string, caller_affiliation string, args []string) ([]byte, error) {
	// Song ID is args[0]
	//var User_rating int = int(args[1])
	var User_Id string = args[2]
	User_rating, err := strconv.Atoi(args[1])
	if err != nil {
		// handle error
		return nil, err
	}
	//	// User = caller
	//	var User_role string = AUDIENCE
	//	var new_user_ratings []string
	//	match := false
	//	new_index := 0
	var sum int
	if s.Obsolete != true {

		s.User_rating[User_Id] = User_rating
		for _, value := range s.User_rating {
			sum = sum + value
		}
		s.AVG_Rating = float32((float32(sum)) / float32(len(s.User_rating)))
		//		for index, user := range s.User_Id {
		//			if user == User_Id {
		//
		//			}
		//
		//		}

	} else {

		return nil, errors.New(fmt.Sprint("Permission denied to set a rating."))

	}

	_, err2 := t.save_changes(stub, s)

	if err2 != nil {
		fmt.Printf("UPDATE_MAKE: Error saving changes: %s", err2)
		return nil, errors.New("Error saving changes")
	}

	// Here we create our composite keys

	//	user_rating_indexName := "user_rating~name"
	//
	//	user_rating_IndexKey, err := stub.CreateCompositeKey(user_rating_indexName, []string{s.User_Id, s.Song_ID})
	//	if err != nil {
	//		return shim.Error(err.Error())
	//	}
	//	//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the marble.
	//	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	//	value := []byte{0x00}
	//	stub.PutState(user_rating_IndexKey, value)

	return nil, nil

}

//=================================================================================================================================
//	 set_contract - The CA may provide a contract to a singer for a particular song. Only 1 contract allowed per copyright authority
//=================================================================================================================================
func (t *SimpleChaincode) set_contract(stub shim.ChaincodeStubInterface, caller string, caller_affiliation string, args []string) ([]byte, error) {

	if len(args) != 9 {
		return nil, errors.New("Not enough arguments passed to function. Cannot store contract")
	}

	Singer_ID := args[5]
	Copyright_Id := args[1]
	Copyright_date_created := args[2]
	Copyright_Institution_Id := args[3]
	Copyright_Institution_Name := args[4]
	Contract_date_from := args[6]
	Contract_date_to := args[7]
	SmartContract_ID := args[8]

	if Contract_date_from < Contract_date_to {
		return nil, errors.New("Contract Start date after contract end date. Cannot store contract")
	}
	//	from_date := strings.Split(Contract_date_from, ".")
	//	to_date := strings.Split(Contract_date_to, ".")
	var c Contract
	c.Contract_date_from = Contract_date_from
	c.Contract_date_to = Contract_date_to
	c.Copyright_Id = Copyright_Id
	c.Copyright_date_created = Copyright_date_created
	c.Copyright_Institution_Id = Copyright_Institution_Id
	c.Copyright_Institution_Name = Copyright_Institution_Name
	c.SmartContract_ID = SmartContract_ID

	//	if len(from_date) != 3 || len(to_date) != 3 {
	//		return nil, errors.New("Wrong date format. Cannot store contract")
	//	}

	//	from_year, err := strconv.Atoi(from_date[2])
	//	from_month, err := strconv.Atoi(from_date[1])
	//	from_day, err := strconv.Atoi(from_date[0])
	//
	//	to_year, err := strconv.Atoi(to_date[2])
	//	to_month, err := strconv.Atoi(to_date[1])
	//	to_day, err := strconv.Atoi(to_date[0])
	//
	//	from_time := time.Date(from_year, from_month, from_day, 0, 0, 0, 0, time.UTC)
	//	to_time := time.Date(to_year, to_month, to_day, 0, 0, 0, 0, time.UTC)

	var contracts Contract_holder

	bytes, err := stub.GetState(Singer_ID)

	if err != nil {
		// Here we already have a contract and need to check if the new one is valid and can be added
		err = json.Unmarshal(bytes, &contracts)

		if err != nil {
			return nil, errors.New("Error getting Contract Holder")
		}
		// We check if there is a collision with a contract that was already signed
		//		for _, current_c := range contracts.Contracts {
		//
		//			temp_from_date := strings.Split(current_c.Contract_date_from, ".")
		//			temp_to_date := strings.Split(current_c.Contract_date_to, ".")
		//			temp_from_year, err := strconv.Atoi(temp_from_date[2])
		//			temp_from_month, err := strconv.Atoi(temp_from_date[1])
		//			temp_from_day, err := strconv.Atoi(temp_from_date[0])
		//
		//			temp_to_year, err := strconv.Atoi(temp_to_date[2])
		//			temp_to_month, err := strconv.Atoi(temp_to_date[1])
		//			temp_to_day, err := strconv.Atoi(temp_to_date[0])
		//			temp_from_time := time.Date(temp_from_year, temp_from_month, temp_from_day, 0, 0, 0, 0, time.UTC)
		//			temp_to_time := time.Date(temp_to_year, temp_to_month, temp_to_day, 0, 0, 0, 0, time.UTC)
		//
		//			if current_c.Copyright_State == STATE_CONTRACT_ACCEPTED {
		//				//if current_c.Contract_date_from < c.Contract_date_from && current_c.Contract_date_to > c.Contract_date_from {
		//				if temp_from_time < from_time && temp_to_time > from_time {
		//					return nil, errors.New("Contract already accepeted in that date range. Contract ends in a period within an active contract")
		//				}
		//
		//				if temp_from_time > from_time && temp_from_time < to_time {
		//					return nil, errors.New("Contract already accepeted in that date range. Contract starts in a period within an active contract")
		//				}
		//
		//			}
		//
		//		}
		contracts.Contracts = append(contracts.Contracts, c)
		bytes, err = json.Marshal(contracts)
		err = stub.PutState(Singer_ID, bytes)
		if err != nil {
			return nil, err
		}
	} else {
		contracts.Contracts = append(contracts.Contracts, c)
		bytes, err = json.Marshal(contracts)
		err = stub.PutState(Singer_ID, bytes)
		if err != nil {
			return nil, err
		}

	}

	//	if s.Obsolete != true {
	//
	//		s.Copyright_Id = Copyright_Id
	//		s.Copyright_date_created = Copyright_date_created
	//		s.Copyright_Institution_Id = Copyright_Institution_Id
	//		s.Copyright_Institution_Name = Copyright_Institution_Name
	//		s.Copyright_State = STATE_CONTRACT_PROVIDED
	//		s.SmartContract_Unique_ID = SmartContract_Unique_ID
	//	} else {
	//
	//		return nil, errors.New(fmt.Sprint("Permission denied to set a contract."))
	//
	//	}
	//
	//	_, err := t.save_changes(stub, s)
	//
	//	if err != nil {
	//		fmt.Printf("UPDATE_MAKE: Error saving changes: %s", err)
	//		return nil, errors.New("Error saving changes")
	//	}

	// Here we create our composite keys

	//	copyright_indexName := "copyright~name"
	//
	//	copyrightIndexKey, err := stub.CreateCompositeKey(copyright_indexName, []string{s.Copyright_Institution_Id, s.Song_ID})
	//	if err != nil {
	//		return shim.Error(err.Error())
	//	}
	//	//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the marble.
	//	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	//	value := []byte{0x00}
	//	stub.PutState(copyrightIndexKey, value)
	//
	//	return nil, nil

	return nil, nil

}

//=================================================================================================================================
//	 set_contract_response - The singer can either accept or reject a contract.
//=================================================================================================================================
func (t *SimpleChaincode) set_contract_response(stub shim.ChaincodeStubInterface, caller string, caller_affiliation string, args []string) ([]byte, error) {

	if len(args) != 4 {
		return nil, errors.New("Not enough arguments passed to function. Cannot store contract")
	}

	Singer_ID := args[0]
	//	Copyright_Id := args[1]
	//	Copyright_date_created := args[2]
	//	Copyright_Institution_Id := args[3]
	//	Copyright_Institution_Name := args[4]
	//	Contract_date_from := args[6]
	//	Contract_date_to := args[7]
	SmartContract_ID := args[1]

	Copyright_decision := args[2]
	Copyright_date_decision := args[3]
	match := false

	//	var c Contract
	var contracts Contract_holder
	var new_contract_h Contract_holder

	bytes, err := stub.GetState(Singer_ID)

	if err == nil {
		// Here we already have a contract and need to check if the new one is valid and can be added
		err = json.Unmarshal(bytes, &contracts)

		if err != nil {
			return nil, errors.New("Error getting Contract Holder")
		}
		// We check if there is a collision with a contract that was already signed
		for _, current_c := range contracts.Contracts {

			if current_c.SmartContract_ID == SmartContract_ID { // we found the contract
				if current_c.Copyright_State != STATE_CONTRACT_ACCEPTED {
					if Copyright_decision == "true" {
						current_c.Copyright_State = STATE_CONTRACT_ACCEPTED
						current_c.Copyright_date_accepted = Copyright_date_decision
						match = true
					}

					if Copyright_decision == "false" {
						current_c.Copyright_date_rejected = Copyright_date_decision
						current_c.Copyright_State = STATE_CONTRACT_REJECTED
						match = true
					}
				}

			}
			var new_contract Contract
			new_contract = current_c
			new_contract_h.Contracts = append(new_contract_h.Contracts, new_contract)

		}

		if match == true {

			bytes, err = json.Marshal(new_contract_h)
			err = stub.PutState(Singer_ID, bytes)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, errors.New("Not able to update a contract. Either wrong Smartcontract ID or decision")
		}

	} else {
		return nil, errors.New("Error unmarshaling")
	}

	return nil, nil
}

//=================================================================================================================================
//	 Song Obsolete - not sure if we need this function. I have just implemented it if we want to make songs obsolete
//=================================================================================================================================
func (t *SimpleChaincode) set_obsolete(stub shim.ChaincodeStubInterface, s Song, caller string, caller_affiliation string) ([]byte, error) {

	if s.Obsolete == false {

		return nil, errors.New("Cannot make song obsolete. Song is already obsolete")
	} else {
		s.Obsolete = false
		_, err := t.save_changes(stub, s)

		if err != nil {
			fmt.Printf("SONG_OBSOLETE: Error saving changes: %s", err)
			return nil, errors.New("SONG_OBSOLETE: Error saving changes")
		} else {
			return nil, nil

		}
	}

	return nil, nil
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
	//Get the list of all the song IDs
	bytes, err := stub.GetState(SongKey)

	if err != nil {
		return nil, errors.New("Unable to get Songs")
	}
	return bytes, nil

	//	var Song_IDs Song_Holder
	//
	//	err = json.Unmarshal(bytes, &Song_IDs)
	//
	//	if err != nil {
	//		return nil, errors.New("Corrupt Song")
	//	}

	//	result := "["
	//

	//	var temp []byte
	//	var s Song
	//loop through song IDs and get the song structures from the ledger
	//	for index, song := range Song_IDs.Songs {
	//
	//		//s, err = t.retrieve_Song_ID(stub, songId)
	//
	//		//		if err != nil {
	//		//			return nil, errors.New("Failed to retrieve Song")
	//		//		}
	//
	//		//temp, err = t.get_song_details(stub, s, caller, caller_affiliation)
	//
	//		//		if err == nil {
	//		//
	//		//		}
	//
	//		result += string(temp) + ","
	//	}
	//
	//	if len(result) == 1 {
	//		result = "[]"
	//	} else {
	//		result = result[:len(result)-1] + "]"
	//	}

	//	return []byte(result), nil

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
func (t *SimpleChaincode) get_rating(stub shim.ChaincodeStubInterface, caller string, caller_affiliation string, args []string) ([]byte, error) {
	// to be implemented
	var Song_ID = args[0]
	var User_Id string = args[1]

	s, err := t.retrieve_Song_ID(stub, Song_ID)

	if err != nil {
		fmt.Printf("INVOKE: Error retrieving Song: %s", err)
		return nil, errors.New("Error retrieving Song")
	}
	if s.Obsolete != true {
		rating := s.User_rating[User_Id]
		if rating == 0 {
			return nil, errors.New("No rating found")
		} else {
			bytes, err2 := json.Marshal(rating)
			if err2 != nil {
				return nil, errors.New("Error marshaling rating")
			}
			return bytes, nil
		}

	} else {

		return nil, errors.New(fmt.Sprint("Permission denied, song is obslete."))

	}

	return nil, nil
}

//=================================================================================================================================
//	 Get_overall_rating - Returns the average rating of a song
//=================================================================================================================================
func (t *SimpleChaincode) get_overall_rating(stub shim.ChaincodeStubInterface, Song_ID string, caller string, caller_affiliation string) ([]byte, error) {

	s, err := t.retrieve_Song_ID(stub, Song_ID)

	if err != nil {
		fmt.Printf("INVOKE: Error retrieving Song: %s", err)
		return nil, errors.New("Error retrieving Song")
	}
	if s.Obsolete != true {
		overall_rating := s.AVG_Rating

		bytes, err2 := json.Marshal(overall_rating)
		if err2 != nil {
			return nil, errors.New("Error marshaling rating")
		}
		return bytes, nil

	} else {

		return nil, errors.New(fmt.Sprint("Permission denied, song is obslete."))

	}

	return nil, nil
}

//=================================================================================================================================
//	 Get_contract - Provides a specific contract for a specific singer
//=================================================================================================================================
func (t *SimpleChaincode) get_contract(stub shim.ChaincodeStubInterface, caller string, caller_affiliation string, args []string) ([]byte, error) {
	Singer_ID := args[0]
	SmartContract_ID := args[1]
	bytes, err := stub.GetState(Singer_ID)
	var contracts Contract_holder

	if err != nil {
		return nil, errors.New("Unable to get contracts for singer")
	}
	err = json.Unmarshal(bytes, &contracts)

	for _, current_c := range contracts.Contracts {

		if current_c.SmartContract_ID == SmartContract_ID {
			bytes, err := json.Marshal(current_c)
			if err != nil {
				return nil, errors.New("Unable to marshal contract")
			}
			return bytes, nil

		}

	}
	return nil, errors.New("Unable to get contracts for singer")

}

//=================================================================================================================================
//	 Get_contracts - Provides all contracts of a singer
//=================================================================================================================================
func (t *SimpleChaincode) get_contracts(stub shim.ChaincodeStubInterface, Singer_ID string, caller string, caller_affiliation string) ([]byte, error) {
	// to be implemented

	bytes, err := stub.GetState(Singer_ID)

	if err != nil {
		return nil, errors.New("Unable to get contracts for singer")
	}
	return bytes, nil
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
