package adapters

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/pkg/errors"

	clientEntities "reimagined_eureka/internal/client/entities"
	"reimagined_eureka/internal/common"
)

const urlPrefix = "/api"
const pathSignIn = "/user/login"
const pathSignUp = "/user/register"

type ServerProxy struct {
	serverURL *url.URL
}

func NewServerProxy(serverURL *url.URL) *ServerProxy {
	return &ServerProxy{serverURL: serverURL}
}

func (p *ServerProxy) LogIn(login string, password string) (*clientEntities.UserDataResponse, error) {
	return p.signInOrUp(login, password, pathSignIn, nil)
}

func (p *ServerProxy) Register(
	login string, password string, entropy *common.EncryptionResult,
) (*clientEntities.UserDataResponse, error) {
	return p.signInOrUp(login, password, pathSignUp, entropy)
}

func (p *ServerProxy) signInOrUp(
	login, password, path string,
	entropy *common.EncryptionResult,
) (*clientEntities.UserDataResponse, error) {
	fullURL := url.URL{
		Scheme: p.serverURL.Scheme,
		Host:   p.serverURL.Host,
		Path:   urlPrefix + path,
	}
	authReqBody := common.Credentials{Login: login, Password: password}
	if entropy != nil {
		authReqBody.Entropy = entropy
	}
	payload, err := json.Marshal(authReqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	var buf bytes.Buffer
	if _, err := buf.Write(payload); err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req, err := http.NewRequest("POST", fullURL.String(), &buf)
	if err != nil {
		return nil, fmt.Errorf("failed to create a request: %v", err)
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := (&http.Client{}).Do(req)
	if resp == nil {
		return nil, fmt.Errorf("no response from the server")
	}
	if err != nil || resp.StatusCode != http.StatusOK {
		body, bodyErr := io.ReadAll(resp.Body)
		defer resp.Body.Close()
		var respErr error
		if bodyErr == nil {
			respErr = fmt.Errorf("server error response: %s", string(body))
		} else {
			respErr = errors.Wrap(err, fmt.Errorf("failed to read body: %v", bodyErr).Error())
		}
		if err == nil {
			err = respErr
		} else {
			err = errors.Wrap(err, respErr.Error())
		}
		return nil, fmt.Errorf("request failed: %v", err)
	}

	for _, cookie := range resp.Cookies() {
		if cookie.Name == common.SessionCookieName {
			return &clientEntities.UserDataResponse{SessionCookie: cookie.Value}, nil
		}
	}
	return nil, fmt.Errorf("incorrect response from server: no session cookie set")
}
