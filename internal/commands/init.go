package commands

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Cagangedik/cli-tool/internal/api"
	"github.com/Cagangedik/cli-tool/internal/config"
	"github.com/spf13/cobra"
)

func NewInitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize Hopsule in the current directory",
		Long: `Initialize Hopsule in the current directory by creating a .hopsule file.

This connects your local project directory to a Hopsule project.
You can either select an existing project or create a new one.

The .hopsule file contains:
  - Project ID and slug
  - Organization ID and slug

After initialization, CLI commands will automatically use this project.`,
		RunE: runInit,
	}

	cmd.Flags().String("project", "", "Project ID to use (skip interactive selection)")
	cmd.Flags().String("org", "", "Organization ID to use (skip interactive selection)")
	cmd.Flags().Bool("force", false, "Overwrite existing .hopsule file")

	return cmd
}

func runInit(cmd *cobra.Command, args []string) error {
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

	// Check for existing .hopsule file
	force, _ := cmd.Flags().GetBool("force")
	if config.ProjectConfigExists() && !force {
		existingCfg, path, _ := config.LoadProjectConfig()
		if existingCfg != nil {
			fmt.Println("┌─────────────────────────────────────────┐")
			fmt.Println("│     Project Already Initialized         │")
			fmt.Println("└─────────────────────────────────────────┘")
			fmt.Println()
			fmt.Printf("Found existing .hopsule file:\n")
			fmt.Printf("  Project: %s (%s)\n", existingCfg.Project.Name, existingCfg.Project.Slug)
			fmt.Printf("  Org:     %s (%s)\n", existingCfg.Project.Organization.Name, existingCfg.Project.Organization.Slug)
			fmt.Printf("  Path:    %s\n", path)
			fmt.Println()
			fmt.Println("Use --force to reinitialize.")
			return nil
		}
	}

	client := api.NewClient(cfg)

	fmt.Println("┌─────────────────────────────────────────┐")
	fmt.Println("│         Initialize Hopsule Project      │")
	fmt.Println("└─────────────────────────────────────────┘")
	fmt.Println()

	// Fetch user data
	meResp, err := client.GetMe()
	if err != nil {
		return fmt.Errorf("failed to fetch user data: %w", err)
	}

	if len(meResp.Organizations) == 0 {
		fmt.Println("You don't have any organizations yet.")
		fmt.Println("Create one at https://hopsule.com/onboarding first.")
		return nil
	}

	// Build maps for lookups
	orgMap := make(map[string]*api.Organization)
	for _, org := range meResp.Organizations {
		orgMap[org.ID] = org
	}

	// Select organization
	var selectedOrg *api.Organization
	orgID, _ := cmd.Flags().GetString("org")

	if orgID != "" {
		selectedOrg = orgMap[orgID]
		if selectedOrg == nil {
			return fmt.Errorf("organization %s not found", orgID)
		}
	} else if len(meResp.Organizations) == 1 {
		selectedOrg = meResp.Organizations[0]
		fmt.Printf("Using organization: %s (@%s)\n", selectedOrg.Name, selectedOrg.Slug)
	} else {
		fmt.Println("Select an organization:")
		fmt.Println()
		for i, org := range meResp.Organizations {
			fmt.Printf("  [%d] %s (@%s)\n", i+1, org.Name, org.Slug)
		}
		fmt.Println()
		fmt.Print("Enter number: ")

		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		idx, err := strconv.Atoi(input)
		if err != nil || idx < 1 || idx > len(meResp.Organizations) {
			return fmt.Errorf("invalid selection")
		}
		selectedOrg = meResp.Organizations[idx-1]
	}

	fmt.Println()

	// Filter projects by organization
	var orgProjects []*api.Project
	for _, proj := range meResp.Projects {
		if proj.OrganizationID == selectedOrg.ID {
			orgProjects = append(orgProjects, proj)
		}
	}

	// Select or create project
	var selectedProject *api.Project
	projectID, _ := cmd.Flags().GetString("project")

	if projectID != "" {
		for _, proj := range orgProjects {
			if proj.ID == projectID {
				selectedProject = proj
				break
			}
		}
		if selectedProject == nil {
			return fmt.Errorf("project %s not found in organization %s", projectID, selectedOrg.Name)
		}
	} else {
		fmt.Println("Select a project:")
		fmt.Println()
		for i, proj := range orgProjects {
			desc := ""
			if proj.Description != "" {
				desc = " - " + truncate(proj.Description, 30)
			}
			fmt.Printf("  [%d] %s%s\n", i+1, proj.Name, desc)
		}
		fmt.Println()
		fmt.Printf("  [%d] Create new project\n", len(orgProjects)+1)
		fmt.Println()
		fmt.Print("Enter number: ")

		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		idx, err := strconv.Atoi(input)
		if err != nil || idx < 1 || idx > len(orgProjects)+1 {
			return fmt.Errorf("invalid selection")
		}

		if idx == len(orgProjects)+1 {
			// Create new project - redirect to web
			fmt.Println()
			fmt.Println("To create a new project, visit:")
			fmt.Printf("  https://hopsule.com/workspace/%s/new-project\n", selectedOrg.Slug)
			fmt.Println()
			fmt.Println("Then run 'hopsule init' again to connect to it.")
			return nil
		}

		selectedProject = orgProjects[idx-1]
	}

	// Create .hopsule file
	projectCfg := &config.ProjectConfig{
		Version: config.HopsuleFileVersion,
		Project: config.ProjectInfo{
			ID:   selectedProject.ID,
			Slug: selectedProject.Slug,
			Name: selectedProject.Name,
			Organization: config.OrganizationInfo{
				ID:   selectedOrg.ID,
				Slug: selectedOrg.Slug,
				Name: selectedOrg.Name,
			},
		},
	}

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	if err := config.SaveProjectConfig(cwd, projectCfg); err != nil {
		return fmt.Errorf("failed to save .hopsule file: %w", err)
	}

	// Also update global config with this project
	cfg.Project = selectedProject.ID
	cfg.Organization = selectedOrg.ID
	config.SaveConfig(cfg)

	fmt.Println()
	fmt.Println("┌─────────────────────────────────────────┐")
	fmt.Println("│     ✓ Project Initialized!              │")
	fmt.Println("└─────────────────────────────────────────┘")
	fmt.Println()
	fmt.Printf("Project:      %s\n", selectedProject.Name)
	fmt.Printf("Organization: %s\n", selectedOrg.Name)
	fmt.Printf("Config file:  .hopsule\n")
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("  • Run 'hopsule list' to see decisions")
	fmt.Println("  • Run 'hopsule create' to create a new decision")
	fmt.Println("  • Run 'hopsule status' to see project statistics")
	fmt.Println()
	fmt.Println("Tip: Add .hopsule to your .gitignore if you don't want to share project settings.")

	return nil
}
