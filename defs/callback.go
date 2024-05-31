package defs

import "encoding/json"

type Command struct {
	Command string `json:"Command"`
	Argstr  string `json:"Argstr"`
}

func (c *Command) ToJSON() string {
	ans, _ := json.Marshal(c)
	return string(ans)
}

const (
	CMD_UNKNOWN = "unknown"
	CMD_APPROVE = "approve"
	CMD_REJECT  = "reject"
)

func NewApproveCommand(mcName string) *Command {
	return &Command{
		Command: CMD_APPROVE,
		Argstr:  mcName,
	}
}

func NewRejectCommand(mcName string) *Command {
	return &Command{
		Command: CMD_REJECT,
		Argstr:  mcName,
	}
}

func NewCommandFromJSON(jsonStr string) (*Command, error) {
	var cmd Command
	err := json.Unmarshal([]byte(jsonStr), &cmd)
	return &cmd, err
}
