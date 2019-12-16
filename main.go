package main

import (
	"fmt"
	"os"

	"github.com/serverlessresearch/srk/pkg/srkmgr"
)

func main() {
	mgrArgs := map[string]interface{}{}
	mgrArgs["config-file"] = "./srk.yaml"
	mgr, err := srkmgr.NewManager(mgrArgs)
	if err != nil {
		fmt.Printf("Failed to initialize: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Initialized SRK manager\n")

	if err = mgr.CreateRaw("./lambdas/echo/", "echo", nil); err != nil {
		fmt.Printf("Failed to create raw directory for lambda: %v\n", err)
		os.Exit(1)
	}
	rawDir := mgr.GetRawPath("echo")
	fmt.Printf("Created raw: %v\n", rawDir)

	pkgPath, err := mgr.Provider.Faas.Package(rawDir)
	if err != nil {
		fmt.Printf("Packaging failed: %v\n")
		os.Exit(1)
	}
	fmt.Printf("Function packaged to %v\n", pkgPath)

	if err := mgr.Provider.Faas.Install(rawDir); err != nil {
		fmt.Printf("Installation failed: %v\n")
		os.Exit(1)
	}
	fmt.Printf("Function installed\n")

	resp, err := mgr.Provider.Faas.Invoke("echo", "{\"hello\" : \"World\"}")
	if err != nil {
		fmt.Printf("Failed to invoke function: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Function Response:\n%v\n", resp)

	if err := mgr.Provider.Faas.Remove("echo"); err != nil {
		fmt.Printf("Service removal failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Successfully removed function \"echo\"")

	mgr.Destroy()
}
