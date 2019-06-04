package auth

import (
	"crypto/md5"
	"fmt"
	"github.com/kataras/iris"
	"io"
	"strconv"
	"time"
)

// token : username
var tokens = make(map[string]string)

func NewToken(username string) string {

	token := makeToken()

	tokens[token] = username

	return token
}

// 检查token是否正确
func CheckToken(ctx iris.Context) {

	token := ctx.GetCookie("token")
	username := ctx.GetCookie("username")

	if un, ok := tokens[token]; !ok || username != un { // 如果不存在token或token不正确
		ctx.Redirect("/", 302)
	}
}

func CheckAdmin(ctx iris.Context) {
	if ctx.GetCookie("username") != "admin" {
		ctx.Redirect("/", 302)
	}
}

func RemoveToken(token string) {
	delete(tokens, token)
}

func RemoveTokenByUsername(username string) {
	if t, ok := findTokenByUsername(username); ok {
		RemoveToken(t)
	}
}

func findTokenByUsername(username string) (string, bool) {

	for t, u := range tokens {
		if u == username {
			return t, true
		}
	}

	return "", false
}

func makeToken() string {

	currentTime := time.Now().Unix()

	h := md5.New()

	_, _ = io.WriteString(h, strconv.FormatInt(currentTime, 10))

	token := fmt.Sprintf("%x", h.Sum(nil))

	return token
}
