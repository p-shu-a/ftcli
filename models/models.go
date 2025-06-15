package models

type Header struct {
	FileName string `json:"file_name"`
	CheckSum string `json:"check_sum"`
	Nonce 	 []byte `json:"nonce,omitempty"`
	Salt 	 []byte `json:"salt,omitempty"`
	IV       []byte `json:"iv,omitempty"`
}