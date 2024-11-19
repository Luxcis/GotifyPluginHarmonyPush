# Gotify Plugin Harmony Push

This Gotify plugin forwards all received messages to Harmony Client.

## Prerequisite

- A Service Account Config File. You can get that information by following
  this [doc](https://developer.huawei.com/consumer/cn/doc/start/api-0000001062522591#section11695162765311).
- A Harmony Application Client Token. You can get that information by following
  this [app](https://github.com/Luxcis/Gotify_Next).
- Golang, Docker, wget (If you want to build the binary from source).

## Installation

* **By shared object**
    1. Get the compatible shared object from [release](https://github.com/Luxcis/GotifyPluginHarmonyPush/releases).
    2. Put it into Gotify plugin folder.
    3. Set secrets via environment variables (List of mandatory secrets is in [Appendix](#appendix)).
    4. Restart gotify.

* **Build from source**
    1. Change GOTIFY_VERSION in Makefile.
    2. Build the binary. `make build`
    3. Follow instructions from step 2 in the shared object installation.

## Appendix

Mandatory secrets.

```(shell)
GOTIFY_SERVER_PORT=YOUR_SERVER_PORT (depending on your setup, "80" will likely work by default)
GOTIFY_CLIENT_TOKEN=YOUR_CLIENT_TOKEN (create a new Client in Gotify and use the Token from there, or you can use an existing client)
HARMONY_CLIENT_TOKEN=YOUR_HARMONY_CLIEN_TOKEN (API token provided by Harmony App)
```