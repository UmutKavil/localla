package cmd

import (
"fmt"
)

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
fmt.Println("[TODO] Network scanning fonksiyonu implement edilecek")
}

func ScanPorts(ip string) {
fmt.Printf("🔍 %s taraması başlıyor...\n", ip)
fmt.Println("[TODO] Port scanning fonksiyonu implement edilecek")
}

func ListServices() {
fmt.Println("📋 Bulunan servislerin listesi:")
fmt.Println("[TODO] Service listing fonksiyonu implement edilecek")
}
