package servers

import (
	"asterism/caches"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"
	"path"
	"strconv"
)

type HttpServer struct {
	cache *caches.Cache
}

func NewHTTPServer(cache *caches.Cache) *HttpServer {
	return &HttpServer{
		cache: cache,
	}
}

func (hs *HttpServer) Run(address string) error {
	return http.ListenAndServe(address, hs.routerHandler())
}

func wrapUriWithVersion(uri string) string {
	return path.Join("/", APIVersion, uri)
}
func (hs *HttpServer) routerHandler() http.Handler {
	router := httprouter.New()
	router.GET(wrapUriWithVersion("/cache/:key"), hs.getHandler)
	router.PUT(wrapUriWithVersion("/cache/:key"), hs.SetHandler)
	router.DELETE(wrapUriWithVersion("/cache/:key"), hs.deleteHandler)
	router.GET(wrapUriWithVersion("/status"), hs.statusHandler)
	return router
}

func (hs *HttpServer) getHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	key := params.ByName("key")
	fmt.Println("Start get ", key)
	value, ok := hs.cache.Get(key)
	if !ok {
		writer.WriteHeader(http.StatusNotFound)
		return
	}
	writer.Write(value)
}

func (hs *HttpServer) SetHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	key := params.ByName("key")
	fmt.Println("Start Set ", key)
	value, err := ioutil.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	ttl, err := ttlOf(request)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
	}

	err = hs.cache.SetWithTTL(key, value, ttl)
	if err != nil {
		writer.WriteHeader(http.StatusRequestEntityTooLarge)
		writer.Write([]byte("Error :" + err.Error()))
		return
	}
	writer.WriteHeader(http.StatusCreated)
}

func ttlOf(request *http.Request) (int64, error) {
	ttls, ok := request.Header["Ttl"]
	if !ok || len(ttls) < 1 {
		return caches.NeverDie, nil
	}
	return strconv.ParseInt(ttls[0], 10, 64)
}

func (hs *HttpServer) deleteHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	key := params.ByName("key")
	fmt.Println("Start delete ", key)
	hs.cache.Delete(key)
}

func (hs *HttpServer) statusHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	status, err := json.Marshal(hs.cache.Status())
	fmt.Println("Start get status ", status)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.Write(status)
}
