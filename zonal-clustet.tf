resource "google_container_cluster" "cluster-zonal-1" {
  allow_net_admin                          = null
  cluster_ipv4_cidr                        = "10.52.0.0/14"
  datapath_provider                        = "LEGACY_DATAPATH"
  default_max_pods_per_node                = 110
  deletion_protection                      = true
  description                              = null
  disable_l4_lb_firewall_reconciliation    = false
  enable_autopilot                         = false
  enable_cilium_clusterwide_network_policy = false
  enable_fqdn_network_policy               = false
  enable_intranode_visibility              = false
  enable_kubernetes_alpha                  = false
  enable_l4_ilb_subsetting                 = false
  enable_legacy_abac                       = false
  enable_multi_networking                  = false
  enable_shielded_nodes                    = true
  enable_tpu                               = false
  initial_node_count                       = 0
  location                                 = "us-central1-c"
  logging_service                          = "logging.googleapis.com/kubernetes"
  min_master_version                       = null
  monitoring_service                       = "monitoring.googleapis.com/kubernetes"
  name                                     = "cluster-zonal-1"
  network                                  = "projects/cluster-converter-1/global/networks/default"
  networking_mode                          = "VPC_NATIVE"
  node_locations                           = []
  node_version                             = "1.31.5-gke.1233001"
  private_ipv6_google_access               = null
  project                                  = "cluster-converter-1"
  remove_default_node_pool                 = null
  resource_labels                          = {}
  subnetwork                               = "projects/cluster-converter-1/regions/us-central1/subnetworks/default"
  addons_config {
    dns_cache_config {
      enabled = false
    }
    gce_persistent_disk_csi_driver_config {
      enabled = true
    }
    gcs_fuse_csi_driver_config {
      enabled = false
    }
    horizontal_pod_autoscaling {
      disabled = false
    }
    http_load_balancing {
      disabled = false
    }
    network_policy_config {
      disabled = true
    }
    ray_operator_config {
      enabled = false
    }
  }
  authenticator_groups_config {
    security_group = ""
  }
  binary_authorization {
    evaluation_mode = "DISABLED"
  }
  cluster_autoscaling {
    auto_provisioning_locations = []
    autoscaling_profile         = "BALANCED"
    enabled                     = false
  }
  control_plane_endpoints_config {
    dns_endpoint_config {
      allow_external_traffic = false
      endpoint               = "gke-de672f388d3b425692cdd083e24ee12096f8-1089806473266.us-central1-c.gke.goog"
    }
    ip_endpoints_config {
      enabled = true
    }
  }
  database_encryption {
    key_name = null
    state    = "DECRYPTED"
  }
  default_snat_status {
    disabled = false
  }
  enterprise_config {
    desired_tier = "STANDARD"
  }
  fleet {
    project = "cluster-converter-1"
  }
  ip_allocation_policy {
    cluster_ipv4_cidr_block       = "10.52.0.0/14"
    cluster_secondary_range_name  = "gke-cluster-zonal-1-pods-de672f38"
    services_ipv4_cidr_block      = "34.118.224.0/20"
    services_secondary_range_name = null
    stack_type                    = "IPV4"
    pod_cidr_overprovision_config {
      disabled = false
    }
  }
  logging_config {
    enable_components = ["SYSTEM_COMPONENTS", "WORKLOADS"]
  }
  master_auth {
    client_certificate_config {
      issue_client_certificate = false
    }
  }
  monitoring_config {
    enable_components = ["SYSTEM_COMPONENTS", "STORAGE", "POD", "DEPLOYMENT", "STATEFULSET", "DAEMONSET", "HPA", "CADVISOR", "KUBELET"]
    advanced_datapath_observability_config {
      enable_metrics = false
      enable_relay   = false
    }
    managed_prometheus {
      enabled = true
      auto_monitoring_config {
        scope = "NONE"
      }
    }
  }
  network_policy {
    enabled  = false
    provider = "PROVIDER_UNSPECIFIED"
  }
  node_config {
    boot_disk_kms_key           = null
    disk_size_gb                = 100
    disk_type                   = "pd-balanced"
    enable_confidential_storage = false
    image_type                  = "COS_CONTAINERD"
    labels                      = {}
    local_ssd_count             = 0
    local_ssd_encryption_mode   = null
    logging_variant             = "DEFAULT"
    machine_type                = "e2-medium"
    max_run_duration            = null
    metadata = {
      disable-legacy-endpoints = "true"
    }
    min_cpu_platform = null
    node_group       = null
    oauth_scopes     = ["https://www.googleapis.com/auth/devstorage.read_only", "https://www.googleapis.com/auth/logging.write", "https://www.googleapis.com/auth/monitoring", "https://www.googleapis.com/auth/service.management.readonly", "https://www.googleapis.com/auth/servicecontrol", "https://www.googleapis.com/auth/trace.append"]
    preemptible      = false
    resource_labels = {
      goog-gke-node-pool-provisioning-model = "on-demand"
    }
    resource_manager_tags = {}
    service_account       = "default"
    spot                  = false
    storage_pools         = []
    tags                  = []
    advanced_machine_features {
      enable_nested_virtualization = false
      threads_per_core             = 0
    }
    kubelet_config {
      allowed_unsafe_sysctls                 = []
      container_log_max_files                = 0
      container_log_max_size                 = null
      cpu_cfs_quota                          = false
      cpu_cfs_quota_period                   = null
      cpu_manager_policy                     = null
      image_gc_high_threshold_percent        = 0
      image_gc_low_threshold_percent         = 0
      image_maximum_gc_age                   = null
      image_minimum_gc_age                   = null
      insecure_kubelet_readonly_port_enabled = "TRUE"
      pod_pids_limit                         = 0
    }
    shielded_instance_config {
      enable_integrity_monitoring = true
      enable_secure_boot          = false
    }
    windows_node_config {
      osversion = null
    }
  }
  node_pool {
    initial_node_count = 1
    max_pods_per_node  = 110
    name               = "pool-name"
    name_prefix        = null
    node_count         = 1
    node_locations     = ["us-central1-c"]
    version            = "1.31.5-gke.1233001"
    management {
      auto_repair  = true
      auto_upgrade = false
    }
    network_config {
      create_pod_range     = false
      enable_private_nodes = false
      pod_ipv4_cidr_block  = "10.52.0.0/14"
      pod_range            = "gke-cluster-zonal-1-pods-de672f38"
    }
    node_config {
      boot_disk_kms_key           = null
      disk_size_gb                = 100
      disk_type                   = "pd-balanced"
      enable_confidential_storage = false
      image_type                  = "COS_CONTAINERD"
      labels                      = {}
      local_ssd_count             = 0
      local_ssd_encryption_mode   = null
      logging_variant             = "DEFAULT"
      machine_type                = "e2-medium"
      max_run_duration            = null
      metadata = {
        disable-legacy-endpoints = "true"
      }
      min_cpu_platform = null
      node_group       = null
      oauth_scopes     = ["https://www.googleapis.com/auth/devstorage.read_only", "https://www.googleapis.com/auth/logging.write", "https://www.googleapis.com/auth/monitoring", "https://www.googleapis.com/auth/service.management.readonly", "https://www.googleapis.com/auth/servicecontrol", "https://www.googleapis.com/auth/trace.append"]
      preemptible      = false
      resource_labels = {
        goog-gke-node-pool-provisioning-model = "on-demand"
      }
      resource_manager_tags = {}
      service_account       = "default"
      spot                  = false
      storage_pools         = []
      tags                  = []
      advanced_machine_features {
        enable_nested_virtualization = false
        threads_per_core             = 0
      }
      kubelet_config {
        allowed_unsafe_sysctls                 = []
        container_log_max_files                = 0
        container_log_max_size                 = null
        cpu_cfs_quota                          = false
        cpu_cfs_quota_period                   = null
        cpu_manager_policy                     = null
        image_gc_high_threshold_percent        = 0
        image_gc_low_threshold_percent         = 0
        image_maximum_gc_age                   = null
        image_minimum_gc_age                   = null
        insecure_kubelet_readonly_port_enabled = "TRUE"
        pod_pids_limit                         = 0
      }
      shielded_instance_config {
        enable_integrity_monitoring = true
        enable_secure_boot          = false
      }
      windows_node_config {
        osversion = null
      }
    }
    queued_provisioning {
      enabled = false
    }
    upgrade_settings {
      max_surge       = 1
      max_unavailable = 0
      strategy        = "SURGE"
    }
  }
  node_pool_auto_config {
    resource_manager_tags = null
  }
  node_pool_defaults {
    node_config_defaults {
      insecure_kubelet_readonly_port_enabled = "FALSE"
      logging_variant                        = "DEFAULT"
    }
  }
  notification_config {
    pubsub {
      enabled = false
      topic   = null
    }
  }
  private_cluster_config {
    enable_private_endpoint     = false
    enable_private_nodes        = false
    master_ipv4_cidr_block      = null
    private_endpoint_subnetwork = null
    master_global_access_config {
      enabled = false
    }
  }
  release_channel {
    channel = "UNSPECIFIED"
  }
  secret_manager_config {
    enabled = false
  }
  security_posture_config {
    mode               = "BASIC"
    vulnerability_mode = "VULNERABILITY_DISABLED"
  }
  service_external_ips_config {
    enabled = false
  }
}