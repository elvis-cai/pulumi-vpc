package main

import (

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi-gcp/sdk/v5/go/gcp/compute"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		network, err := compute.NewNetwork(ctx, "my-vpc", &compute.NetworkArgs{
			AutoCreateSubnetworks: pulumi.Bool(false),
		})
		if err != nil {
			return err
		}

		// Export the VPC name.
		ctx.Export("vpcName", network.Name)
		ctx.Export("subnetwork", network.ID())
		return nil
	})
}
