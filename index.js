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
  "box": {
    "_id": "58233fe540198a00100b5792",
    "createdAt": "2016-11-09T15:25:25.160Z",
    "updatedAt": "2017-01-23T14:01:42.849Z",
    "name": "Marientalstraße",
    "boxType": "fixed",
    "grouptag": "",
    "exposure": "outdoor",
    "model": "homeEthernet",
    "sensors": [
      {
        "title": "Temperatur",
        "unit": "°C",
        "sensorType": "HDC1008",
        "icon": "osem-thermometer",
        "_id": "58233fe540198a00100b5797",
        "lastMeasurement": {
          "value": "1.52",
          "createdAt": "2017-01-23T14:01:37.903Z"
        }
      },
      {
        "title": "rel. Luftfeuchte",
        "unit": "%",
        "sensorType": "HDC1008",
        "icon": "osem-humidity",
        "_id": "58233fe540198a00100b5796",
        "lastMeasurement": {
          "value": "82.5",
          "createdAt": "2017-01-23T14:01:39.133Z"
        }
      },
      {
        "title": "Luftdruck",
        "unit": "hPa",
        "sensorType": "BMP280",
        "icon": "osem-barometer",
        "_id": "58233fe540198a00100b5795",
        "lastMeasurement": {
          "value": "1019.28",
          "createdAt": "2017-01-23T14:01:40.370Z"
        }
      },
      {
        "title": "Beleuchtungsstärke",
        "unit": "lx",
        "sensorType": "TSL45315",
        "icon": "osem-brightness",
        "_id": "58233fe540198a00100b5794",
        "lastMeasurement": {
          "value": "2068",
          "createdAt": "2017-01-23T14:01:41.602Z"
        }
      },
      {
        "title": "UV-Intensität",
        "unit": "μW/cm²",
        "sensorType": "VEML6070",
        "icon": "osem-brightness",
        "_id": "58233fe540198a00100b5793",
        "lastMeasurement": {
          "value": "90",
          "createdAt": "2017-01-23T14:01:42.847Z"
        }
      }
    ],
    "loc": [
      {
        "geometry": {
          "coordinates": [
            7.614847,
            51.972261
          ],
          "type": "Point"
        },
        "type": "feature"
      }
    ],
    "image": "58233fe540198a00100b5792.jpg?1478874872593"
  }
}]


req.write(JSON.stringify(payload));

req.end();

req.on('error', function(e) {
    console.error(e);
});