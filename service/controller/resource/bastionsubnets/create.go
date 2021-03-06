package bastionsubnets

import (
	"context"
	"fmt"
	"net"

	"github.com/apparentlymart/go-cidr/cidr"
	"github.com/giantswarm/microerror"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/installation-operator/service/controller/key"
)

func (r *Resource) readExistingSubnets() (map[string][]*net.IPNet, error) {
	installations, err := r.k8sClient.G8sClient().CoreV1alpha1().Installations("").List(metav1.ListOptions{})
	if err != nil {
		return nil, microerror.Mask(err)
	}

	existing := map[string][]*net.IPNet{}
	for _, installation := range installations.Items {
		for _, subnetString := range installation.Status.BastionSubnets {
			_, subnet, err := net.ParseCIDR(subnetString)
			if err != nil {
				return nil, microerror.Mask(err)
			}
			existing[installation.Name] = append(existing[installation.Name], subnet)
		}
	}

	return existing, nil
}

func (r *Resource) findAvailableSubnets(ctx context.Context) ([]*net.IPNet, error) {
	existing, err := r.readExistingSubnets()
	if err != nil {
		return nil, microerror.Mask(err)
	}
	_, candidate, err := net.ParseCIDR("10.0.220.0/27")
	if err != nil {
		return nil, microerror.Mask(err)
	}
	var combined []*net.IPNet
	for _, subnets := range existing {
		for _, subnet := range subnets {
			combined = append(combined, subnet)
		}
	}
	for {
		err := cidr.VerifyNoOverlap(combined, candidate)
		if err == nil {
			break
		}
		var end bool
		candidate, end = cidr.NextSubnet(candidate, 27)
		if end {
			return nil, microerror.New("no available bastion subnet found")
		}
	}
	var subnets []*net.IPNet
	for i := 0; i < 2; i++ {
		subnet, err := cidr.Subnet(candidate, 1, i)
		if err != nil {
			return nil, microerror.Mask(err)
		}
		subnets = append(subnets, subnet)
	}
	return subnets, nil
}

func terraformEncodeSlice(items []*net.IPNet) string {
	encoded := ""
	for i, item := range items {
		encoded += fmt.Sprintf("\"%s\"", item.String())
		if i < len(items)-1 {
			encoded += ","
		}
	}
	return fmt.Sprintf("[%s]", encoded)
}

func (r *Resource) EnsureCreated(ctx context.Context, obj interface{}) error {
	cr, err := key.ToInstallation(obj)
	if err != nil {
		return microerror.Mask(err)
	}
	if cr.Spec.Provider != "aws" {
		r.logger.LogCtx(ctx, "message", "provider not supported", "provider", cr.Spec.Provider)
		return nil
	}

	secret, err := r.k8sClient.K8sClient().CoreV1().Secrets(cr.Namespace).Get(cr.Name, metav1.GetOptions{})
	if err != nil {
		return microerror.Mask(err)
	}

	if len(secret.Data["subnets_bastion"]) > 0 {
		return nil
	}

	subnets, err := r.findAvailableSubnets(ctx)
	if err != nil {
		return microerror.Mask(err)
	}
	secret.Data["subnets_bastion"] = []byte(terraformEncodeSlice(subnets))
	_, err = r.k8sClient.K8sClient().CoreV1().Secrets(cr.Namespace).Update(secret)

	return nil
}
