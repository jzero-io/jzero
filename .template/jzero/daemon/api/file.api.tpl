syntax = "v1"

type DownloadRequest {
  File string `path:"file"`
}

type UploadResponse {
  Code int `json:"code"`
}

@server (
	prefix: /api/v1
	group:  file
)
service {{ .APP }} {
	@handler DownloadHandler
	get /static/:file(DownloadRequest)

	@handler UploadHandler
  	post /upload returns (UploadResponse)
}