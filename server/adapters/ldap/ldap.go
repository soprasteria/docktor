package ldap

import (
	"errors"
	"fmt"
	"strings"

	ldapV2 "gopkg.in/ldap.v2"

	log "github.com/Sirupsen/logrus"
	"github.com/soprasteria/docktor/server/types"
)

// ErrInvalidCredentials is an error message when credentials are invalid
var ErrInvalidCredentials = errors.New("Invalid Username or Password")

// Client interface for LDAP
type Client interface {
	Search(username string) (*UserInfo, error)
	Login(query types.UserQuery) (*UserInfo, error)
}

// Config contains data used to connect to a LDAP directory service
type Config struct {
	LdapServer   string
	BaseDN       string
	BindDN       string
	BindPassword string
	SearchFilter string
	Attr         Attributes
}

// Attributes list all LDAP attributes names in the LDAP
type Attributes struct {
	Username  string
	Firstname string
	Lastname  string
	Realname  string
	Email     string
}

//UserInfo contains LDAP attributes values for a user
type UserInfo struct {
	DN        string
	Username  string
	FirstName string
	LastName  string
	RealName  string
	Email     string
}

// DefaultClient is an DefaultClient entry point allowing authentication with a DefaultClient server
type DefaultClient struct {
	conf *Config
}

//NewClient create a LDAP entrypoint
func NewClient(conf *Config) Client {
	return &DefaultClient{conf: conf}
}

// Search search existence of username in LDAP
// Returns the info about the user if found, error otherwize
func (a *DefaultClient) Search(username string) (*UserInfo, error) {
	// Reach the ldap server
	conn, err := ldapV2.Dial("tcp", a.conf.LdapServer)
	if err != nil {
		log.WithError(err).Error("LDAP dialing failed")
		return nil, err
	}
	defer conn.Close()

	// perform initial authentication
	if err := a.bind(conn, a.conf.BindDN, a.conf.BindPassword); err != nil {
		log.WithError(err).Error("LDAP binding failed")
		return nil, err
	}

	// find user entry & attributes
	ldapUser, err := a.searchForUser(conn, username)
	if err != nil {
		log.WithError(err).WithField("username", username).Error("Error looking for user in AD")
		return nil, err
	}

	return ldapUser, nil
}

//Login log the user in
func (a *DefaultClient) Login(query types.UserQuery) (*UserInfo, error) {
	// Reach the ldap server
	conn, err := ldapV2.Dial("tcp", a.conf.LdapServer)
	if err != nil {
		log.WithError(err).Error("LDAP dialing failed")
		return nil, err
	}
	defer conn.Close()

	// perform initial authentication
	if err := a.bind(conn, a.conf.BindDN, a.conf.BindPassword); err != nil {
		log.WithError(err).Error("LDAP binding failed")
		return nil, err
	}

	// find user entry & attributes
	ldapUser, err := a.searchForUser(conn, query.Username)
	if err != nil {
		log.WithError(err).WithField("username", query.Username).Error("Error looking for user in AD")
		return nil, err
	}

	// Authenticate user with password
	return ldapUser, a.bind(conn, ldapUser.DN, query.Password)

}

// bind creates the first connexion with readonly user
func (a *DefaultClient) bind(conn *ldapV2.Conn, dn, password string) error {

	if dn == "" || password == "" {
		return fmt.Errorf("Password or DN to bind is empty")
	}

	// LDAP Bind
	if err := conn.Bind(dn, password); err != nil {
		if ldapErr, ok := err.(*ldapV2.Error); ok {
			if ldapErr.ResultCode == ldapV2.LDAPResultInvalidCredentials {
				return ErrInvalidCredentials
			}
		}
		return err
	}

	return nil
}

// searchForUser search the user in LDAP
func (a *DefaultClient) searchForUser(conn *ldapV2.Conn, username string) (*UserInfo, error) {
	var searchResult *ldapV2.SearchResult
	var err error

	searchReq := ldapV2.SearchRequest{
		BaseDN:       a.conf.BaseDN,
		Scope:        ldapV2.ScopeWholeSubtree,
		DerefAliases: ldapV2.NeverDerefAliases,
		Attributes: []string{
			a.conf.Attr.Username,
			a.conf.Attr.Firstname,
			a.conf.Attr.Lastname,
			a.conf.Attr.Realname,
			a.conf.Attr.Email,
			"dn",
		},
		Filter: strings.Replace(a.conf.SearchFilter, "%s", ldapV2.EscapeFilter(username), -1),
	}

	searchResult, err = conn.Search(&searchReq)
	if err != nil {
		return nil, err
	}

	if len(searchResult.Entries) == 0 {
		return nil, ErrInvalidCredentials
	}

	if len(searchResult.Entries) > 1 {
		return nil, errors.New("Ldap search matched more than one entry, please review your filter setting")
	}

	return &UserInfo{
		DN:        searchResult.Entries[0].DN,
		Username:  getLdapAttr(a.conf.Attr.Username, searchResult),
		FirstName: getLdapAttr(a.conf.Attr.Firstname, searchResult),
		LastName:  getLdapAttr(a.conf.Attr.Lastname, searchResult),
		RealName:  getLdapAttr(a.conf.Attr.Realname, searchResult),
		Email:     getLdapAttr(a.conf.Attr.Email, searchResult),
	}, nil
}

func getLdapAttr(name string, result *ldapV2.SearchResult) string {
	for _, attr := range result.Entries[0].Attributes {
		if attr.Name == name && len(attr.Values) > 0 {
			return attr.Values[0]
		}
	}
	return ""
}
