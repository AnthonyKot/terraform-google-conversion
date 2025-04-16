Certainly! Here is a **strict list of fields** (including nested blocks) that must be handled by your converter to produce a Terraform resource block for `google_container_cluster` that matches your zonal-clustet.tf output **100% accurately** from the provided CAI asset.

---

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

---

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

#### node_config
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

#### node_pool (repeated block)
- initial_node_count
- max_pods_per_node
- name
- name_prefix
- node_count
- node_locations
- version
- management.auto_repair
- management.auto_upgrade
- network_config.create_pod_range
- network_config.enable_private_nodes
- network_config.pod_ipv4_cidr_block
- network_config.pod_range
- node_config (same as above)
- queued_provisioning.enabled
- upgrade_settings.max_surge
- upgrade_settings.max_unavailable
- upgrade_settings.strategy

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
- Some fields may be `null` or empty in the CAI asset, but must still be handled for strict fidelity.
- Some blocks (like `node_pool`) are repeated/nested and must be handled as lists.
- You may need to map CAI field names to Terraform field names (they sometimes differ).

---

**To achieve 100% accurate conversion, your converter must extract and map all of the above fields and nested blocks from the CAI asset to the Terraform resource.**





1. Missing Top-Level Arguments
The following fields are present in zonal-clustet.tf but missing in converted.tf:

allow_net_admin
datapath_provider
default_max_pods_per_node
deletion_protection
description
disable_l4_lb_firewall_reconciliation
enable_autopilot
enable_cilium_clusterwide_network_policy
enable_fqdn_network_policy
enable_intranode_visibility
enable_kubernetes_alpha
enable_l4_ilb_subsetting
enable_legacy_abac
enable_multi_networking
enable_shielded_nodes
enable_tpu
initial_node_count
location
logging_service
min_master_version
monitoring_service
networking_mode
node_locations
node_version
private_ipv6_google_access
remove_default_node_pool
resource_labels
2. Block/Nested Structure Differences
add-ons config
converted.tf is missing:
gcs_fuse_csi_driver_config
ray_operator_config
Some blocks may have different default values or missing fields.
authenticator_groups_config
Present in zonal-clustet.tf, missing in converted.tf.
binary_authorization
Present in zonal-clustet.tf, missing in converted.tf.
cluster_autoscaling
Present in zonal-clustet.tf, missing in converted.tf.
control_plane_endpoints_config
Present in zonal-clustet.tf, missing in converted.tf.
database_encryption
Present in zonal-clustet.tf, missing in converted.tf.
default_snat_status
Present in zonal-clustet.tf, missing in converted.tf.
enterprise_config
Present in zonal-clustet.tf, missing in converted.tf.
fleet
Present in zonal-clustet.tf, missing in converted.tf.
logging_config
Present in zonal-clustet.tf, missing in converted.tf.
master_auth
Present in zonal-clustet.tf, missing in converted.tf.
monitoring_config
Present in zonal-clustet.tf, missing in converted.tf.
network_policy
Present in zonal-clustet.tf, missing in converted.tf.
node_pool and related sub-blocks
Present in zonal-clustet.tf, missing in converted.tf.
node_pool_auto_config
Present in zonal-clustet.tf, missing in converted.tf.
node_pool_defaults
Present in zonal-clustet.tf, missing in converted.tf.
notification_config
Present in zonal-clustet.tf, missing in converted.tf.
private_cluster_config
Present in zonal-clustet.tf, missing in converted.tf.
release_channel
Present in zonal-clustet.tf, missing in converted.tf.
secret_manager_config
Present in zonal-clustet.tf, missing in converted.tf.
security_posture_config
Present in zonal-clustet.tf, missing in converted.tf.
service_external_ips_config
Present in zonal-clustet.tf, missing in converted.tf.
3. Field Value Differences
Some fields in converted.tf use short names for network/subnetwork ("default") instead of full resource paths as in zonal-clustet.tf.
Some blocks in converted.tf have fewer or defaulted values compared to the more explicit values in zonal-clustet.tf.
4. Instructions for LLM
Here’s a sample instruction you can give to an LLM to improve the converter:

Instruction Example:

Update the GKE cluster CAI-to-Terraform converter so that it produces a google_container_cluster resource with all fields and nested blocks present in the real Terraform output (zonal-clustet.tf).
Specifically:

Map all top-level arguments from the CAI asset to their Terraform equivalents, including those currently missing (see list above).
For each nested block (e.g., addons_config, authenticator_groups_config, binary_authorization, etc.), ensure all subfields and sub-blocks are mapped and populated as in the reference TF.
Use full resource paths for fields like network and subnetwork, not just short names.
Ensure all default and explicit values are set to match the real TF output, even if they are default or null.
If a field or block is missing in the CAI asset, set it to null, {}, or the default value as appropriate for Terraform.
The goal is for the converter output to be structurally and semantically identical to the hand-written zonal-clustet.tf.
Summary Table:

Section/Field	In Reference TF	In Converted TF	Action Needed
Top-level fields	Yes	No	Add all missing fields
Nested blocks (e.g. addons)	Yes	Partially	Add all missing sub-blocks
Full resource paths	Yes	No	Use full paths
Default/null values	Yes	Partially	Set as in reference TF
