# Minter Gate

The official repository of Minter Gate service.

Minter Gate is a service which provides to clients publish prepared transactions to Minter Network

_NOTE: This project in active development stage so feel free to send us questions, issues, and wishes_


## Related services:
- [explorer-extender](https://github.com/MinterTeam/minter-explorer-extender)
- [explorer-api](https://github.com/MinterTeam/minter-explorer-api)
- [explorer-validators](https://github.com/MinterTeam/minter-explorer-validators) - API for validators meta
- [explorer-tools](https://github.com/MinterTeam/minter-explorer-tools) - common packages
- [explorer-genesis-uploader](https://github.com/MinterTeam/explorer-genesis-uploader)

## How to use this image

```bash
docker run -d --name gate  \
    -e GATE_DEBUG=true \
    -e GATE_PORT=9000 \
    -e BASE_COIN=MNT \
    -e NODE_API=https://texasnet.node-api.minter.network/ \
    -e NODE_API_TIMEOUT=30
```

## ... via docker-compose


Example ```docker-compose.yml``` for Minter Explorer Genesis Uploader:


```yml
version: '3.6'

services:
  app:
    image: minterteam/explorer-gate:latest
  ports:
      - 9000:9000
  environment:
      GATE_DEBUG: true
      GATE_PORT: 9000
      BASE_COIN: MNT
      NODE_API: https://minter-node-1.testnet.minter.network:8841/
      NODE_API_TIMEOUT: 30
```
