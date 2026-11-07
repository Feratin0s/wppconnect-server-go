# wppconnect-server-go: armv7 build

Fork of [wppconnect-team/wppconnect-server-go](https://github.com/wppconnect-team/wppconnect-server-go)
with changes to run on 32-bit ARM (armv7) devices.

## Tested hardware

- Device: TV Box **MXQ Pro 4k** (SoC **RK3228A**, ~1 GB RAM)
- OS: Armbian bookworm (armv7l), kernel 6.18
- Runtime: Docker, image built for `linux/arm/v7`

## Changes from upstream

- Fetches the current WhatsApp Web client version on startup
  (`cmd/server/waversion.go`) and falls back to the compiled version on failure.
- Dockerfile uses `golang:1.27.2-alpine`.
- `docker-compose.yml` targets `linux/arm/v7`.

## Build and run

    docker buildx build --platform linux/arm/v7 -t wppconnect-server-go:armv7 --load .
    docker compose up -d
    docker logs -f wpp-server-go

A native build on the TV box takes roughly 15 to 20 minutes. To build faster,
cross-compile on another machine and transfer the image:

    docker save wppconnect-server-go:armv7 | ssh user@tvbox 'docker load'

On startup you should see a log line like
`WhatsApp Web client version updated: 2.3000.xxxxxxxxxx`.
