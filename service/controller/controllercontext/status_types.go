package controllercontext

import (
	"github.com/aws/aws-sdk-go/service/ec2"
)

type ContextStatus struct {
	AWSAccountID string
	NATGateway   ContextStatusNATGateway
	RouteTables  []*ec2.RouteTable
	PeerRole     ContextStatusPeerRole
	VPC          ContextStatusVPC
}

type ContextStatusNATGateway struct {
	Addresses []*ec2.Address
}

type ContextStatusPeerRole struct {
	ARN string
}

type ContextStatusVPC struct {
	CIDR string
	ID   string
}
