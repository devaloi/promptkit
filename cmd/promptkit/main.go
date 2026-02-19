// Package main is the entry point for the promptkit CLI.
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/devaloi/promptkit/internal/chain"
	"github.com/devaloi/promptkit/internal/config"
	"github.com/devaloi/promptkit/internal/engine"
	"github.com/devaloi/promptkit/internal/registry"
	"github.com/devaloi/promptkit/internal/validator"
)

func main() {
	if err := rootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}

func rootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "promptkit",
		Short: "LLM prompt template engine",
		Long:  "A template engine for LLM prompts with variable injection, validation, includes, and chaining.",
	}

	cmd.AddCommand(renderCmd(), validateCmd(), listCmd(), chainCmd())
	return cmd
}

func renderCmd() *cobra.Command {
	var (
		dir     string
		varFlag []string
	)

	cmd := &cobra.Command{
		Use:   "render <template>",
		Short: "Render a prompt template",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			reg := registry.New()
			if err := reg.LoadDir(dir); err != nil {
				return fmt.Errorf("loading templates: %w", err)
			}

			tmpl, err := reg.Get(args[0])
			if err != nil {
				return err
			}

			vars := parseVars(varFlag)

			if len(tmpl.Meta.RequiredVars) > 0 {
				if err := validator.Validate(tmpl.Meta.RequiredVars, vars); err != nil {
					return err
				}
			}

			result, err := engine.Render(tmpl.Content, vars, reg.Includes())
			if err != nil {
				return err
			}

			fmt.Print(result.Output)
			return nil
		},
	}

	cmd.Flags().StringVarP(&dir, "dir", "d", config.DefaultTemplateDir, "template directory")
	cmd.Flags().StringArrayVar(&varFlag, "var", nil, "variable in key=value format")

	return cmd
}

func validateCmd() *cobra.Command {
	var dir string

	cmd := &cobra.Command{
		Use:   "validate <template>",
		Short: "Validate required variables for a template",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			reg := registry.New()
			if err := reg.LoadDir(dir); err != nil {
				return fmt.Errorf("loading templates: %w", err)
			}

			tmpl, err := reg.Get(args[0])
			if err != nil {
				return err
			}

			if len(tmpl.Meta.RequiredVars) == 0 {
				fmt.Println("No required variables.")
				return nil
			}

			fmt.Printf("Required variables for %q:\n", tmpl.Name)
			for _, v := range tmpl.Meta.RequiredVars {
				fmt.Printf("  - %s\n", v)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&dir, "dir", "d", config.DefaultTemplateDir, "template directory")
	return cmd
}

func listCmd() *cobra.Command {
	var dir string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all available templates",
		RunE: func(_ *cobra.Command, _ []string) error {
			reg := registry.New()
			if err := reg.LoadDir(dir); err != nil {
				return fmt.Errorf("loading templates: %w", err)
			}

			templates := reg.List()
			if len(templates) == 0 {
				fmt.Println("No templates found.")
				return nil
			}

			for _, tmpl := range templates {
				desc := tmpl.Meta.Description
				if desc == "" {
					desc = "(no description)"
				}
				fmt.Printf("%-20s %s\n", tmpl.Name, desc)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&dir, "dir", "d", config.DefaultTemplateDir, "template directory")
	return cmd
}

func chainCmd() *cobra.Command {
	var (
		dir     string
		varFlag []string
	)

	cmd := &cobra.Command{
		Use:   "chain <chain.yaml>",
		Short: "Execute a prompt chain",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			def, err := chain.ParseFile(args[0])
			if err != nil {
				return err
			}

			reg := registry.New()
			if err := reg.LoadDir(dir); err != nil {
				return fmt.Errorf("loading templates: %w", err)
			}

			vars := parseVars(varFlag)

			result, err := chain.Execute(def, reg, vars)
			if err != nil {
				return err
			}

			fmt.Print(result.Final)
			return nil
		},
	}

	cmd.Flags().StringVarP(&dir, "dir", "d", config.DefaultTemplateDir, "template directory")
	cmd.Flags().StringArrayVar(&varFlag, "var", nil, "variable in key=value format")

	return cmd
}

func parseVars(flags []string) map[string]any {
	vars := make(map[string]any, len(flags))
	for _, f := range flags {
		parts := strings.SplitN(f, "=", 2)
		if len(parts) == 2 {
			vars[parts[0]] = parts[1]
		}
	}
	return vars
}
