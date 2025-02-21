package static

import "net/http"

type StaticHandler struct{}

func NewStaticHandler(router *http.ServeMux) {
	handler := &StaticHandler{}

	router.HandleFunc("/", handler.GetStatic)
}

func (handler *StaticHandler) GetStatic(res http.ResponseWriter, req *http.Request) {
	fs := http.FileServer(http.Dir("./web"))
	http.StripPrefix("/", fs).ServeHTTP(res, req)
}
