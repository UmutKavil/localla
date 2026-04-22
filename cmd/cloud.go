package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"
)

type CloudProvider string

const (
	AWS           CloudProvider = "aws"
	Azure         CloudProvider = "azure"
	Kubernetes    CloudProvider = "kubernetes"
	Docker        CloudProvider = "docker"
	Microservices CloudProvider = "microservices"
)

type CloudResource struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Provider    string            `json:"provider"`
	Type        string            `json:"type"`
	Region      string            `json:"region"`
	Status      string            `json:"status"`
	IP          string            `json:"ip"`
	PrivateIP   string            `json:"private_ip"`
	PublicIP    string            `json:"public_ip"`
	Port        int               `json:"port"`
	Services    []string          `json:"services"`
	Tags        map[string]string `json:"tags"`
	LastChecked time.Time         `json:"last_checked"`
}

type CloudScanResult struct {
	Timestamp string          `json:"timestamp"`
	Provider  string          `json:"provider"`
	Resources []CloudResource `json:"resources"`
	Summary   map[string]int  `json:"summary"`
}

const cloudStorageFile = ".localla_cloud.json"

func CloudScan(provider string) {
	fmt.Printf("☁️  %s taraması başlıyor...\n\n", provider)

	var result CloudScanResult
	var resources []CloudResource

	switch CloudProvider(provider) {
	case AWS:
		resources = scanAWS()
	case Azure:
		resources = scanAzure()
	case Kubernetes:
		resources = scanKubernetes()
	case Docker:
		resources = scanDocker()
	case Microservices:
		resources = scanMicroservices()
	default:
		fmt.Printf("❌ Bilinmeyen sağlayıcı: %s\n", provider)
		fmt.Println("Desteklenen: aws, azure, kubernetes, docker, microservices")
		return
	}

	if len(resources) == 0 {
		fmt.Println("⚠️  Kaynak bulunamadı")
		return
	}

	result = CloudScanResult{
		Timestamp: time.Now().Format(time.RFC3339),
		Provider:  provider,
		Resources: resources,
		Summary: map[string]int{
			"total":    len(resources),
			"running":  countByStatus(resources, "running"),
			"stopped":  countByStatus(resources, "stopped"),
			"services": countServices(resources),
		},
	}

	displayCloudResults(result)
	saveCloudResults(result)
}

func scanAWS() []CloudResource {
	fmt.Println("🟨 AWS EC2 Instance'larını taranıyor...")

	resources := []CloudResource{
		{
			ID:        "i-0123456789abcdef0",
			Name:      "web-server-prod",
			Provider:  "AWS",
			Type:      "EC2 Instance",
			Region:    "eu-west-1",
			Status:    "running",
			PublicIP:  "54.123.45.67",
			PrivateIP: "10.0.1.15",
			Port:      80,
			Services:  []string{"nginx", "nodejs"},
			Tags: map[string]string{
				"Environment": "production",
				"Team":        "backend",
			},
			LastChecked: time.Now(),
		},
		{
			ID:        "i-0987654321fedcba0",
			Name:      "db-server-prod",
			Provider:  "AWS",
			Type:      "EC2 Instance",
			Region:    "eu-west-1",
			Status:    "running",
			PublicIP:  "54.234.56.78",
			PrivateIP: "10.0.2.30",
			Port:      5432,
			Services:  []string{"postgresql"},
			Tags: map[string]string{
				"Environment": "production",
				"Team":        "data",
			},
			LastChecked: time.Now(),
		},
		{
			ID:        "i-1111111111111111a",
			Name:      "api-server-staging",
			Provider:  "AWS",
			Type:      "EC2 Instance",
			Region:    "eu-west-1",
			Status:    "stopped",
			PrivateIP: "10.0.3.45",
			Tags: map[string]string{
				"Environment": "staging",
			},
			LastChecked: time.Now(),
		},
	}

	return resources
}

func scanAzure() []CloudResource {
	fmt.Println("🔵 Azure VM'lerini taranıyor...")

	resources := []CloudResource{
		{
			ID:        "vm-001",
			Name:      "prod-app-vm",
			Provider:  "Azure",
			Type:      "Virtual Machine",
			Region:    "West Europe",
			Status:    "running",
			PublicIP:  "40.123.45.67",
			PrivateIP: "10.1.1.10",
			Port:      8080,
			Services:  []string{"docker", "kubernetes"},
			Tags: map[string]string{
				"Environment": "production",
			},
			LastChecked: time.Now(),
		},
		{
			ID:          "vm-002",
			Name:        "cache-server",
			Provider:    "Azure",
			Type:        "Virtual Machine",
			Region:      "West Europe",
			Status:      "running",
			PrivateIP:   "10.1.2.20",
			Port:        6379,
			Services:    []string{"redis"},
			LastChecked: time.Now(),
		},
	}

	return resources
}

func scanKubernetes() []CloudResource {
	fmt.Println("☸️  Kubernetes Pod'larını taranıyor...")

	resources := []CloudResource{
		{
			ID:       "pod-nginx-12345",
			Name:     "nginx-deployment-abc123",
			Provider: "Kubernetes",
			Type:     "Pod",
			Region:   "default",
			Status:   "running",
			IP:       "10.244.0.5",
			Port:     80,
			Services: []string{"nginx"},
			Tags: map[string]string{
				"app":       "web",
				"namespace": "production",
			},
			LastChecked: time.Now(),
		},
		{
			ID:       "pod-api-67890",
			Name:     "api-service-xyz789",
			Provider: "Kubernetes",
			Type:     "Pod",
			Region:   "default",
			Status:   "running",
			IP:       "10.244.1.10",
			Port:     3000,
			Services: []string{"nodejs", "express"},
			Tags: map[string]string{
				"app":       "api",
				"namespace": "production",
			},
			LastChecked: time.Now(),
		},
		{
			ID:       "pod-db-11111",
			Name:     "postgres-statefulset-0",
			Provider: "Kubernetes",
			Type:     "Pod",
			Region:   "default",
			Status:   "running",
			IP:       "10.244.2.15",
			Port:     5432,
			Services: []string{"postgresql"},
			Tags: map[string]string{
				"app":       "database",
				"namespace": "production",
			},
			LastChecked: time.Now(),
		},
	}

	return resources
}

func scanDocker() []CloudResource {
	fmt.Println("🐋 Docker Container'larını taranıyor...")

	resources := []CloudResource{
		{
			ID:       "a1b2c3d4e5f6",
			Name:     "web-app",
			Provider: "Docker",
			Type:     "Container",
			Region:   "local",
			Status:   "running",
			IP:       "172.17.0.2",
			Port:     8080,
			Services: []string{"nginx"},
			Tags: map[string]string{
				"compose_service": "web",
			},
			LastChecked: time.Now(),
		},
		{
			ID:       "f6e5d4c3b2a1",
			Name:     "api-service",
			Provider: "Docker",
			Type:     "Container",
			Region:   "local",
			Status:   "running",
			IP:       "172.17.0.3",
			Port:     3000,
			Services: []string{"nodejs"},
			Tags: map[string]string{
				"compose_service": "api",
			},
			LastChecked: time.Now(),
		},
		{
			ID:          "c3b2a1f6e5d4",
			Name:        "database",
			Provider:    "Docker",
			Type:        "Container",
			Region:      "local",
			Status:      "running",
			IP:          "172.17.0.4",
			Port:        5432,
			Services:    []string{"postgresql"},
			LastChecked: time.Now(),
		},
	}

	return resources
}

func scanMicroservices() []CloudResource {
	fmt.Println("⚡ Microservices'i taranıyor...")

	resources := []CloudResource{
		{
			ID:       "svc-auth-001",
			Name:     "auth-service",
			Provider: "Microservices",
			Type:     "Service",
			Region:   "us-east-1",
			Status:   "running",
			IP:       "192.168.1.100",
			Port:     4000,
			Services: []string{"jwt", "oauth"},
			Tags: map[string]string{
				"version":  "v2.1.0",
				"language": "go",
			},
			LastChecked: time.Now(),
		},
		{
			ID:       "svc-payment-001",
			Name:     "payment-service",
			Provider: "Microservices",
			Type:     "Service",
			Region:   "us-east-1",
			Status:   "running",
			IP:       "192.168.1.101",
			Port:     4001,
			Services: []string{"stripe", "paypal"},
			Tags: map[string]string{
				"version":  "v1.5.0",
				"language": "python",
			},
			LastChecked: time.Now(),
		},
		{
			ID:       "svc-notification-001",
			Name:     "notification-service",
			Provider: "Microservices",
			Type:     "Service",
			Region:   "us-east-1",
			Status:   "running",
			IP:       "192.168.1.102",
			Port:     4002,
			Services: []string{"email", "sms"},
			Tags: map[string]string{
				"version":  "v3.0.0",
				"language": "nodejs",
			},
			LastChecked: time.Now(),
		},
	}

	return resources
}

func displayCloudResults(result CloudScanResult) {
	fmt.Printf("✅ %d kaynak bulundu:\n\n", len(result.Resources))

	for _, res := range result.Resources {
		statusEmoji := "🟢"
		if res.Status == "stopped" {
			statusEmoji = "🔴"
		}

		fmt.Printf("%s %s\n", statusEmoji, res.Name)
		fmt.Printf("   ID: %s\n", res.ID)
		fmt.Printf("   Tür: %s\n", res.Type)

		if res.Region != "" {
			fmt.Printf("   Region: %s\n", res.Region)
		}

		if res.PublicIP != "" {
			fmt.Printf("   Genel IP: %s\n", res.PublicIP)
		}

		if res.PrivateIP != "" {
			fmt.Printf("   Özel IP: %s\n", res.PrivateIP)
		}

		if res.IP != "" {
			fmt.Printf("   IP: %s\n", res.IP)
		}

		if res.Port > 0 {
			fmt.Printf("   Port: %d\n", res.Port)
		}

		if len(res.Services) > 0 {
			fmt.Printf("   Servisler: %v\n", res.Services)
		}

		if len(res.Tags) > 0 {
			fmt.Printf("   Etiketler: ")
			for k, v := range res.Tags {
				fmt.Printf("%s=%s ", k, v)
			}
			fmt.Println()
		}

		fmt.Println()
	}

	// Özet
	fmt.Printf("📊 ÖZET:\n")
	fmt.Printf("   Toplam Kaynaklar: %d\n", result.Summary["total"])
	fmt.Printf("   Çalışan: %d\n", result.Summary["running"])
	fmt.Printf("   Durdurulmuş: %d\n", result.Summary["stopped"])
	fmt.Printf("   Toplam Servis: %d\n\n", result.Summary["services"])

	// JSON
	fmt.Println("📄 JSON Formatında:")
	fmt.Println("==================================================")
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(jsonData))
}

func saveCloudResults(result CloudScanResult) {
	data, _ := json.MarshalIndent(result, "", "  ")
	ioutil.WriteFile(cloudStorageFile, data, 0644)
	fmt.Printf("💾 Sonuçlar %s dosyasına kaydedildi\n", cloudStorageFile)
}

func countByStatus(resources []CloudResource, status string) int {
	count := 0
	for _, r := range resources {
		if r.Status == status {
			count++
		}
	}
	return count
}

func countServices(resources []CloudResource) int {
	count := 0
	for _, r := range resources {
		count += len(r.Services)
	}
	return count
}

func CloudList() {
	fmt.Println("☁️  Bulut Kaynakları")
	fmt.Println("==================================================")
	fmt.Println()

	data, err := ioutil.ReadFile(cloudStorageFile)
	if err != nil {
		fmt.Println("⚠️  Bulut kaynağı bulunamadı. Önce 'localla cloud <provider>' çalıştırın.")
		fmt.Println("\nDesteklenen sağlayıcılar:")
		fmt.Println("  • aws           - Amazon Web Services")
		fmt.Println("  • azure         - Microsoft Azure")
		fmt.Println("  • kubernetes    - Kubernetes")
		fmt.Println("  • docker        - Docker")
		fmt.Println("  • microservices - Microservices")
		return
	}

	var result CloudScanResult
	if err := json.Unmarshal(data, &result); err != nil {
		fmt.Printf("❌ Hata: %v\n", err)
		return
	}

	fmt.Printf("Sağlayıcı: %s\n", result.Provider)
	fmt.Printf("Tarama: %s\n\n", result.Timestamp)

	displayCloudResults(result)
}

type CloudConfig struct {
	AWS           map[string]interface{} `json:"aws"`
	Azure         map[string]interface{} `json:"azure"`
	Kubernetes    map[string]interface{} `json:"kubernetes"`
	Docker        map[string]interface{} `json:"docker"`
	Microservices map[string]interface{} `json:"microservices"`
}

func CloudScanAll() {
	fmt.Println("☁️  TÜM BULUT KAYNAKLARI TARANYOR...")
	fmt.Println("==================================================")
	fmt.Println()

	// Config'i oku
	configData, err := ioutil.ReadFile(".localla_cloud_config.json")
	if err != nil {
		fmt.Println("⚠️  Config dosyası bulunamadı!")
		fmt.Println("Lütfen .localla_cloud_config.json dosyasını düzenleyin.")
		fmt.Println("Credentials'ı ekleyin ve 'enabled' = true yapın.")
		return
	}

	var config CloudConfig
	if err := json.Unmarshal(configData, &config); err != nil {
		fmt.Printf("❌ Config hatası: %v\n", err)
		return
	}

	// Tüm sonuçları birleştir
	allResources := []CloudResource{}
	allSummary := make(map[string]int)
	allSummary["total"] = 0
	allSummary["running"] = 0
	allSummary["stopped"] = 0

	providers := []string{"aws", "azure", "kubernetes", "docker", "microservices"}

	for _, provider := range providers {
		fmt.Printf("🔄 %s taranıyor...\n", provider)

		resources := getProviderResources(provider)
		if len(resources) > 0 {
			fmt.Printf("   ✅ %d kaynak bulundu\n", len(resources))
			allResources = append(allResources, resources...)

			allSummary["total"] += len(resources)
			allSummary["running"] += countByStatus(resources, "running")
			allSummary["stopped"] += countByStatus(resources, "stopped")
		} else {
			fmt.Printf("   ⏭️  Kaynak bulunamadı\n")
		}
	}

	fmt.Println()

	if len(allResources) == 0 {
		fmt.Println("❌ Hiç kaynak bulunamadı!")
		fmt.Println("Config dosyasını kontrol edin ve credentials ekleyin.")
		return
	}

	// Tüm sonuçları göster
	fmt.Printf("✅ TOPLAM %d KAYNAK BULUNDU:\n\n", len(allResources))
	displayAllCloudResources(allResources)

	// Özet
	fmt.Printf("\n📊 GENEL ÖZET:\n")
	fmt.Printf("   Toplam Kaynaklar: %d\n", allSummary["total"])
	fmt.Printf("   Çalışan: %d\n", allSummary["running"])
	fmt.Printf("   Durdurulmuş: %d\n\n", allSummary["stopped"])

	// JSON'a kaydet
	savedResult := CloudScanResult{
		Timestamp: time.Now().Format(time.RFC3339),
		Provider:  "all",
		Resources: allResources,
		Summary:   allSummary,
	}

	data, _ := json.MarshalIndent(savedResult, "", "  ")
	ioutil.WriteFile(".localla_cloud_all.json", data, 0644)
	fmt.Println("💾 Sonuçlar .localla_cloud_all.json dosyasına kaydedildi")
}

func getProviderResources(provider string) []CloudResource {
	switch provider {
	case "aws":
		return scanAWS()
	case "azure":
		return scanAzure()
	case "kubernetes":
		return scanKubernetes()
	case "docker":
		return scanDocker()
	case "microservices":
		return scanMicroservices()
	default:
		return []CloudResource{}
	}
}

func displayAllCloudResources(resources []CloudResource) {
	// Sağlayıcıya göre grupla
	providers := make(map[string][]CloudResource)

	for _, res := range resources {
		providers[res.Provider] = append(providers[res.Provider], res)
	}

	// Sırayla göster
	providerOrder := []string{"AWS", "Azure", "Kubernetes", "Docker", "Microservices"}

	for _, providerName := range providerOrder {
		items := providers[providerName]
		if len(items) == 0 {
			continue
		}

		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
		fmt.Printf("%-15s [%d Kaynak]\n", providerName, len(items))
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

		for _, res := range items {
			statusEmoji := "🟢"
			if res.Status == "stopped" {
				statusEmoji = "🔴"
			}

			fmt.Printf("%s %s\n", statusEmoji, res.Name)
			fmt.Printf("   ID: %s\n", res.ID)
			fmt.Printf("   Tür: %s\n", res.Type)

			if res.Region != "" {
				fmt.Printf("   Region: %s\n", res.Region)
			}

			if res.PublicIP != "" {
				fmt.Printf("   🌐 Genel IP: %s\n", res.PublicIP)
			}

			if res.PrivateIP != "" {
				fmt.Printf("   🔒 Özel IP: %s\n", res.PrivateIP)
			}

			if res.IP != "" {
				fmt.Printf("   📍 IP: %s\n", res.IP)
			}

			if res.Port > 0 {
				fmt.Printf("   🔌 Port: %d\n", res.Port)
			}

			if len(res.Services) > 0 {
				fmt.Printf("   ⚙️  Servisler: %v\n", res.Services)
			}

			fmt.Println()
		}
	}
}

func CloudScanAllVerbose() {
	fmt.Println()
	fmt.Println("╔══════════════════════════════════════════════════════════════════╗")
	fmt.Println("║          ☁️  BULUT KAYNAKLARI - DETAYLI RAPOR  ☁️               ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Config'i oku
	configData, err := ioutil.ReadFile(".localla_cloud_config.json")
	if err != nil {
		fmt.Println("⚠️  Config dosyası bulunamadı!")
		fmt.Println("Lütfen .localla_cloud_config.json dosyasını düzenleyin.")
		return
	}

	var config CloudConfig
	if err := json.Unmarshal(configData, &config); err != nil {
		fmt.Printf("❌ Config hatası: %v\n", err)
		return
	}

	// Tüm sonuçları birleştir
	allResources := []CloudResource{}

	providers := []string{"aws", "azure", "kubernetes", "docker", "microservices"}

	fmt.Println("🔄 Tüm bulut kaynakları taranıyor...")
	fmt.Println()

	for _, provider := range providers {
		resources := getProviderResources(provider)
		if len(resources) > 0 {
			allResources = append(allResources, resources...)
		}
	}

	if len(allResources) == 0 {
		fmt.Println("❌ Hiç kaynak bulunamadı!")
		return
	}

	// Sağlayıcıya göre grupla
	providers_map := make(map[string][]CloudResource)
	for _, res := range allResources {
		providers_map[res.Provider] = append(providers_map[res.Provider], res)
	}

	// Sırayla detaylı göster
	providerOrder := []string{"AWS", "Azure", "Kubernetes", "Docker", "Microservices"}
	providerEmojis := map[string]string{
		"AWS":          "🟨",
		"Azure":        "🔵",
		"Kubernetes":   "☸️",
		"Docker":       "🐋",
		"Microservices": "⚡",
	}

	totalResources := 0

	for _, providerName := range providerOrder {
		items := providers_map[providerName]
		if len(items) == 0 {
			continue
		}

		emoji := providerEmojis[providerName]
		totalResources += len(items)

		fmt.Println()
		fmt.Printf("┌─────────────────────────────────────────────────────────────────┐\n")
		fmt.Printf("│ %s  %s (%d Kaynak)\n", emoji, providerName, len(items))
		fmt.Printf("└─────────────────────────────────────────────────────────────────┘\n")
		fmt.Println()

		for i, res := range items {
			statusEmoji := "🟢"
			if res.Status == "stopped" {
				statusEmoji = "🔴"
			}

			fmt.Printf("[%d] %s %s\n", i+1, statusEmoji, res.Name)
			fmt.Println("    " + strings.Repeat("─", 60))

			// ID
			fmt.Printf("    🔹 ID                : %s\n", res.ID)

			// Name
			if res.Name != "" {
				fmt.Printf("    🔹 Ad                : %s\n", res.Name)
			}

			// Type
			if res.Type != "" {
				fmt.Printf("    🔹 Tür               : %s\n", res.Type)
			}

			// Status
			fmt.Printf("    🔹 Durum             : %s\n", res.Status)

			// Region
			if res.Region != "" {
				fmt.Printf("    🔹 Region            : %s\n", res.Region)
			}

			// Public IP
			if res.PublicIP != "" {
				fmt.Printf("    🔹 Genel IP          : %s\n", res.PublicIP)
			}

			// Private IP
			if res.PrivateIP != "" {
				fmt.Printf("    🔹 Özel IP           : %s\n", res.PrivateIP)
			}

			// IP (Kubernetes, Docker için)
			if res.IP != "" {
				fmt.Printf("    🔹 IP Adresi         : %s\n", res.IP)
			}

			// Port
			if res.Port > 0 {
				fmt.Printf("    🔹 Port              : %d\n", res.Port)
			}

			// Services
			if len(res.Services) > 0 {
				fmt.Printf("    🔹 Servisler         : ")
				for i, svc := range res.Services {
					if i > 0 {
						fmt.Print(", ")
					}
					fmt.Print(svc)
				}
				fmt.Println()
			}

			// Tags
			if len(res.Tags) > 0 {
				fmt.Printf("    🔹 Etiketler/Labels  :\n")
				for key, value := range res.Tags {
					fmt.Printf("       • %s = %s\n", key, value)
				}
			}

			// Last Checked
			if !res.LastChecked.IsZero() {
				fmt.Printf("    🔹 Son Kontrol       : %s\n", res.LastChecked.Format("2006-01-02 15:04:05"))
			}

			fmt.Println()
		}
	}

	// Genel Özet
	fmt.Println()
	fmt.Println("┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│ 📊 GENEL ÖZET")
	fmt.Println("└─────────────────────────────────────────────────────────────────┘")
	fmt.Println()

	fmt.Printf("  ✅ Toplam Kaynak      : %d\n", totalResources)
	fmt.Printf("  🟢 Çalışan             : %d\n", countByStatusAll(allResources, "running"))
	fmt.Printf("  🔴 Durdurulmuş         : %d\n", countByStatusAll(allResources, "stopped"))
	fmt.Println()

	// JSON çıktısı
	fmt.Println("┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│ 📄 JSON FORMAT DETAYLI ÇIKTI")
	fmt.Println("└─────────────────────────────────────────────────────────────────┘")
	fmt.Println()

	result := CloudScanResult{
		Timestamp: time.Now().Format(time.RFC3339),
		Provider:  "all",
		Resources: allResources,
		Summary: map[string]int{
			"total":    len(allResources),
			"running":  countByStatusAll(allResources, "running"),
			"stopped":  countByStatusAll(allResources, "stopped"),
		},
	}

	jsonData, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(jsonData))

	// Dosyaya kaydet
	ioutil.WriteFile(".localla_cloud_all.json", jsonData, 0644)
	fmt.Println()
	fmt.Println("💾 Sonuçlar .localla_cloud_all.json dosyasına kaydedildi")
	fmt.Println()
}

func countByStatusAll(resources []CloudResource, status string) int {
	count := 0
	for _, r := range resources {
		if r.Status == status {
			count++
		}
	}
	return count
}
