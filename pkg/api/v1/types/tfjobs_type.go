package types

type ImportConfig struct {
	UserID    string   `jsob:"userID"`
	Partition string   `json:"partition"`
	Jobid     string   `json:"jobid"`
	Images    string   `json:"images"`
	Command   []string `json:"command"`
}

type StorageSpec struct {
	StorageID             int64  `json:"storageID"`
	MountPath             string `json:"mountPath"`
	ReadOnly              bool   `json:"ReadOnly"`
	PersistentVolumeClaim string `json:persistentVolumeClaim`
}
