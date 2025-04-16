resource "google_container_cluster" "cluster-zonal-1" {
  binary_authorization {
    evaluation_mode = "DISABLED"
  }

  default_max_pods_per_node = 110

  fleet {
    project = "cluster-converter-1"
  }

  ip_allocation_policy {
    cluster_ipv4_cidr_block      = "10.52.0.0/14"
    cluster_secondary_range_name = "gke-cluster-zonal-1-pods-de672f38"
    services_ipv4_cidr_block     = "34.118.224.0/20"
  }

  logging_config {}

  master_version = "1.31.5-gke.1233001"

  monitoring_config {
    managed_prometheus {
      enabled = true
    }
  }

  name                     = "cluster-zonal-1"
  networking_mode          = "VPC_NATIVE"
  node_locations           = ["us-central1-c"]
  project                  = "cluster-converter-1"
  remove_default_node_pool = true
  subnetwork               = "default"
}

resource "google_container_node_pool" "pool-name" {
  cluster = "cluster-zonal-1"

  management {
    auto_upgrade = false
  }

  name = "pool-name"

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

    oauth_scopes = ["https://www.googleapis.com/auth/devstorage.read_only", "https://www.googleapis.com/auth/logging.write", "https://www.googleapis.com/auth/monitoring", "https://www.googleapis.com/auth/service.management.readonly", "https://www.googleapis.com/auth/servicecontrol", "https://www.googleapis.com/auth/trace.append"]

    resource_labels = {
      goog-gke-node-pool-provisioning-model = "on-demand"
    }
  }

  node_count     = 1
  node_locations = ["us-central1-c"]
  project        = "cluster-converter-1"
  version        = "1.31.5-gke.1233001"
}