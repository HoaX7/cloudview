package main

import (
	"bytes"
	"cloudview/agents/exporter/core/logging"
	"cloudview/agents/exporter/core/stats/sysinfo"
	"encoding/json"
	"fmt"
	"net/http"
)

func Recover(cleanups ...func()) {
	for _, cleanup := range cleanups {
		cleanup()
	}

	if p := recover(); p != nil {
		logging.Error(p)
	}
}

func runSafe(fn func()) {
	defer Recover()
	fn()
}

const url = "https://5f08-106-51-172-210.ngrok-free.app/metrics"

type Usage struct {
	Type    string `json:"type"`
	Percent int    `json:"percent"`
}

func reportMetrics(info sysinfo.SysInfo, usage ...Usage) {
	// send data to backend
	data := map[string]any{
		"usage":  usage,
		"sys":    info,
		"config": staticConfig,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("json.Marshal: error:", jsonData)
		return
	}

	req, reqErr := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if reqErr != nil {
		fmt.Println("report metric error:", reqErr)
		logging.Error("Unable to report metrics:", err.Error())
		return
	}
	req.Header.Set("authorization", staticConfig.PublicKey)
	req.Header.Set("content-type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("report metric error:", err)
		logging.Error("Unable to report metrics:", err.Error())
		return
	}

	defer res.Body.Close()
	return
}
