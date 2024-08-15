package main

import (
	"context"
	"strconv"

	"github.com/jalevin/gottl/dagger/internal/dagger"
)

var (
	imagePostgres = "postgres:latest"
	imageGolang   = "golang:latest"
)

type Workflows struct {
	dag *dagger.Client
}

func (g *Workflows) init() {
	if g.dag != nil {
		return
	}

	g.dag = dagger.Connect()
}

func (g *Workflows) TestEnv(source *dagger.Directory) *dagger.Container {
	postgres := g.dag.Container().
		From(imagePostgres).
		WithEnvVariable("POSTGRES_PASSWORD", "postgres").
		WithEnvVariable("POSTGRES_USER", "postgres").
		AsService()

	golang := g.dag.Container().
		From(imageGolang).
		WithServiceBinding("postgres", postgres).
		WithEnvVariable("GOTTL_POSTGRES_HOST", "postgres").
		WithWorkdir("/src/")

	// Mount Go Module Files
	golang = golang.
		WithFile("/src/go.mod", source.File("go.mod")).
		WithFile("/src/go.sum", source.File("go.sum")).
		WithMountedCache("/go/pkg/mod", g.dag.CacheVolume("go-mod-121")).
		WithEnvVariable("GOMODCACHE", "/go/pkg/mod").
		WithExec([]string{"go", "mod", "download"})

	// Mount source code and build cache
	golang = golang.
		WithDirectory("/src/", source, dagger.ContainerWithDirectoryOpts{
			// Exclude already copied files and directories that are not relevant to Go
			Exclude: []string{"go.mod", "go.sum", ".git/**", ".task/**"},
		}).
		WithMountedCache("/go/build-cache", g.dag.CacheVolume("go-build-121")).
		WithEnvVariable("GOCACHE", "/go/build-cache")

	return golang
}

// Test runs the integration tests for the project with the associated services.
func (g *Workflows) Test(
	ctx context.Context,
	source *dagger.Directory,
	// +optional
	// +default=true
	integration bool,
) (string, error) {
	g.init()

	// get the build environment container
	// by calling another Dagger Function
	return g.TestEnv(source).
		WithEnvVariable("UNSAFE_PASSWORD_PROTECTION", "yes_i_am_sure"). // disable password hashing
		WithEnvVariable("TEST_INTEGRATION", strconv.FormatBool(integration)).
		WithExec([]string{"go", "test", "-v", "./..."}).
		Stdout(ctx)
}
