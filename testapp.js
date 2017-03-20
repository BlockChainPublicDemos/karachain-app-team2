'use strict';
/* global process */
/* global __dirname */
/*******************************************************************************
 * Copyright (c) 2015 IBM Corp.
 *
 * All rights reserved. 
 *
 * Contributors:
 *   David Huffman - Initial implementation
 *******************************************************************************/
/////////////////////////////////////////
///////////// Setup Node.js /////////////
/////////////////////////////////////////
var express = require('express');
var session = require('express-session');
var compression = require('compression');
var serve_static = require('serve-static');
var path = require('path');
var morgan = require('morgan');
var cookieParser = require('cookie-parser');
var bodyParser = require('body-parser');
var http = require('http');
var app = express();
var url = require('url');
var setup = require('./setup');
var fs = require('fs');
var cors = require('cors');
var qr = require('qr-image');

//// Set Server Parameters ////
var host = setup.SERVER.HOST;
var port = setup.SERVER.PORT;

////////  Pathing and Module Setup  ////////
app.set('views', path.join(__dirname, 'views'));
app.set('view engine', 'jade');
app.engine('.html', require('jade').__express);
app.use(compression());
app.use(morgan('dev'));
app.use(bodyParser.json());
app.use(bodyParser.urlencoded()); 
app.use(cookieParser());
app.use('/cc/summary', serve_static(path.join(__dirname, 'cc_summaries')) );												//for chaincode investigator
app.use( serve_static(path.join(__dirname, 'public'), {maxAge: '1d', setHeaders: setCustomCC}) );							//1 day cache
//app.use( serve_static(path.join(__dirname, 'public')) );
app.use(session({secret:'Somethignsomething1234!test', resave:true, saveUninitialized:true}));
function setCustomCC(res, path) {
	if (serve_static.mime.lookup(path) === 'image/jpeg')  res.setHeader('Cache-Control', 'public, max-age=2592000');		//30 days cache
	else if (serve_static.mime.lookup(path) === 'image/png') res.setHeader('Cache-Control', 'public, max-age=2592000');
	else if (serve_static.mime.lookup(path) === 'image/x-icon') res.setHeader('Cache-Control', 'public, max-age=2592000');
}
// Enable CORS preflight across the board.
app.options('*', cors());
app.use(cors());

//---------------------
// Cache Busting Hash
//---------------------
var bust_js = require('./busters_js.json');
var bust_css = require('./busters_css.json');
process.env.cachebust_js = bust_js['public/js/singlejshash'];			//i'm just making 1 hash against all js for easier jade implementation
process.env.cachebust_css = bust_css['public/css/singlecsshash'];		//i'm just making 1 hash against all css for easier jade implementation
console.log('cache busting hash js', process.env.cachebust_js, 'css', process.env.cachebust_css);


///////////  Configure Webserver  ///////////
app.use(function(req, res, next){
	var keys;
	console.log('------------------------------------------ incoming request ------------------------------------------');
	console.log('New ' + req.method + ' request for', req.url);
	req.bag = {};																			//create object for my stuff
	req.bag.session = req.session;
	
	var url_parts = url.parse(req.url, true);
	req.parameters = url_parts.query;
	keys = Object.keys(req.parameters);
	if(req.parameters && keys.length > 0) console.log({parameters: req.parameters});		//print request parameters for debug
	keys = Object.keys(req.body);
	if (req.body && keys.length > 0) console.log({body: req.body});							//print request body for debug
	next();
});

//// Router ////
app.use('/', require('./routes/site_router'));

////////////////////////////////////////////
////////////// Error Handling //////////////
////////////////////////////////////////////
app.use(function(req, res, next) {
	var err = new Error('Not Found');
	err.status = 404;
	next(err);
});
app.use(function(err, req, res, next) {														// = development error handler, print stack trace
	console.log('Error Handeler -', req.url);
	var errorCode = err.status || 500;
	res.status(errorCode);
	req.bag.error = {msg:err.stack, status:errorCode};
	if(req.bag.error.status == 404) req.bag.error.msg = 'Sorry, I cannot locate that file';
	res.render('template/error', {bag:req.bag});
});

var qr = require('qr-image');  
var express = require('express');

var app = express();
//Enable CORS preflight across the board.
app.options('*', cors());
app.use(cors());
app.get('/', function(req, res) { 
	/**
	 *  {"Contracts":[{"Copyright_Ids":"cw8017108","Copyright_date_created":"01/30/2017","Copyright_date_accepted":"","Copyright_date_rejected":"","Copyright_Institution_Id":"ci6556271","Copyright_Institution_Name":"institution","Copyright_State":"","Contract_date_from":"01/30/2017","Contract_date_to":"01/30/2017","SmartContract_ID":"ct8981330"},{"Copyright_Ids":"cw2540339","Copyright_date_created":"01/30/2017","Copyright_date_accepted":"","Copyright_date_rejected":"","Copyright_Institution_Id":"ci9358508","Copyright_Institution_Name":"institution","Copyright_State":"","Contract_date_from":"01/30/2017","Contract_date_to":"01/30/2017","SmartContract_ID":"ct2738039"},{"Copyright_Ids":"cw4566778","Copyright_date_created":"03/20/17","Copyright_date_accepted":"","Copyright_date_rejected":"","Copyright_Institution_Id":"ci7177279","Copyright_Institution_Name":"StudioK","Copyright_State":"","Contract_date_from":"04/29/17","Contract_date_to":"04/30/17","SmartContract_ID":"ct2778628"}]}
	 */
	var contractObj =  {"Contracts":[{"Copyright_Ids":"cw8017108","Copyright_date_created":"01/30/2017","Copyright_date_accepted":"","Copyright_date_rejected":"","Copyright_Institution_Id":"ci6556271","Copyright_Institution_Name":"institution","Copyright_State":"","Contract_date_from":"01/30/2017","Contract_date_to":"01/30/2017","SmartContract_ID":"ct8981330"},{"Copyright_Ids":"cw2540339","Copyright_date_created":"01/30/2017","Copyright_date_accepted":"","Copyright_date_rejected":"","Copyright_Institution_Id":"ci9358508","Copyright_Institution_Name":"institution","Copyright_State":"","Contract_date_from":"01/30/2017","Contract_date_to":"01/30/2017","SmartContract_ID":"ct2738039"},{"Copyright_Ids":"cw4566778","Copyright_date_created":"03/20/17","Copyright_date_accepted":"","Copyright_date_rejected":"","Copyright_Institution_Id":"ci7177279","Copyright_Institution_Name":"StudioK","Copyright_State":"","Contract_date_from":"04/29/17","Contract_date_to":"04/30/17","SmartContract_ID":"ct2778628"}]};
	var contracts = contractObj.Contracts;
	console.log("root URL ",contracts);
	var parsed = JSON.parse('{"SONG_ID_001":{"Song_ID":"SONG_ID_001","Date_created":"26.02.2017","Singer_Id":"Singer_ID_123","Singer_Name":"Singer_ANYBODY","Video_Id":"Video_ID_001","Owner":"","Video_Link":"http:123.de","Video_date_created":"26.02.2017","Video_QR_code_Id":"QR_STRING123","Venue_Id":"Venue_ID_001","Venue_Name":"Venue_Name_NY","User_rating":{},"Obsolete":false,"Status":"UNDEFINED","Song_Name":"Song_ANYONE","AVG_Rating":0},"SONG_ID_002":{"Song_ID":"SONG_ID_002","Date_created":"26.02.2017","Singer_Id":"Singer_ID_123","Singer_Name":"Singer_ANYBODY","Video_Id":"Video_ID_001","Owner":"","Video_Link":"http:123.de","Video_date_created":"26.02.2017","Video_QR_code_Id":"QR_STRING123","Venue_Id":"Venue_ID_001","Venue_Name":"Venue_Name_NY","User_rating":{},"Obsolete":false,"Status":"UNDEFINED","Song_Name":"Song_ANYONE","AVG_Rating":0},"kc4059949":{"Song_ID":"kc4059949","Date_created":"http://trotzkowski.de/","Singer_Id":"user_type1_1","Singer_Name":"Carsten","Video_Id":"vd9291200","Owner":"","Video_Link":"https://www.youtube.com/watch?v=Lsty-LgDNxc","Video_date_created":"http://trotzkowski.de/","Video_QR_code_Id":"qr9702995","Venue_Id":"vu6975760","Venue_Name":"internet","User_rating":{"user_type2_0":5},"Obsolete":false,"Status":"UNDEFINED","Song_Name":"RockNRoll","AVG_Rating":5},"kc6508151":{"Song_ID":"kc6508151","Date_created":"01/01/2017","Singer_Id":"user_type1_1","Singer_Na 1me":"Carsten","Video_Id":"vd2326197","Owner":"","Video_Link":"https://www.youtube.com/watch?v=Lsty-LgDNxc","Video_date_created":"01/01/2017","Video_QR_code_Id":"qr746245","Venue_Id":"vu5722680","Venue_Name":"lazy dog","User_rating":{"user_type2_0":5},"Obsolete":false,"Status":"UNDEFINED","Song_Name":"RockNRoll","AVG_Rating":5}}') ;
    console.log("root URL ",Object.keys(parsed) );
   var songarray = new Array();
   var id = 0;
   for (var prop in parsed) {
	   console.log('parsed.' + prop, '=', parsed[prop]);
	   parsed[prop].id = id++;
	   songarray.push(parsed[prop]);
	 }
  var code = qr.image(new Date().toString(), { type: 'svg' });
//  res.type('svg');
//  code.pipe(res);
 // res.header 'Access-Control-Allow-Origin', '*'
  res.send(JSON.stringify(contracts));
});


// ============================================================================================================================
// 														Launch Webserver
// ============================================================================================================================
var server = http.createServer(app).listen(port, function() {});
process.env.NODE_TLS_REJECT_UNAUTHORIZED = '0';
process.env.NODE_ENV = 'production';
server.timeout = 240000;																							// Ta-da.
console.log('------------------------------------------ Server Up - ' + host + ':' + port + ' ------------------------------------------');
if(process.env.PRODUCTION) console.log('Running using Production settings');
else console.log('Running using Developer settings');


// ============================================================================================================================
// 														Deployment Tracking
// ============================================================================================================================
console.log('- Tracking Deployment');
require('cf-deployment-tracker-client').track();		//reports back to us, this helps us judge interest! feel free to remove it


// ============================================================================================================================
// ============================================================================================================================
// ============================================================================================================================
// ============================================================================================================================
// ============================================================================================================================
// ============================================================================================================================

// ============================================================================================================================
// 														Warning
// ============================================================================================================================

// ============================================================================================================================
// 														Entering
// ============================================================================================================================

// ============================================================================================================================
// 														Work Area
// ============================================================================================================================
//var part1 = require('./utils/ws_part1');														//websocket message processing for part 1
//var part2 = require('./utils/ws_part2');														//websocket message processing for part 2
var karachainsvc = require('./utils/ws_karachain');	
var qrManager = require('./utils/qrManager');	
var ws = require('ws');																			//websocket mod
var wss = {};
var Ibc1 = require('ibm-blockchain-js');														//rest based SDK for ibm blockchain
var ibc = new Ibc1();

//
// ==================================
// load peers manually or from VCAP, VCAP will overwrite hardcoded list!
// ==================================
try{
	//this hard coded list is intentionaly left here, feel free to use it when initially starting out
	//please create your own network when you are up and running
	var manual = JSON.parse(fs.readFileSync('mycreds_docker_compose.json', 'utf8'));
	//var manual = JSON.parse(fs.readFileSync('mycreds_bluemix.json', 'utf8'));
	var peers = manual.credentials.peers;
	console.log('loading hardcoded peers');
	var users = null;																			//users are only found if security is on
	if(manual.credentials.users) users = manual.credentials.users;
	console.log('loading hardcoded users');
}
catch(e){
	console.log('Error - could not find hardcoded peers/users, this is okay if running in bluemix');
}

//Rest interface
function genQRpng(singerName, perfName,singerId, perfId,perfDate){
	var qrstring = "{Singer Name:"+singerName+"Performance Name:"+perfName+"Singer ID:"+singerId+"Performance Date:"+perfDate+"}";
	var qr_png = qr.image(qrstring, { type: 'png' });
	return qr_png;
}
//QR image service
app.get('/getqrcode/singername/:singerName/songname/:songName/singerId/:singerId/songId/:songId/perfDate/:perfDate', function(req, res) {  
	console.log("getqrcode",req.params) ; 
	var code = qrManager.genQRcode(req.params.singerName, req.params.songName,req.params.singerId, req.params.songId,req.params.perfDate);
	  //var code = qr.image("Love Shack", { type: 'png' });
	  res.type('png');
	  code.pipe(res);
	});