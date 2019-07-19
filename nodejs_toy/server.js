var fs = require('fs');
var http = require('http')
var WebSocketServer = require('websocket').server;

var start = process.hrtime();

function send_list(path, msg, callback) {
  fs.readdir(path, function (err, files) {
    callback(JSON.stringify({
      "msg": msg,
      "files": files,
    }));
  });
}

var srv = http.createServer(function(request, response) {
  // process HTTP request
});
srv.listen(7777, function() { });

wsSrv = new WebSocketServer({
  httpServer: srv
});

wsSrv.on('request', function(request) {
  var connection = request.accept(null, request.origin);

  connection.on('message', function(message) {
    if (message.type === 'utf8') {
      send_list('.', message.utf8Data, function (fileList) {
        connection.sendUTF(fileList);
      });
    }
  });

  connection.on('close', function(connection) {
    // close user connection
  });
});
