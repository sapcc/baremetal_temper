package diagnostics

import "github.com/sapcc/ironic_temper/pkg/model"

type Diagnostics interface {
	Run(n *model.IronicNode) error
}
