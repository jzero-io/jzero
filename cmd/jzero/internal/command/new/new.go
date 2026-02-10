/*
Copyright © 2024 jaronnie <jaron@jaronnie.com>

*/

package new

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/pkg/errors"
	"github.com/rinchsan/gosimports"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"

	"github.com/jzero-io/jzero/cmd/jzero/internal/command/check"
	"github.com/jzero-io/jzero/cmd/jzero/internal/command/gen/gen"
	"github.com/jzero-io/jzero/cmd/jzero/internal/config"
	"github.com/jzero-io/jzero/cmd/jzero/internal/desc"
	"github.com/jzero-io/jzero/cmd/jzero/internal/embeded"
	"github.com/jzero-io/jzero/cmd/jzero/internal/hooks"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/console"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/filex"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/mod"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/templatex"
)

func IsBase64(base64 string) bool {
	return regexp.MustCompile(`^(?:[A-Za-z0-9+\\/]{4})*(?:[A-Za-z0-9+\\/]{2}==|[A-Za-z0-9+\\/]{3}=|[A-Za-z0-9+\\/]{4})$`).MatchString(base64)
}

type JzeroNew struct {
	TemplateData map[string]any
	nc           config.NewConfig
	base         string
}

type GeneratedFile struct {
	Path    string
	Content bytes.Buffer
	Skip    bool
}

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:          "new",
	Short:        `Used to create project from templates`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		var app string
		if len(args) > 0 {
			app = args[0]
		} else {
			app = config.C.New.Name
		}

		if config.C.New.Serverless && config.C.New.Output != "" {
			return errors.New("serverless mode not support output dir, must be in current project")
		}

		if config.C.New.Output == "" {
			if len(args) > 0 {
				config.C.New.Output = args[0]
			} else {
				config.C.New.Output = config.C.New.Name
			}

			if pathx.FileExists(config.C.New.Output) {
				return errors.Errorf("%s already exists", config.C.New.Output)
			}
		}

		if config.C.New.Serverless {
			config.C.New.Output = filepath.Join("plugins", config.C.New.Output)
		}

		if !config.C.Quiet {
			fmt.Printf("%s project %s in %s dir\n", console.Green("Creating"), app, config.C.New.Output)
		}

		if config.C.New.Module == "" {
			if len(args) > 0 {
				config.C.New.Module = args[0]
			} else {
				config.C.New.Module = config.C.New.Name
			}
		}
		// 在 go.mod 项目下但是项目本身没有 go.mod 文件
		if config.C.New.Mono {
			wd, _ := os.Getwd()
			parentPackage, err := mod.GetParentPackage(wd)
			if err != nil {
				return err
			}
			config.C.New.Module = filepath.ToSlash(filepath.Join(parentPackage, config.C.New.Output))
		}
		gosimports.LocalPrefix = config.C.New.Module

		home, _ := os.UserHomeDir()

		var base string
		switch {
		// 使用内置模板
		case config.C.New.Frame != "":
			base = filepath.Join("frame", config.C.New.Frame, "app")
		// 指定本地路径 ~/.jzero/templates/local 下的某文件夹作为模板
		case config.C.New.Local != "":
			embeded.Home = filepath.Join(home, ".jzero", "templates", "local", config.C.New.Local)
			base = filepath.Join("app")
		// 使用远程仓库模板
		case config.C.New.Remote != "" && config.C.New.Branch != "":
			fp := filepath.Join(home, ".jzero", "templates", "remote", config.C.New.Branch)
			if filex.DirExists(fp) && config.C.New.Cache {
				if !config.C.Quiet {
					fmt.Printf("%s cache templates from '%s', please wait...\n", console.Green("Using"), fp)
				}
			} else {
				_ = os.RemoveAll(fp)
				if !config.C.Quiet {
					fmt.Printf("%s templates into '%s', please wait...\n", console.Green("Cloning"), fp)
				}
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(config.C.New.RemoteTimeout))
				defer cancel()

				if strings.HasPrefix(config.C.New.Remote, "git@") {
					// SSH 协议
					commandContext := exec.CommandContext(ctx, "git", "clone", "--depth", "1", "--branch", config.C.New.Branch, config.C.New.Remote, fp)
					if resp, err := commandContext.CombinedOutput(); err != nil {
						return errors.New(string(resp))
					}
				} else {
					// HTTP 协议
					auth := &http.BasicAuth{
						Username: config.C.New.RemoteAuthUsername,
						Password: config.C.New.RemoteAuthPassword,
					}
					// clone to local
					if _, err := git.PlainCloneContext(ctx, fp, false, &git.CloneOptions{
						SingleBranch:  true,
						URL:           config.C.New.Remote,
						Depth:         0,
						ReferenceName: plumbing.ReferenceName("refs/heads/" + config.C.New.Branch),
						Auth:          auth,
					}); err != nil {
						return err
					}
				}

				_ = os.RemoveAll(filepath.Join(fp, ".git"))
			}
			if !config.C.Quiet {
				fmt.Println(console.Green("Done"))
			}
			embeded.Home = fp
			base = filepath.Join("app")
		// 指定特定路径作为模板
		case config.C.Home != "":
			embeded.Home = config.C.Home
			if config.C.New.Frame != "" {
				base = filepath.Join("frame", config.C.New.Frame, "app")
			} else {
				base = filepath.Join("app")
				if !pathx.FileExists(base) && pathx.FileExists(filepath.Join(embeded.Home, "frame")) {
					base = filepath.Join("frame", "api", "app")
				}
			}
		default:
			// 默认使用 api 模板
			config.C.New.Frame = "api"
			base = filepath.Join("frame", "api", "app")
		}

		if err := Run(app, base); err != nil {
			return err
		}

		if !config.C.New.Gen {
			return nil
		}

		// change dir to project
		if err := os.Chdir(config.C.New.Output); err != nil {
			return err
		}
		defer func() {
			dir, _ := os.Getwd()
			if err := os.Chdir(dir); err != nil {
				cobra.CheckErr(err)
			}
		}()

		frameType, err := desc.GetFrameType()
		if err != nil {
			return err
		}
		if frameType != "" {
			if err := check.RunCheck(false); err != nil {
				return err
			}
		}

		// special dir for jzero
		if !filex.DirExists("desc") {
			return nil
		}
		if !config.C.Quiet {
			fmt.Printf("%s desc dir in %s, auto generate code\n", console.Green("Detected"), config.C.New.Output)
		}

		config.ResetConfig()
		if err := config.InitConfig(cmd.Root()); err != nil {
			return err
		}

		// for gen persistent flags
		if config.C.Style == "" {
			config.C.Style = "gozero"
		}
		if config.C.Home == "" {
			config.C.Home = filepath.Join(config.C.Wd(), ".template")
		}

		// run gen before hooks
		if err := hooks.Run(cmd, "Before", "gen", config.C.Gen.Hooks.Before); err != nil {
			return err
		}
		if err := gen.Run(); err != nil {
			return err
		}
		return hooks.Run(cmd, "After", "gen", config.C.Gen.Hooks.After)
	},
}

func Run(appName, base string) error {
	if err := os.MkdirAll(config.C.New.Output, 0o755); err != nil {
		return err
	}

	goVersion, err := mod.GetGoVersion()
	if err != nil {
		return err
	}

	templateData := map[string]any{
		"GoVersion": goVersion,
		"GoArch":    runtime.GOARCH,
	}
	templateData["Features"] = config.C.New.Features
	templateData["Module"] = config.C.New.Module
	templateData["APP"] = appName
	if abs, err := filepath.Abs(config.C.New.Output); err == nil {
		templateData["DirName"] = filepath.Base(abs)
	} else {
		return err
	}
	templateData["Style"] = config.C.Style
	templateData["Serverless"] = config.C.New.Serverless

	jn := JzeroNew{
		TemplateData: templateData,
		nc:           config.C.New,
		base:         base,
	}

	gfs, err := jn.New(base)
	if err != nil {
		return err
	}

	for _, gf := range gfs {
		if !gf.Skip {
			err = checkWrite(gf.Path, gf.Content.Bytes())
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func checkWrite(path string, bytes []byte) error {
	var err error
	if len(bytes) == 0 {
		return nil
	}
	if !pathx.FileExists(filepath.Dir(path)) {
		err = os.MkdirAll(filepath.Dir(path), 0o755)
		if err != nil {
			return err
		}
	}

	bytesFormat := bytes
	if filepath.Ext(path) == ".go" {
		bytesFormat, err = gosimports.Process("", bytes, nil)
		if err != nil {
			return errors.Wrapf(err, "format %s", path)
		}
	}

	// 增加可执行权限
	if lo.Contains(config.C.New.ExecutableExtensions, filepath.Ext(path)) {
		logx.Debugf("Write executable file: %s", path)
		return os.WriteFile(path, bytesFormat, 0o755)
	}
	return os.WriteFile(path, bytesFormat, 0o644)
}

func (jn *JzeroNew) New(dirname string) ([]*GeneratedFile, error) {
	var gsf []*GeneratedFile

	dir := embeded.ReadTemplateDir(dirname)
	for _, file := range dir {
		if file.IsDir() {
			files, err := jn.New(filepath.Join(dirname, file.Name()))
			if err != nil {
				return nil, err
			}
			gsf = append(gsf, files...)
		}

		filename := file.Name()
		if IsBase64(filename) {
			filenameBytes, _ := base64.StdEncoding.DecodeString(filename)
			filename = string(filenameBytes)
		}

		filename = strings.TrimSuffix(filename, ".tpl")

		rel, err := filepath.Rel(jn.base, filepath.Join(dirname, filename))
		if err != nil {
			return nil, err
		}

		var fileBytes []byte
		if strings.HasSuffix(file.Name(), ".tpl.tpl") {
			// .tpl.tpl suffix means it is a template, do not parse if anymore
			fileBytes = embeded.ReadTemplateFile(filepath.Join(dirname, file.Name()))
		} else {
			fileBytes, err = templatex.ParseTemplate(filepath.Join(dirname, file.Name()), jn.TemplateData, embeded.ReadTemplateFile(filepath.Join(dirname, file.Name())))
			if err != nil {
				return nil, err
			}
		}

		// parse template name
		templatePath := filepath.Join(filepath.Dir(rel), filename)
		stylePathBytes, err := templatex.ParseTemplate(templatePath, jn.TemplateData, []byte(templatePath))
		if err != nil {
			return nil, err
		}

		// specify
		if filename == "go.mod" && jn.nc.Mono {
			continue
		}

		gsf = append(gsf, &GeneratedFile{
			Path:    filepath.Join(jn.nc.Output, string(stylePathBytes)),
			Content: *bytes.NewBuffer(fileBytes),
			Skip: func() bool {
				var ignore []string
				for _, v := range jn.nc.Ignore {
					ignore = append(ignore, filepath.ToSlash(v))
				}
				for _, v := range jn.nc.IgnoreExtra {
					ignore = append(ignore, filepath.ToSlash(v))
				}
				// if ignore is dir
				for _, v := range ignore {
					if config.C.New.Serverless {
						v = filepath.Join(config.C.New.Output, v)
					}
					if stat, err := os.Stat(v); err == nil && stat.IsDir() {
						if filepath.ToSlash(filepath.Dir(string(stylePathBytes))) == filepath.ToSlash(v) {
							return true
						}
					} else {
						if filepath.ToSlash(string(stylePathBytes)) == filepath.ToSlash(v) {
							return true
						}
					}
				}
				return false
			}(),
		})
	}
	return gsf, nil
}

func GetCommand() *cobra.Command {
	newCmd.Flags().StringP("name", "", "", "set project name")
	newCmd.Flags().StringP("module", "m", "", "set go module")
	newCmd.Flags().StringP("output", "o", "", "set output dir with project name")
	newCmd.Flags().StringP("frame", "", "", "set frame such as api/rpc/gateway")
	newCmd.Flags().StringP("remote", "r", "https://github.com/jzero-io/templates", "remote templates repo")
	newCmd.Flags().IntP("remote-timeout", "", 30, "remote templates repo timeout")
	newCmd.Flags().StringP("remote-auth-username", "", "", "remote templates repo auth username")
	newCmd.Flags().StringP("remote-auth-password", "", "", "remote templates repo auth password")
	newCmd.Flags().StringP("branch", "b", "", "use remote template repo branch")
	newCmd.Flags().BoolP("cache", "", false, "remote template using cache")
	newCmd.Flags().StringP("local", "", "", "use local template")
	newCmd.Flags().StringSliceP("features", "", []string{}, "set features such as model/cache/redis")
	newCmd.Flags().BoolP("mono", "", false, "mono project under go mod project")
	newCmd.Flags().BoolP("serverless", "", false, "create serverless project")
	newCmd.Flags().BoolP("gen", "", true, "gen code after new project")
	newCmd.Flags().StringSliceP("ignore", "", []string{}, "set ignore file")
	newCmd.Flags().StringSliceP("ignore-extra", "", []string{}, "set ignore extra file")
	newCmd.Flags().StringSliceP("executable-extensions", "", []string{".sh"}, "select executable extensions")

	return newCmd
}
