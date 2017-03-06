// ==================================
// karachain - incoming messages, look for type
// ==================================
var ibc = {};
var chaincode = {};
var ws ={};
var async = require('async');

module.exports.setup = function(sdk, cc, qrsvc){
	console.log("karachain setup");
	ibc = sdk;
	chaincode = cc;
	qr = qrsvc;
	
};

module.exports.process_msg = function(wssvc, data){
	console.log('karachain svc: process message ',data.type);
	ws = wssvc;
	if(data.v === 1){																						//only look at messages for part 1
	    if(data.type == 'createsinger'){
			console.log('karachain svc: create a singer');
			if(data.name && data.color && data.size && data.user){
				chaincode.invoke.init_singer([data.name, data.color, data.size, data.user], cb_invoked);	//create a new singer
			}
		}
		else if(data.type == 'createperformance'){
			console.log('karachain svc: create performance - singer singing song');
			//"type" : "createperformance","name":"rocknroll", "venue":"lazy dog", "date":"undefined","singer": "bob","v":1}
			//data.name data.date data.singer data.videoid, data.videourl data.videoqrcode data.venuid data.venue 
			var songId = "kc"+Math.round(Math.pow(10,7)*Math.random());
			var qr_png = qr.imageSync(songId, { type: 'png' });
			//qr.image(songId, { type: 'png' });
			data.videoid = "vd"+Math.round(Math.pow(10,7)*Math.random());
			data.videourl = "https://www.youtube.com/watch?v=Lsty-LgDNxc";
			data.venueid = "vu"+Math.round(Math.pow(10,7)*Math.random());
			chaincode.invoke.create_song([songId,data.date,data.videoid, data.videourl, qr_png, data.venueid, data.venue], cb_invoked);	//create a new song		
			console.log('karachain svc: create performance - reading song back ',songId);
			chaincode.query.read([songId], cb_query_response);
			console.log('karachain: create performance - submitted song query ',songId);
			var response ={
					qr:qr_png
			};
			sendMsg(response);
		}
		else if(data.type == 'createvisitor'){
			console.log('karachain svc: create visitor');
			
		}
		else if(data.type == 'createeventmgr'){
			console.log('karachain svc: create event mgr');
			
		} 
		else if(data.type == 'voteperformance'){
			console.log('karachain svc: vote performance');
			
		}
		else if(data.type == 'getmyperformances'){
			console.log('karachain svc: get my performances');
			chaincode.query.read(['_allsongsindex'], cb_got_index);
			
		}
		else if(data.type == 'getmyoffers'){
			console.log('karachain svc: get my offers');
			
		}
		else if(data.type == 'acceptoffer'){
			console.log('karachain svc: accept offer');
			
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