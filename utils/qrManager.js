/**
 * Service module to generate QR codes and maintain a cache indexed by a key provided by the caller
 */
var qr ={};
var async = require('async');
var qrMap = {};
function genQRpng(singerName, perfName,singerId, perfId,perfDate){
	var qrstring = "{Singer Name:"+singerName+"Performance Name:"+perfName+"Singer ID:"+singerId+"Performance Date:"+perfDate+"}";
	var qr_png = qr.image(qrstring, { type: 'png' });
	qrMap[perfId] = qr_png;
	return qr_png;
}
module.exports.setup = function(qrsvc){
	console.log("qrManager setup");
	qr = qrsvc;
	
};
module.exports.genQRcode = function(singerName, perfName,singerId, perfId,perfDate){
	return genQRpng(singerName, perfName,singerId, perfId,perfDate);

	
};
module.exports.getQRCode = function(perfId){
	return qrMap[perfId];

	
};