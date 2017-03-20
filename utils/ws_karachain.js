// ==================================
// karachain - incoming messages, look for type
// ==================================
var ibc = {};
var chaincode = {};
var ws ={};
var qr ={};
var async = require('async');
var performances = [];
var lastSongId = "kc846908"; //placeholder value
var lastSingerId = "user_type1_1";
var lastVisitorId = "user_type2_0";
var lastEvtMgrId = "user_type4_0";
var qrMap = {};
/**
 * admine0e2435d74
WebAppAdmin502f66b484
user_type1_0a01996bb60 //type1_0 = admin
user_type1_16ee2c832e7//singer
user_type1_24035add836
user_type1_31691188bc1
user_type1_4590e3fec92
user_type2_06a80515894//type2 = visitors
user_type2_1d8d80a9f81
user_type2_21b5feb25a0
user_type2_3fd60fb13e3
user_type2_4e5d0c330a6
user_type4_0d0bd04950e //type4 = evt mgrs
user_type4_12f176b4351
user_type4_20903dc8ee1
user_type4_3cb769ac24b
user_type4_492f82c1c16
user_type8_04a8d6f6a59
user_type8_18d582015f2
user_type8_2c0b47a9471
user_type8_3973dcda0b8
user_type8_4ae3d1e2132
 */
function genQRpng(singerName, perfName,singerId, perfId,perfDate){
	var qrstring = "{Singer Name:"+singerName+"Performance Name:"+perfName+"Singer ID:"+singerId+"Performance Date:"+perfDate+"}";
	var qr_png = qr.image(qrstring, { type: 'png' });
	return qr_png;
}
module.exports.setup = function(sdk, cc, qrsvc){
	console.log("karachain setup");
	ibc = sdk;
	chaincode = cc;
	qr = qrsvc;
	
};
module.exports.genQRcode = function(singerName, perfName,singerId, perfId,perfDate){
	return genQRpng(singerName, perfName,singerId, perfId,perfDate);

	
};
/*Go code
 * if function == "create_song" { // we create a song from scratch
		return t.create_song(stub, caller, caller_affiliation, args)
	} else if function == "ping" {
		return t.ping(stub)
	} else { // If the function is not a create then there must be a Song so we need to retrieve the Song.

		Song_ID := args[0]

		s, err := t.retrieve_Song_ID(stub, Song_ID)

		if err != nil {
			fmt.Printf("INVOKE: Error retrieving Song: %s", err)
			return nil, errors.New("Error retrieving Song")
		}

		if function == "Set_Contract" { // Only the copyright authority is allowed to set a contract
			return t.set_contract(stub, s, caller, caller_affiliation, args)
		} else if function == "Set_Rating" { // Rating can be set by anybody, but only once and not by the singer
			return t.set_rating(stub, s, caller, caller_affiliation, args)
		} else if function == "Set_Contract_Response" { // Function my only be called by the singer
			return t.set_contract_response(stub, s, caller, caller_affiliation, args)
		} else if function == "update_song" { // Function may only be called by the singer
			return t.update_song(stub, s, caller, caller_affiliation, args[0])
		} else if function == "update_contract" { // Function may only be called by the copyright institution if the existing contract
			return t.update_contract(stub, s, caller, caller_affiliation, args[0])
		} else if function == "update_rating" { // Rating can be set by anybody, but only once and not by the singer. In this case, a previous rating of the user must exist
			return t.update_rating(stub, s, caller, caller_affiliation, args[0])
		}

 */
module.exports.process_msg = function(wssvc, data){
	console.log('karachain svc: process message ',data.type);
	ws = wssvc;
	if(data.v === 1){																						//only look at messages for part 1
	    if(data.type == 'createsinger'){
			console.log('karachain svc: create a singer');
     		//chaincode.invoke.init_singer([], cb_invoked);	//create a new singer
			
		}
		else if(data.type == 'createperformance'){
			console.log('karachain svc: create performance - singer singing song');
			//{"type" : "createperformance","name":"rocknroll", "venue":"lazy dog", "date":"undefined","singer": "bob","v":1}
			//data.name data.date data.singer data.videoid, data.videourl data.videoqrcode data.venuid data.venue 
			
			var songId = "kc"+Math.round(Math.pow(10,7)*Math.random());
			data.videoid = "vd"+Math.round(Math.pow(10,7)*Math.random());
			//data.videourl = "https://www.youtube.com/watch?v=Lsty-LgDNxc";
			data.venueid = "vu"+Math.round(Math.pow(10,7)*Math.random());
			data.qrid = "qr"+Math.round(Math.pow(10,7)*Math.random());
			data.singerid = lastSingerId;
			//data.singerName = "Carsten";
			//data.perfName = "RockNRoll";
			//data.performancename
			//data.performancevideo
			//data.performancevenue
			//data.videodate
			//data.performancedate
			//data.videourl
			
			
			chaincode.invoke.create_song([songId,data.videodate,data.videoid, data.performancevideo,data.performancedate, data.qrid, data.venueid, data.performancevenue,data.singerid,data.singerName,data.performancename], cb_invoked);	//create a new song		
			
			console.log('karachain svc: create performance - reading song back ',songId);
			lastSongId = songId;
			//chaincode.query.read([songId], cb_query_response);
			
			console.log('karachain: create performance - gen qr code ',songId);
			var qr_png = genQRpng(data.singerName, data.perfName,data.singerid, data.songId,data.date);
			var qrjson ={
					songid:songId,
					qr:qr_png
			};
			sendMsg(qrjson);
		}
		else if(data.type == 'createvisitor'){
			console.log('karachain svc: create visitor');
			var visitorId = "vs"+Math.round(Math.pow(10,7)*Math.random());
			
		}
		else if(data.type == 'createeventmgr'){
			console.log('karachain svc: create event mgr');
			
		} 
		else if(data.type == 'voteperformance'){
			console.log('karachain svc: vote performance ',data.songid,data.rating,data.visitorid);
			/**
			 * Song_IDUser_Rating
			 * "AA1111127", "5"
			 * Set_Rating
			 */
			//TODO get data from visitor client
			//data.rating = 5;
			data.songid = lastSongId;  
			chaincode.invoke.Set_Rating([data.songid,data.rating,lastVisitorId], cb_invoked);	//create a new song		
		}
		else if(data.type == 'viewmyperformances'){
			console.log('karachain svc: get my performances');
			chaincode.query.Get_Songs([], s);
			/*
			 * Song_ID
			 * "AA1111127"
			 */
		}
		else if(data.type == 'viewmyperformance'){
			console.log('karachain svc: get a performances');
			chaincode.query.Get_Song(lastSongId,cb_get_songs);
			/*
			 * Song_ID
			 * "AA1111127"
			 */
		}		
		else if(data.type == 'submitoffer'){
			/*
			* if len(args) != 9 {
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
			 */
			console.log('karachain svc: submit offer');
			/*
			 * "BLANK", "Copyright_Id_001", "26.02.2017", "Copyright_Institution_Id_001", "COPYRIGHT_RECORDS", "SINGER_ID_123", "01.03.2017", "01.04.2017", "SmartContract_ID_001"
			 */
			//TODO .. get data from client
			data.blank = " ";
			//data.songid = lastSongId;
			data.contractid = "ct"+Math.round(Math.pow(10,7)*Math.random());
			data.copywriteid = "cw"+Math.round(Math.pow(10,7)*Math.random());
			data.copywriteinstid = "ci"+Math.round(Math.pow(10,7)*Math.random());
			//data.cwrec = "COPYRIGHT_RECORD";
			data.singerid = lastSingerId;
			//data.date = "01/30/2017";
			//data.copywritedate = "01/30/2017";
			//data.copywritestartdate = "01/30/2017";
			//data.copywriteenddate = "01/30/2017";
			//data.copyright_inst_name = "institution";
			//data.eventmgrid = "em"+Math.round(Math.pow(10,7)*Math.random());
			//data.eventmgrname = lastEvtMgrId;
			/*
			 * karachain team2: received ws msg: {"type" : "submitoffer","copywriteid":"undefined", "copywritedate":"03/17","copywriteinstid": "undefined","copyright_inst_name":"StudioC","singerid":"user_type1_1","songid":"kc6979312","copywritestartdate":"03/17","copywriteenddate":"04/17","contractid":"undefined","contractvalue":"4000","v":1}
			 */
			chaincode.invoke.Set_Contract([data.songid,data.copywriteid,data.copywritedate,data.copywriteinstid,data.copyright_inst_name,data.singerid, data.copywritestartdate,data.copywriteenddate,data.contractid], cb_invoked);	//create a new song		
			console.log('karachain svc:submitted offer ',data.songid);
			
		}
		else if(data.type == 'viewtopperformances'){
			console.log('karachain svc: view performance ratings');
			
		}
		else if(data.type == 'getmyoffers'){
			/*
			 * type Contract struct {
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
			 */
			console.log('karachain svc: get my offers');
			data.singerid = lastSingerId;
			chaincode.query.Get_Contracts([data.singerid], cb_getoffers);
			
		}
		else if(data.type == 'acceptoffer'){
			console.log('karachain svc: accept offer');
			/*
			 * Set_Contract_Response
			 * singerid, contract id, boolean, date
			 */
			data.singerid = lastSingerId;
			data.songid = lastSongId;
			data.contractid = "ct"+Math.round(Math.pow(10,7)*Math.random());
			data.accepted = false;
			data.date = "01/30/2017";
			
			chaincode.invoke.Set_Contract_Response([data.singerid,data.contractid,data.accepted,data.date], cb_invoked);	//create a new song		
			console.log('karachain svc:accept offer ',data.songid);
		}
		else if(data.type == 'transfer'){
			console.log('transfering msg');
			
		}
		else if(data.type == 'remove'){
			console.log('removing msg');
		
		}
		else if(data.type == 'chainstats'){
			console.log('chainstats msg');
			ibc.chain_stats(cb_chainstats);
		}
	}

	//got the songs index, lets get each song
	function cb_got_index(e, index){
		if(e != null) console.log('[ws error] did not get song index:', e);
		else{
			try{
				var json = JSON.parse(index);
				var keys = Object.keys(json);
				var concurrency = 1;

				//serialized version
				async.eachLimit(keys, concurrency, function(key, cb) {
					console.log('!', json[key]);
//					chaincode.query.read([json[key]], function(e, song) {
//						if(e != null) console.log('[ws error] did not get song:', e);
//						else {
//							if(song) sendMsg({msg: 'songs', e: e, song: JSON.parse(song)});
//							cb(null);
//						}
//					});
				}, function() {
					sendMsg({msg: 'action', e: e, status: 'finished'});
				});
			}
			catch(e){
				console.log('[ws error] could not parse response', e);
			}
		}
	}
	//invoke call back
	function cb_invoked(e, a){
		console.log('invoke response: ', e, a);
	}
	//parse chaincode map into json from mobile app
	function mapToJson(jsonStr) {
		//var parsed = JSON.parse('{"SONG_ID_001":{"Song_ID":"SONG_ID_001","Date_created":"26.02.2017","Singer_Id":"Singer_ID_123","Singer_Name":"Singer_ANYBODY","Video_Id":"Video_ID_001","Owner":"","Video_Link":"http:123.de","Video_date_created":"26.02.2017","Video_QR_code_Id":"QR_STRING123","Venue_Id":"Venue_ID_001","Venue_Name":"Venue_Name_NY","User_rating":{},"Obsolete":false,"Status":"UNDEFINED","Song_Name":"Song_ANYONE","AVG_Rating":0},"SONG_ID_002":{"Song_ID":"SONG_ID_002","Date_created":"26.02.2017","Singer_Id":"Singer_ID_123","Singer_Name":"Singer_ANYBODY","Video_Id":"Video_ID_001","Owner":"","Video_Link":"http:123.de","Video_date_created":"26.02.2017","Video_QR_code_Id":"QR_STRING123","Venue_Id":"Venue_ID_001","Venue_Name":"Venue_Name_NY","User_rating":{},"Obsolete":false,"Status":"UNDEFINED","Song_Name":"Song_ANYONE","AVG_Rating":0},"kc4059949":{"Song_ID":"kc4059949","Date_created":"http://trotzkowski.de/","Singer_Id":"user_type1_1","Singer_Name":"Carsten","Video_Id":"vd9291200","Owner":"","Video_Link":"https://www.youtube.com/watch?v=Lsty-LgDNxc","Video_date_created":"http://trotzkowski.de/","Video_QR_code_Id":"qr9702995","Venue_Id":"vu6975760","Venue_Name":"internet","User_rating":{"user_type2_0":5},"Obsolete":false,"Status":"UNDEFINED","Song_Name":"RockNRoll","AVG_Rating":5},"kc6508151":{"Song_ID":"kc6508151","Date_created":"01/01/2017","Singer_Id":"user_type1_1","Singer_Na 1me":"Carsten","Video_Id":"vd2326197","Owner":"","Video_Link":"https://www.youtube.com/watch?v=Lsty-LgDNxc","Video_date_created":"01/01/2017","Video_QR_code_Id":"qr746245","Venue_Id":"vu5722680","Venue_Name":"lazy dog","User_rating":{"user_type2_0":5},"Obsolete":false,"Status":"UNDEFINED","Song_Name":"RockNRoll","AVG_Rating":5}}') ;
	   var parsedSongs = JSON.parse(jsonStr);
	   console.log("songs: ",Object.keys(parsedSongs) );
	   var songarray = new Array();
	   var id = 0;
	   for (var songid in parsedSongs) {
		   console.log('parsed.' + songid, '=', parsedSongs[songid]);
		   parsedSongs[songid].id = id++;
		   songarray.push(parsedSongs[songid]);
		 }
	  return(JSON.stringify(songarray));
	}
	//get songs cb - parse chaincode song map into array of songs for mobile app
	function s(e, songs){
		if(e != null) {
			console.log('[query songs error] did not get query response:', e);
		}else{
			if (songs != null){
				console.log('[query songs] got query song sresponse:', songs);
				var sjmap = mapToJson(songs);
				console.log("[query songs] parsed songs ",sjmap)
			}else{
				console.log('[query songs] NULL query response:');
			}
		}
		
		sendJson(sjmap);
	}
	//cc query callback
	function cb_query_response(e, response){
		if(e != null) {
			console.log('[query error] did not get query response:', e);
		}else{
			if (response != null){
				console.log('[query resonse] got query response:', response);
			}else{
				console.log('[query resonse] NULL query response:');
			}
		}
	}
	//get songs callback
	function cb_get_songs(e, response){
		if(e != null) {
			console.log('[query songs error] did not get query response:', e);
		}else{
			if (response != null){
				console.log('[query songs resonse] got query response:', response);
				//build songs json
			}else{
				console.log('[query songs resonse] NULL query response:');
				//build null response json
			}
		}
	}
	//get songs callback
	function cb_getoffers(e, contracts){
		var contracts ={};
		if(e != null) {
			console.log('[query contracts error] did not get query response:', e);
		}else{
			if (contracts != null){
				console.log('[query contracts response] got query response:', contracts);
				//build songs json
				var contractObj = JSON.parse(contracts);
				contracts = contractObj.Contracts;
				
			}else{
				console.log('[query contracts response] NULL query response:');
				//build null response json
			}
		}
		sendJson(JSON.stringify(contracts));
	}
	//call back for getting the blockchain stats, lets get the block stats now
	function cb_chainstats(e, chain_stats){
		if(chain_stats && chain_stats.height){
			chain_stats.height = chain_stats.height - 1;								//its 1 higher than actual height
			var list = [];
			for(var i = chain_stats.height; i >= 1; i--){								//create a list of heights we need
				list.push(i);
				if(list.length >= 8) break;
			}
			list.reverse();																//flip it so order is correct in UI
			async.eachLimit(list, 1, function(block_height, cb) {						//iter through each one, and send it
				ibc.block_stats(block_height, function(e, stats){
					if(e == null){
						stats.height = block_height;
						sendMsg({msg: 'chainstats', e: e, chainstats: chain_stats, blockstats: stats});
					}
					cb(null);
				});
			}, function() {
			});
		}
	}
	//send a message, socket might be closed...
	function sendMsg(json){
		if(ws){
			try{
				ws.send(JSON.stringify(json));
				console.log('[ws send] msg sent',JSON.stringify(json) );

			}
			catch(e){
				console.log('[ws error] could not send msg', e);
			}
		}
	}
	function sendJson(json){
		if(ws){
			try{
				ws.send(json);
				console.log('[sendJson ws send] msg sent',json );

			}
			catch(e){
				console.log('[sendJson ws error] could not send msg', e);
			}
		}
	}
};