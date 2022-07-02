# Dynamik

Dynamik is a dynamic DNS client written in Go for updating a single record in Cloudflare if the public IP does not match the A record provided.

Best used within a cronjob on a NAS. :) 

## Installation

Grab latest in the artefacts.

## Setup .env

The .env file will require the following to be present alongside `dynamik`:

```
CLOUDFLARE_API_KEY: ""
CLOUDFLARE_ZONE_NAME: ""
CLOUDFLARE_DYNAMIC_RECORD_NAME: ""
```

- `CLOUDFLARE_API_KEY` is a scoped API token from Cloudflare. As dynamik only reads and writes to DNS, a token with `DNS:Edit` is sufficient.
- `CLOUDFLARE_ZONE_NAME` is the name of the primary zone your token is issued for.
- `CLOUDFLARE_DYNAMIC_RECORD_NAME` is the name of the A record that you wish to check against.

## Usage

```shell
./dynamik
```

## To-do's

- Create unit tests
- Add quiet mode

## License

[MIT](https://choosealicense.com/licenses/mit/)
