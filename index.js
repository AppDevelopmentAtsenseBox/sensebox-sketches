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
      "_id": "57000b8745fd40c8196ad04c",
      "sensors": [{
        "_id": "57000b8745fd40c8196ad04e",
        "lastMeasurement": {
          "value": "80",
          "createdAt": "2017-01-19T09:15:39.901Z"
        },
        "sensorType": "VEML6070",
        "title": "UV-Intensität",
        "unit": "μW/cm²"
      }, {
        "_id": "57000b8745fd40c8196ad04f",
        "lastMeasurement": {
          "value": "2456",
          "createdAt": "2017-01-19T09:15:39.901Z"
        },
        "sensorType": "TSL45315",
        "title": "Beleuchtungsstärke",
        "unit": "lx"
      }, {
        "_id": "57000b8745fd40c8196ad050",
        "lastMeasurement": {
          "value": "1032.03",
          "createdAt": "2017-01-19T09:15:39.901Z"
        },
        "sensorType": "BMP280",
        "title": "Luftdruck",
        "unit": "hPa"
      }, {
        "_id": "57000b8745fd40c8196ad051",
        "lastMeasurement": {
          "value": "93.99",
          "createdAt": "2017-01-19T09:15:39.901Z"
        },
        "sensorType": "HDC1008",
        "title": "rel. Luftfeuchte",
        "unit": "%"
      }, {
        "_id": "57000b8745fd40c8196ad052",
        "lastMeasurement": {
          "value": "-1.62",
          "createdAt": "2017-01-19T09:15:39.901Z"
        },
        "sensorType": "HDC1008",
        "title": "Temperatur",
        "unit": "°C"
      }, {
        "_id": "576996be6c521810002479dd",
        "sensorType": "WiFi",
        "unit": "dBm",
        "title": "Wifi-Stärke",
        "lastMeasurement": {
          "value": "-72",
          "createdAt": "2017-01-19T09:15:39.901Z"
        }
      }, {
        "_id": "579f9eae68b4a2120069edc8",
        "sensorType": "VCC",
        "unit": "V",
        "title": "Eingangsspannung",
        "lastMeasurement": {
          "value": "2.73",
          "createdAt": "2017-01-19T09:15:39.901Z"
        },
        "icon": "osem-shock"
      }]
    }
  }
}]


req.write(JSON.stringify(payload));

req.end();

req.on('error', function(e) {
    console.error(e);
});