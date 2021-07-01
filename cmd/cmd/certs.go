package cmd

import (
  "fmt"
  "os"
  "log"
  "os/exec"
  "github.com/spf13/cobra"
)

func init() {
  rootCmd.AddCommand(certsCmd)
}

var certsCmd = &cobra.Command{
  Use:   "certs [ca listener peer]",
  Short: "Create certificates for server and peers",
  Long:  `Create certificates for server and peers`,
  Args: cobra.MinimumNArgs(1),
  Run: func(cmd *cobra.Command, args []string) {

	var name string
	if len(args) == 2 {
		name = args[1]
	}

	switch args[0] {
	case "ca":
		err := exec.Command("openssl", "genrsa", "-out", "/etc/guardmq/certs/ca.key", "4096").Run()
		if err != nil {
			fmt.Printf("%s", err)
		}

		command := exec.Command("openssl", "req", "-new", "-x509", "-key", "/etc/guardmq/certs/ca.key", "-out", "/etc/guardmq/certs/ca.crt")
		command.Stderr = os.Stderr
		command.Stdin = os.Stdin
		command.Stdout = os.Stdout
		
		if err := command.Start(); nil != err {
			log.Fatalf("Error starting program: %s, %s", command.Path, err.Error())
		}
		command.Wait()

	case "listener":
		err := exec.Command("openssl", "genrsa", "-out", "/etc/guardmq/certs/server.key", "4096").Run()
		if err != nil {
			fmt.Printf("%s", err)
		}

		command := exec.Command("openssl", "req", "-new", "-key", "/etc/guardmq/certs/server.key", "-out", "/etc/guardmq/certs/server.csr")
		command.Stderr = os.Stderr
		command.Stdin = os.Stdin
		command.Stdout = os.Stdout
		
		if err := command.Start(); nil != err {
			log.Fatalf("Error starting program: %s, %s", command.Path, err.Error())
		}
		command.Wait()

		err = exec.Command("openssl", "x509", "-req", "-days", "1095", "-in", "/etc/guardmq/certs/server.csr", "-CA", "/etc/guardmq/certs/ca.crt", "-CAkey", "/etc/guardmq/certs/ca.key", "-CAcreateserial", "-out", "/etc/guardmq/certs/server.crt").Run()
		if err != nil {
			fmt.Printf("%s", err)
		}

	case "peer":
		err := exec.Command("openssl", "genrsa", "-out", "/etc/guardmq/certs/peers/"+name+".key", "4096").Run()
		if err != nil {
			fmt.Printf("%s", err)
		}

		command := exec.Command("openssl", "req", "-new", "-key", "/etc/guardmq/certs/peers/"+name+".key", "-out", "/etc/guardmq/certs/peers/"+name+".csr")
		command.Stderr = os.Stderr
		command.Stdin = os.Stdin
		command.Stdout = os.Stdout
		
		if err := command.Start(); nil != err {
			log.Fatalf("Error starting program: %s, %s", command.Path, err.Error())
		}
		command.Wait()

		err = exec.Command("openssl", "x509", "-req", "-days", "1095", "-in", "/etc/guardmq/certs/peers/"+name+".csr", "-CA", "/etc/guardmq/certs/ca.crt", "-CAkey", "/etc/guardmq/certs/ca.key", "-CAcreateserial", "-out", "/etc/guardmq/certs/peers/"+name+".crt").Run()
		if err != nil {
			fmt.Printf("%s", err)
		}
	}

  },
}