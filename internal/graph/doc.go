// Package graph computes pairwise similarity between environment files based
// on their key overlap.
//
// Given a set of named environments and their key maps, Build returns a slice
// of Edge values — one for every unique pair — each carrying the count of
// shared keys, the size of the union, and a normalised similarity score in
// the range [0, 1].
//
// A score of 1.0 means both environments share exactly the same set of keys;
// a score of 0.0 means they have no keys in common.
//
// Typical usage:
//
//	envs := map[string]map[string]string{
//		"dev":  devKeys,
//		"prod": prodKeys,
//	}
//	edges := graph.Build(envs)
//	for _, e := range edges {
//		fmt.Printf("%s <-> %s: %.2f\n", e.Left, e.Right, e.Similarity)
//	}
package graph
