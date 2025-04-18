# GKE Cluster Converter - Plan for April 17th, 2025

## Comparison: Exported (Desired) vs. Converted (Actual) HCL

### `google_container_cluster` Differences

| Attribute/Block             | Exported (Desired)                                                                                                                                    | Converted (Actual)                                                                       | Analysis                                                                                                                                                                                                                                                             |
| :-------------------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `location`                  | `"us-central1-c"`                                                                                                                                     | *Missing* | **ISSUE:** Required field missing in Converted.                                                                                                                                                                                                                  |
| `node_version`              | `"1.31.5-gke.1233001"`                                                                                                                                | *Missing* (has `master_version` instead)                                                   | **ISSUE:** Incorrect version field used/omitted. Cluster-level versioning likely shouldn't be set when `remove_default_node_pool=true`.                                                                                                                                 |
| `master_version`            | *Missing* | `"1.31.5-gke.1233001"`                                                                     | **ISSUE:** See `node_version`. Added when it shouldn't be if `remove_default_node_pool=true`.                                                                                                                                                                   |
| `logging_config`            | `enable_components = [...]`                                                                                                                           | `{}` (Empty block)                                                                       | **ISSUE:** Block content missing (`enable_components`).                                                                                                                                                                                                         |
| `monitoring_config`         | `enable_components = [...]`, `managed_prometheus { ..., auto_monitoring_config { ... }}`                                                              | `managed_prometheus { enabled = true }`                                                  | **ISSUE:** Missing `enable_components`. Missing `auto_monitoring_config` (Known v1 API limitation).                                                                                                                                                               |
| `security_posture_config`   | `mode = "BASIC", ...`                                                                                                                                 | *Missing* | **ISSUE:** Non-default block incorrectly omitted.                                                                                                                                                                                                               |
| `subnetwork`                | `"projects/.../subnetworks/default"` (Full path)                                                                                                      | `"default"` (Short name)                                                                 | **DIFFERENCE/Potential Issue:** Format mismatch. Need consistent approach.                                                                                                                                                                                      |
| `networking_mode`           | *Missing* (Implied default)                                                                                                                           | `"VPC_NATIVE"`                                                                           | **DIFFERENCE:** Explicitly added default/implied value.                                                                                                                                                                                                          |
| `node_locations`            | *Missing* (Default)                                                                                                                                   | `["us-central1-c"]`                                                                      | **DIFFERENCE:** Explicitly added potentially default value.                                                                                                                                                                                                     |
| `fleet`                     | *Missing* (Default)                                                                                                                                   | `fleet { project = "..." }`                                                              | **DIFFERENCE:** Explicitly added potentially default block.                                                                                                                                                                                                     |
| `node_pool` block           | Present (Inline definition)                                                                                                                           | *Missing* (Converter handled as separate resource)                                       | **STRUCTURAL DIFFERENCE:** Expected and correct converter behavior.                                                                                                                                                                                             |

### `google_container_node_pool` Differences

(Comparing Exported inline pool attributes to Converted resource)

| Attribute/Block          | Exported (Desired - inline)      | Converted (Actual - resource) | Analysis                                                            |
| :----------------------- | :------------------------------- | :---------------------------- | :------------------------------------------------------------------ |
| `location`               | (Implicit from cluster)        | *Missing* | **ISSUE:** Required field missing in Converted resource.            |
| `max_pods_per_node`      | `110`                            | *Missing* | **ISSUE:** Non-default attribute incorrectly omitted.               |
| `node_config.oauth_scopes` | *Missing* (Default)            | `[...]`                       | **ISSUE:** Default attribute incorrectly included.                  |
| `node_config.kubelet_config` | `{ insecure... = true }`       | *Missing* | **ISSUE:** Non-default block incorrectly omitted.                   |
| Other fields             | Match                            | Match                         | Most other fields align correctly between the two representations. |

---

## Plan for Tomorrow (Thursday, April 17th, 2025)

### High Priority (Fix Core Functionality & Missing Data)

1.  **Add `location` to Resources:**
    * Ensure `convertClusterData` adds the `location` attribute.
    * Ensure `convertNodePoolData` adds the `location` attribute.
2.  **Fix Config Block Content (`logging_config`, `monitoring_config`):**
    * Debug `flattenLoggingConfig` and `flattenMonitoringConfig` to ensure `enable_components` list is correctly placed within the nested `component_config` map.
    * Verify `MapToCtyValWithSchema` handles the `map[string][]interface{map[string]interface{}}` structure.
3.  **Fix Missing `security_posture_config`:**
    * Implement/debug `flattenSecurityPostureConfig` to return non-nil if `mode` or `vulnerability_mode` are non-default.
4.  **Fix Missing Node Pool `max_pods_per_node`:**
    * Review `flattenMaxPodsConstraint` to ensure it returns non-nil for positive values. Check call site in `convertNodePoolData`.
5.  **Fix Missing Node Pool `kubelet_config`:**
    * Debug `flattenKubeletConfig` to ensure it returns non-nil when non-default fields like `insecure_kubelet_readonly_port_enabled = true` are present. Handle string-to-bool conversion if needed.
6.  **Fix Incorrect Node Pool `oauth_scopes`:**
    * Update `flattenNodeConfig` to identify and omit the default OAuth scopes list (may require hardcoding the list or skipping omission if too complex).
7.  **Fix Cluster Versioning (`master_version`/`node_version`):**
    * Modify `convertClusterData`: If `remove_default_node_pool=true`, **do not** set cluster-level `master_version` or `node_version`.

### Medium Priority (Refine Output & Reduce Redundancy)

8.  **Omit Redundant Cluster `networking_mode`:**
    * Modify `convertClusterData` to skip setting `networking_mode` if `ip_allocation_policy` is present.
9.  **Omit Redundant Cluster `node_locations`:**
    * Modify `convertClusterData` to only add `node_locations` if non-empty and differs from the implicit default based on `location`.
10. **Omit Redundant `fleet` Block:**
    * Modify `flattenFleet` to return `nil` if `fleet.Project` matches cluster project and `fleet.PreRegistered` is false/default.
11. **Standardize `subnetwork` Format:**
    * Decide and implement consistently whether `convertClusterData` should output the short name or the full self-link for `subnetwork`. (Self-link recommended).