package main

import (
	"log"
	"net/http"

	"github.com/emicklei/go-restful"
)

type Workspace struct {
	Id, Name string
}

type wsResource struct {
	// normally one would use DAO (data access object)
	ws map[string]Workspace
}

func (wr wsResource) Register(container *restful.Container) {
	// 创建新的WebService
	wsService := new(restful.WebService)

	// 设定WebService对应的路径("/workspaces")和支持的MIME类型(restful.MIME_XML/ restful.MIME_JSON)
	wsService.
		Path("/workspaces").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML)

	// 添加路由： GET /{ws-id} --> ws.findWorkspace
	wsService.Route(wsService.GET("/{ws-id}").To(wr.findWorkspace))

	// 添加路由： POST / --> ws.updateWorkspace
	wsService.Route(wsService.POST("").To(wr.updateWorkspace))

	// 添加路由： PUT /{ws-id} --> ws.createWorkspace
	wsService.Route(wsService.PUT("/{ws-id}").To(wr.createWorkspace))

	// 添加路由： DELETE /{ws-id} --> ws.removeWorkspace
	wsService.Route(wsService.DELETE("/{ws-id}").To(wr.removeWorkspace))

	// 将初始化好的WebService添加到Container中
	container.Add(wsService)
}

// GET http://localhost:8080/workspaces/1
func (wr wsResource) findWorkspace(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("ws-id")
	workspace := wr.ws[id]
	if len(workspace.Id) == 0 {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "Workspace could not be found.")
	} else {
		response.WriteEntity(workspace)
	}
}

// POST http://localhost:8080/workspaces
func (wr *wsResource) updateWorkspace(request *restful.Request, response *restful.Response) {
	workspace := new(Workspace)
	err := request.ReadEntity(&workspace)
	if err == nil {
		wr.ws[workspace.Id] = *workspace
		response.WriteEntity(workspace)
	} else {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
	}
}

// PUT http://localhost:8080/workspaces/1
func (wr *wsResource) createWorkspace(request *restful.Request, response *restful.Response) {
	workspace := Workspace{Id: request.PathParameter("ws-id")}
	err := request.ReadEntity(&workspace)
	if err == nil {
		wr.ws[workspace.Id] = workspace
		response.WriteHeader(http.StatusCreated)
		response.WriteEntity(workspace)
	} else {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
	}
}

// DELETE http://localhost:8080/workspaces/1
func (wr *wsResource) removeWorkspace(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("ws-id")
	delete(wr.ws, id)
}

func main() {
	// 创建一个空的Container
	wsContainer := restful.NewContainer()

	// 设定路由为CurlyRouter(快速路由)
	wsContainer.Router(restful.CurlyRouter{})

	// 创建自定义的Resource Handle(此处为UserResource)
	wr := wsResource{map[string]Workspace{}}

	// 创建WebService，并将WebService加入到Container中
	wr.Register(wsContainer)

	log.Printf("start listening on localhost:8080")
	server := &http.Server{Addr: ":8080", Handler: wsContainer}

	// 启动服务
	log.Fatal(server.ListenAndServe())
}
