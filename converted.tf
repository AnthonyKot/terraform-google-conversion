resource "google_container_cluster" "cluster-zonal-1" {
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

  datapath_provider         = "LEGACY_DATAPATH"
  default_max_pods_per_node = 110
  enable_shielded_nodes     = true

  ip_allocation_policy {
    cluster_ipv4_cidr_block      = "10.52.0.0/14"
    cluster_secondary_range_name = "gke-cluster-zonal-1-pods-de672f38"
    services_ipv4_cidr_block     = "34.118.224.0/20"
    stack_type                   = "IPV4"
  }

  logging_service          = "logging.googleapis.com/kubernetes"
  master_version           = "1.31.5-gke.1233001"
  monitoring_service       = "monitoring.googleapis.com/kubernetes"
  name                     = "cluster-zonal-1"
  network                  = "default"
  networking_mode          = "VPC_NATIVE"
  node_locations           = ["us-central1-c"]
  node_version             = "1.31.5-gke.1233001"
  project                  = "cluster-converter-1"
  remove_default_node_pool = true
  subnetwork               = "default"
}

resource "google_container_node_pool" "pool-name" {
  cluster = "cluster-zonal-1"

  management {
    auto_repair  = true
    auto_upgrade = false
  }

  name = "pool-name"

  network_config {
    create_pod_range    = false
    pod_ipv4_cidr_block = "10.52.0.0/14"
    pod_range           = "gke-cluster-zonal-1-pods-de672f38"
  }

  node_config {
    disk_size_gb = 100
    disk_type    = "pd-balanced"
    image_type   = "COS_CONTAINERD"

    kubelet_config {
      cpu_cfs_quota = false
    }

    machine_type = "e2-medium"

    metadata = {
      disable-legacy-endpoints = "true"
    }

    oauth_scopes    = ["https://www.googleapis.com/auth/devstorage.read_only", "https://www.googleapis.com/auth/logging.write", "https://www.googleapis.com/auth/monitoring", "https://www.googleapis.com/auth/service.management.readonly", "https://www.googleapis.com/auth/servicecontrol", "https://www.googleapis.com/auth/trace.append"]
    service_account = "default"

    shielded_instance_config {
      enable_integrity_monitoring = true
      enable_secure_boot          = false
    }
  }

  node_count = 1
  project    = "cluster-converter-1"

  upgrade_settings {
    max_surge       = 1
    max_unavailable = 0
    strategy        = "SURGE"
  }

  version = "1.31.5-gke.1233001"
}