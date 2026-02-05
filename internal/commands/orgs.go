package commands

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/Cagangedik/cli-tool/internal/api"
	"github.com/Cagangedik/cli-tool/internal/config"
	"github.com/spf13/cobra"
)

func NewOrgsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "orgs",
		Aliases: []string{"organizations"},
		Short:   "List your organizations",
		Long: `List all organizations you have access to.

Shows organization names, slugs, and your role in each.`,
		RunE: runOrgs,
	}

	cmd.Flags().Bool("json", false, "Output as JSON")

	return cmd
}

func runOrgs(cmd *cobra.Command, args []string) error {
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

	fmt.Println("Fetching organizations...")
	fmt.Println()

	meResp, err := client.GetMe()
	if err != nil {
		return fmt.Errorf("failed to fetch organizations: %w", err)
	}

	if len(meResp.Organizations) == 0 {
		fmt.Println("┌─────────────────────────────────────────┐")
		fmt.Println("│           No Organizations              │")
		fmt.Println("└─────────────────────────────────────────┘")
		fmt.Println()
		fmt.Println("You don't belong to any organizations yet.")
		fmt.Println("Create one at https://hopsule.com/onboarding")
		return nil
	}

	fmt.Println("┌─────────────────────────────────────────┐")
	fmt.Println("│           Your Organizations            │")
	fmt.Println("└─────────────────────────────────────────┘")
	fmt.Println()

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "NAME\tSLUG\tID")
	fmt.Fprintln(w, "────\t────\t──")

	for _, org := range meResp.Organizations {
		fmt.Fprintf(w, "%s\t@%s\t%s\n",
			org.Name,
			org.Slug,
			truncateID(org.ID),
		)
	}
	w.Flush()

	fmt.Println()
	fmt.Printf("Total: %d organization(s)\n", len(meResp.Organizations))

	return nil
}

func truncateID(id string) string {
	if len(id) <= 12 {
		return id
	}
	return id[:8] + "..."
}
