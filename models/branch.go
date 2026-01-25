package models

// BranchInfo представляет информацию о ветке
type BranchInfo struct {
	ServiceName  string               `json:"serviceName"`
	Leaves       []Leaf               `json:"leaves"`
	AllowedUsers map[string]UserRight `json:"allowedUsers"`
}

// BranchMessage представляет сообщение для ветки
type BranchMessage struct {
	BranchName  string                 `json:"branchName"`
	LeafName    string                 `json:"leafName"`
	MessageData map[string]interface{} `json:"messageData"`
}
