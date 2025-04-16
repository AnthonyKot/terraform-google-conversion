package container

import (
	"fmt"
	// "time" // Uncomment if implementing DurationToSeconds
	// "strings" // Uncomment if implementing DurationToSeconds

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/cai2hcl/converters/utils"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/cai2hcl/models"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/caiasset"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/tpgresource"
	"google.golang.org/api/container/v1"
)

const ContainerClusterAssetType string = "container.googleapis.com/Cluster"
const ContainerClusterSchemaName string = "google_container_cluster"
const ContainerNodePoolSchemaName string = "google_container_node_pool"

type ContainerClusterConverter struct {
	clusterName    string
	clusterSchema  map[string]*schema.Schema
	nodePoolSchema map[string]*schema.Schema
}

func NewContainerClusterConverter(provider *schema.Provider) models.Converter {
	clusterSchema, clusterOk := provider.ResourcesMap[ContainerClusterSchemaName]
	nodePoolSchema, nodePoolOk := provider.ResourcesMap[ContainerNodePoolSchemaName]

	if !clusterOk {
		fmt.Printf("Warning: Schema for %s not found in provider\n", ContainerClusterSchemaName)
	}
	if !nodePoolOk {
		fmt.Printf("Warning: Schema for %s not found in provider\n", ContainerNodePoolSchemaName)
	}

	var cs map[string]*schema.Schema
	if clusterOk && clusterSchema != nil {
		cs = clusterSchema.Schema
	}

	var nps map[string]*schema.Schema
	if nodePoolOk && nodePoolSchema != nil {
		nps = nodePoolSchema.Schema
	}

	return &ContainerClusterConverter{
		clusterName:    ContainerClusterSchemaName,
		clusterSchema:  cs,
		nodePoolSchema: nps,
	}
}

func (c *ContainerClusterConverter) Convert(asset *caiasset.Asset) ([]*models.TerraformResourceBlock, error) {
	if asset == nil || asset.Resource == nil || asset.Resource.Data == nil {
		return nil, fmt.Errorf("asset or asset resource data is nil")
	}

	project := utils.ParseFieldValue(asset.Name, "projects")
	location := utils.ParseFieldValue(asset.Name, "locations")
	clusterName := utils.ParseFieldValue(asset.Name, "clusters")

	var cluster *container.Cluster
	if err := utils.DecodeJSON(asset.Resource.Data, &cluster); err != nil {
		return nil, fmt.Errorf("failed to decode cluster JSON: %w", err)
	}

	clusterBlock, err := c.convertClusterData(cluster, project, location, clusterName)
	if err != nil {
		return nil, fmt.Errorf("failed to convert cluster data: %w", err)
	}
	if clusterBlock == nil && err == nil {
		return nil, fmt.Errorf("cluster conversion returned nil block without error")
	}

	blocks := []*models.TerraformResourceBlock{clusterBlock}

	if len(cluster.NodePools) > 0 {
		for _, nodePool := range cluster.NodePools {
			if nodePool == nil {
				continue
			}
			nodePoolBlock, err := c.convertNodePoolData(nodePool, cluster, project, location)
			if err != nil {
				fmt.Printf("Warning: Failed to convert node pool %s: %v. Skipping.\n", nodePool.Name, err)
				continue
			}
			if nodePoolBlock != nil {
				blocks = append(blocks, nodePoolBlock)
			}
		}
	}

	return blocks, nil
}

func (c *ContainerClusterConverter) convertClusterData(cluster *container.Cluster, project, location, clusterName string) (*models.TerraformResourceBlock, error) {
	if cluster == nil {
		return nil, fmt.Errorf("cluster data is nil")
	}
	if c.clusterSchema == nil {
		return nil, fmt.Errorf("cluster schema is nil in converter")
	}

	hclData := make(map[string]interface{})

	hclData["name"] = clusterName
	hclData["location"] = location
	if project != "" {
		hclData["project"] = project
	}

	if cluster.Description != "" {
		hclData["description"] = cluster.Description
	}

	if len(cluster.NodePools) > 0 {
		hclData["remove_default_node_pool"] = true
	} else {
		if cluster.InitialNodeCount > 0 {
			hclData["initial_node_count"] = cluster.InitialNodeCount
		}
		if cluster.NodeConfig != nil {
			hclData["node_config"] = flattenNodeConfig(cluster.NodeConfig)
		}
	}

	if cluster.Network != "" {
		hclData["network"] = tpgresource.GetResourceNameFromSelfLink(cluster.Network)
	}
	if cluster.Subnetwork != "" {
		hclData["subnetwork"] = tpgresource.GetResourceNameFromSelfLink(cluster.Subnetwork)
	}

	if cluster.IpAllocationPolicy != nil {
		hclData["networking_mode"] = "VPC_NATIVE"
		hclData["ip_allocation_policy"] = flattenIPAllocationPolicy(cluster.IpAllocationPolicy)
	}

	if cluster.NetworkConfig != nil {
		if cluster.NetworkConfig.DatapathProvider != "" && cluster.NetworkConfig.DatapathProvider != "DATAPATH_PROVIDER_UNSPECIFIED" {
			hclData["datapath_provider"] = cluster.NetworkConfig.DatapathProvider
		}
		// Commenting out potentially missing field:
		// if cluster.NetworkConfig.DisableL4IlbFwReconciliation {
		// 	hclData["disable_l4_lb_firewall_reconciliation"] = true
		// }
		if cluster.NetworkConfig.EnableFqdnNetworkPolicy {
			hclData["enable_fqdn_network_policy"] = true
		}
		if cluster.NetworkConfig.EnableL4ilbSubsetting {
			hclData["enable_l4_ilb_subsetting"] = true
		}
		if cluster.NetworkConfig.EnableMultiNetworking {
			hclData["enable_multi_networking"] = true
		}
	}

	if cluster.DefaultMaxPodsConstraint != nil {
		hclData["default_max_pods_per_node"] = cluster.DefaultMaxPodsConstraint.MaxPodsPerNode
	}

	if cluster.NetworkPolicy != nil {
		flattenedPolicy := flattenNetworkPolicy(cluster.NetworkPolicy)
		if flattenedPolicy != nil {
			hclData["network_policy"] = flattenedPolicy
		}
	}

	if cluster.LoggingService != "" && cluster.LoggingService != "logging.googleapis.com" {
		hclData["logging_service"] = cluster.LoggingService
	}
	if cluster.MonitoringService != "" && cluster.MonitoringService != "monitoring.googleapis.com" {
		hclData["monitoring_service"] = cluster.MonitoringService
	}

	if cluster.AddonsConfig != nil {
		hclData["addons_config"] = flattenAddonsConfig(cluster.AddonsConfig)
	}

	if len(cluster.Locations) > 0 {
		hclData["node_locations"] = cluster.Locations
	}

	if cluster.ResourceLabels != nil {
		hclData["resource_labels"] = cluster.ResourceLabels
	}

	if cluster.ReleaseChannel == nil || cluster.ReleaseChannel.Channel == "" || cluster.ReleaseChannel.Channel == "UNSPECIFIED" {
		if cluster.CurrentMasterVersion != "" {
			hclData["master_version"] = cluster.CurrentMasterVersion
		}
	}

	if cluster.CurrentNodeVersion != "" {
		hclData["node_version"] = cluster.CurrentNodeVersion
	}

	// Commenting out potentially missing field:
	// if cluster.DeletionProtection {
	// 	hclData["deletion_protection"] = true
	// }

	if cluster.Autopilot != nil && cluster.Autopilot.Enabled {
		hclData["enable_autopilot"] = true
	}
	if cluster.EnableKubernetesAlpha {
		hclData["enable_kubernetes_alpha"] = true
	}
	if cluster.EnableTpu {
		hclData["enable_tpu"] = true
	}
	// Commenting out potentially missing field:
	// if cluster.IntraNodeVisibilityConfig != nil && cluster.IntraNodeVisibilityConfig.Enabled {
	//     hclData["enable_intranode_visibility"] = true
	// }
	if cluster.LegacyAbac != nil && cluster.LegacyAbac.Enabled {
		hclData["enable_legacy_abac"] = true
	}
	if cluster.ShieldedNodes != nil && cluster.ShieldedNodes.Enabled {
		hclData["enable_shielded_nodes"] = true
	}

	// Commenting out potentially missing field:
	// if cluster.PrivateIpv6GoogleAccess != "" && cluster.PrivateIpv6GoogleAccess != "PRIVATE_IPV6_GOOGLE_ACCESS_UNSPECIFIED" {
	//     hclData["private_ipv6_google_access"] = cluster.PrivateIpv6GoogleAccess
	// }

	ctyVal, err := utils.MapToCtyValWithSchema(hclData, c.clusterSchema)
	if err != nil {
		return nil, fmt.Errorf("error converting cluster data to cty.Value: %w", err)
	}

	return &models.TerraformResourceBlock{
		Labels: []string{c.clusterName, clusterName},
		Value:  ctyVal,
	}, nil
}

func (c *ContainerClusterConverter) convertNodePoolData(nodePool *container.NodePool, cluster *container.Cluster, project, location string) (*models.TerraformResourceBlock, error) {
	if nodePool == nil {
		return nil, fmt.Errorf("node pool data is nil")
	}
	if cluster == nil {
		return nil, fmt.Errorf("cluster context is nil for node pool")
	}
	if c.nodePoolSchema == nil {
		return nil, fmt.Errorf("node pool schema is nil in converter")
	}

	hclData := make(map[string]interface{})

	hclData["name"] = nodePool.Name
	hclData["cluster"] = cluster.Name
	if location != "" {
		hclData["location"] = location
	}
	if project != "" {
		hclData["project"] = project
	}

	if nodePool.Autoscaling != nil && nodePool.Autoscaling.Enabled {
		hclData["autoscaling"] = flattenNodePoolAutoscaling(nodePool.Autoscaling)
		if nodePool.InitialNodeCount > 0 {
			hclData["initial_node_count"] = nodePool.InitialNodeCount
		}
	} else {
		if nodePool.InitialNodeCount > 0 {
			hclData["node_count"] = nodePool.InitialNodeCount
		}
	}

	if nodePool.Config != nil {
		hclData["node_config"] = flattenNodeConfig(nodePool.Config)
	}

	if nodePool.Management != nil {
		hclData["management"] = flattenNodeManagement(nodePool.Management)
	}

	if nodePool.MaxPodsConstraint != nil {
		hclData["max_pods_constraint"] = flattenMaxPodsConstraint(nodePool.MaxPodsConstraint)
	}

	if nodePool.NetworkConfig != nil {
		hclData["network_config"] = flattenNodeNetworkConfig(nodePool.NetworkConfig)
	}

	if nodePool.UpgradeSettings != nil {
		hclData["upgrade_settings"] = flattenNodePoolUpgradeSettings(nodePool.UpgradeSettings)
	}

	if nodePool.Version != "" {
		hclData["version"] = nodePool.Version
	}

	if nodePool.PlacementPolicy != nil {
		hclData["placement_policy"] = flattenPlacementPolicy(nodePool.PlacementPolicy)
	}

	ctyVal, err := utils.MapToCtyValWithSchema(hclData, c.nodePoolSchema)
	if err != nil {
		return nil, fmt.Errorf("error converting node pool %s data to cty.Value: %w", nodePool.Name, err)
	}

	return &models.TerraformResourceBlock{
		Labels: []string{ContainerNodePoolSchemaName, nodePool.Name},
		Value:  ctyVal,
	}, nil
}

// --- ALL FLATTEN FUNCTIONS (Multiline Format) ---

func flattenNodeConfig(config *container.NodeConfig) []interface{} {
	if config == nil {
		return nil
	}
	nodeConfig := make(map[string]interface{})
	if config.MachineType != "" {
		nodeConfig["machine_type"] = config.MachineType
	}
	if config.DiskSizeGb > 0 {
		nodeConfig["disk_size_gb"] = config.DiskSizeGb
	}
	if config.DiskType != "" {
		nodeConfig["disk_type"] = config.DiskType
	}
	if len(config.OauthScopes) > 0 {
		nodeConfig["oauth_scopes"] = config.OauthScopes
	}
	if config.ServiceAccount != "" {
		nodeConfig["service_account"] = config.ServiceAccount
	}
	if len(config.Metadata) > 0 {
		nodeConfig["metadata"] = config.Metadata
	}
	if config.ImageType != "" {
		nodeConfig["image_type"] = config.ImageType
	}
	if len(config.Labels) > 0 {
		nodeConfig["labels"] = config.Labels
	}
	if config.LocalSsdCount != 0 {
		nodeConfig["local_ssd_count"] = config.LocalSsdCount
	}
	if len(config.Tags) > 0 {
		nodeConfig["tags"] = config.Tags
	}
	if config.Preemptible {
		nodeConfig["preemptible"] = config.Preemptible
	}
	if config.Spot {
		nodeConfig["spot"] = config.Spot
	}
	if config.MinCpuPlatform != "" {
		nodeConfig["min_cpu_platform"] = config.MinCpuPlatform
	}
	if config.WorkloadMetadataConfig != nil {
		nodeConfig["workload_metadata_config"] = flattenWorkloadMetadataConfig(config.WorkloadMetadataConfig)
	}
	if config.ShieldedInstanceConfig != nil {
		nodeConfig["shielded_instance_config"] = flattenShieldedInstanceConfig(config.ShieldedInstanceConfig)
	}
	if config.BootDiskKmsKey != "" {
		nodeConfig["boot_disk_kms_key"] = config.BootDiskKmsKey
	}
	if len(config.Accelerators) > 0 {
		nodeConfig["guest_accelerator"] = flattenAccelerators(config.Accelerators)
	}
	if config.ReservationAffinity != nil {
		nodeConfig["reservation_affinity"] = flattenReservationAffinity(config.ReservationAffinity)
	}
	if config.ConfidentialNodes != nil {
		nodeConfig["confidential_nodes"] = flattenConfidentialNodes(config.ConfidentialNodes)
	}
	if config.KubeletConfig != nil {
		nodeConfig["kubelet_config"] = flattenKubeletConfig(config.KubeletConfig)
	}
	if config.LinuxNodeConfig != nil {
		nodeConfig["linux_node_config"] = flattenLinuxNodeConfig(config.LinuxNodeConfig)
	}
	if len(config.Taints) > 0 {
		nodeConfig["taint"] = flattenNodeTaints(config.Taints)
	}
	return []interface{}{nodeConfig}
}

func flattenWorkloadMetadataConfig(config *container.WorkloadMetadataConfig) []interface{} {
	if config == nil {
		return nil
	}
	data := make(map[string]interface{})
	data["mode"] = config.Mode
	return []interface{}{data}
}

func flattenShieldedInstanceConfig(config *container.ShieldedInstanceConfig) []interface{} {
	if config == nil {
		return nil
	}
	data := make(map[string]interface{})
	data["enable_secure_boot"] = config.EnableSecureBoot
	data["enable_integrity_monitoring"] = config.EnableIntegrityMonitoring
	return []interface{}{data}
}

func flattenAccelerators(accelerators []*container.AcceleratorConfig) []interface{} {
	if len(accelerators) == 0 {
		return nil
	}
	result := make([]interface{}, len(accelerators))
	for i, acc := range accelerators {
		data := make(map[string]interface{})
		data["accelerator_type"] = acc.AcceleratorType
		data["accelerator_count"] = acc.AcceleratorCount
		data["gpu_partition_size"] = acc.GpuPartitionSize // Add only if not empty? Check schema defaults
		result[i] = data
	}
	return result
}

func flattenReservationAffinity(config *container.ReservationAffinity) []interface{} {
	if config == nil {
		return nil
	}
	ra := make(map[string]interface{})
	ra["consume_reservation_type"] = config.ConsumeReservationType
	if len(config.Key) > 0 && len(config.Values) > 0 {
		ra["key"] = config.Key
		ra["values"] = config.Values
	}
	return []interface{}{ra}
}

func flattenConfidentialNodes(config *container.ConfidentialNodes) []interface{} {
	if config == nil {
		return nil
	}
	data := make(map[string]interface{})
	data["enabled"] = config.Enabled
	return []interface{}{data}
}

func flattenKubeletConfig(config *container.NodeKubeletConfig) []interface{} {
	if config == nil {
		return nil
	}
	kc := make(map[string]interface{})
	if config.CpuManagerPolicy != "" {
		kc["cpu_manager_policy"] = config.CpuManagerPolicy
	}
	// Explicitly set bools as TF schema might require them
	kc["cpu_cfs_quota"] = config.CpuCfsQuota
	if config.CpuCfsQuotaPeriod != "" {
		kc["cpu_cfs_quota_period"] = config.CpuCfsQuotaPeriod
	}
	if config.PodPidsLimit != 0 {
		kc["pod_pids_limit"] = config.PodPidsLimit
	}
	return []interface{}{kc}
}

func flattenLinuxNodeConfig(config *container.LinuxNodeConfig) []interface{} {
	if config == nil {
		return nil
	}
	lnc := make(map[string]interface{})
	if config.Sysctls != nil {
		lnc["sysctls"] = config.Sysctls
	}
	return []interface{}{lnc}
}

func flattenNodeTaints(taints []*container.NodeTaint) []interface{} {
	if len(taints) == 0 {
		return nil
	}
	result := make([]interface{}, len(taints))
	for i, taint := range taints {
		t := make(map[string]interface{})
		t["key"] = taint.Key
		t["value"] = taint.Value
		t["effect"] = taint.Effect
		result[i] = t
	}
	return result
}

func flattenIPAllocationPolicy(policy *container.IPAllocationPolicy) []interface{} {
	if policy == nil {
		return nil
	}
	ipa := make(map[string]interface{})
	ipa["use_ip_aliases"] = policy.UseIpAliases // This is often true if the block exists
	if policy.ClusterSecondaryRangeName != "" {
		ipa["cluster_secondary_range_name"] = policy.ClusterSecondaryRangeName
	}
	if policy.ServicesSecondaryRangeName != "" {
		ipa["services_secondary_range_name"] = policy.ServicesSecondaryRangeName
	}
	if policy.ClusterIpv4CidrBlock != "" {
		ipa["cluster_ipv4_cidr_block"] = policy.ClusterIpv4CidrBlock
	}
	if policy.ServicesIpv4CidrBlock != "" {
		ipa["services_ipv4_cidr_block"] = policy.ServicesIpv4CidrBlock
	}
	if policy.NodeIpv4CidrBlock != "" {
		ipa["node_ipv4_cidr_block"] = policy.NodeIpv4CidrBlock
	}
	if policy.CreateSubnetwork {
		ipa["create_subnetwork"] = policy.CreateSubnetwork
	}
	if policy.SubnetworkName != "" {
		ipa["subnetwork_name"] = policy.SubnetworkName
	}
	if policy.TpuIpv4CidrBlock != "" {
		ipa["tpu_ipv4_cidr_block"] = policy.TpuIpv4CidrBlock
	}
	if policy.StackType != "" {
		ipa["stack_type"] = policy.StackType
	}
	if policy.Ipv6AccessType != "" {
		ipa["ipv6_access_type"] = policy.Ipv6AccessType
	}
	return []interface{}{ipa}
}

func flattenAddonsConfig(config *container.AddonsConfig) []interface{} {
	if config == nil {
		return nil
	}
	addons := make(map[string]interface{})

	if config.HttpLoadBalancing != nil {
		httpLb := make(map[string]interface{})
		httpLb["disabled"] = config.HttpLoadBalancing.Disabled
		addons["http_load_balancing"] = []interface{}{httpLb}
	}
	if config.HorizontalPodAutoscaling != nil {
		hpa := make(map[string]interface{})
		hpa["disabled"] = config.HorizontalPodAutoscaling.Disabled
		addons["horizontal_pod_autoscaling"] = []interface{}{hpa}
	}
	if config.NetworkPolicyConfig != nil {
		npc := make(map[string]interface{})
		npc["disabled"] = config.NetworkPolicyConfig.Disabled
		addons["network_policy_config"] = []interface{}{npc}
	}
	if config.CloudRunConfig != nil {
		crc := make(map[string]interface{})
		crc["disabled"] = config.CloudRunConfig.Disabled
		if config.CloudRunConfig.LoadBalancerType != "" {
			crc["load_balancer_type"] = config.CloudRunConfig.LoadBalancerType
		}
		addons["cloudrun_config"] = []interface{}{crc}
	}
	if config.DnsCacheConfig != nil {
		dcc := make(map[string]interface{})
		dcc["enabled"] = config.DnsCacheConfig.Enabled
		addons["dns_cache_config"] = []interface{}{dcc}
	}
	if config.GcePersistentDiskCsiDriverConfig != nil {
		gpd := make(map[string]interface{})
		gpd["enabled"] = config.GcePersistentDiskCsiDriverConfig.Enabled
		addons["gce_persistent_disk_csi_driver_config"] = []interface{}{gpd}
	}
	if config.GcpFilestoreCsiDriverConfig != nil {
		gfs := make(map[string]interface{})
		gfs["enabled"] = config.GcpFilestoreCsiDriverConfig.Enabled
		addons["gcp_filestore_csi_driver_config"] = []interface{}{gfs}
	}
	if config.ConfigConnectorConfig != nil {
		ccc := make(map[string]interface{})
		ccc["enabled"] = config.ConfigConnectorConfig.Enabled
		addons["config_connector_config"] = []interface{}{ccc}
	}
	if config.GkeBackupAgentConfig != nil {
		gba := make(map[string]interface{})
		gba["enabled"] = config.GkeBackupAgentConfig.Enabled
		addons["gke_backup_agent_config"] = []interface{}{gba}
	}
	if config.GcsFuseCsiDriverConfig != nil {
		gcsFuseConfig := make(map[string]interface{})
		gcsFuseConfig["enabled"] = config.GcsFuseCsiDriverConfig.Enabled
		addons["gcs_fuse_csi_driver_config"] = []interface{}{gcsFuseConfig}
	}
	if config.RayOperatorConfig != nil {
		rayConfig := make(map[string]interface{})
		rayConfig["enabled"] = config.RayOperatorConfig.Enabled
		addons["ray_operator_config"] = []interface{}{rayConfig}
	}

	if len(addons) == 0 {
		return nil
	}
	return []interface{}{addons}
}

func flattenNetworkPolicy(policy *container.NetworkPolicy) []interface{} {
	if policy == nil {
		return nil
	}

	hasProvider := policy.Provider != "" && policy.Provider != "PROVIDER_UNSPECIFIED"
	// Commenting out potentially missing fields:
	// hasAllowNetAdmin := policy.AllowNetAdmin
	// hasEnableCilium := policy.EnableCiliumClusterwideNetworkPolicy

	// Adjust logic based only on provider presence for now
	if !hasProvider { // && !hasAllowNetAdmin && !hasEnableCilium {
		return nil
	}

	data := make(map[string]interface{})

	if hasProvider {
		data["provider"] = policy.Provider
	}
	// Commenting out potentially missing fields:
	// else {
	//     if hasAllowNetAdmin || hasEnableCilium {
	//         fmt.Println("Warning: NetworkPolicy boolean flags set but no Provider specified in API data. Omitting network_policy block.")
	//         return nil
	//     }
	// }

	// Commenting out potentially missing fields:
	// if hasProvider {
	//     data["allow_net_admin"] = policy.AllowNetAdmin
	// }

	// Commenting out potentially missing fields:
	// if policy.EnableCiliumClusterwideNetworkPolicy {
	//     data["enable_cilium_clusterwide_network_policy"] = true
	// }

	if len(data) == 0 {
		return nil
	}
	return []interface{}{data}
}

func flattenNodePoolAutoscaling(autoscaling *container.NodePoolAutoscaling) []interface{} {
	if autoscaling == nil {
		return nil
	}
	data := make(map[string]interface{})
	// These are required in TF block if autoscaling is enabled
	data["min_node_count"] = autoscaling.MinNodeCount
	data["max_node_count"] = autoscaling.MaxNodeCount

	// These are optional
	if autoscaling.TotalMinNodeCount > 0 {
		data["total_min_node_count"] = autoscaling.TotalMinNodeCount
	}
	if autoscaling.TotalMaxNodeCount > 0 {
		data["total_max_node_count"] = autoscaling.TotalMaxNodeCount
	}
	if autoscaling.LocationPolicy != "" {
		data["location_policy"] = autoscaling.LocationPolicy
	}
	return []interface{}{data}
}

func flattenNodeManagement(mgmt *container.NodeManagement) []interface{} {
	if mgmt == nil {
		return nil
	}
	data := make(map[string]interface{})
	// TF requires these bools, map them directly
	data["auto_repair"] = mgmt.AutoRepair
	data["auto_upgrade"] = mgmt.AutoUpgrade
	return []interface{}{data}
}

func flattenNodePoolUpgradeSettings(settings *container.UpgradeSettings) []interface{} {
	if settings == nil {
		return nil
	}
	data := make(map[string]interface{})
	// TF requires these integers if block is present
	data["max_surge"] = settings.MaxSurge
	data["max_unavailable"] = settings.MaxUnavailable

	if settings.Strategy != "" {
		data["strategy"] = settings.Strategy
	}

	if settings.BlueGreenSettings != nil {
		data["blue_green_settings"] = flattenBlueGreenSettings(settings.BlueGreenSettings)
	}

	return []interface{}{data}
}

func flattenBlueGreenSettings(settings *container.BlueGreenSettings) []interface{} {
	if settings == nil || settings.StandardRolloutPolicy == nil {
		return nil
	}
	data := make(map[string]interface{})

	standardRolloutPolicy := make(map[string]interface{})
	standardRolloutPolicy["batch_percentage"] = settings.StandardRolloutPolicy.BatchPercentage
	standardRolloutPolicy["batch_node_count"] = settings.StandardRolloutPolicy.BatchNodeCount
	// Add duration parsing here if needed for settings.StandardRolloutPolicy.BatchSoakDuration -> batch_soak_duration_sec

	data["standard_rollout_policy"] = []interface{}{standardRolloutPolicy}

	// Add duration parsing here if needed for settings.NodePoolSoakDuration -> node_pool_soak_duration_sec

	return []interface{}{data}
}

func flattenMaxPodsConstraint(constraint *container.MaxPodsConstraint) []interface{} {
	if constraint == nil {
		return nil
	}
	data := make(map[string]interface{})
	data["max_pods_per_node"] = constraint.MaxPodsPerNode
	return []interface{}{data}
}

func flattenNodeNetworkConfig(config *container.NodeNetworkConfig) []interface{} {
	if config == nil {
		return nil
	}
	data := make(map[string]interface{})
	if config.PodRange != "" {
		data["pod_range"] = config.PodRange
	}
	if config.PodIpv4CidrBlock != "" {
		data["pod_ipv4_cidr_block"] = config.PodIpv4CidrBlock
	}
	data["create_pod_range"] = config.CreatePodRange // Map bool directly

	if config.NetworkPerformanceConfig != nil && config.NetworkPerformanceConfig.TotalEgressBandwidthTier != "" {
		perfConfig := map[string]interface{}{
			"total_egress_bandwidth_tier": config.NetworkPerformanceConfig.TotalEgressBandwidthTier,
		}
		data["network_performance_config"] = []interface{}{perfConfig}
	}

	return []interface{}{data}
}

func flattenPlacementPolicy(policy *container.PlacementPolicy) []interface{} {
	if policy == nil {
		return nil
	}
	data := make(map[string]interface{})
	if policy.Type != "" {
		data["type"] = policy.Type
	}
	// Add other fields like tpu_topology if needed
	return []interface{}{data}
}