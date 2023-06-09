package main

import (
	"gorfapi/apps/hello"

	"github.com/gin-gonic/gin"
	"github.com/go-gorf/admin"
	"github.com/go-gorf/gorf"
)

// add all the apps
var apps = []gorf.GorfApp{
	&hello.HelloApp,
	&admin.AdminApp,
}

func LoadSettings() {
	// jwt secret key
	gorf.Settings.SecretKey = "GOo8Rs8ht7qdxv6uUAjkQuopRGnql2zWJu08YleBx6pEv0cQ09a"
	// app settings
	// auth.AuthSettings.NewUserState = auth.AuthState(true)
}

// bootstrap server
func BootstrapRouter() *gin.Engine {
	gorf.Apps = append(apps)
	LoadSettings()
	_ = gorf.InitializeDatabase()

	gorf.SetupApps()
	r := gin.Default()
	gorf.RegisterApps(r)
	return r
}
