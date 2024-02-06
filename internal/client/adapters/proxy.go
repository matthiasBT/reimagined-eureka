package adapters

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/pkg/errors"

	clientEntities "reimagined_eureka/internal/client/entities"
	"reimagined_eureka/internal/common"
)

const urlPrefix = "/api"
const pathSignIn = "/user/login"
const pathSignUp = "/user/register"
const pathCredentials = "/secrets/credentials"
const pathNote = "/secrets/notes"
const pathFile = "/secrets/files"
const pathCard = "/secrets/cards"

var ErrNoSessionCookie = errors.New("no session cookie set")

type ServerProxy struct {
	serverURL     *url.URL
	sessionCookie string
}

func NewServerProxy(serverURL *url.URL) *ServerProxy {
	return &ServerProxy{serverURL: serverURL}
}

func (p *ServerProxy) SetSessionCookie(cookie string) {
	p.sessionCookie = cookie
}

func (p *ServerProxy) LogIn(login string, password string) (*clientEntities.UserDataResponse, error) {
	return p.signInOrUp(login, password, pathSignIn, nil)
}

func (p *ServerProxy) Register(
	login string, password string, entropy *common.Entropy,
) (*clientEntities.UserDataResponse, error) {
	return p.signInOrUp(login, password, pathSignUp, entropy)
}

func (p *ServerProxy) AddCredentials(creds *common.CredentialsReq) (int, error) {
	if p.sessionCookie == "" {
		return 0, ErrNoSessionCookie
	}
	fullURL := url.URL{
		Scheme: p.serverURL.Scheme,
		Host:   p.serverURL.Host,
		Path:   urlPrefix + pathCredentials,
	}
	payload, err := json.Marshal(creds)
	return p.addSecret(false, fullURL.String(), payload, err)
}

func (p *ServerProxy) ReadCredentials(startID int, batchSize int) ([]*common.CredentialsReq, error) {
	if p.sessionCookie == "" {
		return nil, ErrNoSessionCookie
	}
	params := url.Values{}
	params.Add("startID", strconv.Itoa(startID))
	params.Add("batchSize", strconv.Itoa(batchSize))
	fullURL := url.URL{
		Scheme: p.serverURL.Scheme,
		Host:   p.serverURL.Host,
		Path:   urlPrefix + pathCredentials,
	}
	req, err := http.NewRequest("GET", fullURL.String()+"?"+params.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	p.addSessionCookie(req)
	body, _, err := getResponse(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	var result []*common.CredentialsReq
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse server response: %v", err)
	}
	return result, nil
}

func (p *ServerProxy) AddNote(note *common.NoteReq) (int, error) {
	if p.sessionCookie == "" {
		return 0, ErrNoSessionCookie
	}
	fullURL := url.URL{
		Scheme: p.serverURL.Scheme,
		Host:   p.serverURL.Host,
		Path:   urlPrefix + pathNote,
	}
	payload, err := json.Marshal(note)
	return p.addSecret(false, fullURL.String(), payload, err)
}

func (p *ServerProxy) UpdateNote(note *common.NoteReq) error {
	if p.sessionCookie == "" {
		return ErrNoSessionCookie
	}
	fullURL := url.URL{
		Scheme: p.serverURL.Scheme,
		Host:   p.serverURL.Host,
		Path:   urlPrefix + pathNote,
	}
	payload, err := json.Marshal(note)
	_, err = p.addSecret(true, fullURL.String(), payload, err)
	return err
}

func (p *ServerProxy) UpdateCredentials(creds *common.CredentialsReq) error {
	if p.sessionCookie == "" {
		return ErrNoSessionCookie
	}
	fullURL := url.URL{
		Scheme: p.serverURL.Scheme,
		Host:   p.serverURL.Host,
		Path:   urlPrefix + pathCredentials,
	}
	payload, err := json.Marshal(creds)
	_, err = p.addSecret(true, fullURL.String(), payload, err)
	return err
}

func (p *ServerProxy) DeleteCredentials(rowID int) error {
	if p.sessionCookie == "" {
		return ErrNoSessionCookie
	}
	fullURL := url.URL{
		Scheme: p.serverURL.Scheme,
		Host:   p.serverURL.Host,
		Path:   urlPrefix + pathCredentials + fmt.Sprintf("/%d", rowID),
	}
	if err := p.deleteSecret(fullURL.String()); err != nil {
		return err
	}
	return nil
}

func (p *ServerProxy) DeleteNote(rowID int) error {
	if p.sessionCookie == "" {
		return ErrNoSessionCookie
	}
	fullURL := url.URL{
		Scheme: p.serverURL.Scheme,
		Host:   p.serverURL.Host,
		Path:   urlPrefix + pathNote + fmt.Sprintf("/%d", rowID),
	}
	if err := p.deleteSecret(fullURL.String()); err != nil {
		return err
	}
	return nil
}

func (p *ServerProxy) DeleteCard(rowID int) error {
	if p.sessionCookie == "" {
		return ErrNoSessionCookie
	}
	fullURL := url.URL{
		Scheme: p.serverURL.Scheme,
		Host:   p.serverURL.Host,
		Path:   urlPrefix + pathCard + fmt.Sprintf("/%d", rowID),
	}
	if err := p.deleteSecret(fullURL.String()); err != nil {
		return err
	}
	return nil
}

func (p *ServerProxy) DeleteFile(rowID int) error {
	if p.sessionCookie == "" {
		return ErrNoSessionCookie
	}
	fullURL := url.URL{
		Scheme: p.serverURL.Scheme,
		Host:   p.serverURL.Host,
		Path:   urlPrefix + pathFile + fmt.Sprintf("/%d", rowID),
	}
	if err := p.deleteSecret(fullURL.String()); err != nil {
		return err
	}
	return nil
}

func (p *ServerProxy) AddFile(file *common.FileReq) (int, error) {
	if p.sessionCookie == "" {
		return 0, ErrNoSessionCookie
	}
	fullURL := url.URL{
		Scheme: p.serverURL.Scheme,
		Host:   p.serverURL.Host,
		Path:   urlPrefix + pathFile,
	}
	payload, err := json.Marshal(file)
	return p.addSecret(false, fullURL.String(), payload, err)
}

func (p *ServerProxy) UpdateFile(file *common.FileReq) error {
	if p.sessionCookie == "" {
		return ErrNoSessionCookie
	}
	fullURL := url.URL{
		Scheme: p.serverURL.Scheme,
		Host:   p.serverURL.Host,
		Path:   urlPrefix + pathFile,
	}
	payload, err := json.Marshal(file)
	_, err = p.addSecret(true, fullURL.String(), payload, err)
	return err
}

func (p *ServerProxy) AddCard(card *common.CardReq) (int, error) {
	if p.sessionCookie == "" {
		return 0, ErrNoSessionCookie
	}
	fullURL := url.URL{
		Scheme: p.serverURL.Scheme,
		Host:   p.serverURL.Host,
		Path:   urlPrefix + pathCard,
	}
	payload, err := json.Marshal(card)
	return p.addSecret(false, fullURL.String(), payload, err)
}

func (p *ServerProxy) UpdateCard(card *common.CardReq) error {
	if p.sessionCookie == "" {
		return ErrNoSessionCookie
	}
	fullURL := url.URL{
		Scheme: p.serverURL.Scheme,
		Host:   p.serverURL.Host,
		Path:   urlPrefix + pathCard,
	}
	payload, err := json.Marshal(card)
	_, err = p.addSecret(true, fullURL.String(), payload, err)
	return err
}

func (p *ServerProxy) signInOrUp(
	login, password, path string,
	entropy *common.Entropy,
) (*clientEntities.UserDataResponse, error) {
	fullURL := url.URL{
		Scheme: p.serverURL.Scheme,
		Host:   p.serverURL.Host,
		Path:   urlPrefix + path,
	}
	authReqBody := common.UserCredentials{Login: login, Password: password}
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
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Add("Content-Type", "application/json")

	body, cookies, err := getResponse(req)
	if err != nil {
		return nil, err
	}
	var userEntropy *common.Entropy
	if entropy == nil { // for sign-up requests
		if err := json.Unmarshal(body, &userEntropy); err != nil {
			return nil, fmt.Errorf("failed to read server response: %v", err)
		}
	}
	for _, cookie := range cookies {
		if cookie.Name == common.SessionCookieName {
			return &clientEntities.UserDataResponse{
				Entropy:       userEntropy,
				SessionCookie: cookie.Value,
			}, nil
		}
	}
	return nil, fmt.Errorf("incorrect response from server: no session cookie set")
}

func (p *ServerProxy) addSessionCookie(req *http.Request) {
	cookie := &http.Cookie{
		Name:  common.SessionCookieName,
		Value: p.sessionCookie,
		Path:  "/",
	}
	req.AddCookie(cookie)
}

func (p *ServerProxy) addSecret(exists bool, url string, payload []byte, err error) (int, error) {
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %v", err)
	}

	var buf bytes.Buffer
	if _, err := buf.Write(payload); err != nil {
		return 0, fmt.Errorf("failed to create request: %v", err)
	}
	var method string
	if exists {
		method = "PUT"
	} else {
		method = "POST"
	}
	req, err := http.NewRequest(method, url, &buf)
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Add("Content-Type", "application/json")
	p.addSessionCookie(req)

	body, _, err := getResponse(req)
	if err != nil {
		return 0, err
	}
	rowID, err := strconv.Atoi(string(body))
	if err != nil {
		return 0, fmt.Errorf("incorrect response from server: can't interpret response as row ID")
	}
	return rowID, nil
}

func (p *ServerProxy) deleteSecret(url string) error {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	p.addSessionCookie(req)
	resp, err := (&http.Client{}).Do(req)
	if resp == nil || err != nil || resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed: %v", err)
	}
	return nil
}

func getResponse(req *http.Request) ([]byte, []*http.Cookie, error) {
	resp, err := (&http.Client{}).Do(req)
	if resp == nil {
		return nil, nil, fmt.Errorf("no response from the server")
	}
	body, bodyErr := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil || bodyErr != nil || resp.StatusCode != http.StatusOK {
		var respErr error
		if bodyErr == nil {
			respErr = fmt.Errorf("server error response: %s", string(body))
		} else {
			respErr = fmt.Errorf("failed to read body: %v", bodyErr)
		}
		if err == nil {
			err = respErr
		} else {
			err = errors.Wrap(err, respErr.Error())
		}
		return nil, nil, fmt.Errorf("request failed: %v", err)
	}
	return body, resp.Cookies(), nil
}
