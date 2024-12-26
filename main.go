package main

import (
    "fmt"

    "github.com/gin-gonic/gin"
    "github.com/go-errors/errors"

    "WireguardAPI/router"
    "WireguardAPI/middlewares/token"
    "WireguardAPI/utils/config"
    "WireguardAPI/utils/errutil"
)

func main() {
    if !config.Debug {
        gin.SetMode(gin.ReleaseMode)
    }
    backend := gin.Default()
    backend.Use(errorHandler)
    backend.Use(gin.CustomRecovery(panicHandler))
    backend.Use(token.AddMeta)
    router.Init(&backend.RouterGroup)
    backend.Run(fmt.Sprintf("%s:%s", config.Address, config.Port))
}

func panicHandler(c *gin.Context, err any) {
    goErr := errors.Wrap(err, 2)
    errutil.AbortAndError(c, &errutil.Err{
        Code: 500,
        Msg: "Internal server error",
        Data: goErr.Error(),
    })
}

func errorHandler(c *gin.Context) {
    c.Next()

    for _, e := range c.Errors {
        err := e.Err
        if myErr, ok := err.(*errutil.Err); ok {
            if myErr.Msg != nil {
                c.JSON(myErr.Code, myErr.ToH())
            } else {
                c.Status(myErr.Code)
            }
        } else {
            c.JSON(500, gin.H{
                "code": 500,
                "msg": "Internal server error",
                //"data": err.Error(),
            })
        }
        return
    }
}
