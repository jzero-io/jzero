package vars

type (
	Scope string

	Resource string

	ResourceHTTPInterfaceMap = map[Resource][]*HTTPInterface
)

// HTTPInterface parse grpc http options, go-zero api file
type HTTPInterface struct {
	Resource Resource

	Method     string
	URL        string
	MethodName string

	// body
	Request *Request

	Response *Response

	// param
	PathParams []*PathParam

	QueryParams []*QueryParam

	// comments
	Comments string

	IsStreamClient bool
	IsStreamServer bool

	WrapCodeMsg        bool
	WrapCodeMsgMapping *WrapCodeMsgMapping
}

type WrapCodeMsgMapping struct {
	Code string
	Data string
	Msg  string
}

type Request struct {
	Body         string // if proto. it takes effect. * or others
	RealBodyName string // if proto and body is not *. use it
	Name         string // request type name
	Type         string // proto or api
	Package      string // for example. types. *types.HelloParamRequest
	FullName     string
}

type Response struct {
	FullName string
	Package  string
}

type PathParam struct {
	Index  int
	Name   string
	GoName string
}

type QueryParam struct {
	GoName string
	Name   string
}
