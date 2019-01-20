const connect = require('connect');
const http = require('http');
const urlParser = require('url').parse;
const bodyParser = require('body-parser');
const app = connect();
const responseTime = require('response-time');
const morgan = require('morgan');

const { compileHandler, payloadValidator } = require('./builder');
const { downloadHandler } = require('./download');
const { HTTPError } = require('./utils');

const defaultHeaders = {
  'Content-Type': 'application/json',
  'Accept': 'application/json',
  'Allow': 'GET,POST',
  'X-Backend-Server': (require('os').hostname())
};

const preflight = function preflight (req, res, next) {
  // preflight POST request https://gist.github.com/balupton/3696140
  res.setHeader('Access-Control-Allow-Origin', '*');
  res.setHeader('Access-Control-Request-Method', '*');
  res.setHeader('Access-Control-Allow-Methods', 'GET,POST');
  res.setHeader('Access-Control-Allow-Headers', 'content-type');
  res.setHeader('Access-Control-Expose-Headers', 'x-backend-server, x-response-time');
  if (req.method === 'OPTIONS') {
    res.writeHead(204);
    res.end();
    return;
  }

  next();
}

const preRequestValidator = function preRequestValidator (req, res, next) {
  // set some headers, just in case
  for (const [k, v] of Object.entries(defaultHeaders)) {
    res.setHeader(k, v);
  }

  // Parse URL to get query
  const url = urlParser(req.url, true);
  req._url = url;

  // reject everything not coming through /compile or /download
  if (url.pathname !== '/compile' && url.pathname !== '/download') {
    return next(new HTTPError({ code: 404, error: `Cannot serve ${req.url}` }));
  }

  // reject all non POST request
  if (req.method !== 'POST' && req.method !== 'GET') {
    return next(new HTTPError({ code: 405, error: 'Invalid HTTP method. Only GET or POST requests allowed.' }));
  }

  next();
};

const errorHandler = function errorHandler (err, req, res, next) {
  if (res.headersSent) {
    return next(err);
  }
  res.statusCode = (err.statusCode ? err.statusCode : 500);
  res.setHeader('Content-Type', 'application/json');
  res.end(JSON.stringify({
    code: http.STATUS_CODES[res.statusCode],
    message: err.message
  }));
};

const startServer = function startServer () {
  app.use(morgan(':date[iso] :res[x-backend-server] :remote-addr :req[x-real-ip] :method :url :response-time[0] :status'));
  app.use(responseTime());
  app.use(preflight);
  app.use(preRequestValidator);
  app.use('/compile', bodyParser.json());
  app.use('/compile', payloadValidator);
  app.use('/compile', compileHandler);
  app.use('/download', downloadHandler);
  app.use(errorHandler);

  http.createServer(app).listen(3000);
  console.log('Compiler started and listening on port 3000!');
};

startServer();

module.exports = app; // for testing
