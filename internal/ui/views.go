package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ============================================================================
// VIEW — main dispatch
// ============================================================================

func (m model) View() string {
	if m.currentView == viewLogin {
		return m.renderLoginView()
	}
	return m.renderPageView()
}

// ============================================================================
// LOGIN VIEW
// ============================================================================

func (m model) renderLoginView() string {
	var s strings.Builder
	border := versionStyle.Render("──────────────────────────────────────────────────────────────────")

	s.WriteString("\n  " + border + "\n\n")

	s.WriteString(logoStyle.Render(`    ⢠⣶⣶⣶⣶⣶⣶⣶⣆          ⣴⣶⣶⣶⣶⣶⣶⣶⡄
    ⢸⣿⣯⣷⣿⢿⣾⣟⣿          ⣿⣿⣽⣾⡿⣷⡿⣯⡇
    ⢸⣿⣾⣟⣿⣟⣯⣿⣿          ⣿⣷⣿⣻⡿⣟⣿⣿⡇
    ⠘⢿⣻⣿⣽⣾⣿⣻⣽          ⣿⣷⣿⣿⢿⣷⣿⠿⠃
             ⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣷⡀
             ⣿⣿⣻⣷⣿⣿⣻⣿⣯⣿⡿⣯⣿⣻⣿⣄
    ⢠⣾⣿⣿⢿⣿⡿⣿⣿          ⣿⣽⣿⣽⣷⣿⡿⣟⡇
    ⢸⣿⢿⣻⣾⣿⣻⣿⣻          ⣿⣷⣿⣯⣿⣷⣿⣾⡇
    ⠘⠿⠿⠻⠿⠻⠽⠿⠾          ⠻⠷⠿⠿⠻⠾⠟⠿⠁`) + "\n\n")

	s.WriteString("  " + titleStyle.Render("Hopsule") + "  " + dimStyle.Render("Decision & Memory Layer") + "\n")
	s.WriteString("  " + dimStyle.Render("for AI teams & coding tools") + "  " + versionStyle.Render("v0.8.0") + "\n\n")

	if m.loading && !m.loginPolling {
		s.WriteString("  " + infoStyle.Render("Loading...") + "\n")
	} else if m.errorMsg != "" {
		s.WriteString("  " + lipgloss.NewStyle().Foreground(redColor).Render("Error: "+m.errorMsg) + "\n\n")
	}

	if m.loginPolling {
		spinner := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		s.WriteString("  " + accentStyle.Render(spinner[m.loginSpinnerIdx]) + " Waiting for browser authentication\n\n")
		if m.loginDeviceCode != "" {
			s.WriteString("  Device Code: " + selectedStyle.Render(m.loginDeviceCode) + "\n\n")
		}
		if m.loginAuthURL != "" {
			s.WriteString("  " + dimStyle.Render("If browser didn't open:") + "\n")
			s.WriteString("  " + logoStyle.Render(m.loginAuthURL) + "\n\n")
		}
		s.WriteString("  " + dimStyle.Render(fmt.Sprintf("Polling... %ds", m.loginPollCount*2)) + "\n")
	} else {
		if m.selected == 0 {
			s.WriteString("  " + selectedStyle.Render("> Login") + "  " + infoStyle.Render("(authenticate)") + "\n")
		} else {
			s.WriteString("    " + normalStyle.Render("Login") + "  " + infoStyle.Render("(authenticate)") + "\n")
		}
		s.WriteString("\n  " + dimStyle.Render("Press Enter to continue") + "\n")
	}

	s.WriteString("\n  " + border + "\n")
	s.WriteString("  " + helpStyle.Render("enter login  •  q quit") + "\n")
	return s.String()
}

// ============================================================================
// PAGE VIEW — unified frame for all non-login views
// ============================================================================

func (m model) renderPageView() string {
	var s strings.Builder
	border := versionStyle.Render("──────────────────────────────────────────────────────────────────")

	s.WriteString("  " + border + "\n")
	s.WriteString(m.renderBreadcrumb())
	s.WriteString("  " + border + "\n\n")

	if m.confirmDelete {
		action := "delete"
		if m.confirmDeleteType == "decision" {
			action = "deprecate"
		}
		s.WriteString("  " + lipgloss.NewStyle().Foreground(redColor).Bold(true).Render("Are you sure you want to "+action+" this "+m.confirmDeleteType+"?") + "\n\n")
		s.WriteString("  " + dimStyle.Render("Press [y] to confirm, any other key to cancel") + "\n")
	} else if m.inputMode {
		s.WriteString("  " + titleStyle.Render(m.inputPrompt) + "\n\n")
		s.WriteString("  " + m.inputValue + "\u2588\n\n")
		s.WriteString("  " + dimStyle.Render("Enter to confirm  •  Esc to cancel") + "\n")
	} else if m.loading {
		s.WriteString("  " + infoStyle.Render("Loading...") + "\n")
	} else {
		if m.errorMsg != "" {
			s.WriteString("  " + lipgloss.NewStyle().Foreground(redColor).Render("Error: "+m.errorMsg) + "\n")
		}
		switch m.currentView {
		case viewOrganizations:
			s.WriteString(m.renderOrganizationsContent())
		case viewProjects:
			s.WriteString(m.renderProjectsContent())
		case viewProjectMenu:
			s.WriteString(m.renderProjectMenuContent())
		case viewDashboard:
			s.WriteString(m.renderDashboardContent())
		case viewDecisions:
			s.WriteString(m.renderDecisionsContent())
		case viewMemories:
			s.WriteString(m.renderMemoriesContent())
		case viewCapsules:
			s.WriteString(m.renderCapsulesContent())
		case viewTasks:
			s.WriteString(m.renderTasksContent())
		case viewHopper:
			s.WriteString(m.renderHopperContent())
		case viewMCP:
			s.WriteString(m.renderMCPContent())
		}
	}

	s.WriteString("\n  " + border + "\n")
	s.WriteString(m.renderFooter())
	return s.String()
}

// ============================================================================
// BREADCRUMB
// ============================================================================

func (m model) renderBreadcrumb() string {
	var parts []string

	switch m.currentView {
	case viewOrganizations:
		parts = []string{"Hopsule", "Organizations"}
	case viewProjects:
		orgName := "..."
		if m.currentOrg != nil {
			orgName = m.currentOrg.Name
		}
		parts = []string{orgName, "Projects"}
	case viewProjectMenu:
		orgName := "..."
		projName := "..."
		if m.currentOrg != nil {
			orgName = m.currentOrg.Name
		}
		if m.currentProj != nil {
			projName = m.currentProj.Name
		}
		parts = []string{orgName, projName}
	default:
		orgName := "..."
		projName := "..."
		pageName := ""
		if m.currentOrg != nil {
			orgName = m.currentOrg.Name
		}
		if m.currentProj != nil {
			projName = m.currentProj.Name
		}
		switch m.currentView {
		case viewDashboard:
			pageName = "Dashboard"
		case viewDecisions:
			pageName = "Decisions"
		case viewMemories:
			pageName = "Memories"
		case viewCapsules:
			pageName = "Capsules"
		case viewTasks:
			pageName = "Tasks"
		case viewHopper:
			pageName = "Hopper"
		case viewMCP:
			pageName = "MCP"
		}
		parts = []string{orgName, projName, pageName}
	}

	var s strings.Builder
	s.WriteString("  ")
	for i, part := range parts {
		if i > 0 {
			s.WriteString(dimStyle.Render(" / "))
		}
		if i == len(parts)-1 {
			s.WriteString(titleStyle.Render(part))
		} else {
			s.WriteString(dimStyle.Render(part))
		}
	}

	if m.cfg != nil && m.cfg.IsAuthenticated() && m.cfg.User != nil {
		s.WriteString("  " + dimStyle.Render("|") + "  " + statusOnStyle.Render("●") + " " + dimStyle.Render(m.cfg.User.Name))
	}

	s.WriteString("\n")
	return s.String()
}

// ============================================================================
// NAVIGATION CONTENT VIEWS
// ============================================================================

func (m model) renderOrganizationsContent() string {
	var s strings.Builder
	if len(m.organizations) == 0 {
		s.WriteString("  " + dimStyle.Render("No organizations found.") + "\n")
		s.WriteString("  " + dimStyle.Render("Create one at app.hopsule.com") + "\n")
	} else {
		for i, org := range m.organizations {
			if i == m.selected {
				s.WriteString("  " + selectedStyle.Render("> "+org.Name) + " " + dimStyle.Render("@"+org.Slug) + "\n")
			} else {
				s.WriteString("    " + normalStyle.Render(org.Name) + " " + dimStyle.Render("@"+org.Slug) + "\n")
			}
		}
		s.WriteString("\n")
		if m.selected == len(m.organizations) {
			s.WriteString("  " + selectedStyle.Render("> Logout") + "\n")
		} else {
			s.WriteString("    " + dimStyle.Render("Logout") + "\n")
		}
	}
	return s.String()
}

func (m model) renderProjectsContent() string {
	var s strings.Builder
	projects := m.getOrgProjects()
	if len(projects) == 0 {
		s.WriteString("  " + dimStyle.Render("No projects found.") + "\n")
		s.WriteString("  " + dimStyle.Render("Create one at app.hopsule.com") + "\n")
		return s.String()
	}

	greenDot := statusOnStyle.Render("●")
	grayDot := dimStyle.Render("●")

	for i, proj := range projects {
		desc := proj.Description
		if desc == "" {
			desc = "No description"
		}
		if len(desc) > 40 {
			desc = desc[:37] + "..."
		}

		isInitialized := proj.ID == m.initializedProjectID

		if isInitialized {
			// Active project — green dot, full brightness
			if i == m.selected {
				s.WriteString("  " + selectedStyle.Render("> "+proj.Name) + " " + greenDot + "  " + dimStyle.Render(desc) + "\n")
			} else {
				s.WriteString("    " + normalStyle.Render(proj.Name) + " " + greenDot + "  " + dimStyle.Render(desc) + "\n")
			}
		} else {
			// Not in this directory — gray dot, dimmed name
			if i == m.selected {
				s.WriteString("  " + dimStyle.Render("> "+proj.Name) + " " + grayDot + "  " + dimStyle.Render(desc) + "\n")
			} else {
				s.WriteString("    " + dimStyle.Render(proj.Name) + " " + grayDot + "  " + dimStyle.Render(desc) + "\n")
			}
		}
	}
	return s.String()
}

func (m model) renderProjectMenuContent() string {
	var s strings.Builder
	for i, item := range m.menuItems {
		if item.name == "" {
			s.WriteString("\n")
			continue
		}
		name := item.name
		for len(name) < 12 {
			name += " "
		}
		if i == m.selected {
			s.WriteString("  " + selectedStyle.Render("> "+item.icon+" "+name) + "  " + dimStyle.Render(item.description) + "\n")
		} else {
			s.WriteString("    " + normalStyle.Render(item.icon+" "+name) + "  " + dimStyle.Render(item.description) + "\n")
		}
	}
	return s.String()
}

func (m model) renderMCPContent() string {
	var s strings.Builder

	s.WriteString("  " + titleStyle.Render("AI Tool Integrations") + "\n")
	s.WriteString("  " + dimStyle.Render("Connect Hopsule to your AI coding tools via MCP") + "\n\n")

	greenDot := statusOnStyle.Render("●")
	grayDot := dimStyle.Render("●")
	dimDot := dimStyle.Render("○")

	for i, ide := range m.mcpIDEs {
		name := ide.name
		for len(name) < 18 {
			name += " "
		}

		var dot, status string
		if !ide.detected {
			dot = dimDot
			status = dimStyle.Render("Not detected")
		} else if ide.configured {
			dot = greenDot
			status = statusOnStyle.Render("Configured")
		} else {
			dot = grayDot
			status = dimStyle.Render("Not configured")
		}

		if i == m.selected {
			s.WriteString("  " + selectedStyle.Render("> "+name) + " " + dot + "  " + status + "\n")
		} else {
			s.WriteString("    " + normalStyle.Render(name) + " " + dot + "  " + status + "\n")
		}
	}

	s.WriteString("\n")
	if m.selected < len(m.mcpIDEs) {
		ide := m.mcpIDEs[m.selected]
		if !ide.detected && m.mcpConfirmManual {
			s.WriteString("  " + warningStyle.Render(ide.name+" not detected.") + " " + dimStyle.Render("Press Enter again to setup manually") + "\n")
		} else if !ide.detected {
			s.WriteString("  " + dimStyle.Render(ide.name+" not detected — press Enter to setup manually") + "\n")
		} else if ide.configured {
			s.WriteString("  " + dimStyle.Render("Press Enter to reinstall/update MCP configuration") + "\n")
		} else {
			s.WriteString("  " + dimStyle.Render("Press Enter to configure Hopsule MCP for "+ide.name) + "\n")
		}
	}

	return s.String()
}

// ============================================================================
// FOOTER
// ============================================================================

func (m model) renderFooter() string {
	if m.confirmDelete {
		action := "delete"
		if m.confirmDeleteType == "decision" {
			action = "deprecate"
		}
		return "  " + helpStyle.Render("[y] confirm "+action+"  •  any key to cancel") + "\n"
	}

	if m.inputMode {
		return "  " + helpStyle.Render("type  •  enter confirm  •  esc cancel") + "\n"
	}

	var help string
	viewToggle := ""
	if m.currentView == viewDecisions || m.currentView == viewMemories || m.currentView == viewCapsules {
		if m.currentViewMode == viewModeCard {
			viewToggle = "  •  [v] list view"
		} else {
			viewToggle = "  •  [v] card view  •  h/l prev/next page"
		}
	}

	switch m.currentView {
	case viewLogin:
		help = "enter login  •  q quit"
	case viewOrganizations:
		help = "j/k navigate  •  enter select  •  q quit"
	case viewProjects:
		help = "j/k navigate  •  enter select & init  •  " + statusOnStyle.Render("●") + " initialized  •  q/esc back"
	case viewProjectMenu:
		help = "j/k navigate  •  enter select  •  q/esc back"
	case viewMCP:
		if m.mcpConfirmManual {
			help = "enter confirm manual setup  •  j/k cancel  •  q/esc back"
		} else {
			help = "j/k navigate  •  enter install/setup  •  q/esc back"
		}
	case viewDashboard:
		help = "q/esc back"
	case viewDecisions:
		if m.detailView {
			help = "[a]ccept  •  [x]deprecate  •  [d]eprec  •  q/esc back"
		} else if m.searchMode {
			help = "type to search  •  enter confirm  •  esc cancel"
		} else {
			help = "hjkl navigate  •  enter detail  •  [/]search  •  tab filter  •  [n]ew  •  [a]ccept  •  [d]eprec" + viewToggle + "  •  q/esc back"
		}
	case viewMemories:
		if m.detailView {
			help = "[e]dit  •  [d]elete  •  q/esc back"
		} else if m.searchMode {
			help = "type to search  •  enter confirm  •  esc cancel"
		} else {
			help = "hjkl navigate  •  enter detail  •  [/]search  •  [n]ew  •  [e]dit  •  [d]elete" + viewToggle + "  •  q/esc back"
		}
	case viewCapsules:
		if m.detailView {
			help = "q/esc back"
		} else if m.searchMode {
			help = "type to search  •  enter confirm  •  esc cancel"
		} else {
			help = "hjkl navigate  •  enter detail  •  [/]search" + viewToggle + "  •  q/esc back"
		}
	case viewTasks:
		if m.searchMode {
			help = "type to search  •  enter confirm  •  esc cancel"
		} else {
			help = "j/k  •  tab view  •  [/]search  •  [n]ew  •  [e]dit  •  [t]oggle  •  [d]elete  •  q/esc back"
		}
	case viewHopper:
		if m.hopperSidebarFocus {
			help = "j/k navigate  •  enter open  •  ctrl+n new chat  •  d delete  •  tab chat  •  esc back"
		} else {
			help = "type message  •  enter send  •  tab history  •  ctrl+n new  •  ctrl+u/d scroll  •  esc back"
		}
	}
	return "  " + helpStyle.Render(help) + "\n"
}
