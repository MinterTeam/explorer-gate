<p align="center" background="black"><img src="minter-logo.svg" width="400"></p>

<p align="center" style="text-align: center;">
    <a href="https://github.com/daniildulin/explorer-gate/blob/master/LICENSE">
        <img src="https://img.shields.io/packagist/l/doctrine/orm.svg" alt="License">
    </a>
    <img alt="undefined" src="https://img.shields.io/github/last-commit/MinterTeam/explorer-gate.svg">
    <a href="https://hub.docker.com/r/minterteam/explorer-gate">
        <img alt="Docker Pulls" src="https://img.shields.io/docker/pulls/minterteam/explorer-gate">
   </a>
</p>

# Minter Gate

The official repository of Minter Gate service.

Minter Gate is a service which provides to clients publish prepared transactions to Minter Network

_NOTE: This project in active development stage so feel free to send us questions, issues, and wishes_

<p align="center" background="black"><img src="minter-explorer.jpeg" width="400"></p>


## Related services:
- [explorer-extender](https://github.com/MinterTeam/minter-explorer-extender)
- [explorer-api](https://github.com/MinterTeam/minter-explorer-api)
- [explorer-validators](https://github.com/MinterTeam/minter-explorer-validators) - API for validators meta
- [explorer-tools](https://github.com/MinterTeam/minter-explorer-tools) - common packages
- [explorer-genesis-uploader](https://github.com/MinterTeam/explorer-genesis-uploader)

## API Docs

Don't forget to read the [documentation](https://minterteam.github.io/minter-gate-docs/)

## BUILD

- run `go mod download`

- run `make build`

## USE

### Setup

- build and move the compiled file to the directory e.g. `/opt/minter/gate`

- copy .env.dist to .env file in extender's directory and fill with own values

#### Run

./gate

## Docker

Change the port in `docker-compose.yml` if you want to use a different port and run  `docker-compose up`
