
# Proxy

A plain Proxy server, for development and testing

- HTTP, HTTPS, ~~WS, WSS, H2C, H2~~
- serves proxy.pac
- ~~Basic Authentication~~
- Works behind another proxy
    - static proxy
    - ~~pac file~~
    - ~~WPAD~~


options

    --port <port>       specify the port
    --pac <file.pac>    specify local pac file
    --proxy <url>       proxy
    --proxy-env <name>  proxy url (env)
    --proxy-pac <url>   proxy (pac URL)
