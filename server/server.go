package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"syscall"

	log "github.com/Sirupsen/logrus"
	"github.com/fishworks/api"
	"github.com/fishworks/api/auth"
	"github.com/julienschmidt/httprouter"
)

var (
	Apps   []*api.App
	Builds []*api.Build
	Users  []*auth.User
)

// HTTPServer is an API Server which listens and responds to HTTP requests.
type HTTPServer struct {
	srv *http.Server
	l   net.Listener
}

// Serve starts the HTTP server, accepting all new connections.
func (s *HTTPServer) Serve() error {
	return s.srv.Serve(s.l)
}

// Close shuts down the HTTP server, dropping all current connections.
func (s *HTTPServer) Close() error {
	return s.l.Close()
}

// ServeRequest processes a single HTTP request.
func (s *HTTPServer) ServeRequest(w http.ResponseWriter, req *http.Request) {
	s.srv.Handler.ServeHTTP(w, req)
}

// New sets up the required Server and does protocol specific checking.
func New(proto, addr string) (*HTTPServer, error) {
	switch proto {
	case "tcp":
		return setupTCPHTTP(addr)
	case "unix":
		return setupUnixHTTP(addr)
	default:
		return nil, fmt.Errorf("Invalid protocol format.")
	}
}

func setupTCPHTTP(addr string) (*HTTPServer, error) {
	r := createRouter()

	l, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &HTTPServer{&http.Server{Addr: addr, Handler: r}, l}, nil
}

func setupUnixHTTP(addr string) (*HTTPServer, error) {
	r := createRouter()

	if err := syscall.Unlink(addr); err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	mask := syscall.Umask(0777)
	defer syscall.Umask(mask)

	l, err := net.Listen("unix", addr)
	if err != nil {
		return nil, err
	}

	if err := os.Chmod(addr, 0660); err != nil {
		return nil, err
	}

	return &HTTPServer{&http.Server{Addr: addr, Handler: r}, l}, nil
}

func getApp(id string) *api.App {
	for _, app := range Apps {
		if app.ID == id {
			return app
		}
	}
	return nil
}

func createRouter() *httprouter.Router {
	r := httprouter.New()

	authNotRequiredMap := map[string]map[string]httprouter.Handle{
		"GET": {
			"/_ping": ping,
		},
	}

	for method, routes := range authNotRequiredMap {
		for route, funct := range routes {
			r.Handle(method, route, logRequestMiddleware(funct))
		}
	}

	authRequiredMap := map[string]map[string]httprouter.Handle{
		"GET": {
			"/apps":            getAppsJSON,
			"/apps/:id":        getAppJSON,
			"/apps/:id/builds": getAppBuildsJSON,
			"/apps/:id/logs":   getAppLogs,
		},
		"POST": {
			"/apps":            createApp,
			"/apps/:id/builds": createBuild,
		},
		"DELETE": {
			"/apps/:id": deleteApp,
		},
	}

	for method, routes := range authRequiredMap {
		for route, funct := range routes {
			r.Handle(method, route, logRequestMiddleware(funct))
		}
	}

	return r
}

func logRequestMiddleware(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		log.Infof("%s %s", r.Method, r.RequestURI)
		// Delegate request to the given handle
		h(w, r, p)
	}
}

// WriteJSON writes the value v to the http response stream as json with standard
// json encoding.
func WriteJSON(w http.ResponseWriter, v interface{}, code int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(v)
}

func ping(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Write([]byte{'P', 'O', 'N', 'G'})
}

func getAppsJSON(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if len(Apps) == 0 {
		w.WriteHeader(http.StatusNoContent)
	} else {
		if err := WriteJSON(w, Apps, http.StatusOK); err != nil {
			log.Error(err)
		}
	}
}

func getAppJSON(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if app := getApp(p.ByName("id")); app != nil {
		if err := WriteJSON(w, app, http.StatusOK); err != nil {
			log.Error(err)
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func getAppBuildsJSON(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var builds []*api.Build
	if app := getApp(p.ByName("id")); app != nil {
		for _, build := range Builds {
			if build.App == app {
				builds = append(builds, build)
			}
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if len(builds) == 0 {
		w.WriteHeader(http.StatusNoContent)
	} else {
		if err := WriteJSON(w, builds, http.StatusOK); err != nil {
			log.Error(err)
		}
	}
}

func createApp(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var (
		app  *api.App
		form struct {
			ID string `json:"id"`
		}
	)
	if r.Body != nil {
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&form); err != nil {
			// the request body is always non-nil (except in tests) but will return EOF immediately when no body is present.
			// http://golang.org/pkg/net/http/#Request
			if err != io.EOF {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("could not decode request: " + err.Error()))
				return
			}
		}
		app = api.NewApp(form.ID)
	} else {
		app = api.NewApp("")
	}
	Apps = append(Apps, app)
	w.WriteHeader(http.StatusCreated)
}

func createBuild(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var build *api.Build
	if r.Body != nil {
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&build); err != nil {
			// the request body is always non-nil (except in tests) but will return EOF immediately when no body is present.
			// http://golang.org/pkg/net/http/#Request
			if err != io.EOF {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("could not decode request: " + err.Error()))
				return
			}
		}
		app := getApp(p.ByName("id"))
		if app == nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("could not find app with id " + p.ByName("id")))
			return
		}
		// attach app to build
		build.App = app
		// add build to in-memory list
		Builds = append(Builds, build)
		release := app.NewRelease(build, nil)
		if err := release.Publish(); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(fmt.Sprintf("there was an error deploying this release: %v", err)))
			return
		}
		app.Log("released " + release.String())
	} else {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func getAppLogs(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	follow := r.URL.Query().Get("follow")
	for _, app := range Apps {
		if app.ID == p.ByName("id") {
			// hijack the connection if we want to "follow" the logs
			if follow == "true" {
				hj, ok := w.(http.Hijacker)
				if !ok {
					http.Error(w, "webserver doesn't support hijacking", http.StatusInternalServerError)
					return
				}
				conn, bufrw, err := hj.Hijack()
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				defer conn.Close()
				f, err := os.Open(app.LogPath)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				defer f.Close()
				go func() {
					for {
						io.Copy(conn, f)
					}
				}()
				bufrw.ReadString('\n')
			} else {
				http.ServeFile(w, r, app.LogPath)
			}
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}

func deleteApp(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	for i, app := range Apps {
		if app.ID == p.ByName("id") {
			Apps = append(Apps[:i], Apps[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
		}
	}
	w.WriteHeader(http.StatusNotFound)
}
