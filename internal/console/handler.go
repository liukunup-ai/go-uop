package console

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (s *Server) handleDevices(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED")
		return
	}

	devices, err := s.deviceMgr.ListDevices()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "LIST_DEVICES_FAILED")
		return
	}

	writeJSON(w, map[string]interface{}{
		"devices": devices,
	})
}

func (s *Server) handleConnect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED")
		return
	}

	var req Device
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST")
		return
	}

	device, err := s.deviceMgr.ConnectDevice(&req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "CONNECT_FAILED")
		return
	}

	info, _ := device.Info()
	writeJSON(w, map[string]interface{}{
		"success": true,
		"device":  req,
		"info":    info,
	})
}

func (s *Server) handleDeviceOps(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/devices/")

	parts := strings.Split(path, "/")
	if len(parts) < 2 {
		writeError(w, http.StatusBadRequest, "INVALID_PATH")
		return
	}

	deviceID := parts[0]
	operation := parts[1]

	switch operation {
	case "screenshot":
		s.handleScreenshot(w, r, deviceID)
	case "info":
		s.handleDeviceInfo(w, r, deviceID)
	case "commands":
		s.handleCommands(w, r, deviceID)
	default:
		writeError(w, http.StatusNotFound, "OPERATION_NOT_FOUND")
	}
}

func (s *Server) handleScreenshot(w http.ResponseWriter, r *http.Request, deviceID string) {
	device, err := s.deviceMgr.GetConnected(deviceID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "DEVICE_NOT_CONNECTED")
		return
	}

	img, err := device.Screenshot()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "SCREENSHOT_FAILED")
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Write(img)
}

func (s *Server) handleDeviceInfo(w http.ResponseWriter, r *http.Request, deviceID string) {
	device, err := s.deviceMgr.GetConnected(deviceID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "DEVICE_NOT_CONNECTED")
		return
	}

	info, err := device.Info()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "INFO_FAILED")
		return
	}

	writeJSON(w, info)
}

func (s *Server) handleCommands(w http.ResponseWriter, r *http.Request, deviceID string) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED")
		return
	}

	var req CommandRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST")
		return
	}

	record, err := s.deviceMgr.ExecuteCommand(deviceID, req.Command, req.Params)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "COMMAND_FAILED")
		return
	}

	s.historyMgr.Add(record)

	writeJSON(w, record)
}

func (s *Server) handleHistory(w http.ResponseWriter, r *http.Request) {
	history := s.historyMgr.GetAll()
	writeJSON(w, map[string]interface{}{
		"history": history,
	})
}

func (s *Server) handleYamlExport(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	ids := query["ids"]

	var records []CommandRecord
	if len(ids) > 0 {
		records = s.historyMgr.GetSelected(ids)
	} else {
		records = s.historyMgr.GetAll()
	}

	name := query.Get("name")
	if name == "" {
		name = "debug-session"
	}

	yamlData, err := ExportToYaml(records, name)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "EXPORT_FAILED")
		return
	}

	w.Header().Set("Content-Type", "text/yaml")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.yaml", name))
	w.Write(yamlData)
}

func writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, code string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": map[string]interface{}{
			"code": code,
		},
	})
}
