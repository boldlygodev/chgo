package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go/build"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
)

func latest(ctx context.Context, host string) ([]version, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, host+"/dl/?mode=json", nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code was not 200 OK")
	}

	if resp.Header.Get("Content-Type") != "application/json" {
		return nil, fmt.Errorf("content-type was not application/json")
	}

	var list []struct {
		Version version
		Stable  bool
	}

	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		return nil, err
	}

	var vers versions

	for _, l := range list {
		if !l.Stable {
			continue
		}

		vers = append(vers, l.Version)
	}

	if len(vers) == 0 {
		return nil, fmt.Errorf("no stable versions were listed")
	}

	sort.Sort(sort.Reverse(vers))

	return vers, nil
}

func install(ctx context.Context, whichgo string, stdout, stderr io.Writer, ver version) error {
	cmd := exec.CommandContext(ctx, whichgo, "install", fmt.Sprintf("golang.org/dl/%s@latest", ver))
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	return cmd.Run()
}

func download(ctx context.Context, whichgo string, stdout, stderr io.Writer) error {
	cmd := exec.CommandContext(ctx, whichgo, "download")
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	return cmd.Run()
}

func symlink(ctx context.Context, homedir, gobin string, ver version) error {
	var (
		oldname = filepath.Join(gobin, ver.String())
		newname = filepath.Join(homedir, "bin", "go")
	)

	if err := os.Symlink(oldname, newname); err != nil {
		if !errors.Is(err, os.ErrExist) {
			return err
		}

		os.Remove(newname)

		return symlink(ctx, homedir, gobin, ver)
	}

	return nil
}

func chgo(ctx context.Context, stdout, stderr io.Writer, gobin, input string) error {
	var ver version

	switch input {
	case "":
		// nop
		return nil
	case "latest":
		vers, err := latest(ctx, "https://go.dev")
		if err != nil {
			return err
		}

		ver = vers[0]
	default:
		if input[:2] != "go" {
			input = "go" + input
		}

		v, err := parse(input)
		if err != nil {
			return fmt.Errorf("invalid version %q", input)
		}

		ver = v
	}

	// install Go version if not present
	if _, err := os.Stat(filepath.Join(gobin, ver.String())); err != nil {
		fmt.Fprintf(stdout, "go install golang.org/dl/%s@latest\n", ver)
		if err = install(ctx, "/usr/local/bin/go", stdout, stderr, ver); err != nil {
			return fmt.Errorf("install go %s: %w", input, err)
		}
	}

	// SDK is always downloaded into $HOME/sdk/$VERSION
	homedir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("homedir: %w", err)
	}

	// download Go version SDK if not present
	if _, err := os.Stat(filepath.Join(homedir, "sdk", input)); err != nil {
		fmt.Fprintf(stdout, "%s/%s download\n", gobin, ver)
		if err := download(ctx, filepath.Join(gobin, ver.String()), stdout, stderr); err != nil {
			return fmt.Errorf("download %s: %w", input, err)
		}
	}

	// update link
	fmt.Fprintf(stdout, "linking %s/%s as %s/bin/go\n", gobin, ver, homedir)
	if err := symlink(ctx, homedir, gobin, ver); err != nil {
		return fmt.Errorf("symlink %s: %w", input, err)
	}

	return nil
}

func chgoRun(ctx context.Context, stdout, stderr io.Writer, version string) error {
	gobin := os.ExpandEnv(os.Getenv("GOBIN"))
	if gobin == "" {
		gopath := os.ExpandEnv(os.Getenv("GOPATH"))

		if gopath == "" {
			gopath = build.Default.GOPATH
		}

		gobin = filepath.Join(gopath, "bin")
	}

	return chgo(ctx, stdout, stderr, gobin, version)
}
