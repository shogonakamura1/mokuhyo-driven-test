package model

type TreeResponse struct {
	Project Project `json:"project"`
	Nodes   []Node  `json:"nodes"`
	Edges   []Edge  `json:"edges"`
}
