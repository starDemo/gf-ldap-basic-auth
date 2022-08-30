// Package ldapauth
// @Function :
// @File : auth.go
// @Author : starliu
// @Time : 2021/3/25 11:24 上午
package ldapauth

import (
	"encoding/base64"
	"github.com/go-ldap/ldap/v3"
	"github.com/gogf/gf/v2/net/ghttp"
	"net/http"
	"strconv"
	"strings"
)

type LdapConfig struct {
	LdapUrl      string
	LdapUser     string
	LdapPassword string
	SearchDn     string
}

type LdapAuth struct {
	config     *LdapConfig
	connection *ldap.Conn
}

func NewLdapAuth(config LdapConfig) (*LdapAuth, error) {
	conn, err := getConn(config.LdapUrl, config.LdapUser, config.LdapPassword)
	if err != nil {
		return nil, err
	}
	return &LdapAuth{
		config:     &config,
		connection: conn,
	}, nil
}

func (l *LdapAuth) MiddlewareBasicAuth(r *ghttp.Request) {
	realm := "Basic realm=" + strconv.Quote("Authorization Required")
	authHeader := r.Header.Get("Authorization")
	status, err := l.ldapBasicAuth(authHeader)
	if err != nil {
		r.Response.WriteHeader(http.StatusInternalServerError)
		r.Exit()
		return
	}
	if !status {
		r.Response.Writer.Header().Set("WWW-Authenticate", realm)
		r.Response.WriteHeader(http.StatusUnauthorized)
		return
	}
	r.Middleware.Next()
}

func (l *LdapAuth) ldapBasicAuth(authData string) (bool, error) {
	if authData == "" {
		return false, nil
	}
	decodedStr, err := base64.StdEncoding.DecodeString(authData[6:])
	if err != nil {
		return false, err
	}
	cred := strings.Split(string(decodedStr), ":")
	return l.authLdapUser(cred[0], cred[1])
}
