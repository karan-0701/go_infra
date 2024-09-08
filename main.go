package main

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// create the resources in here
		// we can delete resources using code as well
		// security groups are different from the ec2 instance
		// same for the key pair functionality

		// Ingress means the traffic the server accepts
		// Egress means the traffic the server can send
		sgArgs := &ec2.SecurityGroupArgs{
			Ingress: ec2.SecurityGroupIngressArray{
				ec2.SecurityGroupIngressArgs{
					// accepts tcp requests on port 8080
					Protocol:   pulumi.String("tcp"),
					FromPort:   pulumi.Int(8080),
					ToPort:     pulumi.Int(8080),
					CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
				},
				ec2.SecurityGroupIngressArgs{
					// accepts http requests on port 22
					Protocol:   pulumi.String("tcp"),
					FromPort:   pulumi.Int(22),
					ToPort:     pulumi.Int(22),
					CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
				},
			},
			Egress: ec2.SecurityGroupEgressArray{
				ec2.SecurityGroupEgressArgs{
					// can send all traffics
					Protocol:   pulumi.String("-1"),
					FromPort:   pulumi.Int(0),
					ToPort:     pulumi.Int(0),
					CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
				},
			},
		}

		// Add your security group creation code here using sgArgs
		sg, err := ec2.NewSecurityGroup(ctx, "jenkins-sg", sgArgs)
		if err != nil {
			return err
		}

		kp, err := ec2.NewKeyPair(ctx, "local-ssh", &ec2.KeyPairArgs{
			PublicKey: pulumi.String("ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC9uMe+zc+tuoO5+UKgViCnITIanfP8NDZc20WGBbP6BlH+40CGN8RPmv2/3AQYZTzp5l0MhbK4P3LXqBE5b3wi8GgiEYwiWr9jhi2nflK189XSAKnWO4rhbD7QUy55XUWIwkz1X59VczCsLKjerXJtgAEGROVS8XH1w7O4RhKkLuEWWz2CmriFixBFkUGc1MTbI6Psa9CVICC5n876oRY+WHlT9cwtx+dUMU/i0Q5RFo5Fl/KDxBbb1nHl5v1nboriUe5bEHlO4QXvQk/D5XAsr0vxuE32OfYZqcv5zH8vdQ5TIbYGucTw+Kc38h0uyMkNpvJE7G+esDwyzxbU9kGbtQzvMslSJia4CmzUmNfY7WiuARxLglFn/0LO1v1rb6SiTI3PahpsL7eHlSLfvq9URTVLJUFF9xuJVkP0ue7V1vpSOsiZ39/Fs3RXNa4zLVMbbDcffMvFBo8xPh0ngkgEJiQjCImXz15325gzj5X0f6LilOHojg3/I1QMTv1io4M= karandeepsingh@Karandeeps-MacBook-Air.local"),
		})

		if err != nil {
			return err
		}

		jenkinsServer, err := ec2.NewInstance(ctx, "jenkins-server", &ec2.InstanceArgs{
			InstanceType:        pulumi.String("t2.micro"),
			VpcSecurityGroupIds: pulumi.StringArray{sg.ID()},
			Ami:                 pulumi.String("ami-0888ba30fd446b771"),
			KeyName:             kp.KeyName,
		})

		fmt.Println(jenkinsServer.PublicIp)
		fmt.Println(jenkinsServer.PublicDns)

		ctx.Export("publicIp", jenkinsServer.PublicIp)
		ctx.Export("publicHostName", jenkinsServer.PublicDns)
		return nil
	})
}
