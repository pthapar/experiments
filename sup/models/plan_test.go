package models

import (
	"path"
	"reflect"
	"testing"
)

func TestPlanParse(t *testing.T) {
	p := Plan{}
	err := Parse(path.Join("testdir", "input", "file-copy-systemd.yaml"), &p)
	if err != nil {
		t.Fatal(err)
	}
	// t.Errorf("%v", p)
	if p.User != "admin" {
		t.Fatalf("expected user  to be admin, but  got %s", p.User)
	}

	if p.Name != "upgrade sherlock_configserver" {
		t.Fatalf("expected name to  be upgrade sherlock_configserver, but got  %s", p.Name)
	}

	if !reflect.DeepEqual(p.Target.Hosts, []string{"host1", "host2"}) {
		t.Fatalf("expected target hosts to  be %v, but got  %v", []string{"host1", "host2"}, p.Target.Hosts)
	}

	if p.Target.ExitOnFirstFailure != true {
		t.Fatalf("expected p.Target.ExitOnFirstFailure=%v, but got %v", true, p.Target.ExitOnFirstFailure)
	}

	modules := []string{}
	numFileCopyOps, numSystemDRestartOps, numDockerImageReloadOps := 0, 0, 0
	for _, s := range p.Steps {
		for module, op := range s {
			modules = append(modules, module)
			switch module {
			case "file":
				fileCopy := FileSystemCopyOp{}
				err := EncodeToObj(op.Operation, &fileCopy)
				if err != nil {
					t.Fatal(err)
				}
				numFileCopyOps++
			case "systemd":
				systemdOp := SystemDServiceRestartOp{}
				err := EncodeToObj(op.Operation, &systemdOp)
				if err != nil {
					t.Fatal(err)
				}
				numSystemDRestartOps++
			case "docker-images":
				dockerImageReloadOp := DockerImageReloadOp{}
				err := EncodeToObj(op.Operation, &dockerImageReloadOp)
				if err != nil {
					t.Fatal(err)
				}
				numDockerImageReloadOps++
			}
		}
	}

	if !reflect.DeepEqual(modules, []string{"file", "systemd"}) {
		t.Fatalf("expected modules %v, but  got %v", []string{"file", "systemd"}, modules)
	}

	if numDockerImageReloadOps != 0 {
		t.Fatalf("expected docker reload ops to be %d, but got %v", 0, numDockerImageReloadOps)
	}

	if numFileCopyOps != 1 {
		t.Fatalf("expected file copy ops to be %d, but got %v", 1, numFileCopyOps)
	}

	if numSystemDRestartOps != 1 {
		t.Fatalf("expected systemd restart ops to be %d, but got %v", 1, numSystemDRestartOps)
	}
}
