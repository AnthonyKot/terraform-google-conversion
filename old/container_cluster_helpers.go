package container

import (
	"fmt"
	"strings"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/cai2hcl/converters/utils"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/tpgresource"
	"google.golang.org/api/container/v1"
)

//
// Node Pool Helper Functions
//

// flattenNodeConfig converts NodeConfig to a flattened map for HCL
func flattenNodeConfig(config *container.NodeConfig) []interface{} {
	if config == nil {
		return nil
	}
	nodeConfig := make(map[string]interface{})
	hasNonDefaultConfig := false // Track if any non-default value is set

	// Machine Type: Typically non-default if specified. Include if present.
	if config.MachineType != "" {
		// TODO: Check against implicit default like e2-medium? Hard to determine.
		nodeConfig["machine_type"] = config.MachineType
		hasNonDefaultConfig = true
	}

	// Disk Size: Assume 100GB is a common implicit default. Omit if matches.
	if config.DiskSizeGb > 0 && config.DiskSizeGb != defaultNodeDiskSizeGb {
		nodeConfig["disk_size_gb"] = config.DiskSizeGb
		hasNonDefaultConfig = true
	}

	// Disk Type: Assume 'pd-standard' is a common implicit default. Omit if matches.
	if config.DiskType != "" && config.DiskType != defaultNodeDiskType {
		nodeConfig["disk_type"] = config.DiskType
		hasNonDefaultConfig = true
	}

	// OAuth Scopes: Omit only if list is empty (simplification, actual defaults are complex).
	if len(config.OauthScopes) > 0 {
		// TODO: Consider comparing against expected default scopes for the SA? Very complex.
		nodeConfig["oauth_scopes"] = config.OauthScopes
		hasNonDefaultConfig = true
	}

	// Service Account: Omit if "default" or matches Compute Engine default (if checkable).
	if config.ServiceAccount != "" && config.ServiceAccount != defaultNodeServiceAccount {
		nodeConfig["service_account"] = config.ServiceAccount
		hasNonDefaultConfig = true
	}

	// Boot Disk KMS Key: Include if set (no default expected).
	if config.BootDiskKmsKey != "" {
		nodeConfig["boot_disk_kms_key"] = config.BootDiskKmsKey
		hasNonDefaultConfig = true
	}

	// Image Type: Skip if matches expected default.
	if config.ImageType != "" && config.ImageType != defaultNodeImageType {
		nodeConfig["image_type"] = config.ImageType
		hasNonDefaultConfig = true
	}

	// Labels: Include only if not empty.
	if len(config.Labels) > 0 {
		nodeConfig["labels"] = config.Labels
		hasNonDefaultConfig = true
	}

	// Local SSD Count: Include if non-zero.
	if config.LocalSsdCount > 0 {
		nodeConfig["local_ssd_count"] = config.LocalSsdCount
		hasNonDefaultConfig = true
	}

	// Metadata: Include if not empty.
	if len(config.Metadata) > 0 {
		nodeConfig["metadata"] = config.Metadata
		hasNonDefaultConfig = true
	}

	// Min CPU Platform: Include if set (no default expected).
	if config.MinCpuPlatform != "" {
		nodeConfig["min_cpu_platform"] = config.MinCpuPlatform
		hasNonDefaultConfig = true
	}

	// Preemptible: Include if true (default is false).
	if config.Preemptible {
		nodeConfig["preemptible"] = config.Preemptible
		hasNonDefaultConfig = true
	}

	// Spot: Include if true (default is false).
	if config.Spot {
		nodeConfig["spot"] = config.Spot
		hasNonDefaultConfig = true
	}

	// Tags: Include if not empty.
	if len(config.Tags) > 0 {
		nodeConfig["tags"] = config.Tags
		hasNonDefaultConfig = true
	}

	// Taint: Include if not empty.
	if len(config.Taints) > 0 {
		nodeConfig["taint"] = flattenNodeTaints(config.Taints)
		hasNonDefaultConfig = true
	}

	// Resource Labels: Include if not empty.
	if len(config.ResourceLabels) > 0 {
		nodeConfig["resource_labels"] = config.ResourceLabels
		hasNonDefaultConfig = true
	}

	// Shielded Instance Config: Include if non-default settings.
	if flattened := flattenShieldedInstanceConfig(config.ShieldedInstanceConfig); flattened != nil {
		nodeConfig["shielded_instance_config"] = flattened
		hasNonDefaultConfig = true
	}

	// Workload Metadata Config: Include if set.
	if flattened := flattenWorkloadMetadataConfig(config.WorkloadMetadataConfig); flattened != nil {
		nodeConfig["workload_metadata_config"] = flattened
		hasNonDefaultConfig = true
	}

	// Sandbox Config: Include if type is specified.
	if config.SandboxConfig != nil && config.SandboxConfig.Type != "" {
		nodeConfig["sandbox_config"] = []interface{}{
			map[string]interface{}{
				"sandbox_type": config.SandboxConfig.Type,
			},
		}
		hasNonDefaultConfig = true
	}

	// Node Group: Include if set.
	if config.NodeGroup != "" {
		nodeConfig["node_group"] = config.NodeGroup
		hasNonDefaultConfig = true
	}

	// Reservation Affinity: Include if non-nil and has custom settings.
	if flattened := flattenReservationAffinity(config.ReservationAffinity); flattened != nil {
		nodeConfig["reservation_affinity"] = flattened
		hasNonDefaultConfig = true
	}

	// Linux Node Config: Include if non-nil and has settings.
	if flattened := flattenLinuxNodeConfig(config.LinuxNodeConfig); flattened != nil {
		nodeConfig["linux_node_config"] = flattened
		hasNonDefaultConfig = true
	}

	// Kubelet Config: Include if non-nil and has settings.
	if flattened := flattenKubeletConfig(config.KubeletConfig); flattened != nil {
		nodeConfig["kubelet_config"] = flattened
		hasNonDefaultConfig = true
	}

	// Guest Accelerators: Include if non-empty.
	if len(config.Accelerators) > 0 {
		nodeConfig["guest_accelerator"] = flattenAccelerators(config.Accelerators)
		hasNonDefaultConfig = true
	}

	// GPU Sharing Config: Include if non-nil and has settings.
	if flattened := flattenGpuSharingConfig(config.GpuSharingConfig); flattened != nil {
		nodeConfig["gvnic"] = flattened
		hasNonDefaultConfig = true
	}

	// GPU Driver Installation Config: Include if non-nil and has settings.
	if flattened := flattenGpuDriverInstallationConfig(config.GpuDriverInstallationConfig); flattened != nil {
		nodeConfig["gpu_driver_installation_config"] = flattened
		hasNonDefaultConfig = true
	}

	// GVNIC: Include if non-nil and has settings.
	if flattened := flattenGvnic(config.Gvnic); flattened != nil {
		nodeConfig["gvnic"] = flattened
		hasNonDefaultConfig = true
	}

	// Confidential Nodes: Include if non-nil and has settings.
	if flattened := flattenConfidentialNodes(config.ConfidentialNodes); flattened != nil {
		nodeConfig["confidential_nodes"] = flattened
		hasNonDefaultConfig = true
	}

	// Ephemeral Storage Config: Include if local_ssd_count > 0
	if config.EphemeralStorageConfig != nil && config.EphemeralStorageConfig.LocalSsdCount > 0 {
		nodeConfig["ephemeral_storage_config"] = []interface{}{
			map[string]interface{}{
				"local_ssd_count": config.EphemeralStorageConfig.LocalSsdCount,
			},
		}
		hasNonDefaultConfig = true
	}

	// Fast Socket: Include if enabled
	if config.FastSocket != nil && config.FastSocket.Enabled {
		nodeConfig["fast_socket"] = []interface{}{
			map[string]interface{}{
				"enabled": config.FastSocket.Enabled,
			},
		}
		hasNonDefaultConfig = true
	}

	// Error: If no non-default configurations were found, return nil.
	if !hasNonDefaultConfig {
		return nil
	}

	return []interface{}{nodeConfig}
}

// flattenWorkloadMetadataConfig converts WorkloadMetadataConfig to a map
func flattenWorkloadMetadataConfig(config *container.WorkloadMetadataConfig) []interface{} {
	if config == nil {
		return nil
	}

	// Mode is the only field we need to check
	if config.Mode == defaultWlcMode {
		return nil // Return nil if default
	}

	return []interface{}{
		map[string]interface{}{
			"mode": config.Mode,
		},
	}
}

// flattenShieldedInstanceConfig converts ShieldedInstanceConfig to a map
func flattenShieldedInstanceConfig(config *container.ShieldedInstanceConfig) []interface{} {
	if config == nil {
		return nil
	}

	// If all fields are false or default, return nil
	if !config.EnableSecureBoot && !config.EnableIntegrityMonitoring {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"enable_secure_boot":          config.EnableSecureBoot,
			"enable_integrity_monitoring": config.EnableIntegrityMonitoring,
		},
	}
}

// flattenAccelerators converts AcceleratorConfig list to a list of maps
func flattenAccelerators(accelerators []*container.AcceleratorConfig) []interface{} {
	if len(accelerators) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(accelerators))
	for _, accel := range accelerators {
		if accel.AcceleratorCount > 0 {
			result = append(result, map[string]interface{}{
				"type":  accel.AcceleratorType,
				"count": accel.AcceleratorCount,
				// No need to check for gpuDriverInstallationConfig or gpuSharingConfig
				// as they are handled separately
			})
		}
	}

	return result
}

// flattenGpuSharingConfig converts GPUSharingConfig to a map
func flattenGpuSharingConfig(config *container.GPUSharingConfig) []interface{} {
	if config == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"gpu_sharing_strategy":     config.GpuSharingStrategy,
			"max_shared_clients_per_gpu": config.MaxSharedClientsPerGpu,
		},
	}
}

// flattenGpuDriverInstallationConfig converts GPUDriverInstallationConfig to a map
func flattenGpuDriverInstallationConfig(config *container.GPUDriverInstallationConfig) []interface{} {
	if config == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"gpu_driver_version": config.GpuDriverVersion,
		},
	}
}

// flattenReservationAffinity converts ReservationAffinity to a map
func flattenReservationAffinity(config *container.ReservationAffinity) []interface{} {
	if config == nil || config.ConsumeReservationType == defaultReservationType {
		return nil
	}

	result := map[string]interface{}{
		"consume_reservation_type": config.ConsumeReservationType,
	}

	if config.ConsumeReservationType == "SPECIFIC_RESERVATION" {
		result["key"] = config.Key
		result["values"] = config.Values
	}

	return []interface{}{result}
}

// flattenConfidentialNodes converts ConfidentialNodes to a map
func flattenConfidentialNodes(config *container.ConfidentialNodes) []interface{} {
	if config == nil || !config.Enabled {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"enabled": config.Enabled,
		},
	}
}

// flattenKubeletConfig converts NodeKubeletConfig to a map
func flattenKubeletConfig(config *container.NodeKubeletConfig) []interface{} {
	if config == nil {
		return nil
	}

	kc := make(map[string]interface{})
	hasNonDefaultConfig := false

	// Add CPU Manager Policy if not default
	if config.CpuManagerPolicy != "" && config.CpuManagerPolicy != defaultKubeletCpuPolicy {
		kc["cpu_manager_policy"] = config.CpuManagerPolicy
		hasNonDefaultConfig = true
	}

	// Add CPU cfs quota if not default
	if config.CpuCfsQuota != nil {
		kc["cpu_cfs_quota"] = *config.CpuCfsQuota
		hasNonDefaultConfig = true
	}

	// Add CPU cfs quota period if set
	if config.CpuCfsQuotaPeriod != "" {
		kc["cpu_cfs_quota_period"] = config.CpuCfsQuotaPeriod
		hasNonDefaultConfig = true
	}

	// Add pod PIDS limit if not default
	if config.PodPidsLimit > 0 {
		kc["pod_pids_limit"] = config.PodPidsLimit
		hasNonDefaultConfig = true
	}

	// Add insecure_kubelet_readonly_port_enabled if true
	if config.InsecureKubeletReadonlyPortEnabled {
		kc["insecure_kubelet_readonly_port_enabled"] = true
		hasNonDefaultConfig = true
	}

	// Error: If no non-default configurations were found, return nil.
	if !hasNonDefaultConfig {
		return nil
	}

	return []interface{}{kc}
}

// flattenLinuxNodeConfig converts LinuxNodeConfig to a map
func flattenLinuxNodeConfig(config *container.LinuxNodeConfig) []interface{} {
	if config == nil || config.CgroupMode == defaultLinuxCgroupMode {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"cgroup_mode": config.CgroupMode,
		},
	}
}

// flattenNodeTaints converts NodeTaint list to a list of maps
func flattenNodeTaints(taints []*container.NodeTaint) []interface{} {
	if len(taints) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(taints))
	for _, taint := range taints {
		result = append(result, map[string]interface{}{
			"key":    taint.Key,
			"value":  taint.Value,
			"effect": taint.Effect,
		})
	}
	return result
}

// flattenGvnic converts VirtualNIC to a map
func flattenGvnic(config *container.VirtualNIC) []interface{} {
	if config == nil || !config.Enabled {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"enabled": config.Enabled,
		},
	}
}

// flattenNodePoolAutoscaling converts NodePoolAutoscaling to a map
func flattenNodePoolAutoscaling(autoscaling *container.NodePoolAutoscaling) []interface{} {
	if autoscaling == nil || !autoscaling.Enabled {
		return nil
	}

	result := map[string]interface{}{
		"min_node_count": autoscaling.MinNodeCount,
		"max_node_count": autoscaling.MaxNodeCount,
	}

	// Include optional fields only if they're set
	if autoscaling.LocationPolicy != "" {
		result["location_policy"] = autoscaling.LocationPolicy
	}

	if autoscaling.TotalMinNodeCount > 0 {
		result["total_min_node_count"] = autoscaling.TotalMinNodeCount
	}

	if autoscaling.TotalMaxNodeCount > 0 {
		result["total_max_node_count"] = autoscaling.TotalMaxNodeCount
	}

	return []interface{}{result}
}

// flattenNodeManagement converts NodeManagement to a map
func flattenNodeManagement(mgmt *container.NodeManagement) []interface{} {
	if mgmt == nil {
		return nil
	}

	hasNonDefaultSettings := false
	result := make(map[string]interface{})

	// Include autoRepair only if false (default is true)
	if mgmt.AutoRepair != nil && !*mgmt.AutoRepair {
		result["auto_repair"] = *mgmt.AutoRepair
		hasNonDefaultSettings = true
	}

	// Include autoUpgrade only if false (default is true in most cases)
	// NOTE: This might vary by release channel, so consider adding nuance here
	if mgmt.AutoUpgrade != nil && !*mgmt.AutoUpgrade {
		result["auto_upgrade"] = *mgmt.AutoUpgrade
		hasNonDefaultSettings = true
	}

	if !hasNonDefaultSettings {
		return nil
	}
	return []interface{}{result}
}

// flattenNodePoolUpgradeSettings converts UpgradeSettings to a map
func flattenNodePoolUpgradeSettings(settings *container.UpgradeSettings) []interface{} {
	if settings == nil {
		return nil
	}

	result := make(map[string]interface{})

	// Include only fields that are set
	if settings.MaxSurge > 0 {
		result["max_surge"] = settings.MaxSurge
	}

	if settings.MaxUnavailable > 0 {
		result["max_unavailable"] = settings.MaxUnavailable
	}

	if settings.Strategy != "" {
		result["strategy"] = settings.Strategy
	}

	if blueGreenSettings := flattenBlueGreenSettings(settings.BlueGreenSettings); blueGreenSettings != nil {
		result["blue_green_settings"] = blueGreenSettings
	}

	if len(result) == 0 {
		return nil
	}
	return []interface{}{result}
}

// flattenBlueGreenSettings converts BlueGreenSettings to a map
func flattenBlueGreenSettings(settings *container.BlueGreenSettings) []interface{} {
	if settings == nil {
		return nil
	}

	result := make(map[string]interface{})

	// Include only fields that are set
	if settings.NodePoolSoakDuration != nil && settings.NodePoolSoakDuration.Seconds > 0 {
		result["node_pool_soak_duration"] = fmt.Sprintf("%ds", settings.NodePoolSoakDuration.Seconds)
	}

	if settings.StandardRolloutPolicy != nil {
		standardPolicy := make(map[string]interface{})

		if settings.StandardRolloutPolicy.BatchSoakDuration != nil && settings.StandardRolloutPolicy.BatchSoakDuration.Seconds > 0 {
			standardPolicy["batch_soak_duration"] = fmt.Sprintf("%ds", settings.StandardRolloutPolicy.BatchSoakDuration.Seconds)
		}

		if settings.StandardRolloutPolicy.BatchNodeCount > 0 {
			standardPolicy["batch_node_count"] = settings.StandardRolloutPolicy.BatchNodeCount
		}

		if settings.StandardRolloutPolicy.BatchPercentage > 0 {
			standardPolicy["batch_percentage"] = settings.StandardRolloutPolicy.BatchPercentage
		}

		if len(standardPolicy) > 0 {
			result["standard_rollout_policy"] = []interface{}{standardPolicy}
		}
	}

	if len(result) == 0 {
		return nil
	}
	return []interface{}{result}
}

// flattenMaxPodsConstraint converts MaxPodsConstraint to a map
func flattenMaxPodsConstraint(constraint *container.MaxPodsConstraint) []interface{} {
	if constraint == nil || constraint.MaxPodsPerNode == 0 {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"max_pods_per_node": constraint.MaxPodsPerNode,
		},
	}
}

// flattenNodeNetworkConfig converts NodeNetworkConfig to a map
func flattenNodeNetworkConfig(config *container.NodeNetworkConfig) []interface{} {
	if config == nil {
		return nil
	}

	result := make(map[string]interface{})

	// Include fields that are set and non-default
	if config.PodRange != "" {
		result["pod_range"] = config.PodRange
	}

	if config.PodIpv4CidrBlock != "" {
		result["pod_ipv4_cidr_block"] = config.PodIpv4CidrBlock
	}

	if config.EnablePrivateNodes {
		result["enable_private_nodes"] = config.EnablePrivateNodes
	}

	if config.CreatePodRange {
		result["create_pod_range"] = config.CreatePodRange
	}

	if len(result) == 0 {
		return nil
	}
	return []interface{}{result}
}

// flattenPlacementPolicy converts PlacementPolicy to a map
func flattenPlacementPolicy(policy *container.PlacementPolicy) []interface{} {
	if policy == nil {
		return nil
	}

	result := make(map[string]interface{})

	// Include fields that are set
	if policy.Type != "" {
		result["type"] = policy.Type
	}

	if policy.TpuTopology != "" {
		result["tpu_topology"] = policy.TpuTopology
	}

	if len(result) == 0 {
		return nil
	}
	return []interface{}{result}
}

//
// Cluster Helper Functions
//

// flattenIPAllocationPolicy converts IPAllocationPolicy to a map
func flattenIPAllocationPolicy(policy *container.IPAllocationPolicy) []interface{} {
	if policy == nil {
		return nil
	}

	// Create a map for the allocation policy
	result := make(map[string]interface{})

	// Include only fields that are set and not empty
	if policy.UseIpAliases {
		result["use_ip_aliases"] = policy.UseIpAliases
	}

	if policy.CreateSubnetwork {
		result["create_subnetwork"] = policy.CreateSubnetwork
	}

	if policy.SubnetworkName != "" {
		result["subnetwork_name"] = policy.SubnetworkName
	}

	if policy.ClusterIpv4CidrBlock != "" {
		result["cluster_ipv4_cidr_block"] = policy.ClusterIpv4CidrBlock
	}

	if policy.NodeIpv4CidrBlock != "" {
		result["node_ipv4_cidr_block"] = policy.NodeIpv4CidrBlock
	}

	if policy.ServicesIpv4CidrBlock != "" {
		result["services_ipv4_cidr_block"] = policy.ServicesIpv4CidrBlock
	}

	if policy.ClusterSecondaryRangeName != "" {
		result["cluster_secondary_range_name"] = policy.ClusterSecondaryRangeName
	}

	if policy.ServicesSecondaryRangeName != "" {
		result["services_secondary_range_name"] = policy.ServicesSecondaryRangeName
	}

	if policy.StackType != defaultIpAllocationStackType {
		result["stack_type"] = policy.StackType
	}

	if policy.IpvDualStackType != "" {
		result["ipv_dual_stack_type"] = policy.IpvDualStackType
	}

	if policy.ServicesIpv6CidrBlock != "" {
		result["services_ipv6_cidr_block"] = policy.ServicesIpv6CidrBlock
	}

	if policy.ClusterIpv6CidrBlock != "" {
		result["cluster_ipv6_cidr_block"] = policy.ClusterIpv6CidrBlock
	}

	if policy.SubnetIpv6CidrBlock != "" {
		result["subnet_ipv6_cidr_block"] = policy.SubnetIpv6CidrBlock
	}

	// Return nil if no fields were set
	if len(result) == 0 {
		return nil
	}

	return []interface{}{result}
}

// flattenAddonsConfig converts AddonsConfig to a map
func flattenAddonsConfig(config *container.AddonsConfig) []interface{} {
	if config == nil {
		return nil
	}

	result := make(map[string]interface{})
	hasNonDefaultConfig := false

	// HttpLoadBalancing - defaults to enabled, so only include if disabled
	if config.HttpLoadBalancing != nil && !config.HttpLoadBalancing.Disabled {
		// Skip if it matches default (enabled)
	} else if config.HttpLoadBalancing != nil {
		result["http_load_balancing"] = []interface{}{
			map[string]interface{}{
				"disabled": config.HttpLoadBalancing.Disabled,
			},
		}
		hasNonDefaultConfig = true
	}

	// HorizontalPodAutoscaling - defaults to enabled, so only include if disabled
	if config.HorizontalPodAutoscaling != nil && !config.HorizontalPodAutoscaling.Disabled {
		// Skip if it matches default (enabled)
	} else if config.HorizontalPodAutoscaling != nil {
		result["horizontal_pod_autoscaling"] = []interface{}{
			map[string]interface{}{
				"disabled": config.HorizontalPodAutoscaling.Disabled,
			},
		}
		hasNonDefaultConfig = true
	}

	// NetworkPolicyConfig - defaults to disabled, so include if enabled
	if config.NetworkPolicyConfig != nil && !config.NetworkPolicyConfig.Disabled {
		result["network_policy_config"] = []interface{}{
			map[string]interface{}{
				"disabled": config.NetworkPolicyConfig.Disabled,
			},
		}
		hasNonDefaultConfig = true
	}

	// CloudRunConfig - include if set (no clear default)
	if config.CloudRunConfig != nil {
		cloudRunConfig := map[string]interface{}{
			"disabled": config.CloudRunConfig.Disabled,
		}

		if config.CloudRunConfig.LoadBalancerType != "" {
			cloudRunConfig["load_balancer_type"] = config.CloudRunConfig.LoadBalancerType
		}

		result["cloudrun_config"] = []interface{}{cloudRunConfig}
		hasNonDefaultConfig = true
	}

	// GcpFilestoreCsiDriverConfig - include if set (no clear default)
	if config.GcpFilestoreCsiDriverConfig != nil {
		result["gcp_filestore_csi_driver_config"] = []interface{}{
			map[string]interface{}{
				"enabled": !config.GcpFilestoreCsiDriverConfig.Disabled,
			},
		}
		hasNonDefaultConfig = true
	}

	// GcePersistentDiskCsiDriverConfig - defaults to enabled, include if disabled
	if config.GcePersistentDiskCsiDriverConfig != nil && config.GcePersistentDiskCsiDriverConfig.Disabled {
		result["gce_persistent_disk_csi_driver_config"] = []interface{}{
			map[string]interface{}{
				"enabled": !config.GcePersistentDiskCsiDriverConfig.Disabled,
			},
		}
		hasNonDefaultConfig = true
	}

	// DnsCacheConfig - include if set (no clear default)
	if config.DnsCacheConfig != nil {
		result["dns_cache_config"] = []interface{}{
			map[string]interface{}{
				"enabled": !config.DnsCacheConfig.Disabled,
			},
		}
		hasNonDefaultConfig = true
	}

	// ConfigConnectorConfig - include if set (no clear default)
	if config.ConfigConnectorConfig != nil {
		result["config_connector_config"] = []interface{}{
			map[string]interface{}{
				"enabled": !config.ConfigConnectorConfig.Disabled,
			},
		}
		hasNonDefaultConfig = true
	}

	// GcsFuseCsiDriverConfig - include if set (no clear default)
	if config.GcsFuseCsiDriverConfig != nil {
		result["gcs_fuse_csi_driver_config"] = []interface{}{
			map[string]interface{}{
				"enabled": !config.GcsFuseCsiDriverConfig.Disabled,
			},
		}
		hasNonDefaultConfig = true
	}

	if !hasNonDefaultConfig {
		return nil
	}

	return []interface{}{result}
}

// flattenNetworkPolicy converts NetworkPolicy to a map
func flattenNetworkPolicy(policy *container.NetworkPolicy) []interface{} {
	if policy == nil || policy.Provider == "" {
		return nil
	}

	result := map[string]interface{}{
		"provider": policy.Provider,
	}

	if policy.Enabled {
		result["enabled"] = policy.Enabled
	}

	return []interface{}{result}
}

// flattenMasterAuth converts MasterAuth to a map
func flattenMasterAuth(auth *container.MasterAuth) []interface{} {
	if auth == nil {
		return nil
	}

	// Check for any non-empty fields
	if auth.Username == "" && auth.Password == "" && auth.ClientCertificate == "" && 
	   auth.ClientKey == "" && auth.ClusterCaCertificate == "" {
		return nil
	}

	result := make(map[string]interface{})

	if auth.ClientCertificateConfig != nil {
		result["client_certificate_config"] = []interface{}{
			map[string]interface{}{
				"issue_client_certificate": auth.ClientCertificateConfig.IssueClientCertificate,
			},
		}
	}

	return []interface{}{result}
}

// flattenPrivateClusterConfigAdapted converts PrivateClusterConfig and NetworkConfig to a map
func flattenPrivateClusterConfigAdapted(config *container.PrivateClusterConfig, netConfig *container.NetworkConfig) []interface{} {
	if config == nil && (netConfig == nil || netConfig.PrivateIpv6GoogleAccess == "") {
		return nil
	}

	result := make(map[string]interface{})
	hasNonDefaultConfig := false

	// Handle fields from PrivateClusterConfig
	if config != nil {
		if config.EnablePrivateNodes {
			result["enable_private_nodes"] = config.EnablePrivateNodes
			hasNonDefaultConfig = true
		}

		if config.EnablePrivateEndpoint {
			result["enable_private_endpoint"] = config.EnablePrivateEndpoint
			hasNonDefaultConfig = true
		}

		if config.MasterIpv4CidrBlock != "" {
			result["master_ipv4_cidr_block"] = config.MasterIpv4CidrBlock
			hasNonDefaultConfig = true
		}

		if config.PrivateEndpoint != "" {
			result["private_endpoint"] = config.PrivateEndpoint
			hasNonDefaultConfig = true
		}

		if config.PublicEndpoint != "" {
			result["public_endpoint"] = config.PublicEndpoint
			hasNonDefaultConfig = true
		}

		if config.MasterGlobalAccessConfig != nil {
			result["master_global_access_config"] = []interface{}{
				map[string]interface{}{
					"enabled": config.MasterGlobalAccessConfig.Enabled,
				},
			}
			hasNonDefaultConfig = true
		}
	}

	// Handle fields from NetworkConfig
	if netConfig != nil && netConfig.PrivateIpv6GoogleAccess != "" {
		result["private_ipv6_google_access"] = netConfig.PrivateIpv6GoogleAccess
		hasNonDefaultConfig = true
	}

	if !hasNonDefaultConfig {
		return nil
	}

	return []interface{}{result}
}

// flattenReleaseChannel converts ReleaseChannel to a map
func flattenReleaseChannel(rc *container.ReleaseChannel) []interface{} {
	if rc == nil || rc.Channel == defaultReleaseChannel {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"channel": rc.Channel,
		},
	}
}

// flattenLoggingConfig converts LoggingConfig to a map
func flattenLoggingConfig(config *container.LoggingConfig) []interface{} {
	if config == nil || len(config.ComponentConfig.EnableComponents) == 0 {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"enable_components": config.ComponentConfig.EnableComponents,
		},
	}
}

// flattenMonitoringConfig converts MonitoringConfig to a map
func flattenMonitoringConfig(config *container.MonitoringConfig) []interface{} {
	if config == nil {
		return nil
	}

	result := make(map[string]interface{})
	hasNonDefaultConfig := false

	// Component Config
	if config.ComponentConfig != nil && len(config.ComponentConfig.EnableComponents) > 0 {
		result["enable_components"] = config.ComponentConfig.EnableComponents
		hasNonDefaultConfig = true
	}

	// Managed Prometheus
	if config.ManagedPrometheus != nil && config.ManagedPrometheus.Enabled {
		managedPrometheusConfig := map[string]interface{}{
			"enabled": config.ManagedPrometheus.Enabled,
		}

		// TODO: Add support for AutoMonitoringConfig from v1beta1
		// if config.ManagedPrometheus.AutoMonitoringConfig != nil {
		//     managedPrometheusConfig["auto_monitoring_config"] = []interface{}{
		//         map[string]interface{}{
		//             "scope": config.ManagedPrometheus.AutoMonitoringConfig.Scope,
		//         },
		//     }
		// }

		result["managed_prometheus"] = []interface{}{managedPrometheusConfig}
		hasNonDefaultConfig = true
	}

	if !hasNonDefaultConfig {
		return nil
	}

	return []interface{}{result}
}

// flattenClusterAutoscaling converts ClusterAutoscaling to a map
func flattenClusterAutoscaling(autoscaling *container.ClusterAutoscaling) []interface{} {
	if autoscaling == nil || !autoscaling.EnableNodeAutoprovisioning {
		return nil
	}

	result := map[string]interface{}{
		"enabled": autoscaling.EnableNodeAutoprovisioning,
	}

	// Resource Limits
	if len(autoscaling.ResourceLimits) > 0 {
		result["resource_limits"] = flattenResourceLimits(autoscaling.ResourceLimits)
	}

	// Auto-provisioning Node Pool Defaults
	if autoscaling.AutoprovisioningNodePoolDefaults != nil {
		result["auto_provisioning_defaults"] = flattenAutoprovisioningNodePoolDefaults(autoscaling.AutoprovisioningNodePoolDefaults)
	}

	// Auto-scaling Profile
	if autoscaling.AutoscalingProfile != "" {
		result["autoscaling_profile"] = autoscaling.AutoscalingProfile
	}

	return []interface{}{result}
}

// flattenResourceLimits converts ResourceLimit list to a list of maps
func flattenResourceLimits(limits []*container.ResourceLimit) []interface{} {
	if len(limits) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(limits))
	for _, limit := range limits {
		result = append(result, map[string]interface{}{
			"resource_type": limit.ResourceType,
			"minimum":       limit.Minimum,
			"maximum":       limit.Maximum,
		})
	}
	return result
}

// flattenAutoprovisioningNodePoolDefaults converts AutoprovisioningNodePoolDefaults to a map
func flattenAutoprovisioningNodePoolDefaults(defaults *container.AutoprovisioningNodePoolDefaults) []interface{} {
	if defaults == nil {
		return nil
	}

	result := make(map[string]interface{})

	// Set fields that are present
	if defaults.OauthScopes != nil && len(defaults.OauthScopes) > 0 {
		result["oauth_scopes"] = defaults.OauthScopes
	}

	if defaults.ServiceAccount != "" {
		result["service_account"] = defaults.ServiceAccount
	}

	if defaults.MachineType != "" {
		result["machine_type"] = defaults.MachineType
	}

	if defaults.DiskSizeGb > 0 {
		result["disk_size_gb"] = defaults.DiskSizeGb
	}

	if defaults.DiskType != "" {
		result["disk_type"] = defaults.DiskType
	}

	if defaults.MinCpuPlatform != "" {
		result["min_cpu_platform"] = defaults.MinCpuPlatform
	}

	if defaults.BootDiskKmsKey != "" {
		result["boot_disk_kms_key"] = defaults.BootDiskKmsKey
	}

	if defaults.Shielded != nil {
		shieldedConfig := make(map[string]interface{})
		
		if defaults.Shielded.EnableSecureBoot {
			shieldedConfig["enable_secure_boot"] = defaults.Shielded.EnableSecureBoot
		}
		
		if defaults.Shielded.EnableIntegrityMonitoring {
			shieldedConfig["enable_integrity_monitoring"] = defaults.Shielded.EnableIntegrityMonitoring
		}
		
		if len(shieldedConfig) > 0 {
			result["shielded_instance_config"] = []interface{}{shieldedConfig}
		}
	}

	if defaults.Management != nil {
		managementConfig := make(map[string]interface{})
		
		if defaults.Management.AutoRepair != nil {
			managementConfig["auto_repair"] = *defaults.Management.AutoRepair
		}
		
		if defaults.Management.AutoUpgrade != nil {
			managementConfig["auto_upgrade"] = *defaults.Management.AutoUpgrade
		}
		
		if len(managementConfig) > 0 {
			result["management"] = []interface{}{managementConfig}
		}
	}

	if defaults.UpgradeSettings != nil {
		upgradeConfig := make(map[string]interface{})
		
		if defaults.UpgradeSettings.MaxSurge > 0 {
			upgradeConfig["max_surge"] = defaults.UpgradeSettings.MaxSurge
		}
		
		if defaults.UpgradeSettings.MaxUnavailable > 0 {
			upgradeConfig["max_unavailable"] = defaults.UpgradeSettings.MaxUnavailable
		}
		
		if len(upgradeConfig) > 0 {
			result["upgrade_settings"] = []interface{}{upgradeConfig}
		}
	}

	if len(result) == 0 {
		return nil
	}
	return []interface{}{result}
}

// flattenDatabaseEncryption converts DatabaseEncryption to a map
func flattenDatabaseEncryption(config *container.DatabaseEncryption) []interface{} {
	if config == nil || config.State == "" || config.KeyName == "" {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"state":    config.State,
			"key_name": config.KeyName,
		},
	}
}

// flattenVerticalPodAutoscaling converts VerticalPodAutoscaling to a map
func flattenVerticalPodAutoscaling(vpa *container.VerticalPodAutoscaling) []interface{} {
	if vpa == nil || !vpa.Enabled {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"enabled": vpa.Enabled,
		},
	}
}

// flattenBinaryAuthorization converts BinaryAuthorization to a map
func flattenBinaryAuthorization(ba *container.BinaryAuthorization) []interface{} {
	if ba == nil || ba.EvaluationMode == defaultBinAuthEvalMode {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"evaluation_mode": ba.EvaluationMode,
		},
	}
}

// flattenCostManagementConfig converts CostManagementConfig to a map
func flattenCostManagementConfig(cmc *container.CostManagementConfig) []interface{} {
	if cmc == nil || !cmc.Enabled {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"enabled": cmc.Enabled,
		},
	}
}

// flattenDnsConfig converts DNSConfig to a map
func flattenDnsConfig(dns *container.DNSConfig) []interface{} {
	if dns == nil || (dns.ClusterDns == defaultDnsClusterDns && dns.ClusterDnsScope == defaultDnsClusterScope) {
		return nil
	}

	result := make(map[string]interface{})

	if dns.ClusterDns != defaultDnsClusterDns {
		result["cluster_dns"] = dns.ClusterDns
	}

	if dns.ClusterDnsScope != defaultDnsClusterScope {
		result["cluster_dns_scope"] = dns.ClusterDnsScope
	}

	if dns.ClusterDnsDomain != "" {
		result["cluster_dns_domain"] = dns.ClusterDnsDomain
	}

	if len(result) == 0 {
		return nil
	}
	return []interface{}{result}
}

// flattenSecurityPostureConfig converts SecurityPostureConfig to a map
func flattenSecurityPostureConfig(config *container.SecurityPostureConfig) []interface{} {
	if config == nil || (config.Mode == "" && config.VulnerabilityMode == "") {
		return nil
	}

	result := make(map[string]interface{})

	if config.Mode != "" {
		result["mode"] = config.Mode
	}

	if config.VulnerabilityMode != "" {
		result["vulnerability_mode"] = config.VulnerabilityMode
	}

	if len(result) == 0 {
		return nil
	}
	return []interface{}{result}
}

// flattenIdentityServiceConfig converts IdentityServiceConfig to a map
func flattenIdentityServiceConfig(isc *container.IdentityServiceConfig) []interface{} {
	if isc == nil || !isc.Enabled {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"enabled": isc.Enabled,
		},
	}
}

// flattenMeshCertificates converts MeshCertificates to a map
func flattenMeshCertificates(mc *container.MeshCertificates) []interface{} {
	if mc == nil || !mc.EnableCertificates {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"enable_certificates": mc.EnableCertificates,
		},
	}
}

// flattenResourceUsageExportConfig converts ResourceUsageExportConfig to a map
func flattenResourceUsageExportConfig(ruec *container.ResourceUsageExportConfig) []interface{} {
	if ruec == nil {
		return nil
	}

	result := make(map[string]interface{})
	hasNonDefaultConfig := false

	if ruec.ConsumptionMeteringConfig != nil && ruec.ConsumptionMeteringConfig.Enabled {
		result["enable_network_egress_metering"] = ruec.EnableNetworkEgressMetering
		hasNonDefaultConfig = true
	}

	if ruec.EnableResourceConsumptionMetering != defaultEnableResourceConsumptionMetering {
		result["enable_resource_consumption_metering"] = ruec.EnableResourceConsumptionMetering
		hasNonDefaultConfig = true
	}

	if ruec.BigqueryDestination != nil && ruec.BigqueryDestination.DatasetId != "" {
		result["bigquery_destination"] = []interface{}{
			map[string]interface{}{
				"dataset_id": ruec.BigqueryDestination.DatasetId,
			},
		}
		hasNonDefaultConfig = true
	}

	if !hasNonDefaultConfig {
		return nil
	}
	return []interface{}{result}
}

// flattenSecretManagerConfig converts SecretManagerConfig to a map
func flattenSecretManagerConfig(smc *container.SecretManagerConfig) []interface{} {
	if smc == nil || !smc.Enabled {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"enabled": smc.Enabled,
		},
	}
}

// flattenServiceExternalIPsConfig converts ServiceExternalIPsConfig to a map
func flattenServiceExternalIPsConfig(seic *container.ServiceExternalIPsConfig) []interface{} {
	if seic == nil || !seic.Enabled {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"enabled": seic.Enabled,
		},
	}
}

// flattenWorkloadIdentityConfig converts WorkloadIdentityConfig to a map
func flattenWorkloadIdentityConfig(wic *container.WorkloadIdentityConfig) []interface{} {
	if wic == nil || wic.WorkloadPool == "" {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"workload_pool": wic.WorkloadPool,
		},
	}
}

// flattenGatewayApiConfig converts GatewayAPIConfig to a map
func flattenGatewayApiConfig(gac *container.GatewayAPIConfig) []interface{} {
	if gac == nil || gac.Channel == defaultGatewayApiChannel {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"channel": gac.Channel,
		},
	}
}

// flattenFleet converts Fleet to a map
func flattenFleet(fleet *container.Fleet) []interface{} {
	if fleet == nil {
		return nil
	}

	result := make(map[string]interface{})

	if fleet.Project != "" {
		result["project"] = fleet.Project
	}

	if len(result) == 0 {
		return nil
	}
	return []interface{}{result}
}

// flattenNodePoolDefaults converts NodePoolDefaults to a map
func flattenNodePoolDefaults(defaults *container.NodePoolDefaults) []interface{} {
	if defaults == nil {
		return nil
	}

	result := make(map[string]interface{})

	if defaults.NodeConfigDefaults != nil {
		nodeConfigDefaults := make(map[string]interface{})
		
		if defaults.NodeConfigDefaults.Gcfs != nil {
			nodeConfigDefaults["gcfs_config"] = []interface{}{
				map[string]interface{}{
					"enabled": defaults.NodeConfigDefaults.Gcfs.Enabled,
				},
			}
		}
		
		if len(nodeConfigDefaults) > 0 {
			result["node_config_defaults"] = []interface{}{nodeConfigDefaults}
		}
	}

	if len(result) == 0 {
		return nil
	}
	return []interface{}{result}
}