// package ldapauth
// @Function :
// @File : auth_test.go
// @Author : starliu
// @Time : 2021/3/25 11:37 上午
package ldapauth

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"testing"
)

func TestBasicAuth(t *testing.T) {
	ldapAuth, err := NewLdapAuth(LdapConfig{
		ldapUrl:      "",
		ldapUser:     "",
		ldapPassword: "",
		searchDn:     "",
	})
	if err != nil {
		fmt.Println("ldap auth init failed")
		return
	}
	s := g.Server()
	s.SetPort(9876)
	s.Group("/", func(baseGroup *ghttp.RouterGroup) {
		baseGroup.Middleware(ldapAuth.MiddlewareBasicAuth)
		baseGroup.GET("/hello", func(r *ghttp.Request) {
			r.Response.WriteExit("world")
		})
	})
	s.Run()
}
