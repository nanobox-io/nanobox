package server

import (
	"fmt"

	"github.com/spf13/cobra"

)

type CmdFunc func(ccmd *cobra.Command, args []string)

type AdminCmds struct {
	cmds map[string]CmdFunc
	response *Response
}

type Response struct {
	Output string
	ExitCode int
}

type Request struct {
	Name string
	Args []string
}

var commands = &AdminCmds{
	cmds: map[string]CmdFunc{},
	response: nil,
}

func AddCmd(key string, cmdFunc CmdFunc) {
	commands.cmds[key] = cmdFunc
}

func AddResponses(resp *Response) {
	commands.response = resp 
}

func (comm *AdminCmds) Run(req Request, resp *Response) error {
	cmd, ok := comm.cmds[req.Name]
	if !ok {
		return fmt.Errorf("invalid command %s", req.Name)
	}
	cmd(nil, req.Args)

	resp.Output = comm.response.Output
	resp.ExitCode = comm.response.ExitCode
	return nil
}