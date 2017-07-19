package entity

import (
    "time"
)

type Agent struct {
    UUID string `json:"uuid"`
    Hostname string `json:"hostname"`
    PID int `json:"pid"`
    OS string `json:"os"`
    Arch string `json:"arch"`
    Interfaces map[string][]string `json:"interfaces"`
    Time time.Time `json:"time"`
}
