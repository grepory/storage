package main

import (
	"log"
	"net/http"
	"time"

	simplerest "github.com/grepory/storage/apis/simple/rest"
	"github.com/grepory/storage/runtime/codec"
	"github.com/grepory/storage/server"
	"github.com/grepory/storage/server/routes"
	"github.com/grepory/storage/storage/etcd"

	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful-openapi"
	"github.com/go-openapi/spec"
	"go.etcd.io/etcd/clientv3"
)

var (
	restStrategies = []server.RestStrategy{
		simplerest.SimpleRestStrategy{},
	}
)

/*
func updateSimple(req *restful.Request, resp *restful.Response) {
	s := simple.Simple{}

	if err := req.ReadEntity(&s); err != nil {
		resp.WriteError(http.StatusBadRequest, err)
	}

	// TODO: Update crashes
	if err := s.store.Update(s.GetName(), &s); err != nil {
		resp.WriteError(http.StatusInternalServerError, err)
	}
}
*/

func main() {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"http://localhost:2379"},
		DialTimeout: 2 * time.Second,
	})
	if err != nil {
		log.Fatal("Could not connect to etcd")
	}

	store := etcd.NewStorage(client, codec.UniversalCodec())
	//restful.Add(Simple{store: store}.WebService())

	builder := routes.NewRestServiceBuilder(store)
	for _, strategy := range restStrategies {
		restful.Add(builder.Build(strategy))
	}

	openAPIConfig := restfulspec.Config{
		WebServices: restful.RegisteredWebServices(),
		APIPath:     "/openapi.json",
		PostBuildSwaggerObjectHandler: enrichSwaggerObject,
	}
	restful.Add(restfulspec.NewOpenAPIService(openAPIConfig))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func enrichSwaggerObject(swo *spec.Swagger) {
	swo.Info = &spec.Info{
		InfoProps: spec.InfoProps{
			Title:       "Sensu API",
			Description: "Some description of the API.",
			Version:     "1.0.0",
		},
	}

	swo.Tags = []spec.Tag{
		spec.Tag{
			TagProps: spec.TagProps{
				Name:        "Simple",
				Description: "Manipulating Simple objects",
			},
		},
	}
}
