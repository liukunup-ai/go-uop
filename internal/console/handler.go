package console

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/liukunup/go-uop/internal/report"
	"github.com/liukunup/go-uop/internal/runner"
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

func (s *Server) handleSerialPorts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED")
		return
	}

	ports, err := s.serialMgr.ListPorts()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "LIST_PORTS_FAILED")
		return
	}

	writeJSON(w, map[string]interface{}{
		"ports": ports,
	})
}

func (s *Server) handleSerialConnect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED")
		return
	}

	var req SerialConnectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST")
		return
	}

	conn, err := s.serialMgr.Connect(&req.Config)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "CONNECT_FAILED")
		return
	}

	writeJSON(w, map[string]interface{}{
		"success":    true,
		"connection": conn,
	})
}

func (s *Server) handleSerialOps(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/serial/")

	parts := strings.Split(path, "/")
	if len(parts) < 2 {
		writeError(w, http.StatusBadRequest, "INVALID_PATH")
		return
	}

	connID := parts[0]
	operation := parts[1]

	switch operation {
	case "disconnect":
		s.handleSerialDisconnect(w, r, connID)
	case "send":
		s.handleSerialSend(w, r, connID)
	case "sendByID":
		s.handleSerialSendByID(w, r, connID)
	case "commands":
		s.handleSerialCommands(w, r, connID)
	case "loadTable":
		s.handleSerialLoadTable(w, r, connID)
	default:
		writeError(w, http.StatusNotFound, "OPERATION_NOT_FOUND")
	}
}

func (s *Server) handleSerialDisconnect(w http.ResponseWriter, r *http.Request, connID string) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED")
		return
	}

	if err := s.serialMgr.Disconnect(connID); err != nil {
		writeError(w, http.StatusInternalServerError, "DISCONNECT_FAILED")
		return
	}

	writeJSON(w, map[string]interface{}{
		"success": true,
	})
}

func (s *Server) handleSerialSend(w http.ResponseWriter, r *http.Request, connID string) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED")
		return
	}

	var req SerialSendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST")
		return
	}

	result, err := s.serialMgr.SendRaw(connID, req.Data)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "SEND_FAILED")
		return
	}

	writeJSON(w, result)
}

func (s *Server) handleSerialSendByID(w http.ResponseWriter, r *http.Request, connID string) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED")
		return
	}

	var req SerialSendByIDRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST")
		return
	}

	result, err := s.serialMgr.SendByCommandID(connID, req.CommandID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "SEND_FAILED")
		return
	}

	writeJSON(w, result)
}

func (s *Server) handleSerialCommands(w http.ResponseWriter, r *http.Request, connID string) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED")
		return
	}

	cmds, err := s.serialMgr.ListCommands(connID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "LIST_COMMANDS_FAILED")
		return
	}

	writeJSON(w, map[string]interface{}{
		"commands": cmds,
	})
}

func (s *Server) handleSerialLoadTable(w http.ResponseWriter, r *http.Request, connID string) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED")
		return
	}

	var req SerialLoadTableRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST")
		return
	}

	var err error
	if req.FilePath != "" {
		err = s.serialMgr.LoadCommandTableFromFile(connID, req.FilePath)
	} else {
		err = s.serialMgr.LoadCommandTable(connID, req.YamlContent)
	}

	if err != nil {
		writeError(w, http.StatusInternalServerError, "LOAD_TABLE_FAILED")
		return
	}

	writeJSON(w, map[string]interface{}{
		"success": true,
	})
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

func (s *Server) handleIosDevices(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED")
		return
	}

	devices, err := s.iosMgr.ListDevices()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "LIST_DEVICES_FAILED")
		return
	}

	writeJSON(w, map[string]interface{}{
		"devices": devices,
	})
}

type IosForwardRequest struct {
	UDID       string `json:"udid"`
	DevicePort int    `json:"devicePort"`
	Ports      []int  `json:"ports"`
}

func (s *Server) handleIosForward(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED")
		return
	}

	var req IosForwardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST")
		return
	}

	if req.UDID == "" {
		writeError(w, http.StatusBadRequest, "UDID_REQUIRED")
		return
	}

	if len(req.Ports) == 0 {
		req.Ports = []int{req.DevicePort}
	}

	err := s.iosMgr.SetupPortForwarding(req.UDID, req.Ports)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, map[string]interface{}{
		"success": true,
		"command": s.iosMgr.GetIOSForwardCommand(req.UDID, req.Ports),
	})
}

type IosWdaStartRequest struct {
	UDID     string `json:"udid"`
	BundleID string `json:"bundleId"`
}

func (s *Server) handleIosWdaStart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED")
		return
	}

	var req IosWdaStartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST")
		return
	}

	err := s.iosMgr.LaunchWDA(req.UDID, req.BundleID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "WDA_START_FAILED")
		return
	}

	writeJSON(w, map[string]interface{}{
		"success": true,
	})
}

func (s *Server) handleIosWdaStop(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED")
		return
	}

	var req struct {
		UDID string `json:"udid"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST")
		return
	}

	if err := s.iosMgr.StopWDA(req.UDID); err != nil {
		writeError(w, http.StatusInternalServerError, "WDA_STOP_FAILED")
		return
	}

	writeJSON(w, map[string]interface{}{
		"success": true,
	})
}

type RunnerRunRequest struct {
	YamlContent string `json:"yamlContent"`
	DeviceID    string `json:"deviceId"`
}

func (s *Server) handleRunnerRun(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED")
		return
	}

	var req RunnerRunRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST")
		return
	}

	flow, err := runner.ParseFlowString(req.YamlContent)
	if err != nil {
		writeError(w, http.StatusBadRequest, "PARSE_FLOW_FAILED")
		return
	}

	if req.DeviceID != "" {
		if err := s.devicePool.SwitchDevice(req.DeviceID); err != nil {
			writeError(w, http.StatusBadRequest, "DEVICE_SWITCH_FAILED")
			return
		}
	}

	s.reportGen.StartTest(flow.Name)
	err = s.executor.ExecuteSuite(flow)
	status := "passed"
	if err != nil {
		status = "failed"
	}
	s.reportGen.EndTest(status, err)

	result := s.reportGen.Generate()
	writeJSON(w, map[string]interface{}{
		"success": true,
		"result":  result,
	})
}

func (s *Server) handleRunnerDebug(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED")
		return
	}

	var req RunnerRunRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST")
		return
	}

	flow, err := runner.ParseFlowString(req.YamlContent)
	if err != nil {
		writeError(w, http.StatusBadRequest, "PARSE_FLOW_FAILED")
		return
	}

	if req.DeviceID != "" {
		if err := s.devicePool.SwitchDevice(req.DeviceID); err != nil {
			writeError(w, http.StatusBadRequest, "DEVICE_SWITCH_FAILED")
			return
		}
	}

	debugger := runner.NewDebugger(flow.TestCases)
	stateMap := map[runner.DebugState]string{
		runner.DebugRunning:  "running",
		runner.DebugPaused:   "paused",
		runner.DebugStepping: "stepping",
		runner.DebugStopped:  "stopped",
	}
	stateStr := stateMap[debugger.State()]
	if stateStr == "" {
		stateStr = "unknown"
	}
	tc, step := debugger.CurrentPosition()
	writeJSON(w, map[string]interface{}{
		"success":        true,
		"totalTestCases": len(flow.TestCases),
		"currentTC":      tc,
		"currentStep":    step,
		"state":          stateStr,
	})
}

type RunnerTestRequest struct {
	DeviceType string `json:"deviceType"`
	Serial     string `json:"serial"`
}

func (s *Server) handleRunnerTest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED")
		return
	}

	var req RunnerTestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST")
		return
	}

	deviceID := req.Serial
	if deviceID == "" {
		deviceID = "test-device"
	}

	if err := s.devicePool.AddDevice(deviceID, req.DeviceType, req.Serial); err != nil {
		writeError(w, http.StatusInternalServerError, "ADD_DEVICE_FAILED")
		return
	}

	writeJSON(w, map[string]interface{}{
		"success":    true,
		"deviceId":   deviceID,
		"deviceType": req.DeviceType,
	})
}

type RunnerSuiteRequest struct {
	SuitePath string `json:"suitePath"`
}

func (s *Server) handleRunnerSuite(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED")
		return
	}

	var req RunnerSuiteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST")
		return
	}

	result, err := runner.ParseAndRunSuiteFile(req.SuitePath, s.devicePool, s.reportGen)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "RUN_SUITE_FAILED")
		return
	}

	writeJSON(w, map[string]interface{}{
		"success":     true,
		"totalSteps":  result.TotalSteps,
		"passedSteps": result.PassedSteps,
		"failedSteps": result.FailedSteps,
		"tcResults":   result.TCResults,
	})
}

func (s *Server) handleRunnerDevices(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED")
		return
	}

	devices := s.devicePool.ListDevices()
	var deviceList []map[string]interface{}
	for id, dev := range devices {
		deviceList = append(deviceList, map[string]interface{}{
			"id":      dev.ID,
			"type":    dev.Type,
			"serial":  dev.Serial,
			"current": id == s.devicePool.CurrentDevice().ID,
		})
	}

	writeJSON(w, map[string]interface{}{
		"devices":       deviceList,
		"currentDevice": s.devicePool.CurrentDevice().ID,
		"defaultDevice": s.devicePool.DefaultDevice().ID,
	})
}

type RunnerSwitchRequest struct {
	DeviceID string `json:"deviceId"`
}

func (s *Server) handleRunnerSwitch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED")
		return
	}

	var req RunnerSwitchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST")
		return
	}

	if err := s.devicePool.SwitchDevice(req.DeviceID); err != nil {
		writeError(w, http.StatusBadRequest, "DEVICE_NOT_FOUND")
		return
	}

	writeJSON(w, map[string]interface{}{
		"success":  true,
		"deviceId": req.DeviceID,
	})
}

type RunnerReportsRequest struct {
	Format string `json:"format"`
	Path   string `json:"path"`
}

func (s *Server) handleRunnerReports(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED")
		return
	}

	format := r.URL.Query().Get("format")
	if format == "" {
		format = "json"
	}

	result := s.reportGen.Generate()

	switch format {
	case "html":
		html, err := report.RenderHTML(result)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "RENDER_HTML_FAILED")
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	case "junit":
		junitXML, err := report.RenderJUnitXML(result)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "RENDER_JUNIT_FAILED")
			return
		}
		w.Header().Set("Content-Type", "application/xml")
		w.Write([]byte(junitXML))
	case "json":
		fallthrough
	default:
		jsonData, err := s.reportGen.ToJSON()
		if err != nil {
			writeError(w, http.StatusInternalServerError, "GENERATE_JSON_FAILED")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	}
}
