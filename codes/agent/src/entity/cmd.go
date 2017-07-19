package entity

type Command struct {
    ID int `json:"id"`
    Type int `json:"type"`
    Args interface{} `json:"args"`
}

type CommandResult struct {
    ID int `json:"id"`
    Status int `json:"status"`
    Progress int `json:"progress"`
    Reason string `json:"reason"`
    Result interface{} `json:"result"`
}
