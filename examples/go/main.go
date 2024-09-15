package main

import (
	"github.com/abklabs/pulumi-svm/sdk/go/svm"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		myRandomResource, err := svm.NewRandom(ctx, "myRandomResource", &svm.RandomArgs{
			Length: pulumi.Int(24),
		})
		if err != nil {
			return err
		}
		ctx.Export("output", map[string]interface{}{
			"value": myRandomResource.Result,
		})
		return nil
	})
}
