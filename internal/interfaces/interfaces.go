// interfaces.go
package interfaces

import (
	"github.com/google/wire"
)

// ProviderSet is interfaces providers.
var ProviderSet = wire.NewSet(NewHandler)
