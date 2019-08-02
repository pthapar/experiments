package executors

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"
	"sup/executors/templates"
	"sup/models"
	"text/template"
	"unicode"
)

var (
	execCmd = func(binPath string, args ...string) error {
		cmd := exec.Command(binPath, args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		// TODO: add handling for cancelling a command if it takes longer than certain time
		return cmd.Run()
	}
)

// Executor defines the abstractions for different executors
type Executor interface {
	Run(plan *models.Plan, dryRun bool) error // executes the plan
	AddFileSystemCopyTask(planName, src, dest string) error
	AddSystemDTask(planName, serviceName, targetState string) error
	// DryRun(plan *models.Plan) error // dry runs an executor specific plan if it supports dry run
}

// Ansible implements executor interface for ansible
type Ansible struct {
	binPath               string
	workDir               string // where ansible can store its state, can be replaced with etcd in future, should come from config
	rollforward, rollback *bytes.Buffer
}

// SystemDUpdateTempl defines the plan for systemd
type SystemDUpdateTempl struct {
	// Name is the systemd service name
	Name string

	// Backup is the absoluate path of backup faile
	Backup string

	// Source is the source path of service binary
	Source string

	// Destination is the destination path of service binary
	Destination string
}

type AnsibleHostTempl struct {
	ID     string
	HostIP string
	User   string
}

type CopyTaskTempl struct {
	Name        string
	Source      string
	Destination string
}

type SystemDTaskTempl struct {
	ServiceName string
	TargetState string
}

type PlayTempl struct {
	User             string
	Name             string
	RollForwardTasks string
	RollbackTasks    string
	Local            bool
	Serial           bool
}

// IsBinaryChanged checks if the  binary has changed
// func (ans *Ansible) IsBinaryChanged(source, target string) bool {
// 	return true
// }

// ValidateBinary makes sure that the binary is good to be updated as per user given check
// func (ans *Ansible) ValidateBinary(cmd []string) error {
// 	return nil
// }

func NewAnsibleExecutor(binPath, workDir string) Executor {
	return &Ansible{
		binPath:     binPath,
		workDir:     workDir,
		rollback:    bytes.NewBuffer([]byte{}),
		rollforward: bytes.NewBuffer([]byte{}),
	}
}

func (ans *Ansible) generateHostfile(hosts []string, user string, hostsPath string) error {
	f, err := os.Create(hostsPath)
	if err != nil {
		return err
	}

	f.WriteString("[local]\n")

	t := template.New("exec template")
	t.Parse(templates.AnsibleHost)
	err = t.Execute(f, &AnsibleHostTempl{HostIP: hosts[0], ID: "local-1", User: user})
	if err != nil {
		return err
	}

	if len(hosts) > 1 {
		f.WriteString("\n[remote]\n")
		for i, h := range hosts[1:] {
			err = t.Execute(f, &AnsibleHostTempl{HostIP: h, ID: fmt.Sprintf("remote-%d", i+1), User: "admin"})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (ans *Ansible) planWorkDir(planName string) string {
	return path.Join(ans.workDir, strings.Join(strings.Split(planName, " "), "-"))
}

func (ans *Ansible) AddFileSystemCopyTask(name, src, dest string) error {
	t := template.New("copy task")
	_, err := t.Parse(templates.CopyTask)
	if err != nil {
		return err
	}
	backUpFile := path.Join(ans.planWorkDir(name), path.Base(dest)) + ".backup"
	// back up  file
	err = t.Execute(ans.rollforward, &CopyTaskTempl{Source: dest, Destination: backUpFile, Name: name})
	if err != nil {
		return err
	}

	err = t.Execute(ans.rollforward, &CopyTaskTempl{Source: src, Destination: dest, Name: name})
	if err != nil {
		return err
	}
	// copy backed up file
	return t.Execute(ans.rollback, &CopyTaskTempl{Source: backUpFile, Destination: dest, Name: name})
}

func (ans *Ansible) AddSystemDTask(planName, serviceName, targetState string) error {
	t := template.New("systemd task")
	_, err := t.Parse(templates.SystemdTask)
	if err != nil {
		return err
	}
	err = t.Execute(ans.rollforward, &SystemDTaskTempl{ServiceName: serviceName, TargetState: targetState})
	if err != nil {
		return err
	}
	return t.Execute(ans.rollback, &SystemDTaskTempl{ServiceName: serviceName, TargetState: targetState})
}

func (ans *Ansible) generateTask(name, user, rollforwardTasks, rollbackTasks string, local, serial bool, w io.Writer) error {
	t := template.New("task")
	_, err := t.Parse(templates.LocalPlay)
	if err != nil {
		return err
	}
	return t.Execute(w, &PlayTempl{User: user, Name: name,
		RollForwardTasks: rollforwardTasks, RollbackTasks: rollbackTasks, Local: local,
		Serial: serial,
	})
}

func indentEachLine(input string, n int) string {
	parts := strings.Split(input, "\n")
	rtn := bytes.NewBuffer(make([]byte, 0, len(input)))
	for _, p := range parts {
		for i := 0; i < n; i++ {
			rtn.WriteRune(rune(' '))
		}
		rtn.WriteString(p)
		rtn.WriteRune('\n')
	}
	return strings.TrimRightFunc(rtn.String(), unicode.IsSpace)
}

// ParseAndAddTasks parses step number #n from the list of steps
func ParseAndAddTasks(plan *models.Plan, exec Executor) error {
	for _, s := range plan.Steps {
		for module, op := range s {
			switch module {
			case "file":
				fileCopy := models.FileSystemCopyOp{}
				err := models.EncodeToObj(op.Operation, &fileCopy)
				if err != nil {
					return err
				}
				err = exec.AddFileSystemCopyTask(plan.Name, fileCopy.Copy.Src, fileCopy.Copy.Dst)
				if err != nil {
					return err
				}
			case "systemd":
				systemdOp := models.SystemDServiceRestartOp{}
				err := models.EncodeToObj(op.Operation, &systemdOp)
				if err != nil {
					return err
				}
				err = exec.AddSystemDTask(plan.Name, systemdOp.Restart.ServiceName, "restarted")
				if err != nil {
					return err
				}
			case "docker-images":
				dockerImageReloadOp := models.DockerImageReloadOp{}
				err := models.EncodeToObj(op.Operation, &dockerImageReloadOp)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// GeneratePlaybook generates ansible playbook and writes it to writedir/<plan-name>/<ver-num>/
func (ans *Ansible) GeneratePlaybook(plan *models.Plan) (string, string, error) {
	// rollforwardBuf, rollbackBuf := bytes.NewBuffer([]byte{}), bytes.NewBuffer([]byte{})
	var err error

	// if plan.SystemD == nil {
	// 	return "", "", fmt.Errorf("systemd missing")
	// }

	// if plan.SystemD.Copy != nil {
	// 	// TODO: add a copy task to perform backup
	// 	err = ans.generateCopyTask(plan.SystemD.Copy.Source, plan.SystemD.Copy.Target, plan.Name, rollforwardBuf)
	// 	err = ans.generateCopyTask(plan.SystemD.Copy.Target, plan.SystemD.Copy.Source, plan.Name, rollbackBuf)
	// }

	// if err != nil {
	// 	return "", "", err
	// }

	// err = ans.generateSystemDTask(plan.SystemD.ServiceName, "restarted", rollforwardBuf)
	// if err != nil {
	// 	return "", "", err
	// }

	// err = ans.generateSystemDTask(plan.SystemD.ServiceName, "restarted", rollbackBuf)
	// if err != nil {
	// 	return "", "", err
	// }
	planFSName := strings.Join(strings.Split(plan.Name, " "), "-")
	playBookPath := path.Join(ans.workDir, fmt.Sprintf("%s-plan.yaml", planFSName))
	f, err := os.Create(playBookPath)
	defer f.Close()
	if err != nil {
		return "", "", err
	}
	rollForwardIndented := indentEachLine(ans.rollforward.String(), 6)
	rollBackIndented := indentEachLine(ans.rollback.String(), 6)

	// generate local task
	err = ans.generateTask(plan.Name, plan.User, rollForwardIndented, rollBackIndented, true, true, f)
	if err != nil {
		return "", "", err
	}

	if len(plan.Target.Hosts) > 1 {
		// generate remote tasks
		err = ans.generateTask(plan.Name, plan.User, rollForwardIndented, rollBackIndented, false, true, f)
		if err != nil {
			return "", "", err
		}
	}

	// f, err := os.Create(path.Join(ans.workDir, fmt.Sprintf("%s-plan.yaml", plan.Name)))
	// defer f.Close()
	// if err != nil {
	// 	return "", "", err
	// }

	// f.WriteString(buf.String())

	// // write playbook to workdir under plan-name, prefer versioned
	// if plan.SystemD != nil {
	// 	// TODO: copy target to backup and remove it on successfull completion
	// 	systemDUpdateTempl := &SystemDUpdateTempl{Name: plan.SystemD.ServiceName,
	// 		Source: plan.SystemD.Copy.Source, Destination: plan.SystemD.Copy.Target,
	// 		Backup: "foo",
	// 	}

	// 	err := ans.generatePlaybook(systemDUpdateTempl, playBookPath)
	// 	if err != nil {
	// 		return playBookPath, hostsFilePath, err
	// 	}
	hostsFilePath := path.Join(ans.workDir, fmt.Sprintf("%s-hosts.yaml", planFSName))
	err = ans.generateHostfile(plan.Target.Hosts, plan.User, hostsFilePath)
	if err != nil {
		return "", "", err
	}
	return playBookPath, hostsFilePath, err
}

// Run executes an ansible plan
func (ans *Ansible) Run(plan *models.Plan, dryRun bool) error {
	err := ParseAndAddTasks(plan, ans)
	if err != nil {
		return err
	}

	playbook, hosts, err := ans.GeneratePlaybook(plan)
	if err != nil {
		return err
	}

	if !dryRun {
		return execCmd(ans.binPath, []string{"-i", hosts, playbook}...)
	}
	return nil
}
