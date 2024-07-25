package toolbox

import (
	"fmt"
	"github.com/spf13/cobra"
	"slices"
)

var (
	showCount  bool
	showInMenu bool
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed ToolBox IDEs",
	RunE: func(cmd *cobra.Command, args []string) error {
		tools, err := ListToolboxTools(_ToolBoxDir, showInMenu)
		if err != nil {
			return err
		}
		if showCount {
			fmt.Println(len(tools))
		} else { // show list
			for _, tool := range tools {
				fmt.Printf("%-30s\t%-10s\n", tool.Name, tool.Version)
			}
		}
		return nil
	},
}

func init() {
	listCmd.Flags().BoolVarP(&showCount, "count", "c", false, "count the number of installed tools")
	listCmd.Flags().BoolVar(&showInMenu, "menu", false, "list the tools shown in the context menu")
}

// ListToolboxTools list local tools
func ListToolboxTools(dir string, showInMenu bool) ([]*Tool, error) {
	if !showInMenu {
		toolBox, err := GetAllTools(dir)
		if err != nil {
			return nil, err
		}
		return toolBox.Tools, err
	}

	toolbox, err := GetLatestTools(dir, _SortNames)
	if err != nil {
		return nil, err
	}

	items, exist, err := ReadSubCommands()
	if err != nil {
		return nil, err
	} else if !exist {
		return nil, nil
	}

	var tools []*Tool
	for _, tool := range toolbox.Tools {
		if slices.ContainsFunc(items, func(id string) bool { return tool.Id == id }) {
			tools = append(tools, tool)
		}
	}
	sortTools(tools, _SortNames)
	return tools, nil
}
