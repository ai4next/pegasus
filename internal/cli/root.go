package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/ai4next/pegasus/internal/expert"
	"github.com/ai4next/pegasus/internal/global"
)

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}

// ensureDirs creates all runtime directories required by the agent.
func ensureDirs() error {
	cfg := global.Config()
	dirs := []string{
		cfg.Workspace,
		global.SkillsDir(),
		global.HooksDir(),
		global.BusDir(),
		global.OrchestratorDir(),
		global.PlansDir(),
		global.StateRootDir(),
		global.AgentStateDir("pegasus"),
		global.AgentStateDir("pegasus-evolver"),
		global.AgentStateDir("expert-evolver"),
		global.AgentStateDir("meta-evolver"),
		global.MemoryRootDir(),
		global.AgentMemoryDir("pegasus-evolver"),
		global.AgentMemoryDir("expert-evolver"),
		global.AgentMemoryDir("meta-evolver"),
		global.MemoryDir(),
		global.L2Dir(),
		global.MemoryL2Dir(global.AgentMemoryDir("pegasus-evolver")),
		global.MemoryL2Dir(global.AgentMemoryDir("expert-evolver")),
		global.MemoryL2Dir(global.AgentMemoryDir("meta-evolver")),
		global.SessionsDir(),
		global.AgentSessionsDir("pegasus-evolver"),
		global.AgentSessionsDir("expert-evolver"),
		global.AgentSessionsDir("meta-evolver"),
	}
	for _, d := range dirs {
		if d == "" {
			continue
		}
		if err := os.MkdirAll(d, 0755); err != nil {
			return fmt.Errorf("create runtime dir %s: %w", d, err)
		}
	}
	if err := expert.SeedDefaultCADExperts(global.ExpertsDir()); err != nil {
		return err
	}
	return nil
}

var rootCmd = &cobra.Command{
	Use:   "pegasus",
	Short: "Pegasus - AI-native CAD development agent",
	Long: `Pegasus is an AI-native CAD development agent built with Google ADK.
	It coordinates CAD specialists, modeling tools, persistent sessions,
	layered memory, a TUI interface, and autonomous reflection modes.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		configPath, _ := cmd.Flags().GetString("config")
		if _, err := global.LoadConfig(configPath); err != nil {
			return err
		}
		return ensureDirs()
	},
	RunE: RunServe,
}

func init() {
	rootCmd.PersistentFlags().String("config", "", "path to config file (default: ./config.yaml)")
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(configureCmd)
	rootCmd.AddCommand(reflectCmd)
	rootCmd.AddCommand(toolsetsCmd)
	rootCmd.AddCommand(sessionsCmd)
	rootCmd.AddCommand(runtimeCmd)
	rootCmd.AddCommand(imCmd)
}
