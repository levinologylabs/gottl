package main

import (
	"context"
	"strconv"

	"github.com/jalevin/gottl/workflows/internal/dagger"
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
		WithDirectory("/src/", source).
		WithWorkdir("/src/")

	return golang
}

// Return the result of running unit tests
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
		WithEnvVariable("TEST_INTEGRATION", strconv.FormatBool(integration)).
		WithExec([]string{"go", "test", "./..."}).
		Stdout(ctx)
}
