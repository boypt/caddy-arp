# caddy-arp
caddy plugin retrives client mac if possible

This plugin reads `/proc/net/arp` for IP-Mac mapping info. Thus works only when client is from same layer-2 scope.


The plugin add a single `arp` directive, provices a placeholder `{client_mac}`.

```
:80 {
  tls off
  arp
  root /srv
  gzip
  ext .htm .html
  templates
  header / x-mac {client_mac}
}
```
