package cloudstack

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"

	anywherev1 "github.com/aws/eks-anywhere/pkg/api/v1alpha1"
	"github.com/aws/eks-anywhere/pkg/logger"
)

const controlEndpointDefaultPort = "6443"

type Validator struct {
	cmk            ProviderCmkClient
	machineConfigs map[string]*anywherev1.CloudStackMachineConfig
}

func NewValidator(cmk ProviderCmkClient, machineConfigs map[string]*anywherev1.CloudStackMachineConfig) *Validator {
	return &Validator{
		cmk:            cmk,
		machineConfigs: machineConfigs,
	}
}

type ProviderCmkClient interface {
	ValidateCloudStackConnection(ctx context.Context) error
	ValidateServiceOfferingPresent(ctx context.Context, domain string, zone anywherev1.CloudStackResourceRef, account string, serviceOffering anywherev1.CloudStackResourceRef) error
	ValidateTemplatePresent(ctx context.Context, domain string, zone anywherev1.CloudStackResourceRef, account string, template anywherev1.CloudStackResourceRef) error
	ValidateAffinityGroupsPresent(ctx context.Context, domain string, zone anywherev1.CloudStackResourceRef, account string, affinityGroupIds []string) error
	ValidateZonePresent(ctx context.Context, zone anywherev1.CloudStackResourceRef) error
}

func (v *Validator) validateCloudStackAccess(ctx context.Context) error {
	if err := v.cmk.ValidateCloudStackConnection(ctx); err != nil {
		return fmt.Errorf("failed validating connection to vCenter: %v", err)
	}
	logger.MarkPass("Connected to server")

	return nil
}

func (v *Validator) ValidateCloudStackDatacenterConfig(ctx context.Context, datacenterConfig *anywherev1.CloudStackDatacenterConfig) error {
	if err := v.cmk.ValidateZonePresent(ctx, datacenterConfig.Spec.Zone); err != nil {
		return err
	}
	logger.MarkPass("Datacenter validated")

	return nil
}

// TODO: dry out machine configs validations
func (v *Validator) ValidateClusterMachineConfigs(ctx context.Context, cloudStackClusterSpec *spec) error {
	var etcdMachineConfig *anywherev1.CloudStackMachineConfig

	if len(cloudStackClusterSpec.Cluster.Spec.ControlPlaneConfiguration.Endpoint.Host) <= 0 {
		return fmt.Errorf("cluster controlPlaneConfiguration.Endpoint.Host is not set or is empty")
	}
	if len(cloudStackClusterSpec.datacenterConfig.Spec.Domain) <= 0 {
		return fmt.Errorf("CloudStackDatacenterConfig domain is not set or is empty")
	}
	if len(cloudStackClusterSpec.datacenterConfig.Spec.Zone.Value) <= 0 {
		return fmt.Errorf("CloudStackDatacenterConfig zone is not set or is empty")
	}
	if len(cloudStackClusterSpec.datacenterConfig.Spec.Network.Value) <= 0 {
		return fmt.Errorf("CloudStackDatacenterConfig network is not set or is empty")
	}
	if cloudStackClusterSpec.Cluster.Spec.ControlPlaneConfiguration.MachineGroupRef == nil {
		return fmt.Errorf("must specify machineGroupRef for control plane")
	}

	if cloudStackClusterSpec.Cluster.Spec.WorkerNodeGroupConfigurations[0].MachineGroupRef == nil {
		return fmt.Errorf("must specify machineGroupRef for worker nodes")
	}

	controlPlaneMachineConfig := cloudStackClusterSpec.controlPlaneMachineConfig()
	if controlPlaneMachineConfig == nil {
		return fmt.Errorf("cannot find CloudStackMachineConfig %v for control plane", cloudStackClusterSpec.Cluster.Spec.ControlPlaneConfiguration.MachineGroupRef.Name)
	}

	if cloudStackClusterSpec.Cluster.Spec.ExternalEtcdConfiguration != nil {
		if cloudStackClusterSpec.Cluster.Spec.ExternalEtcdConfiguration.MachineGroupRef == nil {
			return fmt.Errorf("must specify machineGroupRef for etcd machines")
		}
		etcdMachineConfig = cloudStackClusterSpec.etcdMachineConfig()
		if etcdMachineConfig == nil {
			return fmt.Errorf("cannot find CloudStackMachineConfig %v for etcd machines", cloudStackClusterSpec.Cluster.Spec.ExternalEtcdConfiguration.MachineGroupRef.Name)
		}
		if etcdMachineConfig.Spec.Template != controlPlaneMachineConfig.Spec.Template {
			return fmt.Errorf("control plane and etcd machines must have the same template specified")
		}
	}

	if cloudStackClusterSpec.datacenterConfig.Namespace != cloudStackClusterSpec.Namespace {
		return fmt.Errorf(
			"CloudStackDatacenterConfig and Cluster objects must have the same namespace: CloudstackDatacenterConfig namespace=%s; Cluster namespace=%s",
			cloudStackClusterSpec.datacenterConfig.Namespace,
			cloudStackClusterSpec.Namespace,
		)
	}

	workerNodeGroupMachineConfig, ok := v.machineConfigs[cloudStackClusterSpec.Cluster.Spec.WorkerNodeGroupConfigurations[0].MachineGroupRef.Name]
	if !ok {
		return fmt.Errorf("cannot find CloudStackMachineConfig %v for worker nodes", cloudStackClusterSpec.Cluster.Spec.WorkerNodeGroupConfigurations[0].MachineGroupRef.Name)
	}
	if controlPlaneMachineConfig.Spec.Template != workerNodeGroupMachineConfig.Spec.Template {
		if controlPlaneMachineConfig.Spec.Template != workerNodeGroupMachineConfig.Spec.Template {
			return fmt.Errorf("control plane and worker nodes must have the same template specified")
		}
	}

	hostWithPort, err := v.validateControlPlaneHostAndApplyDefaultPort(cloudStackClusterSpec.Cluster.Spec.ControlPlaneConfiguration.Endpoint.Host)
	if err != nil {
		return fmt.Errorf("failed to validate controlPlaneConfiguration.Endpoint.Host: %v", err)
	}
	cloudStackClusterSpec.Cluster.Spec.ControlPlaneConfiguration.Endpoint.Host = hostWithPort

	for _, machineConfig := range v.machineConfigs {
		if machineConfig.Namespace != cloudStackClusterSpec.Namespace {
			return fmt.Errorf(
				"CloudStackMachineConfig and Cluster objects must have the same namespace: CloudStackMachineConfig namespace=%s; Cluster namespace=%s",
				machineConfig.Namespace,
				cloudStackClusterSpec.Namespace,
			)
		}
		if len(machineConfig.Spec.Users) <= 0 {
			machineConfig.Spec.Users = []anywherev1.UserConfiguration{{}}
		}
		if len(machineConfig.Spec.Users[0].SshAuthorizedKeys) <= 0 {
			machineConfig.Spec.Users[0].SshAuthorizedKeys = []string{""}
		}
		if machineConfig.Spec.Template.Value == "" {
			return fmt.Errorf("template is not set for CloudStackMachineConfig %s. Default template is not supported in CloudStack, please provide a template name", machineConfig.Name)
		}
		if err = v.validateMachineConfig(ctx, cloudStackClusterSpec.datacenterConfig.Spec, etcdMachineConfig); err != nil {
			return fmt.Errorf("machine config %s validation failed: %v", machineConfig.Name, err)
		}
	}
	logger.MarkPass("Validated cluster Machine Configs")

	return nil
}

func (v *Validator) validateMachineConfig(ctx context.Context, datacenterConfigSpec anywherev1.CloudStackDatacenterConfigSpec, machineConfig *anywherev1.CloudStackMachineConfig) error {
	domain := datacenterConfigSpec.Domain
	zone := datacenterConfigSpec.Zone
	account := datacenterConfigSpec.Account
	if err := v.cmk.ValidateTemplatePresent(ctx, domain, zone, account, machineConfig.Spec.Template); err != nil {
		return fmt.Errorf("validating template: %v", err)
	}
	if err := v.cmk.ValidateServiceOfferingPresent(ctx, domain, zone, account, machineConfig.Spec.ComputeOffering); err != nil {
		return fmt.Errorf("validating service offering: %v", err)
	}
	if len(machineConfig.Spec.AffinityGroupIds) > 0 {
		if err := v.cmk.ValidateAffinityGroupsPresent(ctx, domain, zone, account, machineConfig.Spec.AffinityGroupIds); err != nil {
			return fmt.Errorf("validating affinity group ids: %v", err)
		}
	}

	return nil
}

func (v *Validator) validateControlPlaneHostAndApplyDefaultPort(pHost string) (string, error) {
	_, port, err := net.SplitHostPort(pHost)
	portWithHost := pHost
	if err != nil {
		if strings.Contains(err.Error(), "missing port") {
			port = controlEndpointDefaultPort
			portWithHost = fmt.Sprintf("%s:%s", pHost, port)
		} else {
			return "", fmt.Errorf("host %s is invalid: %v", pHost, err.Error())
		}
	}
	_, err = strconv.Atoi(port)
	if err != nil {
		return "", fmt.Errorf("host %s has an invalid port: %v", pHost, err.Error())
	}
	return portWithHost, nil
}
