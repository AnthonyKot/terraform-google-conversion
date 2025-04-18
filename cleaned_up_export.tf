resource "google_container_cluster" "cluster-zonal-1" {
  binary_authorization {
    evaluation_mode = "DISABLED"
  }
  default_max_pods_per_node = 110
  ip_allocation_policy {
    cluster_ipv4_cidr_block    = "10.52.0.0/14"
    cluster_secondary_range_name = "gke-cluster-zonal-1-pods-de672f38"
    services_ipv4_cidr_block = "34.118.224.0/20"
  }
  location = "us-central1-c"
  logging_config {
    enable_components = ["SYSTEM_COMPONENTS", "WORKLOADS"]
  }
  monitoring_config {
    enable_components = ["SYSTEM_COMPONENTS", "STORAGE", "POD", "DEPLOYMENT", "STATEFULSET", "DAEMONSET", "HPA", "CADVISOR", "KUBELET"]
    managed_prometheus {
      enabled = true
      auto_monitoring_config {
        scope = "NONE"
      }
    }
  }
  name     = "cluster-zonal-1"
  node_version              = "1.31.5-gke.1233001"
  project                   = "cluster-converter-1"
  remove_default_node_pool  = true
  security_posture_config {
    mode               = "BASIC"
    vulnerability_mode = "VULNERABILITY_DISABLED"
  }
  subnetwork = "projects/cluster-converter-1/regions/us-central1/subnetworks/default"
}

resource "google_container_node_pool" "pool-name" {
  cluster  = "cluster-zonal-1"
  location = "us-central1-c"
  name     = "pool-name"
  node_count = 1
  node_locations = ["us-central1-c"]
  version   = "1.31.5-gke.1233001"
  max_pods_per_node = 110
  
  management {
    auto_upgrade = false
  }
  
  network_config {
    pod_ipv4_cidr_block = "10.52.0.0/14"
    pod_range           = "gke-cluster-zonal-1-pods-de672f38"
  }
  
  node_config {
    disk_type    = "pd-balanced"
    machine_type = "e2-medium"
    metadata = {
      disable-legacy-endpoints = "true"
    }
    resource_labels = {
      goog-gke-node-pool-provisioning-model = "on-demand"
    }
    kubelet_config {
      insecure_kubelet_readonly_port_enabled = true
    }
  }
  project = "cluster-converter-1"
}

# --- Removed Attributes (Defaults or Null/Empty) ---
# allow_net_admin, cluster_ipv4_cidr, datapath_provider, deletion_protection, description,
# disable_l4_lb_firewall_reconciliation, enable_autopilot, enable_cilium_clusterwide_network_policy,
# enable_fqdn_network_policy, enable_intranode_visibility, enable_kubernetes_alpha,
# enable_l4_ilb_subsetting, enable_legacy_abac, enable_multi_networking, enable_shielded_nodes,
# enable_tpu, initial_node_count, logging_service, min_master_version, monitoring_service, network,
# networking_mode, node_locations (cluster level), private_ipv6_google_access, resource_labels (cluster level)

# --- Removed Blocks (All Defaults or Computed/Redundant) ---
# addons_config, authenticator_groups_config, cluster_autoscaling, control_plane_endpoints_config,
# database_encryption, default_snat_status, enterprise_config, fleet, master_auth, network_policy,
# node_config (cluster level), node_pool_auto_config, node_pool_defaults, notification_config,
# private_cluster_config, release_channel, secret_manager_config, service_external_ips_config