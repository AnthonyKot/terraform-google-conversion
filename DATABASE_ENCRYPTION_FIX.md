# Database Encryption Default Value Fix

## Issue

The converter currently omits the `database_encryption` block when it has the state "DECRYPTED" (which is the default value) or when it's not specified in the input. However, we need to include this block with state="DECRYPTED" in the output for better clarity and consistency, even when it's the default.

## Current Implementation

In `container_cluster_helpers.go`, the `flattenDatabaseEncryption` function is implemented as follows:

```go
func flattenDatabaseEncryption(config *container.DatabaseEncryption) []interface{} {
    if config == nil {
        return nil
    }
    // Keep logic to only include block if state is non-default (ENCRYPTED)
    if config.State == "" || config.State == "DECRYPTION_STATE_UNSPECIFIED" || config.State == "DECRYPTED" {
        return nil
    }

    de := make(map[string]interface{})
    de["state"] = config.State // Must be ENCRYPTED

    // Keep check for required key_name when encrypted
    if config.KeyName != "" {
        de["key_name"] = config.KeyName
    } else {
        fmt.Println("Warning: DatabaseEncryption state is ENCRYPTED but key_name is missing. Omitting block.")
        return nil
    }
    return []interface{}{de}
}
```

This function returns nil (thus omitting the block) if:
1. The config is nil
2. The state is empty, "DECRYPTION_STATE_UNSPECIFIED", or "DECRYPTED"
3. The state is ENCRYPTED but no key_name is provided

## Proposed Fix

We need to modify the function to always include the database_encryption block for better clarity, even with default values:

```go
func flattenDatabaseEncryption(config *container.DatabaseEncryption) []interface{} {
    // If config is nil, return a default config with state "DECRYPTED"
    if config == nil {
        de := make(map[string]interface{})
        de["state"] = "DECRYPTED"
        return []interface{}{de}
    }

    de := make(map[string]interface{})
    
    // Always set the state, default to "DECRYPTED" if empty or unspecified
    if config.State == "" || config.State == "DECRYPTION_STATE_UNSPECIFIED" {
        de["state"] = "DECRYPTED"
    } else {
        de["state"] = config.State
    }

    // Only include key_name if state is ENCRYPTED and key_name is provided
    if config.State == "ENCRYPTED" {
        if config.KeyName != "" {
            de["key_name"] = config.KeyName
        } else {
            fmt.Println("Warning: DatabaseEncryption state is ENCRYPTED but key_name is missing. Setting state to DECRYPTED.")
            de["state"] = "DECRYPTED"
        }
    }
    
    return []interface{}{de}
}
```

## Usage in Container Cluster Converter

In the `convertClusterData` function in `container_cluster.go`, we need to ensure this field is always included:

Replace:
```go
if flattened := flattenDatabaseEncryption(cluster.DatabaseEncryption); flattened != nil {
    hclData["database_encryption"] = flattened
}
```

With:
```go
// Always include database_encryption, using default if not specified
hclData["database_encryption"] = flattenDatabaseEncryption(cluster.DatabaseEncryption)
```

## Testing

After implementing this change, we should verify:

1. When input specifies `databaseEncryption` with state="DECRYPTED", it appears in the output
2. When input doesn't specify `databaseEncryption`, it appears in the output with state="DECRYPTED"
3. When input specifies `databaseEncryption` with state="ENCRYPTED" and a key_name, it appears correctly in the output
4. When input specifies `databaseEncryption` with state="ENCRYPTED" but no key_name, it falls back to state="DECRYPTED" with a warning

## Additional Considerations

1. This approach makes `database_encryption` with state="DECRYPTED" appear in all cluster configurations
2. If we decide this is too verbose, we could make this behavior configurable
3. We should consider applying the same pattern to other default values that should be made explicit in the output