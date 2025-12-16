package genmongo

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"

	"github.com/jzero-io/jzero/cmd/jzero/internal/config"
	"github.com/jzero-io/jzero/cmd/jzero/internal/embeded"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/console"
)

type JzeroMongo struct {
	Module string
}

func (jm *JzeroMongo) Gen() error {
	if len(config.C.Gen.MongoType) == 0 {
		return nil
	}

	var goctlHome string

	if !pathx.FileExists(filepath.Join(config.C.Home, "go-zero", "mongo")) {
		tempDir, err := os.MkdirTemp(os.TempDir(), "")
		if err != nil {
			return err
		}
		defer os.RemoveAll(tempDir)
		err = embeded.WriteTemplateDir(filepath.Join("go-zero", "mongo"), filepath.Join(tempDir, "mongo"))
		if err != nil {
			return err
		}
		goctlHome = tempDir
	} else {
		goctlHome = filepath.Join(config.C.Home, "go-zero")
	}
	logx.Debugf("goctl_home = %s", goctlHome)

	if !config.C.Quiet {
		fmt.Printf("%s to generate mongo model code from types.\n", console.Green("Start"))
	}

	for _, mongoType := range config.C.Gen.MongoType {
		if !config.C.Quiet {
			fmt.Printf("%s mongo type %s\n", console.Green("Using"), mongoType)
		}

		// Support MutiModel with dot notation like "ntls_log.user"
		var typeName string
		if strings.Contains(mongoType, ".") {
			typeName = filepath.Join(strings.Split(mongoType, ".")...)
		} else {
			typeName = mongoType
		}

		modelDir := filepath.Join("internal", "mongo", strings.ToLower(typeName))

		// For MutiModel, only pass the part after the dot to goctl
		var goctlType string
		if strings.Contains(mongoType, ".") {
			goctlType = strings.Split(mongoType, ".")[1] // Only pass the part after the dot
		} else {
			goctlType = mongoType
		}

		args := []string{
			"model", "mongo",
			"-t", goctlType,
			"--dir", modelDir,
			"--home", goctlHome,
			"--style", config.C.Style,
		}

		var enableCache bool
		if config.C.Gen.MongoCache {
			if len(config.C.Gen.MongoCacheType) == 1 && config.C.Gen.MongoCacheType[0] == "*" {
				enableCache = true
			} else {
				for _, cacheType := range config.C.Gen.MongoCacheType {
					if cacheType == mongoType {
						enableCache = true
						break
					}
				}
			}
		}

		if enableCache {
			args = append(args, "--cache=true")
			if config.C.Gen.MongoCachePrefix != "" {
				args = append(args, "-p", config.C.Gen.MongoCachePrefix)
			}
		} else {
			args = append(args, "--cache=false")
		}

		// easy 模式
		args = append(args, fmt.Sprintf("--easy=%v", true))

		cmd := exec.Command("goctl", args...)
		logx.Debug(cmd.String())
		resp, err := cmd.CombinedOutput()
		if err != nil {
			return errors.Errorf("gen mongo model code meet error. Err: %s:%s", err.Error(), resp)
		}
	}

	err := jm.GenRegister(config.C.Gen.MongoType)
	if err != nil {
		return err
	}

	if !config.C.Quiet {
		fmt.Println(console.Green("Done"))
	}

	return nil
}
