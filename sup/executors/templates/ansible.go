package templates

const (
	AnsibleHost = `{{.ID}} ansible_host={{.HostIP}} ansible_user={{.User}}
`
	CopyTask = `- name: copy {{.Name}} binary
  copy:
    src: {{.Source}}
    dest: {{.Destination}}
    force: no
    remote_src: yes
`
	SystemdTask = `- name: restart {{.ServiceName}}
  become_user: root
  systemd:
    name: {{.ServiceName}}
    state: {{.TargetState}}
    daemon_reload: yes
`
	LocalPlay = `- hosts: {{ if .Local }}local{{else}}remote{{end}}
  any_errors_fatal: true
  become: yes
  {{- if .Local }}
  connection: local
  {{- end}}
  serial: {{ if .Serial }}yes{{else}}no{{end}}
  tasks:
    - name: {{.Name}}
      block:
{{.RollForwardTasks}}
      rescue:
{{.RollbackTasks}}
      - fail:
          msg: failed to perform {{.Name}}
`
)
