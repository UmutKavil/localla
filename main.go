package main

import (
	"fmt"
	"os"

	"github.com/UmutKavil/localla/cmd"
)

func main() {
	if len(os.Args) < 2 {
		cmd.PrintHelp()
		return
	}

	command := os.Args[1]

	switch command {
	case "scan":
		cmd.ScanNetwork()
	case "ports":
		if len(os.Args) < 3 {
			fmt.Println("Kullanım: localla ports <IP>")
			return
		}
		cmd.ScanPorts(os.Args[2])
	case "list":
		cmd.ListServices()
	case "history":
		cmd.ShowHistory()
	case "compare":
		cmd.CompareScans()
	case "cloud":
		if len(os.Args) < 3 {
			fmt.Println("Kullanım: localla cloud <provider>")
			fmt.Println("Sağlayıcılar: aws, azure, kubernetes, docker, microservices")
			fmt.Println("Tüm sağlayıcıları taramak için: localla cloud-scan-all veya localla bt")
			return
		}
		cmd.CloudScan(os.Args[2])
	case "cloud-scan-all":
		cmd.CloudScanAll()
	case "bt":
		cmd.CloudScanAllVerbose()
	case "cloud-list":
		cmd.CloudList()
	case "demo":
		cmd.DemoMode()
	case "help":
		cmd.PrintHelp()
	default:
		fmt.Printf("Bilinmeyen komut: %s\n", command)
		cmd.PrintHelp()
	}
}
