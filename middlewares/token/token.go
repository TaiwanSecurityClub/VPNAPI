package token

import (
    "fmt"

    "github.com/gin-gonic/gin"

    "WireguardAPI/utils/errutil"
    "WireguardAPI/utils/config"
)

func CheckAuth(c *gin.Context) {
    if isAuth, exist := c.Get("isAuth"); !exist || !isAuth.(bool) {
        errutil.AbortAndStatus(c, 401)
    }
}

func AddMeta(c *gin.Context) {
    token := c.GetHeader("Authorization")
    if fmt.Sprintf("Token %s", config.Token) == token {
        c.Set("isAuth", true)
    }
}
