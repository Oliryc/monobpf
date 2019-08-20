// Import usdt module
const USDT = require("usdt");

// Create a Provider, which we'll use to register probes
const provider = new USDT.USDTProvider("nodeProvider");

// Register your first probe with two arguments by passing its types
const probe1 = provider.addProbe("firstProbe", "int", "char *");

const probeOneBefore = provider.addProbe("OneBefore", "int");
console.log('probeOneBefore');
console.log(probeOneBefore);
const probeOneAfter = provider.addProbe("OneAfter", "int");

const probeTwoBefore = provider.addProbe("TwoBefore", "int");
const probeTwoAfter = provider.addProbe("TwoAfter", "int");


// Enable provider (needed to fire probes later)
provider.enable();

function functionTwoProbed(a, b) {
  // Doing actual work…
  return a + b;
}

function two(a, b) {
  console.time('before_two');
  probeOneBefore.fire(function() { return [a, b]});
  console.timeEnd('before_two');

  console.time('body_two');
  c = functionOneProbed(a, b);
  console.timeEnd('body_two');

  console.time('after_two');
  probeOneAfter.fire(function() { return [c]});
  console.timeEnd('after_two');
  return;
}


function functionOneProbed(a, b) {
  // Doing actual work…
  return Math.sqrt(a);
}

function one(a, b) {
  console.time('nothing');
  console.timeEnd('nothing');

  console.time('before_one');
  probeOneBefore.fire(function() { return [a, b]});
  console.timeEnd('before_one');

  console.time('body_one');
  c = functionOneProbed(a, b);
  console.timeEnd('body_one');

  console.time('after_one');
  probeOneAfter.fire(function() { return [c]});
  console.timeEnd('after_one');
  return;
}

let countdown = 10;
function waiter() {
  one(Date.now(), Math.random()*1000);
  // two(Math.random()*10000, Math.random()*1000);
  // console.log("Trying to fire probe...");

  /*
  if(countdown <= 0) {
    console.log("Disable provider");
    provider.disable();
  }
  // Try to fire probe
    probe1.fire(function() {
    // This function will only run if the probe was enabled by an external tool
    // console.log("Probe fired!");
    countdown = countdown - 1;
    // Returning values will be passed as arguments to the probe
    return [countdown, "My little string"];
  });
  */
}

setInterval(waiter, 1000);
