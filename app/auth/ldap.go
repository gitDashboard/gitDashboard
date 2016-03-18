package auth

import (
	"errors"
	"github.com/gitDashboard/gitDashboard/app/config"
	ldap "github.com/mqu/openldap"
	"github.com/revel/revel"
	"strings"
)

func Connect() (*ldap.Ldap, error) {
	ldapServer, found := config.LdapConnection()
	if !found {
		revel.ERROR.Fatalf("Error initializing LDAP: configuration not found \"ldap.connection\"\n")
	}
	ldapConn, err := ldap.Initialize(ldapServer)
	if err == nil {
		ldapConn.SetOption(ldap.LDAP_OPT_PROTOCOL_VERSION, ldap.LDAP_VERSION2)
	} else {
		revel.ERROR.Printf("LDAP Error:%s\n", err.Error())
	}
	return ldapConn, err
}

func Login(username, password string) error {
	ldapConn, err := Connect()
	if err != nil {
		revel.ERROR.Printf("LDAP Login Error:%s\n", err.Error())
		return err
	}

	baseDn, found := config.LdapBaseDn()
	if !found {
		err = errors.New("LDAP Configuration \"ldap.baseDn\" not found")
		revel.ERROR.Printf("LDAP Login Error:%s\n", err.Error())
		return err
	}

	//check username/password
	ldapUsername := "uid=" + username + "," + baseDn
	err = ldapConn.Bind(ldapUsername, password)
	if err != nil {
		revel.ERROR.Printf("LDAP Login Error:%s\n", err.Error())
		return err
	}
	defer ldapConn.Close()
	//retrieve user information

	revel.INFO.Printf("LDAP BaseDn:%s\n", baseDn)
	//var userAttrs []string = []string{"cn"}

	userSearchResult, err := ldapConn.SearchAll(baseDn, ldap.LDAP_SCOPE_SUBTREE, "uid="+username, nil)
	if err != nil {
		revel.ERROR.Printf("LDAP Login Error:%s\n", err.Error())
		return err
	}
	revel.INFO.Printf("%v\n", userSearchResult)

	// ldapsearch -x -h vlxioldap01.intra.infocamere.it -p389 -b "ou=gruppi,o=sistema camerale,c=it" "(uniqueMember= uid=YYI3842,ou=utenti,o=Sistema Camerale,c=It)"

	groupsSearchResult, err := ldapConn.SearchAll(baseDn, ldap.LDAP_SCOPE_SUBTREE, "uid="+username, nil)
	if err != nil {
		revel.ERROR.Printf("LDAP Login Error:%s\n", err.Error())
		return err
	}
	revel.INFO.Printf("%v\n", groupsSearchResult)

	return nil
}

func Search(username string) (map[string]string, error) {
	ldapConn, err := Connect()
	if err != nil {
		revel.ERROR.Printf("LDAP Login Error:%s\n", err.Error())
		return nil, err
	}
	defer ldapConn.Close()

	baseDn, found := config.LdapBaseDn()
	if !found {
		err = errors.New("LDAP Configuration \"ldap.baseDn\" not found")
		revel.ERROR.Printf("LDAP Login Error:%s\n", err.Error())
		return nil, err
	}
	scope := ldap.LDAP_SCOPE_SUBTREE // LDAP_SCOPE_BASE, LDAP_SCOPE_ONELEVEL, LDAP_SCOPE_SUBTREE
	filter := "uid=" + username
	attributes := []string{"cn", "mail"} // leave empty for all attributes

	result, err := ldapConn.SearchAll(baseDn, scope, filter, attributes)
	if err != nil {
		revel.ERROR.Println(err)
		return nil, err
	}
	userData := make(map[string]string)
	for _, entry := range result.Entries() {
		for _, attr := range entry.Attributes() {
			userData[attr.Name()] = strings.Join(attr.Values(), ", ")
		}
	}
	return userData, nil
}
