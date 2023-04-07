package http

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/google/martian/v3"
	"github.com/sirupsen/logrus"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type ScfModifier struct {
	apis   []string
	length int
	port   string
}

type httpRequest struct {
	Method string            `json:"method"`
	Url    string            `json:"url"`
	Header map[string]string `json:"headers"`
	Body   string            `json:"body"`
}

type httpResponse struct {
	Url    string            `json:"url"`
	Code   int               `json:"status_code"`
	Header map[string]string `json:"headers"`
	Body   string            `json:"content"`
}

func NewScfModifier(apis []string, lport string) (*ScfModifier, error) {
	length := len(apis)
	return &ScfModifier{apis: apis, length: length, port: lport}, nil
}

func (m *ScfModifier) ModifyRequest(req *http.Request) error {
	// Prevent scfproxy from recursively connecting to itself.
	remoteIp, _, _ := net.SplitHostPort(req.RemoteAddr)
	if remoteIp == req.URL.Hostname() && m.port == req.URL.Port() {
		ctx := martian.NewContext(req)
		ctx.SkipRoundTrip()
		return errors.New("Detecting recursive connection")
	}

	if req.Method == http.MethodConnect {
		return nil
	}

	headers := make(map[string]string)
	for k := range req.Header {
		headers[k] = strings.Join(req.Header.Values(k), ",")
	}

	rawBody, err := io.ReadAll(req.Body)
	if err != nil {
		logrus.Debugf("Error reading request body")
		return err
	}
	req.Body.Close()
	base64Body := base64.StdEncoding.EncodeToString(rawBody)

	hr := httpRequest{Method: req.Method, Url: req.URL.String(), Header: headers, Body: base64Body}
	data, err := json.Marshal(hr)
	if err != nil {
		return err
	}

	scfApi := m.pickRandomApi()
	logrus.Debugf("%s - %s", req.URL, scfApi)
	scfReq, err := http.NewRequest("POST", scfApi, bytes.NewReader(data))
	*req = *scfReq

	return nil
}

func (m *ScfModifier) ModifyResponse(res *http.Response) error {
	if res.Request.Method == http.MethodConnect {
		return nil
	}

	rawBody, err := io.ReadAll(res.Body)
	res.Body.Close()

	var hr httpResponse
	err = json.Unmarshal(rawBody, &hr)
	if err != nil {
		logrus.Debugf("Error Unmarshaling %s", string(rawBody))
		return err
	}

	res.StatusCode = hr.Code
	res.Status = fmt.Sprintf("%d %s", hr.Code, http.StatusText(hr.Code))

	res.Header = http.Header{}
	for k, v := range hr.Header {
		res.Header.Set(k, v)
	}

	body, err := base64.StdEncoding.DecodeString(hr.Body)
	if err != nil {
		logrus.Debugf("Error decoding base64 %s", hr.Body)
		return err
	}
	res.Body = io.NopCloser(bytes.NewReader(body))
	res.ContentLength = int64(len(body))

	return nil
}

func (m *ScfModifier) pickRandomApi() string {
	n := rand.Intn(m.length)
	return m.apis[n]
}
