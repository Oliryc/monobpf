// Run with `node client.js ws://localhost:7777/`

// From https://github.com/theturtle32/WebSocket-Node#client-example
var WebSocketClient = require('websocket').client;

var client = new WebSocketClient();

client.on('connectFailed', function(error) {
    console.log('Connect Error: ' + error.toString());
});

client.on('connect', function(connection) {
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
        }
    });
    
    function sendNumber() {
        if (connection.connected) {
            var times = 10;
            console.log("out");
            for (var t = 1; t <= times; t++) {
              console.log("in");
              // var number = Math.round(Math.random() * 0xFFFFFF);
              connection.sendUTF(t.toString());
              console.log(t, "/", times, " sent");
            }
            // setTimeout(sendNumber, 1000);
        }
    }
    sendNumber();
});

var wsUrl = process.argv[2];
console.log("Connecting to", wsUrl);
client.connect(wsUrl, null);

