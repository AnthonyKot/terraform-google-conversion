resource "google_container_cluster" "cluster-zonal-1" {
  addons_config {
    dns_cache_config {
      enabled = false
    }

    gce_persistent_disk_csi_driver_config {
      enabled = true
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
  }

  ip_allocation_policy {
    cluster_ipv4_cidr_block      = "10.52.0.0/14"
    cluster_secondary_range_name = "gke-cluster-zonal-1-pods-de672f38"
    services_ipv4_cidr_block     = "34.118.224.0/20"
    stack_type                   = "IPV4"
  }

  name    = "cluster-zonal-1"
  network = "default"

  node_config {
    disk_size_gb   = 100
    disk_type      = "pd-balanced"
    image_type     = "COS_CONTAINERD"
    kubelet_config = {}

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

  project    = "cluster-converter-1"
  subnetwork = "default"
}