# Node Pool Default Values Fix

## Issue
GKE node pools have several default configuration blocks that should be included in the Terraform output, even when not specified in the CAI data:

1. The `queued_provisioning` field defaults to `enabled = false`
2. The `node_config.advanced_machine_features` block defaults to `enable_nested_virtualization = false`
3. The `node_config.resource_manager_tags` field defaults to an empty map
4. The `initial_node_count` field (with autoscaling) or `node_count` field (without autoscaling) should default to `1`

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

## Expected Behavior

After these changes:
1. All GKE node pools will have the `queued_provisioning { enabled = false }` block explicitly included
2. All node_config blocks will include `advanced_machine_features { enable_nested_virtualization = false }`
3. All node_config blocks will include `resource_manager_tags = {}`
4. All node pools will include either `initial_node_count = 1` (with autoscaling) or `node_count = 1` (without autoscaling) when not specified

This ensures that these fields are consistently represented in the generated Terraform, making it more explicit and easier to understand the configuration.