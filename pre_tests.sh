#!/usr/bin/env bash

#Add config file
CFG_STR1='{ "name": "Minter Gate", "debug": true, "baseCoin": "MNT", "singleNode": true, "database": { "url": "" }, "minterApi": { "isSecure": false, "link": "'
CFG_STR2='", "port": "'
CFG_STR3='" }, "gateApi": { "isSecure" : false, "link" : "", "port" : "9000" }, "wsServer":{ "isSecure" : false, "link" : "127.0.0.1", "port" : "8800", "key" : "secret key" } }'
CFG="$CFG_STR1$GT_NODE_API_LINK$CFG_STR2$GT_NODE_API_PORT$CFG_STR3"

echo -e "$CFG" > "config.json"
