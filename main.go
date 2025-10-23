package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()
	case "exit":
		exit()
	case "0":
		exit()
	default:
		panic("bad command")
	}
}

func run() {
	fmt.Printf("Running... %v as %d\n", os.Args[2:], os.Getpid())
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS,
	}

	must(cmd.Run())
}
func child() {
	fmt.Printf("Running... %v as %d\n", os.Args[2:], os.Getpid())
	syscall.Sethostname([]byte("container"))
	syscall.Chroot("/vagrant/ubuntu-fs")
	syscall.Chdir("/")
	syscall.Mount("proc", "proc", "proc", 0, "")
	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	must(cmd.Run())
}

func must(err error) {
	if err == nil {
		return
	}
	panic(err)
}

func exit() {
	os.Exit(0)
}
