package app

import (
	"testing"

	"github.com/matthewmcnew/archtest"
)

func TestArchitecture(t *testing.T) {
	// Delivery layer should not depend on anything except service layer
	archtest.Package(t, "github.com/vadimpk/gses-2023/core/internal/controller").
		ShouldNotDependOn("github.com/vadimpk/gses-2023/core/internal/api",
			"github.com/vadimpk/gses-2023/core/internal/storage/...")

	// Service layer should not depend on API or Storage layers
	archtest.Package(t, "github.com/vadimpk/gses-2023/core/internal/service").
		ShouldNotDependOn("github.com/vadimpk/gses-2023/core/internal/api",
			"github.com/vadimpk/gses-2023/core/internal/storage/...")

	// Entity layer should not depend on anything
	archtest.Package(t, "github.com/vadimpk/gses-2023/core/internal/entity").
		ShouldNotDependOn("github.com/vadimpk/gses-2023/core/internal/...", "github.com/vadimpk/gses-2023/core/pkg/...")
}
