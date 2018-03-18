package types

type ImportConfig struct {
	Partition string   `json:"partition"`
	Jobid     string   `json:"jobid"`
	Images    string   `json:"images"`
	Command   []string `json:"command"`
}
