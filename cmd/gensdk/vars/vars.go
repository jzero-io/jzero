package vars

type (
	Scope                    string
	Resource                 string
	ResourceHTTPInterfaceMap map[Resource][]*HTTPInterface
)

// HTTPInterface parse grpc http options, go-zero api file
type HTTPInterface struct {
	Resource Resource

	Method     string
	URL        string
	MethodName string

	// body
	RequestBody  *RequestBody
	ResponseBody *ResponseBody

	// param
	PathParams  []*PathParam
	QueryParams []*QueryParam

	// comments
	Comments string

	IsStreamClient bool
	IsStreamServer bool
}

type RequestBody struct {
	MessageName string
	BodyName    string
	Name        string
	Type        string // proto or api
}

type ResponseBody struct {
	Name         string
	GoImportPath string
	RootPath     string
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
