package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
)

func listAll(ctx context.Context, url string) ([]version, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	var tree struct {
		Entries []struct {
			Path string `json:"path"`
			Type string `json:"type"`
		} `json:"tree"`
		Truncated bool `json:"truncated"`
	}

	if err = json.NewDecoder(resp.Body).Decode(&tree); err != nil {
		return nil, fmt.Errorf("list: %w", err)
	}

	var vers versions

	for _, entry := range tree.Entries {
		if entry.Type != "tree" {
			continue
		}

		if entry.Path[:2] != "go" {
			continue
		}

		var ver version

		ver, err = parse(entry.Path)
		if err != nil {
			return nil, err
		}

		vers = append(vers, ver)
	}

	return vers, nil
}

func list(ctx context.Context, stdout io.Writer, filter string) error {
	var (
		min  version
		vers versions
		err  error
	)

	switch filter {
	case "latest":
		vers, err = latest(ctx, "https://go.dev")
	default:
		if filter != "" {
			min, err = parse(filter)
			if err != nil {
				return err
			}
		}

		vers, err = listAll(ctx, "https://api.github.com/repos/golang/dl/git/trees/master")
	}

	if err != nil {
		return err
	}

	sort.Sort(sort.Reverse(vers))

	for _, ver := range vers {
		if ver.Minor != min.Minor && min != (version{}) {
			continue
		}

		fmt.Fprintln(stdout, ver)
	}

	return nil
}

func listRun(ctx context.Context, stdout, _ io.Writer, filter string) error {
	return list(ctx, stdout, filter)
}
