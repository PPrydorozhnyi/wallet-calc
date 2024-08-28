package model

const (
	DEBIT  string = "DEBIT"
	CREDIT string = "CREDIT"
)

type Reason struct {
	Id        string            `json:"id"`
	Name      string            `json:"name"`
	Reference string            `json:"reference"`
	Meta      map[string]string `json:"metadata"`
}
