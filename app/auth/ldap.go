package auth

import (
	"errors"
	ldap "github.com/mqu/openldap"
	"github.com/revel/revel"
)

func Connect() (*ldap.Ldap, error) {
	ldapServer, found := revel.Config.String("ldap.connection")
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

	baseDn, found := revel.Config.String("ldap.baseDn")
	if !found {
		err = errors.New("LDAP Configuration \"ldap.baseDn\" not found")
		revel.ERROR.Printf("LDAP Login Error:%s\n", err.Error())
		return err
	}

	//check username/password
	ldapUsername := "uid=" + username + "," + baseDn
	revel.INFO.Println("username and password", ldapUsername, password)
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
