package trading

import (
	"github.com/commonpool/backend/pkg/graph"
)

type Module struct {
	Store   Store
	Service Service
}

func NewModule(driver graph.Driver) *Module {
	return &Module{}
}
