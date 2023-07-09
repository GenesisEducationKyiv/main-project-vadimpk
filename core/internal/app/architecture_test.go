package app

import (
	"testing"

	"github.com/matthewmcnew/archtest"
)

func TestArchitecture(t *testing.T) {
	// Delivery layer should not depend on anything except service layer
	archtest.Package(t, "github.com/vadimpk/gses-2023/internal/controller").
		ShouldNotDependOn("github.com/vadimpk/gses-2023/internal/api",
			"github.com/vadimpk/gses-2023/internal/storage/...")

	// Service layer should not depend on API or Storage layers
	archtest.Package(t, "github.com/vadimpk/gses-2023/internal/service").
		ShouldNotDependOn("github.com/vadimpk/gses-2023/internal/api",
			"github.com/vadimpk/gses-2023/internal/storage/...")

	// Entity layer should not depend on anything
	archtest.Package(t, "github.com/vadimpk/gses-2023/internal/entity").
		ShouldNotDependOn("github.com/vadimpk/gses-2023/internal/...", "github.com/vadimpk/gses-2023/pkg/...")
}
