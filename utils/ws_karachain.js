// ==================================
// Part 1 - incoming messages, look for type
// ==================================
var ibc = {};
var chaincode = {};
var async = require('async');

module.exports.setup = function(sdk, cc){
	ibc = sdk;
	chaincode = cc;
};

module.exports.process_msg = function(ws, data){
	console.log('karachain: process message ',data.type);
	if(data.v === 1){																						//only look at messages for part 1
	    if(data.type == 'createsinger'){
			console.log('karachain: create a song');
			if(data.name && data.color && data.size && data.user){
				chaincode.invoke.init_singer([data.name, data.color, data.size, data.user], cb_invoked);	//create a new singer
			}
		}
		else if(data.type == 'createperformance'){
			console.log('karachain: create performance - singer singing song');
			var qr_png = qr.image('performance', { type: 'png' });
			var songId = "sb1234567";
			if(data.name ){
				chaincode.invoke.create_song([songId, qr_png], cb_invoked);	//create a new song
			}
			chaincode.query.read([songId], cb_query_response);
		}
		else if(data.type == 'createvisitor'){
			console.log('karachain: create visitor');
			
		}
		else if(data.type == 'createeventmgr'){
			console.log('karachain: create event mgr');
			
		} 
		else if(data.type == 'voteperformance'){
			console.log('karachain: vote performance');
			
		}
		else if(data.type == 'getmyperformances'){
			console.log('karachain: get my performances');
			chaincode.query.read(['_allsongsindex'], cb_got_index);
			
		}
		else if(data.type == 'getmyoffers'){
			console.log('karachain: get my offers');
			
		}
		else if(data.type == 'acceptoffer'){
			console.log('karachain: accept offer');
			
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
	
	function cb_invoked(e, a){
		console.log('response: ', e, a);
	}
	//cc query callback
	function cb_query_response(e, resonse){
		if(e != null) {
			console.log('[query error] did not get query response:', e);
		}else{
			if (resonse != null){
				console.log('[query resonse] got query response:', response);
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