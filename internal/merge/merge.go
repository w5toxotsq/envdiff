package merge

// MergeOptions controls how env maps are merged.
type MergeOptions struct {
	// Overwrite determines whether values from src overwrite existing keys in dst.
	Overwrite bool
}

// DefaultOptions returns sensible defaults for MergeOptions.
func DefaultOptions() MergeOptions {
	return MergeOptions{
		Overwrite: false,
	}
}

// Result holds the outcome of a merge operation.
type Result struct {
	// Merged is the combined env map.
	Merged map[string]string
	// Overwritten contains keys whose values were replaced from src.
	Overwritten []string
	// Added contains keys that were new from src.
	Added []string
	// Skipped contains keys from src that were skipped because Overwrite was false.
	Skipped []string
}

// Merge combines src into dst according to opts.
// dst is not modified; a new map is returned inside Result.
func Merge(dst, src map[string]string, opts MergeOptions) Result {
	result := Result{
		Merged: make(map[string]string, len(dst)),
	}

	// Copy dst into merged.
	for k, v := range dst {
		result.Merged[k] = v
	}

	for k, v := range src {
		if existing, exists := result.Merged[k]; exists {
			if opts.Overwrite && existing != v {
				result.Merged[k] = v
				result.Overwritten = append(result.Overwritten, k)
			} else if !opts.Overwrite {
				result.Skipped = append(result.Skipped, k)
			}
		} else {
			result.Merged[k] = v
			result.Added = append(result.Added, k)
		}
	}

	sortStrings(result.Overwritten)
	sortStrings(result.Added)
	sortStrings(result.Skipped)

	return result
}
