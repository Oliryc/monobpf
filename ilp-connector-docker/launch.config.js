'use strict'
const path = require('path')
const myaddress = 'rpcf2MMWqR5da7hhJ8X6PwRHGPA2nuZXDj'//<RIPPLE ADDRESS>
const mysecret = 'sss7GHxuR4paftg5pyq531zSbR3Xz'//<RIPPLE SECRET>

// The main differences between plugins are whether they are multi-user and which settlement ledger they use.
const moneydGui = {
  relation: 'child',
  plugin: 'ilp-plugin-mini-accounts',
  assetCode: 'XRP',
  assetScale: 6,//https://interledger.org/rfcs/0031-dynamic-configuration-protocol/
  options: {
    port: 7768
  }
}

// parent&child (one node is the parent of another node that is relatively its child), or, peer&peer (two nodes are peered with one another)
const peercca={
    relation: 'peer',
//a direct peering relationship with another connector over XRP
    plugin: 'ilp-plugin-xrp-paychan',
    assetCode: 'XRP',
    assetScale: 9,
    balance: {
        maximum: '1000000',
        settleThreshold: '-5000000000',
        settleTo: '0'
    },
    options: {
		listener: {
			port: 51666,
			secret: 'bW9uZXlkZ3V5Y2Fzc2FnbmVzY3lyaWwK' // this is the token that your peer must authenticate with.
			},
// If you wish to connect to your peer as a ws client, specify the server option.
// You may specify both the server and client options; in that case it is not deterministic
// which peer will end up as the ws client.
//https://interledger.org/rfcs/0010-connector-to-connector-protocol/
//		server: 'btp+ws://cyrilRcassagnes:6lr+gS6tcm0/AmX/MYgawpAAYZu0lTZeHbjIIVDHPSo=@192.172.1.213:51666', 
//TOKEN should be 64 bits represented in uppercase hex string
// Specify the server that you submit XRP transactions to.
		xrpServer: 'ws://192.172.1.207:51233',
// XRP address and secret
		address: myaddress,
		secret: mysecret,
// Peer's XRP address
		peerAddress: 'rMntKvpragJ6sAqkcq71a6LjBWmqB5SBbe'
	}
}

const ilspServer = {
    relation: 'child',
    plugin: 'ilp-plugin-xrp-asym-server',
    assetCode: 'XRP',
    assetScale: 6,
    options: {
        port: 7443,
//PORT 51233 - https://github.com/interledgerjs/ilp-plugin-xrp-asym-server
        xrpServer: 'ws://192.172.1.207:51233',
        address: myaddress,
		secret: mysecret,
    }
}

const connectorApp = {
    name: 'connector',
    env: {
        DEBUG: 'ilp*,connector*',
        CONNECTOR_ENV: 'production',
        CONNECTOR_ADMIN_API: true,
        CONNECTOR_ADMIN_API_PORT: 7769,
        CONNECTOR_ILP_ADDRESS: 'g.cca16.cojs',    //<YOUR ILP ADDRESS>
        CONNECTOR_BACKEND: 'one-to-one',
        CONNECTOR_SPREAD: '0',
        CONNECTOR_STORE: 'memdown',
        CONNECTOR_ACCOUNTS: JSON.stringify({
            classic_cca16_cnn: peercca,
            ilsp_clients: ilspServer,
            moneyd_GUI: moneydGui
        })
    },
    script: path.resolve(__dirname, 'src/index.js')
}

module.exports = { apps: [ connectorApp ] }

