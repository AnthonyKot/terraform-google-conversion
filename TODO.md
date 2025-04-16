# GKE Cluster Converter - Past TODOs and Discrepancies

This document consolidates previously identified areas for improvement, discrepancies found during testing, and specific field requirements for the GKE cluster converter.

---

## General Improvements / Missing Features (Initial List)

### `google_container_cluster` Resource Improvements:

* **Add Missing Blocks:** The converter was missing several configuration blocks entirely. Adding these would significantly improve coverage:
    * `master_auth`: Crucial for authentication configuration. Needs a `flattenMasterAuth` function.
    * `private_cluster_config`: Essential for configuring private clusters. Needs `flattenPrivateClusterConfig`.
    * `release_channel`: Very common way to manage versions. Needs `flattenReleaseChannel`.
    * `logging_config` / `monitoring_config`: Provide finer-grained control than just `logging_service` / `monitoring_service`. Need `flattenLoggingConfig` and `flattenMonitoringConfig`.
    * `cluster_autoscaling`: Configuration for node pool auto-provisioning. Needs `flattenClusterAutoscaling`.
    * `database_encryption`: For enabling CMEK for secrets. Needs `flattenDatabaseEncryption`.
    * `vertical_pod_autoscaling`: Although not in the specific TF output example, this is another common addon.
    * Others from TF output like `binary_authorization`, `default_snat_status`, `notification_config`, etc., could be added based on priority.
* **Fix/Uncomment Missing Top-Level Fields:** Revisit fields that caused build errors. Find the correct field names/locations in the `container/v1` library version:
    * `deletion_protection` (API: `DeletionProtection bool`?)
    * `enable_intranode_visibility` (API: `IntraNodeVisibilityConfig.Enabled bool`?)
    * `private_ipv6_google_access` (API: `PrivateIpv6GoogleAccess string` enum?)
* **Fix/Uncomment `network_policy` Fields:** Revisit commented-out fields within `flattenNetworkPolicy`. Find correct API field names for:
    * `allow_net_admin` (API: `AllowNetAdmin bool`?)
    * `enable_cilium_clusterwide_network_policy` (API: `EnableCiliumClusterwideNetworkPolicy bool`?)
    * Restructure `flattenNetworkPolicy` if needed to better match TF schema (`enabled` and `provider`).
* **Refine `ip_allocation_policy`:** Add the nested `pod_cidr_overprovision_config` block within `flattenIPAllocationPolicy`.
* **Handle `initial_node_count` Correctly:** When `remove_default_node_pool` is true, Terraform often requires `initial_node_count` to be set (e.g., to 1). The converter currently omits it. It should probably default to setting `initial_node_count = 1` in this scenario. (Note: Later analysis suggested omitting cluster-level versioning/counts if `remove_default_node_pool=true`).
* **Refine `network`/`subnetwork` Output:** Output full self-link or path instead of just the name ("default").
* **Verify `node_locations` Mapping:** Double-check if `cluster.Locations` is the correct API source field for cluster-level `node_locations`.

### `google_container_node_pool` Resource Improvements:

* **Add Missing Top-Level Fields:**
    * `max_pods_per_node`: Map from `nodePool.MaxPodsConstraint.MaxPodsPerNode`.
    * `node_locations`: Map from `nodePool.Locations`.
* **Add Missing Blocks:**
    * `placement_policy`: Ensure `flattenPlacementPolicy` is complete and called.
    * `queued_provisioning`: Add `flattenQueuedProvisioning` if needed.
* **Expand `node_config` Flattening (`flattenNodeConfig`):**
    * Add mappings for simple fields: `labels`, `resource_labels`, `tags`, `logging_variant`, `boot_disk_kms_key`.
    * Consider explicitly setting booleans like `preemptible` and `spot` to `false` to match state, or stick to omitting defaults.
    * Expand `flattenKubeletConfig` for more detail (log/GC settings, etc.).
    * Add flatteners for `advanced_machine_features` and `windows_node_config`.
* **Expand `network_config` Flattening (`flattenNodeNetworkConfig`):**
    * Add mapping for `enable_private_nodes`.

---

## Compiler Error Fixes and Refactoring Notes (from File 2 Analysis)

* **Fix Errors:** Address fields causing compiler errors (`LoggingVariant`, `CpuCfsQuota`, `AutoscalingProfile` (x2), `MaxSurge`, `MaxUnavailable`, `EnableCiliumClusterwideNetworkPolicy`) by finding the correct field/type or removing access. Use provider code as a strong hint for correct types (`bool`, `int64`).
* **Refactor `flattenPrivateClusterConfig`:** Rewrite based on provider's approach using CPEC, PCC, NC inputs.
* **Update `flattenNetworkPolicy`:** Remove CALICO default, remove Cilium field access. Keep Enabled check and nil return.
* **Update `flattenNodePoolUpgradeSettings`:** Change `MaxSurge`/`Unavailable` handling to value types (`int64`).
* **Review `flattenBlueGreenSettings`:** Simplify standard policy logic? Implement duration parsing?
* **Update `flattenReleaseChannel`:** Consider outputting `UNSPECIFIED` explicitly (Note: Later decided against this to omit defaults).
* **Check Addons/Monitoring:** Add missing `ParallelstoreCsiDriverConfig`, `auto_monitoring_config`. Verify `RayOperatorConfig` detail & `AdvancedDatapathObservabilityConfig` fields.
* **Review `flattenIPAllocationPolicy`:** Simplify CIDR logic? Align `StackType` default?
* **(Low Priority/Consistency):** Align return types (`[]map` vs `[]interface`)? Standardize nil/default block handling?

---

## Discrepancies Found vs. `zonal-clustet.tf` Example (Based on File 2/3 Analysis)

### 1. Missing Top-Level Arguments in Converted Output
The following fields were present in the reference `zonal-clustet.tf` but missing in the converter's output for that specific example:

* `allow_net_admin`
* `datapath_provider` (Note: Was default, removed during cleanup)
* `default_max_pods_per_node`
* `deletion_protection` (Note: Was default, removed during cleanup)
* `description`
* `disable_l4_lb_firewall_reconciliation` (Note: Was default, removed during cleanup)
* `enable_autopilot` (Note: Was default, removed during cleanup)
* `enable_cilium_clusterwide_network_policy` (Note: Was default, removed during cleanup)
* `enable_fqdn_network_policy` (Note: Was default, removed during cleanup)
* `enable_intranode_visibility` (Note: Was default, removed during cleanup)
* `enable_kubernetes_alpha` (Note: Was default, removed during cleanup)
* `enable_l4_ilb_subsetting` (Note: Was default, removed during cleanup)
* `enable_legacy_abac` (Note: Was default, removed during cleanup)
* `enable_multi_networking` (Note: Was default, removed during cleanup)
* `enable_shielded_nodes` (Note: Was default, removed during cleanup)
* `enable_tpu` (Note: Was default, removed during cleanup)
* `initial_node_count` (Note: Should be removed if `remove_default_node_pool=true`)
* `location` (**CRITICAL MISSING FIELD**)
* `logging_service` (Note: Was default, removed during cleanup)
* `min_master_version`
* `monitoring_service` (Note: Was default, removed during cleanup)
* `networking_mode` (Note: Was default/implied, removed during cleanup)
* `node_locations` (Note: Was default/empty, removed during cleanup)
* `node_version`
* `private_ipv6_google_access`
* `remove_default_node_pool`
* `resource_labels` (Note: Was default/empty, removed during cleanup)

### 2. Block/Nested Structure Differences in Converted Output

* **Missing Blocks:** `authenticator_groups_config`, `binary_authorization`, `cluster_autoscaling`, `control_plane_endpoints_config`, `database_encryption`, `default_snat_status`, `enterprise_config`, `fleet`, `logging_config`, `master_auth`, `monitoring_config`, `network_policy`, `node_pool` (inline version), `node_pool_auto_config`, `node_pool_defaults`, `notification_config`, `private_cluster_config`, `release_channel`, `secret_manager_config`, `security_posture_config`, `service_external_ips_config`.
* **Partially Missing Blocks:** `addons_config` was missing `gcs_fuse_csi_driver_config` and `ray_operator_config`. (Note: Many blocks were correctly removed during cleanup as they only contained defaults).

### 3. Field Value Differences

* Network/Subnetwork paths used short names ("default") instead of full resource paths.
* Some blocks had fewer fields due to default value omission (which is desired, but need verification).

---

## Strict Field List for 100% `zonal-clustet.tf` Accuracy (from CAI)

*To achieve 100% accurate conversion **for the specific `zonal-clustet.tf` example provided earlier**, the converter must handle the extraction and mapping of all fields listed below from the corresponding CAI asset.*

*(Note: Many of these were likely set to default values in the original `zonal-clustet.tf` export and were removed in the cleaned-up version. This list represents the potential fields present in the CAI export.)*

### Top-level Fields
- allow_net_admin
- cluster_ipv4_cidr
- datapath_provider
- default_max_pods_per_node
- deletion_protection
- description
- disable_l4_lb_firewall_reconciliation
- enable_autopilot
- enable_cilium_clusterwide_network_policy
- enable_fqdn_network_policy
- enable_intranode_visibility
- enable_kubernetes_alpha
- enable_l4_ilb_subsetting
- enable_legacy_abac
- enable_multi_networking
- enable_shielded_nodes
- enable_tpu
- initial_node_count
- location
- logging_service
- min_master_version
- monitoring_service
- name
- network
- networking_mode
- node_locations
- node_version
- private_ipv6_google_access
- project
- remove_default_node_pool
- resource_labels
- subnetwork

### Nested Blocks

#### addons_config
- dns_cache_config.enabled
- gce_persistent_disk_csi_driver_config.enabled
- gcs_fuse_csi_driver_config.enabled
- horizontal_pod_autoscaling.disabled
- http_load_balancing.disabled
- network_policy_config.disabled
- ray_operator_config.enabled

#### authenticator_groups_config
- security_group

#### binary_authorization
- evaluation_mode

#### cluster_autoscaling
- auto_provisioning_locations
- autoscaling_profile
- enabled

#### control_plane_endpoints_config
- dns_endpoint_config.allow_external_traffic
- dns_endpoint_config.endpoint
- ip_endpoints_config.enabled

#### database_encryption
- key_name
- state

#### default_snat_status
- disabled

#### enterprise_config
- desired_tier

#### fleet
- project

#### ip_allocation_policy
- cluster_ipv4_cidr_block
- cluster_secondary_range_name
- services_ipv4_cidr_block
- services_secondary_range_name
- stack_type
- pod_cidr_overprovision_config.disabled

#### logging_config
- enable_components

#### master_auth
- client_certificate_config.issue_client_certificate

#### monitoring_config
- enable_components
- advanced_datapath_observability_config.enable_metrics
- advanced_datapath_observability_config.enable_relay
- managed_prometheus.enabled
- managed_prometheus.auto_monitoring_config.scope

#### network_policy
- enabled
- provider

#### node_config (Cluster Level - Only if default pool not removed)
- boot_disk_kms_key
- disk_size_gb
- disk_type
- enable_confidential_storage
- image_type
- labels
- local_ssd_count
- local_ssd_encryption_mode
- logging_variant
- machine_type
- max_run_duration
- metadata
- min_cpu_platform
- node_group
- oauth_scopes
- preemptible
- resource_labels
- resource_manager_tags
- service_account
- spot
- storage_pools
- tags
- advanced_machine_features.enable_nested_virtualization
- advanced_machine_features.threads_per_core
- kubelet_config.allowed_unsafe_sysctls
- kubelet_config.container_log_max_files
- kubelet_config.container_log_max_size
- kubelet_config.cpu_cfs_quota
- kubelet_config.cpu_cfs_quota_period
- kubelet_config.cpu_manager_policy
- kubelet_config.image_gc_high_threshold_percent
- kubelet_config.image_gc_low_threshold_percent
- kubelet_config.image_maximum_gc_age
- kubelet_config.image_minimum_gc_age
- kubelet_config.insecure_kubelet_readonly_port_enabled
- kubelet_config.pod_pids_limit
- shielded_instance_config.enable_integrity_monitoring
- shielded_instance_config.enable_secure_boot
- windows_node_config.osversion

#### node_pool (Repeated Block - Handled by separate resource converter)
- *Fields within this block are handled by the `google_container_node_pool` converter*

#### node_pool_auto_config
- resource_manager_tags

#### node_pool_defaults.node_config_defaults
- insecure_kubelet_readonly_port_enabled
- logging_variant

#### notification_config.pubsub
- enabled
- topic

#### private_cluster_config
- enable_private_endpoint
- enable_private_nodes
- master_ipv4_cidr_block
- private_endpoint_subnetwork
- master_global_access_config.enabled

#### release_channel
- channel

#### secret_manager_config
- enabled

#### security_posture_config
- mode
- vulnerability_mode

#### service_external_ips_config
- enabled

---
**Note:**
- Some fields might be `null` or empty in the CAI asset but may need specific handling for Terraform representation.
- CAI field names might need mapping to Terraform field names.