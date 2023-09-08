package meters

type Process struct {
	Name string
	Pid  int
	Ppid int
}

func NewProcess(pid int) *Process {

	p := new(Process)

	p.Name = ""
	p.Pid = pid
	p.Ppid = 0

	return p
}
