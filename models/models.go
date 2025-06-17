package models

type Header struct {
	FileName string `json:"file_name,omitempty"`
	CheckSum string `json:"check_sum,omitempty"`
	Nonce    []byte `json:"nonce,omitempty"`
	Salt     []byte `json:"salt,omitempty"`
	IV       []byte `json:"iv,omitempty"`
}
