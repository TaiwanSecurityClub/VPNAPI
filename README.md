# WireguardAPI

Append your wireguard peer by api.

## Build
1. Install golang 1.19
2. Clone this repo and cd into WireguardAPI
``` bash
git clone https://github.com/TaiwanSecurityClub/WireguardAPI
cd WireguardAPI
```
3. Run make
``` bash
make clean && make
```

## Install
1. CD into repo dir
``` bash
cd WireguardAPI
```
2. Run make
```
make install
```

## Uninstall
1. CD into repo dir
``` bash
cd WireguardAPI
```
2. Run make
```
make uninstall
```

## Run
1. Set `/usr/local/share/wgapi/.env` config
```bash
vi /usr/local/share/wgapi/.env
```

2. Run wgapi
```bash
wgapi
```
