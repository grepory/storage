package routes

import (
	"fmt"
	"net/http"
	"path"
	"strings"

	restful "github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful-openapi"
	"github.com/grepory/storage/apis/meta"
	"github.com/grepory/storage/server"
	"github.com/grepory/storage/storage"
)

const (
	apiRoot = "/api"
)

type ServiceBuilder interface {
	Build(obj meta.Object, strategy server.RestStrategy) *restful.WebService
}

func NewRestServiceBuilder(store storage.Store) *RestServiceBuilder {
	return &RestServiceBuilder{
		store: store,
	}
}

type RestServiceBuilder struct {
	store storage.Store
}

func (b RestServiceBuilder) Build(strategy server.RestStrategy) (ws *restful.WebService) {
	ws = new(restful.WebService)

	obj := strategy.PrepareForUpdate()

	gvk := meta.GetGroupVersionKind(obj)
	if gvk == nil {
		return nil
	}

	path := path.Join(apiRoot, strings.ToLower(gvk.GetGroup()))
	ws.Path(path).
		ApiVersion(gvk.GetVersion()).
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	tags := []string{gvk.GetKind()}

	ws.Route(ws.GET("/{version}/{kind}/{name}").To(b.getObject(strategy)).
		Doc(fmt.Sprintf("get a %s", gvk.GetKind())).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.PathParameter("version", "the api version").DataType("string")).
		Param(ws.PathParameter("kind", "the kind obj object").DataType("string")).
		Param(ws.PathParameter("name", fmt.Sprintf("name of the %s", gvk.GetKind())).DataType("string")).
		Writes(obj).
		Returns(http.StatusOK, http.StatusText(http.StatusOK), obj).
		Returns(http.StatusNotFound, http.StatusText(http.StatusNotFound), nil))

	ws.Route(ws.POST("").To(b.createObject(strategy)).
		Doc(fmt.Sprintf("create a %s", gvk.GetKind())).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(obj))

	return ws
}

func (b RestServiceBuilder) getObject(strategy server.RestStrategy) restful.RouteFunction {
	return func(req *restful.Request, resp *restful.Response) {
		obj := strategy.PrepareForUpdate()
		name := req.PathParameter("name")

		if err := b.store.Get(name, obj); err != nil {
			if err == storage.ErrNotFound {
				resp.WriteError(http.StatusNotFound, err)
			}
			resp.WriteError(http.StatusInternalServerError, err)
		}

		resp.WriteEntity(obj)
	}
}

func (b RestServiceBuilder) createObject(strategy server.RestStrategy) restful.RouteFunction {
	return func(req *restful.Request, resp *restful.Response) {
		// returns meta.Object
		obj := strategy.PrepareForUpdate()

		if err := req.ReadEntity(obj); err != nil {
			resp.WriteError(http.StatusBadRequest, err)
		}

		if err := strategy.Validate(obj); err != nil {
			resp.WriteError(http.StatusBadRequest, err)
		}

		if err := b.store.Create(obj.GetName(), obj); err != nil {
			resp.WriteError(http.StatusInternalServerError, err)
		}

		resp.WriteEntity(obj)
	}
}
