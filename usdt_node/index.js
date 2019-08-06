// Import usdt module
const USDT = require("usdt");

// Create a Provider, which we'll use to register probes
const provider = new USDT.USDTProvider("nodeProvider");

// Register your first probe with two arguments by passing its types
const probe1 = provider.addProbe("firstProbe", "int", "char *");

// Enable provider (needed to fire probes later)
provider.enable();

let countdown = 10;
function waiter() {
  console.log("Trying to fire probe...");
  if(countdown <= 0) {
    console.log("Disable provider");
    provider.disable();
  }
  // Try to fire probe
    probe1.fire(function() {
    // This function will only run if the probe was enabled by an external tool
    console.log("Probe fired!");
    countdown = countdown - 1;
    // Returning values will be passed as arguments to the probe
    return [countdown, "My little string"];
  });
}

setInterval(waiter, 1000);
