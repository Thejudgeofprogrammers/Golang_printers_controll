package interfaces

type PrintStats struct {
	Filename string `json:"filename"`
}

type QueryResponse struct {
	Result struct {
		Status struct {
			PrintStats PrintStats `json:"print_stats"`
		} `json:"status"`
	} `json:"result"`
}

type MetadataResponse struct {
	Result struct {
		EstimatedTime  float64 `json:"estimated_time"`
		PrintStartTime float64 `json:"print_start_time"`
	} `json:"result"`
}

type PrinterInfo struct {
	Success       string `json:"success"`
	DateEnd       string `json:"date_end"`
	EstimatedTime string `json:"estimated_time"`
}

// Error implements error.
func (p PrinterInfo) Error() string {
	panic("unimplemented")
}
