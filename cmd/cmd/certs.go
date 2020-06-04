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
		err := exec.Command("openssl", "genrsa", "-out", "certs/ca.key", "4096").Run()
		if err != nil {
			fmt.Printf("%s", err)
		}

		command := exec.Command("openssl", "req", "-new", "-x509", "-key", "certs/ca.key", "-out", "certs/ca.crt")
		command.Stderr = os.Stderr
		command.Stdin = os.Stdin
		command.Stdout = os.Stdout
		
		if err := command.Start(); nil != err {
			log.Fatalf("Error starting program: %s, %s", command.Path, err.Error())
		}
		command.Wait()

	case "listener":
		err := exec.Command("openssl", "genrsa", "-out", "certs/server.key", "4096").Run()
		if err != nil {
			fmt.Printf("%s", err)
		}

		command := exec.Command("openssl", "req", "-new", "-key", "certs/server.key", "-out", "certs/server.csr")
		command.Stderr = os.Stderr
		command.Stdin = os.Stdin
		command.Stdout = os.Stdout
		
		if err := command.Start(); nil != err {
			log.Fatalf("Error starting program: %s, %s", command.Path, err.Error())
		}
		command.Wait()

		err = exec.Command("openssl", "x509", "-req", "-in", "certs/server.csr", "-CA", "certs/ca.crt", "-CAkey", "certs/ca.key", "-CAcreateserial", "-out", "certs/server.crt").Run()
		if err != nil {
			fmt.Printf("%s", err)
		}

	case "peer":
		err := exec.Command("openssl", "genrsa", "-out", "certs/peers/"+name+".key", "4096").Run()
		if err != nil {
			fmt.Printf("%s", err)
		}

		command := exec.Command("openssl", "req", "-new", "-key", "certs/peers/"+name+".key", "-out", "certs/peers/"+name+".csr")
		command.Stderr = os.Stderr
		command.Stdin = os.Stdin
		command.Stdout = os.Stdout
		
		if err := command.Start(); nil != err {
			log.Fatalf("Error starting program: %s, %s", command.Path, err.Error())
		}
		command.Wait()

		err = exec.Command("openssl", "x509", "-req", "-in", "certs/peers/"+name+".csr", "-CA", "certs/ca.crt", "-CAkey", "certs/ca.key", "-CAcreateserial", "-out", "certs/peers/"+name+".crt").Run()
		if err != nil {
			fmt.Printf("%s", err)
		}
	}

  },
}