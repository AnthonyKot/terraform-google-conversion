# Default Values Fix for GKE Resources

## Issue
GKE clusters and node pools have several default configuration blocks that should be included in the Terraform output, even when not specified in the CAI data:

### Node Pool Defaults
1. The `queued_provisioning` field defaults to `enabled = false`
2. The `node_config.advanced_machine_features` block defaults to `enable_nested_virtualization = false`
3. The `node_config.resource_manager_tags` field defaults to an empty map
4. The `initial_node_count` field (with autoscaling) or `node_count` field (without autoscaling) should default to `1`

### Cluster Defaults
5. The `monitoring_config` block should include default values:
   - `enable_components = ["SYSTEM_COMPONENTS", "STORAGE", "POD", "DEPLOYMENT", "STATEFULSET", "DAEMONSET", "HPA", "CADVISOR", "KUBELET"]`
   - `managed_prometheus { enabled = true }`
   - `managed_prometheus.auto_monitoring_config { scope = "NONE" }`

6. The `master_auth` block should include default values:
   - `client_certificate_config { issue_client_certificate = false }`

7. The `control_plane_endpoints_config` block should include default values:
   - `dns_endpoint_config { allow_external_traffic = false }`
   - `ip_endpoints_config { enabled = true }`

## Changes Made

1. Added a new `flattenQueuedProvisioning` function:
```go
func flattenQueuedProvisioning(config *container.QueuedProvisioning) []interface{} {
    qp := make(map[string]interface{})
    
    // Always set enabled to false by default
    qp["enabled"] = false
    
    // If config exists and enabled is true, override the default
    if config != nil && config.Enabled {
        qp["enabled"] = true
    }
    
    return []interface{}{qp}
}
```

2. Added a new `flattenAdvancedMachineFeatures` function:
```go
func flattenAdvancedMachineFeatures(config *container.AdvancedMachineFeatures) []interface{} {
    amf := make(map[string]interface{})
    
    // Always set enable_nested_virtualization to false by default
    amf["enable_nested_virtualization"] = false
    
    // If config exists and enable_nested_virtualization is true, override the default
    if config != nil && config.EnableNestedVirtualization {
        amf["enable_nested_virtualization"] = true
    }
    
    return []interface{}{amf}
}
```

3. Updated the `convertNodePoolData` function to always include the queued_provisioning field:
```go
// Always include queued_provisioning with enabled=false as the default
hclData["queued_provisioning"] = flattenQueuedProvisioning(nodePool.QueuedProvisioning)
```

4. Modified the node count handling to always include either initial_node_count or node_count with default value of 1:
```go
// Node count / Autoscaling - Handle initial_node_count and node_count appropriately
autoscalingBlock := flattenNodePoolAutoscaling(nodePool.Autoscaling)
if autoscalingBlock != nil {
    hclData["autoscaling"] = autoscalingBlock
    // With autoscaling, always include initial_node_count
    if nodePool.InitialNodeCount > 0 {
        hclData["initial_node_count"] = nodePool.InitialNodeCount
    } else {
        // Default to 1 if not specified
        hclData["initial_node_count"] = 1
    }
} else {
    // Autoscaling is disabled or default. Always include node_count.
    if nodePool.InitialNodeCount > 0 {
        hclData["node_count"] = nodePool.InitialNodeCount
    } else {
        // Default to 1 if not specified
        hclData["node_count"] = 1
    }
}
```

4. Modified the `flattenNodeConfig` function to:
   - Always create default node config with advanced_machine_features and resource_manager_tags when config is nil
   - Always include advanced_machine_features with default values
   - Always include resource_manager_tags, even if it's an empty map

5. Updated the resource_manager_tags section to always include it:
```go
// resource_manager_tags: Always include empty map as default
if config.ResourceManagerTags != nil && len(config.ResourceManagerTags.Tags) > 0 {
    nodeConfig["resource_manager_tags"] = config.ResourceManagerTags.Tags // Assign the actual map
} else {
    nodeConfig["resource_manager_tags"] = make(map[string]string) // Empty map as default
}
hasNonDefaultConfig = true
```

6. Completely rewrote the `flattenMonitoringConfig` function to always include defaults:
```go
// flattenMonitoringConfig: Always include default values
func flattenMonitoringConfig(config *container.MonitoringConfig) []interface{} {
    mc := make(map[string]interface{})
    
    // Default components to include if not specified
    defaultComponents := []string{
        "SYSTEM_COMPONENTS", 
        "STORAGE", 
        "POD", 
        "DEPLOYMENT", 
        "STATEFULSET", 
        "DAEMONSET", 
        "HPA", 
        "CADVISOR", 
        "KUBELET",
    }

    // Set components - either from config or defaults
    if config != nil && config.ComponentConfig != nil && len(config.ComponentConfig.EnableComponents) > 0 {
        mc["enable_components"] = config.ComponentConfig.EnableComponents
    } else {
        mc["enable_components"] = defaultComponents
    }
    
    // Managed Prometheus: Always include with default enabled=true
    mpc := make(map[string]interface{})
    
    // Default to enabled=true
    mpc["enabled"] = true
    
    // Override if explicitly set to false
    if config != nil && config.ManagedPrometheusConfig != nil && !config.ManagedPrometheusConfig.Enabled {
        mpc["enabled"] = false
    }

    // Always include auto_monitoring_config with default scope=NONE
    amc := map[string]interface{}{"scope": "NONE"}
    
    // Override if explicitly set
    if config != nil && config.ManagedPrometheusConfig != nil && 
       config.ManagedPrometheusConfig.AutoMonitoringConfig != nil && 
       config.ManagedPrometheusConfig.AutoMonitoringConfig.Scope != "" {
        amc["scope"] = config.ManagedPrometheusConfig.AutoMonitoringConfig.Scope
    }
    
    // Add auto_monitoring_config to managed_prometheus
    mpc["auto_monitoring_config"] = []interface{}{amc}
    
    // Add managed_prometheus to monitoring_config
    mc["managed_prometheus"] = []interface{}{mpc}
    
    return []interface{}{mc}
}
```

7. Updated the convertClusterData function to always include monitoring_config:
```go
// Always include monitoring_config
hclData["monitoring_config"] = monitoringConfigBlock
```

8. Rewrote the `flattenMasterAuth` function to always include defaults:
```go
// flattenMasterAuth: Always include client_certificate_config with default issue_client_certificate = false
func flattenMasterAuth(auth *container.MasterAuth) []interface{} {
    // Create default master_auth block
    ma := make(map[string]interface{})
    
    // Create client_certificate_config block
    ccc := make(map[string]interface{})
    
    // Default to issue_client_certificate = false
    ccc["issue_client_certificate"] = false
    
    // Override if explicitly set to true
    if auth != nil && auth.ClientCertificateConfig != nil && auth.ClientCertificateConfig.IssueClientCertificate {
        ccc["issue_client_certificate"] = true
    }
    
    // Add client_certificate_config to master_auth
    ma["client_certificate_config"] = []interface{}{ccc}
    
    // Return master_auth block
    return []interface{}{ma}
}
```

9. Updated the convertClusterData function to always include master_auth:
```go
// Always include master_auth block with client_certificate_config.issue_client_certificate = false
hclData["master_auth"] = flattenMasterAuth(cluster.MasterAuth)
```

10. Created a new `flattenControlPlaneEndpointsConfig` function to always include defaults:
```go
// flattenControlPlaneEndpointsConfig: Always include with default values
func flattenControlPlaneEndpointsConfig(config *container.ControlPlaneEndpointsConfig) []interface{} {
    result := make(map[string]interface{})
    
    // DNS Endpoint Config - default to allow_external_traffic = false
    dnsConfig := make(map[string]interface{})
    dnsConfig["allow_external_traffic"] = false
    
    // Override if explicitly set to true
    if config != nil && config.DnsEndpointConfig != nil && config.DnsEndpointConfig.AllowExternalTraffic {
        dnsConfig["allow_external_traffic"] = true
    }
    
    result["dns_endpoint_config"] = []interface{}{dnsConfig}
    
    // IP Endpoints Config - default to enabled = true
    ipConfig := make(map[string]interface{})
    ipConfig["enabled"] = true
    
    // Override if explicitly set to false
    if config != nil && config.IpEndpointsConfig != nil && !config.IpEndpointsConfig.Enabled {
        ipConfig["enabled"] = false
    }
    
    result["ip_endpoints_config"] = []interface{}{ipConfig}
    
    return []interface{}{result}
}
```

11. Updated the convertClusterData function to always include control_plane_endpoints_config:
```go
// Always include control_plane_endpoints_config with defaults
hclData["control_plane_endpoints_config"] = flattenControlPlaneEndpointsConfig(cluster.ControlPlaneEndpointsConfig)
```

## Expected Behavior

After these changes:

### Node Pool Behavior
1. All GKE node pools will have the `queued_provisioning { enabled = false }` block explicitly included
2. All node_config blocks will include `advanced_machine_features { enable_nested_virtualization = false }`
3. All node_config blocks will include `resource_manager_tags = {}`
4. All node pools will include either `initial_node_count = 1` (with autoscaling) or `node_count = 1` (without autoscaling) when not specified

### Cluster Behavior
5. All GKE clusters will have the `monitoring_config` block explicitly included with these defaults:
   ```hcl
   monitoring_config {
     enable_components = ["SYSTEM_COMPONENTS", "STORAGE", "POD", "DEPLOYMENT", "STATEFULSET", "DAEMONSET", "HPA", "CADVISOR", "KUBELET"]
     managed_prometheus {
       enabled = true
       auto_monitoring_config {
         scope = "NONE"
       }
     }
   }
   ```

6. All GKE clusters will have the `master_auth` block explicitly included with these defaults:
   ```hcl
   master_auth {
     client_certificate_config {
       issue_client_certificate = false
     }
   }
   ```

7. All GKE clusters will have the `control_plane_endpoints_config` block explicitly included with these defaults:
   ```hcl
   control_plane_endpoints_config {
     dns_endpoint_config {
       allow_external_traffic = false
     }
     ip_endpoints_config {
       enabled = true
     }
   }
   ```

This ensures that these fields are consistently represented in the generated Terraform, making it more explicit and easier to understand the configuration.