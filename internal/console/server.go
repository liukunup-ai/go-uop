package console

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"
	"time"
)

// ConsoleAssets holds the embedded console frontend files.
//
//go:embed all:_out
var ConsoleAssets embed.FS

type Server struct {
	mux        *http.ServeMux
	addr       string
	devMode    bool
	deviceMgr  *DeviceManager
	historyMgr *HistoryManager
	serialMgr  *SerialManager
	iosMgr     *IOSManager
}

func NewServer(addr string, devMode bool) (*Server, error) {
	iosMgr := NewIOSManager()

	s := &Server{
		mux:        http.NewServeMux(),
		addr:       addr,
		devMode:    devMode,
		deviceMgr:  NewDeviceManager(),
		historyMgr: NewHistoryManager(100),
		serialMgr:  NewSerialManager(),
		iosMgr:     iosMgr,
	}
	s.setupRoutes()
	return s, nil
}

func (s *Server) Start() error {
	srv := &http.Server{
		Addr:         s.addr,
		Handler:      s.mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	return srv.ListenAndServe()
}

func (s *Server) setupRoutes() {
	// API routes - devices
	s.mux.HandleFunc("/api/devices", s.handleDevices)
	s.mux.HandleFunc("/api/devices/connect", s.handleConnect)
	s.mux.HandleFunc("/api/devices/", s.handleDeviceOps)

	// API routes - serial
	s.mux.HandleFunc("/api/serial/ports", s.handleSerialPorts)
	s.mux.HandleFunc("/api/serial/connect", s.handleSerialConnect)
	s.mux.HandleFunc("/api/serial/", s.handleSerialOps)

	// API routes - commands & history
	s.mux.HandleFunc("/api/commands/history", s.handleHistory)
	s.mux.HandleFunc("/api/export/yaml", s.handleYamlExport)

	// API routes - iOS
	s.mux.HandleFunc("/api/ios/devices", s.handleIosDevices)
	s.mux.HandleFunc("/api/ios/forward", s.handleIosForward)
	s.mux.HandleFunc("/api/ios/wda/start", s.handleIosWdaStart)
	s.mux.HandleFunc("/api/ios/wda/stop", s.handleIosWdaStop)

	// Frontend static files (SPA)
	if !s.devMode {
		s.mux.HandleFunc("/", s.serveFrontend)
	}
}

// serveFrontend serves the frontend SPA
func (s *Server) serveFrontend(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if path == "/" {
		path = "/index.html"
	}

	// Strip leading slash since embed.FS doesn't use leading slash
	path = strings.TrimPrefix(path, "/")

	data, err := fs.ReadFile(ConsoleAssets, "_out/"+path)
	if err != nil {
		// File not found, return index.html (SPA routing)
		data, err = fs.ReadFile(ConsoleAssets, "_out/index.html")
		if err != nil {
			http.NotFound(w, r)
			return
		}
	}

	// Set Content-Type
	contentType := getContentType(path)
	w.Header().Set("Content-Type", contentType)
	w.Write(data)
}

func getContentType(path string) string {
	ext := strings.ToLower(path[strings.LastIndex(path, "."):])
	switch ext {
	case ".html":
		return "text/html"
	case ".js":
		return "application/javascript"
	case ".css":
		return "text/css"
	case ".png":
		return "image/png"
	case ".svg":
		return "image/svg+xml"
	case ".json":
		return "application/json"
	default:
		return "text/plain"
	}
}
