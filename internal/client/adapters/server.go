package adapters

import "net/url"

type ServerProxy struct {
	serverURL *url.URL
}

func NewServerProxy(serverURL *url.URL) *ServerProxy {
	return &ServerProxy{serverURL: serverURL}
}

func (p *ServerProxy) SignIn(string, string) error {
	return nil
}

func (p *ServerProxy) SignUp(string, string) error {
	return nil
}
