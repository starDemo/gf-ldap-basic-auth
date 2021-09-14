# GoFrame Ldap BasicAuth中间件

## Example
```go
package main

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/stardemo/gf-ldap-basic-auth"
)

func main() {
	fmt.Println("Starting server...")
	ldapAuthClient, err := ldapauth.NewLdapAuth(ldapauth.LdapConfig{
		LdapUrl:      "",
		LdapUser:     "",
		LdapPassword: "",
		SearchDn:     "",
	})
	if err != nil {
		fmt.Println("ldap auth init failed")
		return
	}
	s := g.Server()
	s.SetPort(9876)
	s.Group("/", func(baseGroup *ghttp.RouterGroup) {
		baseGroup.Middleware(ldapAuthClient.MiddlewareBasicAuth)
		baseGroup.GET("/hello", func(r *ghttp.Request) {
			r.Response.WriteExit("world")
		})
	})
	s.Run()

}
```
