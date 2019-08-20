console.log("==========MAIN")
const USDT = require("usdt");
const provider = new USDT.USDTProvider("nodeProvider");
const probeOneBefore = provider.addProbe("OneMain", "int");
console.log("==========/MAIN")

setInterval( () => { console.log("Sleeping") }, 1000);
