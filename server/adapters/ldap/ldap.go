package ldap

import (
	"crypto/tls"
	"errors"
	"fmt"
	"strings"

	ldapV2 "gopkg.in/ldap.v2"

	log "github.com/sirupsen/logrus"
	"github.com/soprasteria/docktor/server/types"
)

// ErrInvalidCredentials is an error message  which appears when credentials are invalid
var ErrInvalidCredentials = errors.New("Invalid Username or Password")

// ErrSearchFailed is an error message which appears when search for user fail
var ErrSearchFailed = errors.New("Search for user failed")

// Client interface for LDAP
type Client interface {
	Open() error
	Search(username string) (*UserInfo, error)
	Login(query types.UserQuery) error
	Close()
}

// Config contains data used to connect to a LDAP directory service
type Config struct {
	SSL          bool
	TLS          tls.Config
	LdapServer   string
	BaseDN       string
	BindDN       string
	BindPassword string
	SearchDN     string
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
	conn *ldapV2.Conn
}

//NewClient create a LDAP entrypoint
func NewClient(conf *Config) Client {
	return &DefaultClient{conf: conf}
}

// Open LDAP connection
func (dc *DefaultClient) Open() error {
	var err error
	if !dc.conf.SSL {
		log.Info("connecting to insecured LDAP server")
		dc.conn, err = ldapV2.Dial("tcp", dc.conf.LdapServer)
		if err != nil {
			log.WithError(err).Errorf("LDAP dialing failed at %v", dc.conf.LdapServer)
			return err
		}

		// Reconnect with TLS
		if dc.conf.TLS.InsecureSkipVerify {
			err = dc.conn.StartTLS(&tls.Config{InsecureSkipVerify: true})
			if err != nil {
				log.WithError(err).Errorf("cannot start TLS at %s", dc.conf.LdapServer)
				return err
			}
		}
	} else {
		log.Info("connecting to secured LDAP server")
		dc.conn, err = ldapV2.DialTLS("tcp", dc.conf.LdapServer, &dc.conf.TLS)
		if err != nil {
			log.WithError(err).Errorf("LDAP TLS dialing failed at %v", dc.conf.LdapServer)
			return err
		}
	}
	return nil
}

// Search search existence of username in LDAP
// Returns the info about the user if found, error otherwize
func (dc *DefaultClient) Search(username string) (*UserInfo, error) {

	// perform initial authentication
	if err := dc.bind(dc.conf.BindDN, dc.conf.BindPassword); err != nil {
		log.WithError(err).Error("LDAP binding failed")
		return nil, err
	}

	// find user entry & attributes
	ldapUser, err := dc.searchForUser(username)
	if err != nil {
		log.WithError(err).WithField("username", username).Error("Error looking for user in AD")
		return nil, err
	}

	return ldapUser, nil
}

//Login log the user in
func (dc *DefaultClient) Login(query types.UserQuery) error {

	username := query.Username
	dn := dc.conf.SearchDN

	// perform initial authentication if needed
	if dn == "" {
		if err := dc.bind(dc.conf.BindDN, dc.conf.BindPassword); err != nil {
			log.WithError(err).Error("LDAP binding failed")
			return err
		}
		// find user entry & attributes
		ldapUser, err := dc.searchForUser(username)
		if err != nil {
			log.WithError(err).WithField("username", username).Error("Error looking for user in AD")
			return err
		}
		dn = ldapUser.DN
	} else {
		dn = fmt.Sprintf(dn, username)
		dn = strings.Replace(dn, "{{.ldapBase}}", dc.conf.BaseDN, -1)
	}

	// Authenticate user with password
	return dc.bind(dn, query.Password)
}

// bind creates the first connexion with readonly user
func (dc *DefaultClient) bind(dn, password string) error {

	if dc.conn == nil {
		return fmt.Errorf("LDAP connction is not open")
	}

	if dn == "" || password == "" {
		return fmt.Errorf("Password or DN to bind is empty")
	}

	// LDAP Bind
	if err := dc.conn.Bind(dn, password); err != nil {
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
func (dc *DefaultClient) searchForUser(username string) (*UserInfo, error) {
	if dc.conn == nil {
		return nil, fmt.Errorf("LDAP connection is not open")
	}

	var searchResult *ldapV2.SearchResult
	var err error

	searchReq := ldapV2.SearchRequest{
		BaseDN:       dc.conf.BaseDN,
		Scope:        ldapV2.ScopeWholeSubtree,
		DerefAliases: ldapV2.NeverDerefAliases,
		Attributes: []string{
			dc.conf.Attr.Username,
			dc.conf.Attr.Firstname,
			dc.conf.Attr.Lastname,
			dc.conf.Attr.Realname,
			dc.conf.Attr.Email,
			"dn",
		},
		Filter: strings.Replace(dc.conf.SearchFilter, "%s", ldapV2.EscapeFilter(username), -1),
	}

	searchResult, err = dc.conn.Search(&searchReq)
	if err != nil {
		return nil, err
	}

	if len(searchResult.Entries) == 0 {
		return nil, ErrInvalidCredentials
	}

	if len(searchResult.Entries) > 1 {
		return nil, errors.New("Ldap search matched more than one entry, please review your filter setting")
	}

	getLdapAttr := func(name string, result *ldapV2.SearchResult) string {
		for _, attr := range result.Entries[0].Attributes {
			if attr.Name == name && len(attr.Values) > 0 {
				return attr.Values[0]
			}
		}
		return ""
	}

	return &UserInfo{
		DN:        searchResult.Entries[0].DN,
		Username:  getLdapAttr(dc.conf.Attr.Username, searchResult),
		FirstName: getLdapAttr(dc.conf.Attr.Firstname, searchResult),
		LastName:  getLdapAttr(dc.conf.Attr.Lastname, searchResult),
		RealName:  getLdapAttr(dc.conf.Attr.Realname, searchResult),
		Email:     getLdapAttr(dc.conf.Attr.Email, searchResult),
	}, nil
}

// Close LDAP connection
func (dc *DefaultClient) Close() {
	if dc.conn != nil {
		dc.conn.Close()
		dc.conn = nil
	}
}
