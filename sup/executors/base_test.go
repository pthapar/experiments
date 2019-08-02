package executors

import (
	"os"
	"path"
	"sup/models"
	"testing"
)

func TestGeneratePlaybook(t *testing.T) {
	// ans := &Ansible{binPath: "/usr/local/bin/ansible", workDir: "testdir"}
	plan := &models.Plan{}
	models.Parse(path.Join("testdir", "input", "copy-systemd-plan.yaml"), plan)
	os.MkdirAll(path.Join("testdir", "ansible"), 0755)
	// defer os.RemoveAll(path.Join("testdir", "ansible"))
	exec := NewAnsibleExecutor("/usr/local/bin/ansible", "testdir/ansible")
	// err := ParseAndAddTasks(plan, exec)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	err := exec.Run(plan, true)
	if err != nil {
		t.Fatal(err)
	}

	// plan := &models.Plan{Name: "test systemd copy", Hosts: []string{"127.0.0.1", "10.40.6.77"}, User: "admin",
	// 	SystemD: &models.SystemD{ServiceName: "test_service", Copy: &models.Copy{Source: "foo", Target: "bar"}},
	// }

	// _, _, err := ans.GeneratePlaybook(plan)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// plan = &models.Plan{Name: "test systemd one host", Target: models.Target{Hosts: []string{"127.0.0.1"}}, User: "admin",
	// 	SystemD: &models.SystemD{ServiceName: "test_service", Copy: &models.Copy{Source: "foo", Target: "bar"}},
	// }

	// _, _, err = ans.GeneratePlaybook(plan)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// plan = &models.Plan{Name: "test systemd no copy", Hosts: []string{"127.0.0.1"}, User: "admin",
	// 	SystemD: &models.SystemD{ServiceName: "test_service"},
	// }

	// _, _, err = ans.GeneratePlaybook(plan)
	// if err != nil {
	// 	t.Fatal(err)
	// }
}
