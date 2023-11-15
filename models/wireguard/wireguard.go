package wireguard

import (
    "fmt"
    "log"
    "os"
    "os/exec"
    "bufio"
    "strings"
    "path/filepath"
    "regexp"
    "net/netip"
    "sync"

    "github.com/google/uuid"

    "WireguardAPI/utils/config"
    "WireguardAPI/utils/ipcalc"
    "WireguardAPI/models/privatekey"
)

var lock *sync.RWMutex

type PeerData struct {
    Index int64 `json:"index"`
    Key string `json:"key"`
    Name string `json:"name"`
}

func init() {
    lock = new(sync.RWMutex)
    _, err := exec.LookPath("wg")
    if err != nil {
        log.Panicln("Please install wireguard.")
    }
}

func getConfig(name string) (conf string, data map[string]string) {
    f, err := os.Open(filepath.Join(config.WGpath, fmt.Sprintf("%s.conf", name)))
    if err != nil {
        log.Panicln(err)
    }
    defer f.Close()

    scanner := bufio.NewScanner(f)
    conf = ""
    data = make(map[string]string)
    ininterface := false
    rmcomment := regexp.MustCompile(`^#\s*`)
    spliteq := regexp.MustCompile(`\s*=\s*`)
    matchend := regexp.MustCompile(`^# BEGIN .*$`)
    for scanner.Scan() {
        now := scanner.Text()
        if now == "" {
            continue
        } else if now == "[Interface]" {
            ininterface = true
        } else if !ininterface {
            nowvar := rmcomment.ReplaceAllString(now, "")
            if spliteq.MatchString(nowvar) {
                data[spliteq.Split(nowvar, 2)[0]] = spliteq.Split(nowvar, 2)[1]
            }
        } else if matchend.MatchString(now) {
            break
        } else {
            if spliteq.MatchString(now) {
                data[spliteq.Split(now, 2)[0]] = spliteq.Split(now, 2)[1]
            }
        }
        conf += fmt.Sprintf("\n%s", now)
    }

    if err := scanner.Err(); err != nil {
        log.Panicln(err)
    }

    return
}

func setConfig(name, conf string) {
    f, err := os.OpenFile(filepath.Join(config.WGpath, fmt.Sprintf("%s.conf", name)), os.O_RDWR|os.O_TRUNC, 0600)
    if err != nil {
        log.Panicln(err)
    }
    _, err = f.WriteString(conf)
    if err != nil {
        log.Panicln(err)
    }
    f.Close()
    out, _ := exec.Command("wg-quick", "strip", name).Output()

    tmpname := uuid.New().String()
    f, err = os.Create(fmt.Sprintf("/tmp/%s", tmpname))
    if err != nil {
        log.Panicln(err)
    }
    _, err = f.Write(out)
    if err != nil {
        log.Panicln(err)
    }
    f.Close()
    defer os.Remove(fmt.Sprintf("/tmp/%s", tmpname))
    exec.Command("wg", "syncconf", name, fmt.Sprintf("/tmp/%s", tmpname)).Run()
}

func Reload(name string, datas []PeerData) bool {
    lock.Lock()
    defer lock.Unlock()

    splitcom := regexp.MustCompile(`\s*,\s*`)
    
    conf, servervar := getConfig(name)
    conf += "\n"
    addressesarr := splitcom.Split(servervar["Address"], -1)
    for _, data := range datas {
        addresses := []string{}
        for _, address := range addressesarr {
            nowipaddr := ipcalc.PrefixIPGet(netip.MustParsePrefix(address), data.Index)
            addresses = append(addresses, netip.PrefixFrom(nowipaddr, nowipaddr.BitLen()).String())
        }
        conf += fmt.Sprintf(
`
# BEGIN %s
[Peer]
AllowedIPs = %s
PublicKey = %s
PersistentKeepalive = %s
# END %s`, 
            data.Name,
            strings.Join(addresses, ","),
            privatekey.Pubkey(data.Key),
            servervar["PersistentKeepalive"],
            data.Name,
        )
    }
    conf += "\n"
    setConfig(name, conf)
    return true
}

func GetPeerConfig(name string, data PeerData) string {
    lock.RLock()
    defer lock.RUnlock()

    _, servervar := getConfig(name)
    splitcom := regexp.MustCompile(`\s*,\s*`)
    addressesarr := splitcom.Split(servervar["Address"], -1)
    addresses := []string{}
    for _, address := range addressesarr {
        nowipaddr := ipcalc.PrefixIPGet(netip.MustParsePrefix(address), data.Index)
        addresses = append(addresses, netip.PrefixFrom(nowipaddr, netip.MustParsePrefix(address).Bits()).String())
    }
    conf := fmt.Sprintf(
`[Interface]
Address = %s
PrivateKey = %s
`,
        strings.Join(addresses, ","),
        data.Key,
    )
    if _, ok := servervar["DNS"]; ok {
        conf += fmt.Sprintf(`DNS = %s
`,
            servervar["DNS"],
        )
    }
    conf += fmt.Sprintf(`
[Peer]
Endpoint = %s
AllowedIPs = %s
PublicKey = %s
`, 
        fmt.Sprintf("%s:%s", servervar["Host"], servervar["ListenPort"]),
        servervar["AllowedIPs"],
        privatekey.Pubkey(servervar["PrivateKey"]),
    )
    if _, ok := servervar["PersistentKeepalive"]; ok {
        conf += fmt.Sprintf(`PersistentKeepalive = %s
`,
            servervar["PersistentKeepalive"],
        )
    }

    return conf
}
