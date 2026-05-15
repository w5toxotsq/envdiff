// Package baseline provides functionality for saving and loading a reference
// set of diff results so that subsequent runs can highlight only newly
// introduced differences.
//
// A baseline is stored as a JSON file containing the captured diff results
// along with metadata such as a timestamp and an optional label. Use Save to
// record the current state and Load to retrieve it. Compare can then be used
// to determine which differences are new since the baseline was taken.
package baseline
