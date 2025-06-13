package models

type Header struct {
	FileName string `json:"file_name"`
	CheckSum string `json:"check_sum"`
	IV       []byte `json:"iv"`
}
