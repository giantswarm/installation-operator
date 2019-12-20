package bootstrap

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
	HiveRepoDir               = "hive"
	InventoryDir              = "hosts_inventory"
	EtcdBackupAwsAccessKeyEnv = "ETCD_BACKUP_AWS_ACCESS_KEY"
	EtcdBackupAwsSecretKeyEnv = "ETCD_BACKUP_AWS_SECRET_KEY"
	OpsctlGithubEnv           = "OPSCTL_GITHUB_TOKEN"
	VaultUnsealTokenEnv       = "VAULT_UNSEAL_TOKEN"
	BootstrapPlaybook         = "boostrap.yml"
	DefaultDC                 = "gridscale"
)

func (r *Resource) ensureBootstrapped(ctx context.Context, dc string) error {
	result := ""
	executor := ansible.New(ansible.Config{
		Env: map[string]string{
			EtcdBackupAwsAccessKeyEnv: "", // TODO
			EtcdBackupAwsSecretKeyEnv: "", // TODO
			OpsctlGithubEnv:           "", // TODO
			VaultUnsealTokenEnv:       "", // TODO
		},
		Logger: r.logger,
		Out:    bytes.NewBufferString(result),
	})
	playbook := &ansibler.AnsiblePlaybookCmd{
		Playbook: path.Join(HiveRepoDir, BootstrapPlaybook),
		Options: &ansibler.AnsiblePlaybookOptions{
			ExtraVars: map[string]interface{}{"dc": dc},
			Inventory: path.Join(HiveRepoDir, InventoryDir, dc),
		},
		Exec: executor,
	}
	err := playbook.Run()
	if err != nil {
		return microerror.Mask(err)
	}
	// TODO: save vault tokens from results
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
	err = r.ensureBootstrapped(ctx, DefaultDC)
	if err != nil {
		return microerror.Mask(err)
	}
	return nil
}
