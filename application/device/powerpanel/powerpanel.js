'use strict';

const NGSI = require('./ngsi/ngsiclient.js');
const fs = require('fs');

// read device profile from the configuration file
var args = process.argv.slice(2);
if(args.length != 1){
    console.log('please specify the device profile');
    return;
}

var cfgFile = args[0];
var profile = JSON.parse(
    fs.readFileSync(cfgFile)
);

var ngsi10client;
var timer;

// find out the nearby IoT Broker according to my location
var discovery = new NGSI.NGSI9Client(profile.discoveryURL)
discovery.findNearbyIoTBroker(profile.location, 1).then( function(brokers) {
    console.log('-------nearbybroker----------');    
    console.log(brokers);    
    console.log('------------end-----------');    
    if(brokers && brokers.length > 0) {
        ngsi10client = new NGSI.NGSI10Client(brokers[0]);

        // generating data observations periodically
        timer = setInterval(function(){ 
            updateContext();
        }, 2000);    

        // register my device profile by sending a device update
        registerDevice();
    }
}).catch(function(error) {
    console.log(error);
});

// register device with its device profile
function registerDevice() 
{
    var ctxObj = {};
    ctxObj.entityId = {
        id: 'Device.' + profile.type + '.' + profile.id,
        type: profile.type,
        isPattern: false
    };
    
    ctxObj.attributes = {};
    
    var degree = Math.floor((Math.random() * 100) + 1);        
    ctxObj.attributes.usage = {
        type: 'integer',
        value: degree
    };   
    ctxObj.attributes.shop = {
        type: 'string',
        value: profile.id
    };       
    ctxObj.attributes.iconURL = {
        type: 'string',
        value: profile.iconURL
    };                   
    
    ctxObj.metadata = {};
    
    ctxObj.metadata.location = {
        type: 'point',
        value: profile.location
    };    
    ctxObj.metadata.shop = {
        type: 'string',
        value: profile.id
    };      
   
    ngsi10client.updateContext(ctxObj).then( function(data) {
        console.log(data);
    }).catch(function(error) {
        console.log('failed to update context');
    });  
}

// update context for streams
function updateContext() 
{
    var ctxObj = {};
    ctxObj.entityId = {
        id: 'Device.' + profile.type + '.' + profile.id,
        type: profile.type,
        isPattern: false
    };
    
    ctxObj.attributes = {};
    
    var degree = Math.floor((Math.random() * 100) + 1);        
    ctxObj.attributes.usage = {
        type: 'integer',
        value: degree
    };       
    
    ngsi10client.updateContext(ctxObj).then( function(data) {
        console.log(data);
    }).catch(function(error) {
        console.log('failed to update context');
    });    
}

process.on('SIGINT', function() 
{    
    if(ngsi10client) {
        clearInterval(timer);
        
        // to delete the device
        var entity = {
            id: 'Device.' + profile.type + '.' + profile.id,
            type: 'Device',
            isPattern: false
        };
        ngsi10client.deleteContext(entity).then( function(data) {
            console.log(data);
        }).catch(function(error) {
            console.log('failed to delete context');
        });        

        ngsi10client.deleteContext(entity).then( function(data) {
            console.log(data);
        }).catch(function(error) {
            console.log('failed to delete context');
        });        
    }
});
