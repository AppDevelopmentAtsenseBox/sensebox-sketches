var fs = require('fs')
var https = require('https')

var options = {
    hostname: 'localhost',
    port: 3924,
    path: '/',
    method: 'POST',
    key: fs.readFileSync('client_key.pem'),
    cert: fs.readFileSync('client_cert.pem'),
    ca: fs.readFileSync('ca_cert.pem')
};

var req = https.request(options, function(res) {
    res.on('data', function(data) {
        process.stdout.write(data);
    });
});

var payload = [
{
  "networktype": "ethernet",
  "payload": {
    "box": {
      "_id": "<some valid senseBox id>",
      "sensors": [
        {
          "title": "<some title>",
          "sensorType": "<some type>",
          "_id": "<some valid senseBox sensor id>"
        },
        {
          "title": "sensor2",
          "sensorType": "<some type2>",
          "_id": "<some valid senseBox sensor id222>"
        }
      ]
    }
  }
}]


req.write(JSON.stringify(payload));

req.end();

req.on('error', function(e) {
    console.error(e);
});