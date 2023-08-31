package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
	"sub-store-manager-cli/docker"
	"sub-store-manager-cli/lib"
)

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "list all sub-store docker containers",
	Run: func(cmd *cobra.Command, args []string) {
		listAllSSMContainer()
	},
}

func listAllSSMContainer() {
	fel, bel, err := docker.GetSSMContainers()
	if err != nil {
		lib.PrintError("Failed to get SSM containers:", err)
	}

	if len(fel) == 0 && len(bel) == 0 {
		fmt.Println("No Sub-Store Manager Front-End Docker Containers found")
		return
	}

	fmt.Println("Sub-Store Manager Docker Containers:")
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("Type", "ID", "Version", "Port", "Status", "Name")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, c := range fel {
		var portStr string
		if p, e := c.GetPortInfo(); e != nil {
			portStr = "none"
		} else {
			portStr = fmt.Sprintf("%s: %s->%s", p.Type, p.Public, p.Private)
		}
		tbl.AddRow(c.ContainerType, c.DockerContainer.ID, c.Version, portStr, c.DockerContainer.Status, c.Name)
	}

	tbl.Print()
}
