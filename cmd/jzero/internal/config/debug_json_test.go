package config

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestMarshalDebugJSONForCommandIncludesOnlyGlobalAndCurrentCommandFlags(t *testing.T) {
	root := &cobra.Command{Use: "jzero"}
	root.PersistentFlags().Bool("debug", false, "")
	root.PersistentFlags().String("style", "gozero", "")
	root.PersistentFlags().String("home", ".template", "")
	root.PersistentFlags().String("config", ".jzero.yaml", "")
	root.PersistentFlags().String("working-dir", ".", "")

	gen := &cobra.Command{Use: "gen"}
	gen.Flags().StringSlice("desc", nil, "")
	gen.Flags().Bool("rpc-client", false, "")
	root.AddCommand(gen)

	cfg := Config{
		Debug: true,
		Style: "gozero",
		Home:  "/tmp/template-home",
		Gen: GenConfig{
			Desc:      []string{"desc/api/user.api"},
			RpcClient: true,
		},
		New: NewConfig{
			Name: "demo",
		},
	}

	data, err := MarshalDebugJSONForCommand(cfg, gen)
	if err != nil {
		t.Fatalf("MarshalDebugJSONForCommand() error = %v", err)
	}

	got := string(data)

	for _, want := range []string{
		`"debug": true`,
		`"style": "gozero"`,
		`"home": "/tmp/template-home"`,
		`"gen": {`,
		`"desc": [`,
		`"rpc-client": true`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("MarshalDebugJSONForCommand() missing %q in output:\n%s", want, got)
		}
	}

	for _, unwanted := range []string{
		`"config"`,
		`"working-dir"`,
		`"new"`,
		`"hooks"`,
	} {
		if strings.Contains(got, unwanted) {
			t.Fatalf("MarshalDebugJSONForCommand() should not contain %q:\n%s", unwanted, got)
		}
	}
}

func TestMarshalDebugJSONForCommandUsesActiveLeafBranchOnly(t *testing.T) {
	root := &cobra.Command{Use: "jzero"}
	root.PersistentFlags().Bool("debug", false, "")

	gen := &cobra.Command{Use: "gen"}
	gen.Flags().Bool("rpc-client", false, "")
	root.AddCommand(gen)

	swagger := &cobra.Command{Use: "swagger"}
	swagger.Flags().String("output", "", "")
	swagger.Flags().Bool("merge", true, "")
	gen.AddCommand(swagger)

	cfg := Config{
		Debug: true,
		Gen: GenConfig{
			RpcClient: true,
			Swagger: GenSwaggerConfig{
				Output: "desc/swagger",
				Merge:  true,
			},
		},
	}

	data, err := MarshalDebugJSONForCommand(cfg, swagger)
	if err != nil {
		t.Fatalf("MarshalDebugJSONForCommand() error = %v", err)
	}

	got := string(data)

	for _, want := range []string{
		`"debug": true`,
		`"gen": {`,
		`"swagger": {`,
		`"output": "desc/swagger"`,
		`"merge": true`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("MarshalDebugJSONForCommand() missing %q in output:\n%s", want, got)
		}
	}

	if strings.Contains(got, `"rpc-client": true`) {
		t.Fatalf("MarshalDebugJSONForCommand() should not contain parent local flags:\n%s", got)
	}
}

func TestMarshalDebugJSONForCommandKeepsInheritedPersistentFlags(t *testing.T) {
	root := &cobra.Command{Use: "jzero"}
	root.PersistentFlags().Bool("debug", false, "")

	migrate := &cobra.Command{Use: "migrate"}
	migrate.PersistentFlags().String("source", "", "")
	migrate.PersistentFlags().String("datasource-url", "", "")
	root.AddCommand(migrate)

	up := &cobra.Command{Use: "up"}
	migrate.AddCommand(up)

	cfg := Config{
		Debug: true,
		Migrate: MigrateConfig{
			Source:        "file://desc/sql_migration",
			DataSourceUrl: "mysql://demo",
		},
	}

	data, err := MarshalDebugJSONForCommand(cfg, up)
	if err != nil {
		t.Fatalf("MarshalDebugJSONForCommand() error = %v", err)
	}

	got := string(data)

	for _, want := range []string{
		`"debug": true`,
		`"migrate": {`,
		`"source": "file://desc/sql_migration"`,
		`"datasource-url": "mysql://demo"`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("MarshalDebugJSONForCommand() missing %q in output:\n%s", want, got)
		}
	}
}
