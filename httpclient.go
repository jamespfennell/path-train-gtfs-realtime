package pathgtfsrt

import "net/http"

type HttpClient interface {
	Get(url string) (resp *http.Response, err error)
}
