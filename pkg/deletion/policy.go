package deletion

import (
	"fmt"

	"github.com/abklabs/svmkit/pkg/runner"
	"github.com/pulumi/pulumi-go-provider/infer"
)

type Policy string

const (
	PolicyKeep                   Policy = "keep"
	PolicyDelete                 Policy = "delete"
	PolicyDeleteAndForceCreation Policy = "delete-and-force-creation"
)

func (Policy) Values() []infer.EnumValue[Policy] {
	return []infer.EnumValue[Policy]{
		{
			Name:        "keep",
			Value:       PolicyKeep,
			Description: "Do not delete files upon destroying their associated component",
		},
		{
			Name:        "delete",
			Value:       PolicyDelete,
			Description: "Delete files upon destroying their associated component; if the files exist before the component is created, will throw an error",
		},
		{
			Name:        "delete-and-force-creation",
			Value:       PolicyDeleteAndForceCreation,
			Description: "Delete files upon destroying their associated component, and overwrite existing files on creation — use with caution",
		},
	}
}

func (p Policy) Check() error {
	switch p {
	case PolicyKeep, PolicyDelete, PolicyDeleteAndForceCreation:
	default:
		return fmt.Errorf("unknown deletion policy '%s'", p)
	}

	return nil
}

type ManagedFiles interface {
	ManagedFiles() []string
}

func (p Policy) Create(files ManagedFiles, b *runner.EnvBuilder) {
	arr := []string{}
	switch p {
	case PolicyKeep, PolicyDeleteAndForceCreation:
	default:
		arr = files.ManagedFiles()
	}

	b.SetArray("DELETION_CHECK_FILES", arr)
}

func (p Policy) Delete(files ManagedFiles, b *runner.EnvBuilder) {
	arr := []string{}
	switch p {
	case PolicyDelete, PolicyDeleteAndForceCreation:
		arr = files.ManagedFiles()
	}

	b.SetArray("DELETION_DELETE_FILES", arr)
}
