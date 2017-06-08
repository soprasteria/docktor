package ldap

import (
	"github.com/soprasteria/docktor/server/types"
	"github.com/stretchr/testify/mock"
)

// Mock for LDAP client
type Mock struct {
	mock.Mock
}

// Search user in ldap
func (m *Mock) Search(username string) (*UserInfo, error) {
	args := m.Called(username)
	return args.Get(0).(*UserInfo), args.Error(1)
}

// Login user
func (m *Mock) Login(query types.UserQuery) (*UserInfo, error) {
	args := m.Called(query)
	return args.Get(0).(*UserInfo), args.Error(1)
}
