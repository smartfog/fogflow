var fs = require('fs'); 
var https = require('https'); 
var options = { 
    hostname: 'host.docker.internal', 
    port: 443, 
    path: '/ngsi9/status', 
    method: 'GET', 
    key: fs.readFileSync('client1-key.pem'), 
    cert: fs.readFileSync('client1-crt.pem'), 
    ca: fs.readFileSync('ca-crt.pem') }; 
var req = https.request(options, function(res) { 
    res.on('data', function(data) { 
        process.stdout.write(data); 
    }); 
}); 
req.end(); 
req.on('error', function(e) { 
    console.error(e); 
});


