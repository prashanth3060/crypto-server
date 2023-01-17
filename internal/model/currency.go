package model

type Currency struct {
	ID          string `json:"id"`
	FullName    string `json:"fullName"`
	Ask         string `json:"ask"`
	Bid         string `json:"bid"`
	Last        string `json:"last"`
	Open        string `json:"open"`
	Low         string `json:"low"`
	High        string `json:"high"`
	FeeCurrency string `json:"feeCurrency"`
}
