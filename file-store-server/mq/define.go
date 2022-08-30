package mq

type TransferData struct {
	// URL string
	FileHash     string `json:"fileHash"`
	CurLocation  string `json:"curLocation"`
	DestLocation string `json:"destLocation"`
}
