package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	storageFile = ".localla_services.json"
	httpTimeout = 3 * time.Second
	commonPorts = "80,443,8000,8080,8443,3000,5000,9000"
)

type Device struct {
	IP    string `json:"ip"`
	MAC   string `json:"mac"`
	Ports []int  `json:"ports"`
}

type Service struct {
	IP       string `json:"ip"`
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
	Title    string `json:"title"`
}

type ScanResult struct {
	Timestamp string    `json:"timestamp"`
	Devices   []Device  `json:"devices"`
	Services  []Service `json:"services"`
}

func PrintHelp() {
	help := `
Localla - Local Ağ Tarama Aracı

Kullanım:
  localla <komut> [argümanlar]

Komutlar:
  scan          - Ağdaki tüm cihazları tara
  ports <IP>    - Belirli IP adresinde açık portları tara
  list          - Bulunmuş tüm servisleri listele
  help          - Bu yardım mesajını göster

Örnekler:
  localla scan              # Tüm ağı tara
  localla ports 192.168.1.1 # Belirli cihazı tara
  localla list              # Bulunan servisleri listele
`
	fmt.Println(help)
}

func ScanNetwork() {
	fmt.Println("🔍 Ağ taraması başlıyor...")

	devices, err := discoverDevices()
	if err != nil {
		fmt.Printf("❌ Hata: %v\n", err)
		return
	}

	if len(devices) == 0 {
		fmt.Println("⚠️  Cihaz bulunamadı")
		return
	}

	fmt.Printf("✅ %d cihaz bulundu:\n\n", len(devices))

	var services []Service
	for _, device := range devices {
		fmt.Printf("📱 %s", device.IP)
		if device.MAC != "" {
			fmt.Printf(" (%s)", device.MAC)
		}
		fmt.Println()

		discoveredServices := discoverServices(device.IP)
		if len(discoveredServices) > 0 {
			for _, svc := range discoveredServices {
				fmt.Printf("   🌐 %s://%s:%d", svc.Protocol, svc.IP, svc.Port)
				if svc.Title != "" {
					fmt.Printf(" - %s", svc.Title)
				}
				fmt.Println()
				services = append(services, svc)
			}
		}
	}

	result := ScanResult{
		Timestamp: time.Now().Format(time.RFC3339),
		Devices:   devices,
		Services:  services,
	}

	saveResults(result)
	fmt.Printf("\n💾 Sonuçlar %s dosyasına kaydedildi\n", storageFile)
}

func ScanPorts(ip string) {
	fmt.Printf("🔍 %s taraması başlıyor...\n", ip)

	portList := parsePortList(commonPorts)
	openPorts := scanPorts(ip, portList)

	if len(openPorts) == 0 {
		fmt.Println("⚠️  Açık port bulunamadı")
		return
	}

	fmt.Printf("✅ %d açık port bulundu:\n\n", len(openPorts))

	var services []Service
	for _, port := range openPorts {
		fmt.Printf("✓ Port %d açık\n", port)

		httpService := discoverService(ip, port, "http")
		httpsService := discoverService(ip, port, "https")

		if httpService != nil {
			fmt.Printf("  🌐 %s://%s:%d", httpService.Protocol, httpService.IP, httpService.Port)
			if httpService.Title != "" {
				fmt.Printf(" - %s", httpService.Title)
			}
			fmt.Println()
			services = append(services, *httpService)
		}

		if httpsService != nil {
			fmt.Printf("  🔒 %s://%s:%d", httpsService.Protocol, httpsService.IP, httpsService.Port)
			if httpsService.Title != "" {
				fmt.Printf(" - %s", httpsService.Title)
			}
			fmt.Println()
			services = append(services, *httpsService)
		}
	}

	result := ScanResult{
		Timestamp: time.Now().Format(time.RFC3339),
		Devices: []Device{
			{IP: ip, Ports: openPorts},
		},
		Services: services,
	}

	saveResults(result)
	fmt.Printf("\n💾 Sonuçlar %s dosyasına kaydedildi\n", storageFile)
}

func ListServices() {
	fmt.Println("📋 Bulunan servislerin listesi:")

	data, err := ioutil.ReadFile(storageFile)
	if err != nil {
		fmt.Println("⚠️  Kayıtlı servis bulunamadı. Önce 'scan' komutunu çalıştırın.")
		return
	}

	var result ScanResult
	if err := json.Unmarshal(data, &result); err != nil {
		fmt.Printf("❌ Hata: %v\n", err)
		return
	}

	if len(result.Services) == 0 {
		fmt.Println("⚠️  Servis bulunamadı")
		return
	}

	fmt.Printf("\n⏰ Tarama zamanı: %s\n\n", result.Timestamp)

	for _, svc := range result.Services {
		fmt.Printf("🌐 %s://%s:%d", svc.Protocol, svc.IP, svc.Port)
		if svc.Title != "" {
			fmt.Printf(" - %s", svc.Title)
		}
		fmt.Println()
	}

	jsonData, _ := json.MarshalIndent(result, "", "  ")
	fmt.Printf("\n📄 JSON çıktı:\n%s\n", string(jsonData))
}

func discoverDevices() ([]Device, error) {
	ifaceWithAddr, err := getDefaultInterface()
	if err != nil {
		return nil, err
	}

	network := getNetworkCIDR(ifaceWithAddr.IPNet)
	hosts := getHostsInNetwork(network)

	var devices []Device
	var mu sync.Mutex
	var wg sync.WaitGroup

	semaphore := make(chan struct{}, 10)

	for _, host := range hosts {
		wg.Add(1)
		go func(ip string) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			if isHostAlive(ip) {
				mac := getMAC(ip)
				mu.Lock()
				devices = append(devices, Device{
					IP:  ip,
					MAC: mac,
				})
				mu.Unlock()
			}
		}(host)
	}

	wg.Wait()
	return devices, nil
}

func scanPorts(ip string, ports []int) []int {
	var openPorts []int
	var mu sync.Mutex
	var wg sync.WaitGroup

	semaphore := make(chan struct{}, 20)

	for _, port := range ports {
		wg.Add(1)
		go func(p int) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			if isPortOpen(ip, p) {
				mu.Lock()
				openPorts = append(openPorts, p)
				mu.Unlock()
			}
		}(port)
	}

	wg.Wait()
	return openPorts
}

func discoverServices(ip string) []Service {
	ports := parsePortList(commonPorts)
	openPorts := scanPorts(ip, ports)

	var services []Service
	for _, port := range openPorts {
		if svc := discoverService(ip, port, "http"); svc != nil {
			services = append(services, *svc)
		}
		if svc := discoverService(ip, port, "https"); svc != nil {
			services = append(services, *svc)
		}
	}
	return services
}

func discoverService(ip string, port int, protocol string) *Service {
	scheme := "http"
	if protocol == "https" {
		scheme = "https"
	}

	url := fmt.Sprintf("%s://%s:%d", scheme, ip, port)

	client := &http.Client{
		Timeout: httpTimeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	title := extractTitle(resp)

	return &Service{
		IP:       ip,
		Port:     port,
		Protocol: protocol,
		Title:    title,
	}
}

func isPortOpen(ip string, port int) bool {
	address := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("tcp", address, 2*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func isHostAlive(ip string) bool {
	conn, err := net.DialTimeout("tcp", ip+":22", 1*time.Second)
	if err == nil {
		conn.Close()
		return true
	}

	conn, err = net.DialTimeout("tcp", ip+":80", 1*time.Second)
	if err == nil {
		conn.Close()
		return true
	}

	conn, err = net.DialTimeout("tcp", ip+":443", 1*time.Second)
	if err == nil {
		conn.Close()
		return true
	}

	return false
}

type InterfaceWithAddr struct {
	Interface *net.Interface
	IPNet     *net.IPNet
}

func getDefaultInterface() (*InterfaceWithAddr, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && ipnet.IP.To4() != nil {
				return &InterfaceWithAddr{
					Interface: &iface,
					IPNet:     ipnet,
				}, nil
			}
		}
	}

	return nil, fmt.Errorf("ağ arayüzü bulunamadı")
}

func getNetworkCIDR(ipnet *net.IPNet) string {
	return ipnet.Network()
}

func getHostsInNetwork(cidr string) []string {
	_, network, err := net.ParseCIDR(cidr)
	if err != nil {
		return []string{}
	}

	var hosts []string
	for ip := network.IP.Mask(network.Mask); network.Contains(ip); incrementIP(ip) {
		if ip.String() != network.IP.String() && !isNetworkBroadcast(ip, network) {
			hosts = append(hosts, ip.String())
		}
	}
	return hosts
}

func incrementIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func isNetworkBroadcast(ip net.IP, network *net.IPNet) bool {
	broadcastIP := make(net.IP, len(network.IP))
	copy(broadcastIP, network.IP)
	for i := range network.Mask {
		broadcastIP[i] = network.IP[i] | ^network.Mask[i]
	}
	return ip.Equal(broadcastIP)
}

func getMAC(ip string) string {
	netInterface, err := net.InterfaceByName("en0")
	if err != nil {
		netInterface, err = net.InterfaceByName("eth0")
		if err != nil {
			return ""
		}
	}
	return netInterface.HardwareAddr.String()
}

func extractTitle(resp *http.Response) string {
	if resp == nil || resp.Body == nil {
		return ""
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	content := string(body)
	start := strings.Index(strings.ToLower(content), "<title>")
	if start == -1 {
		return ""
	}

	start += 7
	end := strings.Index(strings.ToLower(content[start:]), "</title>")
	if end == -1 {
		return ""
	}

	title := strings.TrimSpace(content[start : start+end])
	if len(title) > 100 {
		title = title[:100] + "..."
	}
	return title
}

func parsePortList(portStr string) []int {
	var ports []int
	parts := strings.Split(portStr, ",")
	for _, part := range parts {
		var port int
		fmt.Sscanf(strings.TrimSpace(part), "%d", &port)
		if port > 0 && port < 65536 {
			ports = append(ports, port)
		}
	}
	return ports
}

func saveResults(result ScanResult) {
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Printf("❌ Kaydedilirken hata: %v\n", err)
		return
	}

	path, err := filepath.Abs(storageFile)
	if err != nil {
		path = storageFile
	}

	err = ioutil.WriteFile(path, data, 0644)
	if err != nil {
		fmt.Printf("❌ Dosya yazılırken hata: %v\n", err)
	}
}
