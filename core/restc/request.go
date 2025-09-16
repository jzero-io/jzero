package restc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var (
	DefaultCodeField    = "code"
	DefaultDataField    = "data"
	DefaultMessageField = "msg"
)

// Request allows for building up a request to a server in a chained fashion.
// Any errors are stored until the end of your call, so you only have to
// check once.
type Request struct {
	c *client

	verb string
	path string

	queryParams []QueryParam

	// output
	err error

	// headers
	headers http.Header

	// body
	body io.Reader
}

func NewRequest(c *client) *Request {
	r := &Request{
		c:       c,
		headers: c.headers.Clone(),
	}
	return r
}

func (r *Request) Verb(verb string) *Request {
	r.verb = verb
	return r
}

func (r *Request) GetBody() io.Reader {
	return r.body
}

func (r *Request) GetParams() []QueryParam {
	return r.queryParams
}

func (r *Request) parseParam() string {
	if len(r.queryParams) == 0 {
		return ""
	}

	var queryParams strings.Builder
	queryParams.WriteString("?")
	for i, v := range r.queryParams {
		val := reflect.ValueOf(v.Value)
		kind := val.Kind()
		if kind == reflect.Slice || kind == reflect.Array {
			length := val.Len()
			for j := 0; j < length; j++ {
				value := val.Index(j).Interface()
				if cast.ToString(value) == "" {
					continue
				}
				va := url.QueryEscape(cast.ToString(value))
				if i == len(r.queryParams)-1 && j == length-1 {
					queryParams.WriteString(fmt.Sprintf("%s=%s", v.Name, va))
				} else {
					queryParams.WriteString(fmt.Sprintf("%s=%s&", v.Name, va))
				}
			}
		} else {
			if cast.ToString(v.Value) == "" {
				continue
			}
			va := url.QueryEscape(cast.ToString(v.Value))
			if i == len(r.queryParams)-1 {
				queryParams.WriteString(fmt.Sprintf("%s=%s", v.Name, va))
			} else {
				queryParams.WriteString(fmt.Sprintf("%s=%s&", v.Name, va))
			}
		}
	}
	return queryParams.String()
}

func (r *Request) AddHeader(key, value string) {
	r.c.lock.Lock()
	defer r.c.lock.Unlock()
	r.headers.Set(key, value)
}

type PathParam struct {
	Name  string
	Value any
}

// Path set path
func (r *Request) Path(path string, args ...PathParam) *Request {
	r.path = path
	for _, v := range args {
		val := reflect.ValueOf(v.Value)
		kind := val.Kind()
		if kind == reflect.Slice || kind == reflect.Array {
			js, err := json.Marshal(v.Value)
			if err != nil {
				panic(err)
			}
			path = strings.ReplaceAll(path, "{"+v.Name+"}", cast.ToString(js[1:len(js)-1]))
			path = strings.ReplaceAll(path, ":"+v.Name, cast.ToString(js[1:len(js)-1]))
		} else {
			path = strings.ReplaceAll(path, "{"+v.Name+"}", cast.ToString(v.Value))
			path = strings.ReplaceAll(path, ":"+v.Name, cast.ToString(v.Value))
		}
	}
	return r
}

type QueryParam struct {
	Name  string
	Value any
}

func (r *Request) Params(args ...QueryParam) *Request {
	if len(args) == 0 {
		return r
	}
	r.queryParams = args
	return r
}

// getUrl get url for request
func (r *Request) getUrl() (string, error) {
	if r.c.protocol == "" || r.c.addr == "" {
		return "", errors.New("invalid url, please check")
	}

	if r.c.protocol == "https" && r.c.port == "" {
		r.c.port = "443"
	} else if r.c.protocol == "http" && r.c.port == "" {
		r.c.port = "80"
	}

	return fmt.Sprintf("%s://%s:%s", r.c.protocol, r.c.addr, r.c.port+r.path+r.parseParam()), nil
}

// wsUrl get websocket url for request
func (r *Request) getWsUrl() (string, error) {
	if r.c.protocol == "" || r.c.addr == "" || r.c.port == "" {
		return "", errors.New("invalid url, you may not login")
	}

	// upgrade http to websocket proto
	if r.c.protocol == "https" {
		r.c.protocol = "wss"
	} else {
		r.c.protocol = "ws"
	}

	return fmt.Sprintf("%s://%s:%s", r.c.protocol, r.c.addr, r.c.port+r.path+r.parseParam()), nil
}

// Body makes the request use obj as the body. Optional.
// If obj is a string, try to read a file of that name.
// If obj is a []byte, send it directly.
// default marshal it
func (r *Request) Body(obj any) *Request {
	if r.err != nil {
		return r
	}

	switch t := obj.(type) {
	case io.Reader:
		r.body = t
	case io.ReadCloser:
		r.body = t
	case string:
		r.body = bytes.NewReader([]byte(t))
	case []byte:
		r.body = bytes.NewReader(t)
	default:
		data, err := json.Marshal(obj)
		if err != nil {
			r.err = err
			return r
		}
		r.body = bytes.NewReader(data)
	}
	return r
}

// Result contains the result of calling Request.Do().
type Result struct {
	body       []byte
	err        error
	statusCode int
	status     string
}

// Do format and executes the request. Returns a Result object for easy response
// processing.
//
// Error type:
// http.Client.Do errors are returned directly.
func (r *Request) Do(ctx context.Context) Result {
	if err := r.c.executeRequestMiddlewares(r); err != nil {
		return Result{err: err}
	}

	defaultUrl, err := r.getUrl()
	if err != nil {
		return Result{err: err}
	}

	request, err := http.NewRequestWithContext(ctx, r.verb, defaultUrl, r.body)
	if err != nil {
		return Result{err: err}
	}

	if r.c.client == nil {
		r.c.client = http.DefaultClient
	}

	if r.c.retryTimes == 0 {
		r.c.retryTimes = 1
	}
	request.Header = r.headers

	var rawResp *http.Response
	// if meet error, retry times that you set
	for k := 0; k < r.c.retryTimes; k++ {
		rawResp, err = r.doRequest(r.c.client, request)
		if err != nil {
			// sleep retry delay
			time.Sleep(r.c.retryDelay)
			continue
		}
		break
	}

	if err != nil {
		return Result{err: err}
	}

	if rawResp == nil {
		return Result{err: errors.New("http response is nil")}
	}

	data, err := io.ReadAll(rawResp.Body)
	if err != nil {
		return Result{err: err}
	}
	defer rawResp.Body.Close()

	return Result{
		body:       data,
		err:        err,
		statusCode: rawResp.StatusCode,
		status:     rawResp.Status,
	}
}

func (r *Request) WsConn(ctx context.Context) (*websocket.Conn, *http.Response, error) {
	wsUrl, err := r.getWsUrl()
	if err != nil {
		return nil, nil, err
	}
	return websocket.DefaultDialer.DialContext(ctx, wsUrl, r.c.headers)
}

func (r *Request) doRequest(client *http.Client, request *http.Request) (*http.Response, error) {
	res, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, errors.New("response is nil")
	}
	return res, nil
}

type WrapCodeMsgMapping struct {
	Code string
	Data string
	Msg  string
}

type IntoOptions struct {
	WrapCodeMsg        bool
	WrapCodeMsgMapping WrapCodeMsgMapping
}

// Into stores the result into obj, if possible. If obj is nil it is ignored.
func (r Result) Into(obj any, options *IntoOptions) error {
	if reflect.TypeOf(obj).Kind() != reflect.Ptr {
		return errors.New("object is not a ptr")
	}

	if r.err != nil {
		return r.err
	}

	if options != nil {
		if options.WrapCodeMsg && options.WrapCodeMsgMapping.Code == "" {
			options.WrapCodeMsgMapping.Code = DefaultCodeField
		}
		if options.WrapCodeMsg && options.WrapCodeMsgMapping.Data == "" {
			options.WrapCodeMsgMapping.Data = DefaultDataField
		}
		if options.WrapCodeMsg && options.WrapCodeMsgMapping.Msg == "" {
			options.WrapCodeMsgMapping.Msg = DefaultMessageField
		}
	}

	if r.StatusCode() != http.StatusOK {
		s := string(r.body)

		if len(s) == 0 {
			return fmt.Errorf("empty response body, status code: %d", r.StatusCode())
		}

		if options != nil && options.WrapCodeMsg {
			j, err := simplejson.NewJson(r.body)
			if err != nil {
				return fmt.Errorf("marsher json error: %v, response body: %v", err, r.body)
			}
			message, _ := j.Get(options.WrapCodeMsgMapping.Msg).String()
			return errors.New(message)
		}
		return errors.New(s)
	}

	j, err := simplejson.NewJson(r.body)
	if err != nil {
		return err
	}

	var marshalJSON []byte
	if options.WrapCodeMsg {
		code, err := j.Get(options.WrapCodeMsgMapping.Code).Int()
		if err != nil {
			return err
		}
		if code != http.StatusOK {
			message, _ := j.Get(options.WrapCodeMsgMapping.Msg).String()
			return errors.New(message)
		}
		data := j.Get(options.WrapCodeMsgMapping.Data)
		marshalJSON, err = data.MarshalJSON()
		if err != nil {
			return err
		}
	} else {
		marshalJSON, err = j.MarshalJSON()
		if err != nil {
			return err
		}
	}

	switch v := obj.(type) {
	case proto.Message:
		parser := protojson.UnmarshalOptions{
			DiscardUnknown: true,
		}
		err = parser.Unmarshal(marshalJSON, v)
	default:
		err = json.Unmarshal(marshalJSON, &obj)
	}

	if err != nil {
		return err
	}

	return nil
}

// Stream proto Stream way return io.ReadCloser
func (r *Request) Stream(ctx context.Context) (io.ReadCloser, error) {
	defaultUrl, err := r.getUrl()
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(ctx, r.verb, defaultUrl, r.body)
	if err != nil {
		return nil, err
	}

	request.Header = r.headers

	if r.c.client == nil {
		r.c.client = http.DefaultClient
	}

	if r.c.retryTimes == 0 {
		r.c.retryTimes = 1
	}

	var rawResp *http.Response
	// if meet error, retry times that you set
	for k := 0; k < r.c.retryTimes; k++ {
		rawResp, err = r.doRequest(r.c.client, request)
		if err != nil {
			// sleep retry delay
			time.Sleep(r.c.retryDelay)
			continue
		}
		break
	}

	if err != nil {
		return nil, err
	}

	if rawResp == nil {
		return nil, errors.New("empty resp")
	}

	if rawResp.StatusCode != 200 {
		return nil, errors.Errorf("unhealthy status code: [%d], status message: [%s]", rawResp.StatusCode, rawResp.Status)
	}

	return rawResp.Body, nil
}

func (r Result) RawResponse() ([]byte, error) {
	return r.body, r.err
}

// Error returns the error executing the request, nil if no error occurred.
func (r Result) Error() error {
	return r.err
}

// StatusCode returns the HTTP status code of the request. (Only valid if no
// error was returned.)
func (r Result) StatusCode() int {
	return r.statusCode
}

// Status returns the status executing the request
func (r Result) Status() string {
	return r.status
}
