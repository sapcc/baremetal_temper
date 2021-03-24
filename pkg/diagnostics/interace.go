package diagnostics

import "github.com/sapcc/baremetal_temper/pkg/model"

type Diagnostics interface {
	Run(n *model.Node) error
}
