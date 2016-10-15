package main

import (
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"time"
)

type Levels string

const (
	Info  Levels = "INFO"
	Error        = "ERROR"
)

// ScanResponse object
type ScanResponse struct {
	ScanID       string `json:"scan_id"`
	Sha1         string `json:"sha1"`
	Resource     string `json:"resource"`
	ResponseCode int    `json:"response_code"`
	Sha256       string `json:"sha256"`
	Permalink    string `json:"permalink"`
	Md5          string `json:"md5"`
	VerboseMsg   string `json:"verbose_msg"`
}

// ReportResponse object
type ReportResponse struct {
	ScanResponse
	ScanDate  time.Time `json:"scan_date"`
	Positives int       `json:"positives"`
	Total     int       `json:"total"`
}

// VirusTotal object
type VirusTotal struct {
	APIKey string
}

const (
	scanURI   string = "https://www.virustotal.com/vtapi/v2/file/scan"
	rescanURI string = "https://www.virustotal.com/vtapi/v2/file/rescan"
	reportURI string = "https://www.virustotal.com/vtapi/v2/file/report"
)

// NewVirusTotal setup and return a new VirusTotal instance
func NewVirusTotal(key string) *VirusTotal {
	return &VirusTotal{key}
}

// ScanFile scans a new file
func (v *VirusTotal) ScanFile(filePath string) {
	p := make(map[string]string)
	p["apikey"] = v.APIKey

	SendFormData(scanURI, filePath, p, &ReportResponse{})
}

// ReScanFile scans a new file
func (v *VirusTotal) ReScanFile(filePath string) {
	p := make(map[string]string)
	p["apikey"] = v.APIKey

	SendFormData(rescanURI, filePath, p, &ReportResponse{})
}

// GetReport returns a report from a known hash
func (v *VirusTotal) GetReport(hash string, check bool) (isKnown bool, hasVirus bool) {
	p := make(map[string]string)
	p["apikey"] = v.APIKey
	p["resource"] = hash

	response := &ReportResponse{}

	for response.ResponseCode != 1 {
		SendFormData(reportURI, "", p, response)

		if check {
			return response.ResponseCode == 1, response.Positives > 0
		}

		time.Sleep(200 * time.Millisecond)
	}

	return false, response.Positives > 0
}

func (VirusTotal) GetFileBytesWithHash(filePath string) ([]byte, string) {
	b, _ := ioutil.ReadFile(filePath)
	hash := sha256.New()
	hash.Write(b)
	return b, hex.EncodeToString(hash.Sum(nil))
}
