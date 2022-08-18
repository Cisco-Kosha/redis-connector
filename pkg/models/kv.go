package models

type KV struct {
	Value  string `json:"value,omitempty"`
	Expiry int    `json:"expiry,omitempty"`
}

type Set struct {
	Value []string `json:"value,omitempty"`
}

type Hash struct {
	Value []string `json:"value,omitempty"`
}

type List struct {
	Value []string `json:"value,omitempty"`
}
