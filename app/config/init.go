package config

import (
	"github.com/revel/revel"
	"os"
)

var gitBasePath string
var gitHttpBackEnd string
var ldapConn string
var ldapBaseDn string

func GitBasePath() string {
	if gitBasePath == "" {
		gitBasePath = os.Getenv("GIT_PROJECT_ROOT")
		if gitBasePath == "" {
			return revel.Config.StringDefault("git.baseDir", "/")
		}
	}
	return gitBasePath
}

func GitHttpBackendPath() string {
	if gitHttpBackEnd == "" {
		gitHttpBackEnd = os.Getenv("GIT_HTTP_BACKEND")
		if gitHttpBackEnd == "" {
			return revel.Config.StringDefault("git.httpBackend", "/usr/lib/git-core/git-http-backend")
		}
	}
	return gitHttpBackEnd
}

func LdapConnection() (string, bool) {
	if ldapConn == "" {
		ldapConn = os.Getenv("GIT_LDAP_CONNECTION")
		if ldapConn == "" {
			return revel.Config.String("ldap.connection")
		}
	}
	return ldapConn, ldapConn != ""
}

func LdapBaseDn() (string, bool) {
	if ldapBaseDn == "" {
		ldapBaseDn = os.Getenv("GIT_LDAP_BASEDN")
		if ldapBaseDn == "" {
			return revel.Config.String("ldap.baseDn")
		}
	}
	return ldapBaseDn, ldapBaseDn != ""
}
