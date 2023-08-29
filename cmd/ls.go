package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
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
	// 获取所有 SSM 容器列表
	ssmList := lib.GetSSMContainers()

	if len(ssmList) == 0 {
		fmt.Println("No Sub-Store Manager Docker Containers found")
		return
	}

	fmt.Println("Sub-Store Manager Docker Containers:")
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("ID", "Version", "Port", "Status", "Name")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, container := range ssmList {
		portStr := fmt.Sprintf("%s: %s->%s", container.NetworkType, container.HostPort, container.Port)
		tbl.AddRow(container.Id, container.Version, portStr, container.Status, container.Name)
	}

	tbl.Print()
}
