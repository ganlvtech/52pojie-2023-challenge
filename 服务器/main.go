package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/golang-jwt/jwt/v4"
)

var (
	JWTSecretKey = []byte("80ef070c-cdff-5ca1-b7e8-413b4e27fcfb")
)

type FlagCClaims struct {
	Uid  string `json:"uid"`
	Role string `json:"role"`
}

func (f *FlagCClaims) Valid() error {
	return nil
}

func MD5(s string) string {
	sum := md5.Sum([]byte(s))
	return hex.EncodeToString(sum[:])
}

func CalcFlag(uid string, secret string, t time.Time) (flag string, expiredAt string) {
	flag = MD5(fmt.Sprintf("%s%s%d", uid, secret, t.Unix()/600))[:8]
	expiredAt = time.Unix(((t.Unix()/600)+1)*600, 0).Format(time.RFC3339)
	return
}

func main() {
	var listen string
	flag.StringVar(&listen, "listen", ":8888", "ip:port")
	flag.Parse()

	h := server.Default(server.WithHostPorts(listen))
	h.GET("/", func(ctx context.Context, c *app.RequestContext) {
		// 如果有 flag 则重定向一次，把 flag 去掉
		if c.URI().QueryArgs().Has("flag") {
			c.Redirect(consts.StatusFound, []byte("/"))
			return
		}

		// 根据 Header 中的 X-WuAiPoJie-Uid 提示 flagA
		uid := c.Request.Header.Get("X-52PoJie-Uid")
		if uid == "" {
			c.Header("X-Dynamic-Flag", "flagA{Header X-52PoJie-Uid Not Found}")
		} else {
			flag1, expiredAt := CalcFlag(uid, "_2023challenge_52pojie_cn_", time.Now())
			c.Header("X-Dynamic-Flag", fmt.Sprintf("flagA{%s} ExpiredAt: %s", flag1, expiredAt))
		}

		data, err := os.ReadFile("public/index.html")
		if err != nil {
			c.AbortWithStatus(consts.StatusNotFound)
			return
		}
		c.Data(consts.StatusOK, "text/html; charset=utf-8", data)
	})
	h.GET("/bgm.mp3", func(ctx context.Context, c *app.RequestContext) {
		data, err := os.ReadFile("public/bgm.mp3")
		if err != nil {
			c.AbortWithStatus(consts.StatusNotFound)
			return
		}
		c.Data(consts.StatusOK, "audio/mpeg", data)
	})
	h.GET("/login", func(ctx context.Context, c *app.RequestContext) {
		c.Data(consts.StatusOK, "text/html; charset=utf-8", []byte(`<form method="POST">吾爱破解 UID: <input type="text" name="uid" disabled><input type="submit"></form>`))
	})
	h.POST("/login", func(ctx context.Context, c *app.RequestContext) {
		uid := string(c.FormValue("uid"))
		if uid == "" {
			c.AbortWithStatus(consts.StatusBadRequest)
			return
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, &FlagCClaims{
			Uid:  uid,
			Role: "user",
		})
		tokenString, err := token.SignedString(JWTSecretKey)
		if err != nil {
			_ = c.AbortWithError(consts.StatusInternalServerError, err)
			return
		}
		domain := strings.SplitN(string(c.Host()), ":", 2)[0]
		c.SetCookie("2023_challenge_jwt_token", tokenString, 600, "/", domain, protocol.CookieSameSiteStrictMode, false, true)
		c.Redirect(consts.StatusFound, []byte("/home"))
	})
	h.GET("/home", func(ctx context.Context, c *app.RequestContext) {
		tokenString := string(c.Cookie("2023_challenge_jwt_token"))
		if tokenString == "" {
			c.Redirect(consts.StatusFound, []byte("/login"))
			return
		}
		claims := FlagCClaims{}
		_, _, err := jwt.NewParser().ParseUnverified(tokenString, &claims)
		if err != nil {
			_ = c.AbortWithError(consts.StatusUnauthorized, err)
			return
		}
		if claims.Role != "admin" {
			c.Data(consts.StatusOK, "text/html; charset=utf-8", []byte("您不是 admin，你没有权限获取 flag"))
			return
		}
		flag1, expiredAt := CalcFlag(claims.Uid, "_jwt_must_verify_sign_", time.Now())
		c.Data(consts.StatusOK, "text/html; charset=utf-8", []byte(fmt.Sprintf("欢迎，admin。您的 flag 是 flagC{%s}，过期时间是 %s", flag1, expiredAt)))
	})

	h.Spin()
}
