package createvalidations

import (
	"github.com/aws/eks-anywhere/pkg/validations"
)

func New(opts *validations.Opts) *CreateValidations {
	opts.SetDefaults()
	return &CreateValidations{Opts: opts}
}

type CreateValidations struct {
	Opts *validations.Opts
}
