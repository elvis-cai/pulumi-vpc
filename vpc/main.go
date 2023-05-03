package main

import (
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/compute"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/container"
	serviceAccount "github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/serviceaccount"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		sa, err := serviceAccount.NewAccount(ctx, "default", &serviceAccount.AccountArgs{
			AccountId:   pulumi.String("service-account-id"),
			DisplayName: pulumi.String("Service Account"),
		})
		if err != nil {
			return err
		}

		network, err := compute.NewNetwork(ctx, "my-vpc", &compute.NetworkArgs{
			AutoCreateSubnetworks: pulumi.Bool(false),
		})
		if err != nil {
			return err
		}
		// Create a custom subnet
		subnet, err := compute.NewSubnetwork(ctx, "subnet-poc", &compute.SubnetworkArgs{
			Network:     network.ID(),
			IpCidrRange: pulumi.String("10.2.0.0/16"),
			Region:      pulumi.String("us-east4"),
			SecondaryIpRanges: compute.SubnetworkSecondaryIpRangeArray{
				&compute.SubnetworkSecondaryIpRangeArgs{
					IpCidrRange: pulumi.String("10.3.0.0/20"),
					RangeName:   pulumi.String("secondary-range-1"),
				},
				&compute.SubnetworkSecondaryIpRangeArgs{
					IpCidrRange: pulumi.String("10.3.16.0/20"),
					RangeName:   pulumi.String("secondary-range-2"),
				},
			},
		})

		if err != nil {
			return err
		}

		// Create a GKE cluster
		primary, err := container.NewCluster(ctx, "gke-cluster", &container.ClusterArgs{
			InitialNodeCount:      pulumi.Int(1),
			Location:              pulumi.String("us-east4"),
			RemoveDefaultNodePool: pulumi.Bool(true),
			Network:               network.ID(),
			Subnetwork:            subnet.ID(),
		})
		if err != nil {
			return err
		}

		_, err = container.NewNodePool(ctx, "primaryPreemptibleNodes", &container.NodePoolArgs{
			Location:  pulumi.String("us-east4"),
			Cluster:   primary.Name,
			NodeCount: pulumi.Int(1),
			NodeConfig: &container.NodePoolNodeConfigArgs{
				Preemptible:    pulumi.Bool(true),
				MachineType:    pulumi.String("e2-medium"),
				ServiceAccount: sa.Email,
				OauthScopes: pulumi.StringArray{
					pulumi.String("https://www.googleapis.com/auth/cloud-platform"),
				},
			},
		})
		if err != nil {
			return err
		}

		// Export the VPC name.
		ctx.Export("vpcName", network.Name)
		ctx.Export("subnetwork", network.ID())
		// Export the GKE cluster name and endpoint
		return nil
	})
}

