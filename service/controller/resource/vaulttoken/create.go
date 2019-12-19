package vaulttoken

import (
	"context"

	"github.com/giantswarm/microerror"

	"github.com/giantswarm/installation-operator/service/controller/key"
)

func (r *Resource) EnsureCreated(ctx context.Context, obj interface{}) error {
	cr, err := key.ToInstallation(obj)
	if err != nil {
		return microerror.Mask(err)
	}
	if cr.Spec.Provider != "aws" {
		r.logger.LogCtx(ctx, "message", "provider not supported", "provider", cr.Spec.Provider)
		return nil
	}
	if cr.Status.NodeVaultToken != "" {
		return nil
	}
	token, err := createVaultToken()
	if err != nil {
		return microerror.Mask(err)
	}
	cr.Status.NodeVaultToken = token
	_, err = r.k8sClient.G8sClient().CoreV1alpha1().Installations(cr.Namespace).UpdateStatus(&cr)
	if err != nil {
		return microerror.Mask(err)
	}
	return nil
}
