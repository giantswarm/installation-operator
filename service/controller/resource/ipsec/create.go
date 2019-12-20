package ipsec

import (
	"bytes"
	"context"
	"path"

	"github.com/apenella/go-ansible"
	"github.com/giantswarm/microerror"

	"github.com/giantswarm/installation-operator/service/controller/ansible"
	"github.com/giantswarm/installation-operator/service/controller/key"
)

const (
	PingdomPasswordEnv = "PINGDOM_PASS"
	VpnRepoDir         = "vpn"
	InventoryDir       = "hosts_inventory"
)

func (r *Resource) ensureVpnCreated(ctx context.Context, dc string) error {
	result := ""
	executor := ansible.New(ansible.Config{
		Env: map[string]string{
			PingdomPasswordEnv: "", // TODO
		},
		Logger: r.logger,
		Out:    bytes.NewBufferString(result),
	})
	playbook := &ansibler.AnsiblePlaybookCmd{
		Playbook: path.Join(VpnRepoDir, "vpn.yml"),
		Options: &ansibler.AnsiblePlaybookOptions{
			ExtraVars: map[string]interface{}{"dc": dc},
			Inventory: path.Join(VpnRepoDir, InventoryDir, dc),
		},
		Exec: executor,
	}
	err := playbook.Run()
	if err != nil {
		return microerror.Mask(err)
	}
	return nil
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
	dcs := []string{"gridscale", "vultr"}
	for _, dc := range dcs {
		err := r.ensureVpnCreated(ctx, dc)
		if err != nil {
			return microerror.Mask(err)
		}
	}
	return nil
}
