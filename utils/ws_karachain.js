// ==================================
// karachain - incoming messages, look for type
// ==================================
var ibc = {};
var chaincode = {};
var ws ={};
var async = require('async');
var performances = [];
var lastSongId = "kc846908"; //placeholder value
var lastSingerId = "user_type1_1";
var lastVisitorId = "user_type2_0";
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

module.exports.setup = function(sdk, cc, qrsvc){
	console.log("karachain setup");
	ibc = sdk;
	chaincode = cc;
	qr = qrsvc;
	
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
			var qr_png = qr.imageSync(songId, { type: 'png' });
			var qrjson ={
					songid:songId,
					qr:qr_png
			};
			//qr.image(songId, { type: 'png' });
			data.videoid = "vd"+Math.round(Math.pow(10,7)*Math.random());
			data.videourl = "https://www.youtube.com/watch?v=Lsty-LgDNxc";
			data.venueid = "vu"+Math.round(Math.pow(10,7)*Math.random());
			data.qrid = "qr"+Math.round(Math.pow(10,7)*Math.random());
			data.singerid = lastSingerId;
			
			chaincode.invoke.create_song([songId,data.date,data.videoid, data.videourl,data.date, data.qrid, data.venueid, data.venue,data.singerid], cb_invoked);	//create a new song		
			console.log('karachain svc: create performance - reading song back ',songId);
			lastSongId = songId;
			chaincode.query.read([songId], cb_query_response);
			console.log('karachain: create performance - submitted song query ',songId);
		
			sendMsg(qrjson);
		}
		else if(data.type == 'createvisitor'){
			console.log('karachain svc: create visitor');
			
		}
		else if(data.type == 'createeventmgr'){
			console.log('karachain svc: create event mgr');
			
		} 
		else if(data.type == 'voteperformance'){
			console.log('karachain svc: vote performance ',data.songid,data.rating);
			/**
			 * Song_IDUser_Rating
			 * "AA1111127", "5"
			 * Set_Rating
			 */
			//TODO get data from visitor client
			//data.rating = 5;
			//data.songid = lastSongId;  
			chaincode.invoke.Set_Rating([data.songid,data.rating,lastVisitorId], cb_invoked);	//create a new song		
		}
		else if(data.type == 'viewmyperformances'){
			console.log('karachain svc: get my performances');
			chaincode.query.Get_Songs([], cb_query_songs);
			/*
			 * Song_ID
			 * "AA1111127"
			 */
		}
		else if(data.type == 'viewmyperformance'){
			console.log('karachain svc: get a performances');
			chaincode.query.Get_Song(lastSongId, cb_query_songs);
			/*
			 * Song_ID
			 * "AA1111127"
			 */
		}		
		else if(data.type == 'submitoffer'){
			console.log('karachain svc: submit offer');
			/*
			 * Song_IDCopyright_IDCopyright_Date_createdCopyright_Institution_IDCopyright_Institution_NameSmartContract_Unique_ID
			 * "AA1111127", "12345", "01.01.2017", "Institution_123", "Dr. Dre Records", "56789"
			 */
			//TODO .. get data from client
			data.songid = lastSongId;
			data.contractid = "ct"+Math.round(Math.pow(10,7)*Math.random());
			data.copywriteid = "cw"+Math.round(Math.pow(10,7)*Math.random());
			data.date = "01/30/2017";
			data.copywritedate = "01/30/2017";
			data.eventmgrid = "em"+Math.round(Math.pow(10,7)*Math.random());
			data.eventmgrname = "Bob";
			
			chaincode.invoke.Set_Contract([data.songid,data.copywriteid,data.copywritedate, data.eventmgrid,data.eventmgrname, data.contractid], cb_invoked);	//create a new song		
			console.log('karachain svc:submitted offer ',data.songid);
			
		}
		else if(data.type == 'viewtopperformances'){
			console.log('karachain svc: view performance ratings');
			
		}
		else if(data.type == 'getmyoffers'){
			console.log('karachain svc: get my offers');
			
		}
		else if(data.type == 'acceptoffer'){
			console.log('karachain svc: accept offer');
			/*
			 * Song_IDCopyright_decision (“false” or “true”)Copyright decision date
			 * "AA1111127", "false", "01.01.2017"
			 * Set_Contract_Response
			 */
			data.songid = lastSongId;
			data.contractid = "ct"+Math.round(Math.pow(10,7)*Math.random());
			data.accepted = false;
			data.date = "01/30/2017";
			
			chaincode.invoke.Set_Contract_Response([data.songid,data.accepted,data.date], cb_invoked);	//create a new song		
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
					chaincode.query.read([json[key]], function(e, song) {
						if(e != null) console.log('[ws error] did not get song:', e);
						else {
							if(song) sendMsg({msg: 'songs', e: e, song: JSON.parse(song)});
							cb(null);
						}
					});
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
	//get songs cb
	function cb_query_songs(e, songs){
		if(e != null) {
			console.log('[query songs error] did not get query response:', e);
		}else{
			if (resonse != null){
				console.log('[query songs] got query response:', response);
			}else{
				console.log('[query songs] NULL query response:');
			}
		}
		sendMsg(songs);
	}
	//cc query callback
	function cb_query_response(e, resonse){
		if(e != null) {
			console.log('[query error] did not get query response:', e);
		}else{
			if (resonse != null){
				console.log('[query resonse] got query response:', response);
			}else{
				console.log('[query resonse] NULL query response:');
			}
		}
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
			}
			catch(e){
				console.log('[ws error] could not send msg', e);
			}
		}
	}
};