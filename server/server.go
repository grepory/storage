package main

import (
	"log"
	"net/http"
	"time"

	"github.com/grepory/storage/apis/simple"
	"github.com/grepory/storage/runtime/codec"
	"github.com/grepory/storage/storage"
	"github.com/grepory/storage/storage/etcd"

	"github.com/emicklei/go-restful"
	"go.etcd.io/etcd/clientv3"
)

type Simple struct {
	simple.Simple
	store storage.Store
}

// TODO: Add OpenAPI
func (object Simple) WebService() *restful.WebService {
	ws := new(restful.WebService)

	ws.Path("/simple").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/{name}").To(object.getSimple).
		Doc("get a Simple").
		Param(ws.PathParameter("name", "name of the Simple").DataType("string")).
		Writes(Simple{}).
		Returns(http.StatusOK, http.StatusText(http.StatusOK), Simple{}).
		Returns(http.StatusNotFound, http.StatusText(http.StatusNotFound), nil))

	ws.Route(ws.POST("").To(object.createSimple).
		Doc("create a Simple").
		Reads(Simple{}))

	ws.Route(ws.PUT("").To(object.updateSimple).
		Doc("update a Simple").
		Reads(Simple{}))

	return ws
}

func (object Simple) getSimple(req *restful.Request, resp *restful.Response) {
	s := Simple{}
	name := req.PathParameter("name")

	// TODO: the store crashes if name is not found :(
	if err := object.store.Get(name, &s); err != nil {
		// TODO: switch on different types of errors to give a proper status code.
		resp.WriteError(http.StatusNotFound, err)
	}

	resp.WriteEntity(s)
}

func (object Simple) createSimple(req *restful.Request, resp *restful.Response) {
	s := Simple{}

	if err := req.ReadEntity(&s); err != nil {
		resp.WriteError(http.StatusBadRequest, err)
	}

	if err := object.store.Create(s.GetName(), s); err != nil {
		resp.WriteError(http.StatusInternalServerError, err)
	}

	resp.WriteHeader(http.StatusCreated)
}

func (object Simple) updateSimple(req *restful.Request, resp *restful.Response) {
	s := Simple{}

	if err := req.ReadEntity(&s); err != nil {
		resp.WriteError(http.StatusBadRequest, err)
	}

	// TODO: Update crashes
	if err := s.store.Update(s.GetName(), &s); err != nil {
		resp.WriteError(http.StatusInternalServerError, err)
	}
}

func main() {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"http://localhost:2379"},
		DialTimeout: 2 * time.Second,
	})
	if err != nil {
		log.Fatal("Could not connect to etcd")
	}

	store := etcd.NewStorage(client, codec.UniversalCodec())
	restful.Add(Simple{store: store}.WebService())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
