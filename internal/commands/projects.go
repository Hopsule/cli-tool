package commands

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/Cagangedik/cli-tool/internal/api"
	"github.com/Cagangedik/cli-tool/internal/config"
	"github.com/spf13/cobra"
)

func NewProjectsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "projects",
		Short: "List your projects",
		Long: `List all projects you have access to.

Shows project names, descriptions, and associated organizations.`,
		RunE: runProjects,
	}

	cmd.Flags().Bool("json", false, "Output as JSON")
	cmd.Flags().String("org", "", "Filter by organization ID")

	return cmd
}

func runProjects(cmd *cobra.Command, args []string) error {
	cfg, err := config.GetConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if !cfg.IsAuthenticated() {
		fmt.Println("Not logged in.")
		fmt.Println()
		fmt.Println("Run 'hopsule login' to sign in first.")
		return nil
	}

	client := api.NewClient(cfg)

	fmt.Println("Fetching projects...")
	fmt.Println()

	meResp, err := client.GetMe()
	if err != nil {
		return fmt.Errorf("failed to fetch projects: %w", err)
	}

	// Build org name lookup
	orgNames := make(map[string]string)
	for _, org := range meResp.Organizations {
		orgNames[org.ID] = org.Name
	}

	// Filter by org if specified
	orgFilter, _ := cmd.Flags().GetString("org")
	var filteredProjects []*api.Project
	for _, proj := range meResp.Projects {
		if orgFilter == "" || proj.OrganizationID == orgFilter {
			filteredProjects = append(filteredProjects, proj)
		}
	}

	if len(filteredProjects) == 0 {
		fmt.Println("┌─────────────────────────────────────────┐")
		fmt.Println("│             No Projects                 │")
		fmt.Println("└─────────────────────────────────────────┘")
		fmt.Println()
		if orgFilter != "" {
			fmt.Println("No projects found in this organization.")
		} else {
			fmt.Println("You don't have any projects yet.")
			fmt.Println("Create one at https://hopsule.com or run 'hopsule init'")
		}
		return nil
	}

	fmt.Println("┌─────────────────────────────────────────┐")
	fmt.Println("│              Your Projects              │")
	fmt.Println("└─────────────────────────────────────────┘")
	fmt.Println()

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "NAME\tORGANIZATION\tID")
	fmt.Fprintln(w, "────\t────────────\t──")

	for _, proj := range filteredProjects {
		orgName := orgNames[proj.OrganizationID]
		if orgName == "" {
			orgName = "Unknown"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\n",
			proj.Name,
			orgName,
			truncateID(proj.ID),
		)
	}
	w.Flush()

	fmt.Println()
	fmt.Printf("Total: %d project(s)\n", len(filteredProjects))
	fmt.Println()
	fmt.Println("Tip: Run 'hopsule init' in your project directory to connect it to Hopsule.")

	return nil
}
