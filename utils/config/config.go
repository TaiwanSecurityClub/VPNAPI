package config

import (
    "os"
    "strconv"
    "path/filepath"
    "github.com/joho/godotenv"
)

var WGpath string
var Port string
var Servername string
var Debug bool

func init() {
    err := godotenv.Load()
    ex, err := os.Executable()
    if err == nil {
        exPath := filepath.Dir(ex)
        os.Chdir(exPath)
    }
    err = godotenv.Load()
    exists := false
    Port, exists = os.LookupEnv("PORT")
    if !exists {
        Port = "3000"
    }
    WGpath, exists = os.LookupEnv("WGPATH")
    if !exists {
        WGpath = "/etc/wireguard"
    }
    Servername = os.Getenv("SERVERNAME")
    debugstr, exists := os.LookupEnv("DEBUG")
    if !exists {
        Debug = false
    } else {
        Debug, err = strconv.ParseBool(debugstr)
        if err != nil {
            Debug = false
        }
    }
}
