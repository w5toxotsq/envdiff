package graph_test

import (
	"testing"

	"github.com/nicholasgasior/envdiff/internal/graph"
)

func makeKeys(keys ...string) map[string]string {
	m := make(map[string]string, len(keys))
	for _, k := range keys {
		m[k] = "value"
	}
	return m
}

func TestBuild_TwoEnvs_FullOverlap(t *testing.T) {
	envs := map[string]map[string]string{
		"dev":  makeKeys("A", "B", "C"),
		"prod": makeKeys("A", "B", "C"),
	}
	edges := graph.Build(envs)
	if len(edges) != 1 {
		t.Fatalf("expected 1 edge, got %d", len(edges))
	}
	e := edges[0]
	if e.SharedKeys != 3 {
		t.Errorf("shared keys: want 3, got %d", e.SharedKeys)
	}
	if e.Similarity != 1.0 {
		t.Errorf("similarity: want 1.0, got %f", e.Similarity)
	}
}

func TestBuild_TwoEnvs_NoOverlap(t *testing.T) {
	envs := map[string]map[string]string{
		"dev":  makeKeys("A", "B"),
		"prod": makeKeys("C", "D"),
	}
	edges := graph.Build(envs)
	if len(edges) != 1 {
		t.Fatalf("expected 1 edge, got %d", len(edges))
	}
	e := edges[0]
	if e.SharedKeys != 0 {
		t.Errorf("shared keys: want 0, got %d", e.SharedKeys)
	}
	if e.Similarity != 0.0 {
		t.Errorf("similarity: want 0.0, got %f", e.Similarity)
	}
	if e.TotalKeys != 4 {
		t.Errorf("total keys: want 4, got %d", e.TotalKeys)
	}
}

func TestBuild_ThreeEnvs_EdgeCount(t *testing.T) {
	envs := map[string]map[string]string{
		"dev":     makeKeys("A"),
		"staging": makeKeys("A", "B"),
		"prod":    makeKeys("A", "B", "C"),
	}
	edges := graph.Build(envs)
	// 3 environments → 3 pairs
	if len(edges) != 3 {
		t.Fatalf("expected 3 edges, got %d", len(edges))
	}
}

func TestBuild_EmptyEnvs_ZeroSimilarity(t *testing.T) {
	envs := map[string]map[string]string{
		"dev":  {},
		"prod": {},
	}
	edges := graph.Build(envs)
	if len(edges) != 1 {
		t.Fatalf("expected 1 edge, got %d", len(edges))
	}
	if edges[0].Similarity != 0.0 {
		t.Errorf("similarity: want 0.0, got %f", edges[0].Similarity)
	}
}

func TestBuild_SingleEnv_NoEdges(t *testing.T) {
	envs := map[string]map[string]string{
		"dev": makeKeys("A", "B"),
	}
	edges := graph.Build(envs)
	if len(edges) != 0 {
		t.Errorf("expected 0 edges for single env, got %d", len(edges))
	}
}

func TestBuild_EdgeNamesAreOrdered(t *testing.T) {
	envs := map[string]map[string]string{
		"prod": makeKeys("A"),
		"dev":  makeKeys("A"),
	}
	edges := graph.Build(envs)
	if edges[0].Left != "dev" || edges[0].Right != "prod" {
		t.Errorf("expected left=dev right=prod, got left=%s right=%s", edges[0].Left, edges[0].Right)
	}
}
