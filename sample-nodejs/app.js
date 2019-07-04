const http = require('http');

const hostname = '127.0.0.1';
const port = 3000;

const server = http.createServer((req, res) => {
  res.statusCode = 200;
  res.setHeader('Content-Type', 'text/plain');
  res.end('Hello World\n');

  try {
    if (global.gc) {
      console.log(new Date().toLocaleString());
      console.log("gc running…");
      global.gc();
      console.log("gc finished");
    }
  } catch (e) {
    console.log("`node --expose-gc index.js`");
    process.exit();
}
});

server.listen(port, hostname, () => {
  console.log(`Server running at http://${hostname}:${port}/`);
});
