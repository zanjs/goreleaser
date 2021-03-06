// Package name provides name template logic for the final archive, formulae,
// etc.
package name

import (
	"bytes"
	"text/template"

	"github.com/goreleaser/goreleaser/config"
	"github.com/goreleaser/goreleaser/context"
)

type nameData struct {
	Os          string
	Arch        string
	Arm         string
	Version     string
	Tag         string
	Binary      string // deprecated
	ProjectName string
}

// ForBuild return the name for the given context, goos, goarch, goarm and
// build, using the build.Binary property instead of project_name.
func ForBuild(ctx *context.Context, build config.Build, goos, goarch, goarm string) (string, error) {
	return apply(
		nameData{
			Os:          replace(ctx.Config.Archive.Replacements, goos),
			Arch:        replace(ctx.Config.Archive.Replacements, goarch),
			Arm:         replace(ctx.Config.Archive.Replacements, goarm),
			Version:     ctx.Version,
			Tag:         ctx.Git.CurrentTag,
			Binary:      build.Binary,
			ProjectName: build.Binary,
		},
		ctx.Config.Archive.NameTemplate,
	)
}

// For returns the name for the given context, goos, goarch and goarm.
func For(ctx *context.Context, goos, goarch, goarm string) (string, error) {
	return apply(
		nameData{
			Os:          replace(ctx.Config.Archive.Replacements, goos),
			Arch:        replace(ctx.Config.Archive.Replacements, goarch),
			Arm:         replace(ctx.Config.Archive.Replacements, goarm),
			Version:     ctx.Version,
			Tag:         ctx.Git.CurrentTag,
			Binary:      ctx.Config.ProjectName,
			ProjectName: ctx.Config.ProjectName,
		},
		ctx.Config.Archive.NameTemplate,
	)
}

func apply(data nameData, templateStr string) (string, error) {
	var out bytes.Buffer
	t, err := template.New(data.ProjectName).Parse(templateStr)
	if err != nil {
		return "", err
	}
	err = t.Execute(&out, data)
	return out.String(), err
}

func replace(replacements map[string]string, original string) string {
	result := replacements[original]
	if result == "" {
		return original
	}
	return result
}
