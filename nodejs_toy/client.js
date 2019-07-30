// Run with `node client.js ws://localhost:7777/`

// From https://github.com/theturtle32/WebSocket-Node#client-example
var WebSocketClient = require('websocket').client;

var client = new WebSocketClient();

client.on('connectFailed', function(error) {
    console.log('Connect Error: ' + error.toString());
});

client.on('connect', function(connection) {
    // Number of message to be sent
    var msgNumber = 10;
    // Number of received message
    var msgReceived = 0;

    console.log('WebSocket Client Connected');
    connection.on('error', function(error) {
        console.log("Connection Error: " + error.toString());
    });
    connection.on('close', function() {
        console.log('echo-protocol Connection Closed');
    });
    connection.on('message', function(message) {
        if (message.type === 'utf8') {
            console.log("Received: '" + message.utf8Data + "'");
            msgReceived = msgReceived + 1;
        }
    });
    
    function sendNumber() {
        if (connection.connected) {
            for (var t = 1; t <= msgNumber; t++) {
              // var number = Math.round(Math.random() * 0xFFFFFF);
              connection.sendUTF(t.toString());
              console.log(t, "/", msgNumber, " sent");
            }
        }
    }
    sendNumber();
    console.log("=======");
    function closeWait() {
        if (msgReceived < msgNumber) {
          console.log("Waiting for response (got "+msgReceived+"/"+msgNumber+")")
          setTimeout(closeWait, 1000)
        } else {
          console.log("Got message "+msgReceived+"/"+msgNumber+", closing connection")
          connection.close();
        }
     }
     closeWait();
});

var wsUrl = process.argv[2];
console.log("Connecting to", wsUrl);

var requestOptions = {
  closeTimeout: 10000
};

client.connect(wsUrl, null, null, null, requestOptions);

