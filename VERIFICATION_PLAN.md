# GKE Converter Verification Plan

## Overview

This document outlines a systematic approach to test the GKE converter with various cluster configurations. We'll create multiple test cases covering different GKE features and deployment patterns to ensure the converter produces correct Terraform configurations.

## Test Approach

We will follow a round-trip verification methodology:

1. Create GKE cluster (manually or via Terraform)
2. Export cluster configuration as JSON in Cloud Asset Inventory (CAI) format
3. Convert JSON to Terraform HCL using our converter
4. Apply the Terraform configuration with a different cluster name
5. Export the newly created cluster as JSON
6. Compare original and new JSON field by field

This approach verifies that our converter preserves all essential configurations and properly handles default values.

## Test Scenarios

### 1. Basic Zonal Cluster (Current Example)

**Features:**
- Single-zone cluster (us-central1-c)
- Default network
- Basic node pool with e2-medium machines
- Enabled kubelet_config with insecure_kubelet_readonly_port_enabled
- Basic security posture config
- Managed Prometheus monitoring

**Field Verification Focus:**
- kubelet_config.insecure_kubelet_readonly_port_enabled type handling
- region extraction for subnetwork path
- node_version, max_pods_per_node
- oauth_scopes

### 2. Regional Multi-Zone Cluster

**Features:**
- Regional deployment (us-west1)
- Multiple zones for high availability (us-west1-a, us-west1-b, us-west1-c)
- Custom networking with VPC and subnets
- Multiple node pools with different machine types

**Field Verification Focus:**
- Regional vs. zonal location handling
- node_locations arrays
- VPC networking configuration
- Multiple node pools with different configurations

### 3. Private Cluster with Advanced Networking

**Features:**
- Private cluster (no public endpoint)
- Custom VPC with private service connections
- Network policy enabled
- Custom master authorized networks
- VPC-native networking

**Field Verification Focus:**
- private_cluster_config block
- master_authorized_networks_config
- network_policy
- ip_allocation_policy with custom ranges

### 4. Cluster with Advanced Security Features

**Features:**
- Workload Identity enabled
- Binary Authorization
- GKE Sandbox
- Shielded Nodes
- Customer-managed encryption keys (CMEK)

**Field Verification Focus:**
- workload_identity_config
- binary_authorization
- database_encryption
- node_config.sandbox_config
- shielded_instance_config

### 5. Cluster with Autopilot and Advanced Features

**Features:**
- Autopilot mode
- Release channel (regular)
- Cost management enabled
- Advanced monitoring and logging
- Maintenance window and exclusions

**Field Verification Focus:**
- enable_autopilot
- release_channel
- cost_management_config
- monitoring_config with all components
- maintenance_policy

## Test Data Structure

Each test should have:

1. **Input Files:**
   - `test_name.json` - CAI formatted JSON export of the original cluster
   - `test_name_expected.tf` - Manually validated target Terraform output

2. **Test Process Files:**
   - `test_name_converted.tf` - Output from our converter
   - `test_name_applied.json` - JSON export of the cluster created by applying the converted Terraform

3. **Verification:**
   - Automated diff of key fields between original and recreated cluster
   - Documentation of any expected differences

## Handling Type Differences

The test system should account for known type differences:

1. **String vs. Boolean:**
   - JSON may represent boolean values as strings (`"true"` vs `true`)
   - Document which fields are affected (e.g., kubelet_config.insecure_kubelet_readonly_port_enabled)

2. **String vs. Integer:**
   - Numeric values like maxPodsPerNode may be represented as strings in JSON
   - Test should normalize these for comparison

3. **Empty Objects vs. Nil:**
   - Empty JSON objects ({}) vs. null/nil values
   - Test should treat these as equivalent in many cases

## Next Steps

1. Create test script to automate the verification process
2. Develop cluster templates for each test scenario
3. Document expected field mappings for each test case
4. Implement CI/CD integration for automated testing

## Example Field Mapping Table

For each test, create a table of critical fields to verify:

| Field Path | Original JSON | Expected TF | Converted TF | Applied JSON | Status |
|------------|--------------|------------|-------------|--------------|--------|
| kubeletConfig.insecureKubeletReadonlyPortEnabled | true | true | "true" | true | ⚠️ Type mismatch |
| locations[0] | "us-central1-c" | "us-central1-c" | "us-central1-c" | "us-central1-c" | ✅ Match |
| subnetwork | "default" | "projects/.../subnetworks/default" | "projects/.../subnetworks/default" | "default" | ✅ Format expected |
| ... | ... | ... | ... | ... | ... |