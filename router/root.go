package router

import (
    //"fmt"
    "strconv"
    //"encoding/json"

    "github.com/gin-gonic/gin"
    
    "WireguardAPI/utils/config"
    "WireguardAPI/models/privatekey"
    "WireguardAPI/models/wireguard"
)

var router *gin.RouterGroup

func Init(r *gin.RouterGroup) {
    router = r
    router.GET("/genkey", genkey)
    router.POST("/reload", reload)
    router.POST("/getconfig", getconfig)
}

func genkey(c *gin.Context) {
    c.String(200, privatekey.Generate())
}
func reload(c *gin.Context) {
    var datas []wireguard.PeerData
    c.BindJSON(&datas)
    c.String(200, strconv.FormatBool(wireguard.Reload(config.Servername, datas)))
}
func getconfig(c *gin.Context) {
    var data wireguard.PeerData
    c.BindJSON(&data)
    c.String(200, wireguard.GetPeerConfig(config.Servername, data))
}
