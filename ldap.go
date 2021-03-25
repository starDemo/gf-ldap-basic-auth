// package ldapauth
// @Function :
// @File : ldap.go
// @Author : starliu
// @Time : 2021/3/25 11:56 上午
package ldapauth

import (
	"fmt"
	"github.com/go-ldap/ldap/v3"
)

func (l *LdapAuth) authLdapUser(user, pass string) (bool, error) {
	result, err := l.baseSearch(ldap.NewSearchRequest(
		l.config.SearchDn,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		fmt.Sprintf("(&(objectClass=inetOrgPerson)(uid=%s))", user),
		[]string{"cn"},
		nil,
	))

	if err != nil {
		return false, fmt.Errorf("failed to find user :%s,error %s", user, err.Error())
	}

	// 	用户不存在
	if len(result.Entries) < 1 {
		return false, nil
	}

	// 用户多于一个
	if len(result.Entries) > 1 {
		return false, fmt.Errorf("too many entries returned")
	}

	conn, err := newConn(l.config.LdapUrl)
	if err != nil {
		return false, err
	}
	defer conn.Close()
	if err := conn.Bind(result.Entries[0].DN, pass); err != nil {
		return false, fmt.Errorf("failed to auth. %s", err)
	}
	return true, nil
}

// newConn function
func newConn(ldapURL string) (*ldap.Conn, error) {
	ldapConn, err := ldap.DialURL(ldapURL)
	if err != nil {
		return nil, err
	}
	return ldapConn, nil
}

// getConn function
func getConn(url, bindUser, bindPass string) (*ldap.Conn, error) {
	ldapConn, err := ldap.DialURL(url)
	if err != nil {
		return nil, err
	}
	err = ldapConn.Bind(bindUser, bindPass)
	if err != nil {
		return nil, err
	}
	return ldapConn, nil
}

// checkAlive function
func (l *LdapAuth) checkAlive() error {
	// 重连检测
	if l.connection.IsClosing() {
		if conn, err := getConn(l.config.LdapUrl, l.config.LdapUser, l.config.LdapPassword); err != nil {
			return fmt.Errorf("reconnect Ldap Conn Error:%s", err.Error())
		} else {
			l.connection = conn
		}
	}
	return nil
}

// BaseSearch function
func (l *LdapAuth) baseSearch(searchRequest *ldap.SearchRequest) (*ldap.SearchResult, error) {
	// 检查连接
	if err := l.checkAlive(); err != nil {
		return nil, err
	}
	return l.connection.Search(searchRequest)
}
