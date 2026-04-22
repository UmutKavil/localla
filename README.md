# Localla 🌐

Local ağınızda çalışan web sitelerini ve cihazları keşfetmek için Go ile yazılmış bir CLI aracı.

## 📖 Açıklama

**Localla**, local ağınızda (LAN) bulunan cihazları tarar, açık portları tespit eder ve web servisleri keşfeder. Aynı ağda çalışan tüm servisleri kolayca görüntüleyebilirsiniz.

## ⚙️ Kurulum

### Gereksinimler
- Go 1.21 veya üzeri
- Linux, macOS veya Windows

### Derleme

```bash
git clone https://github.com/UmutKavil/localla.git
cd localla
go build -o localla
```

### Kurulum
```bash
go install
```

## 🚀 Kullanım

### Tüm Ağı Tara
```bash
./localla scan
```
Bu komut local ağınızdaki tüm cihazları bulur ve açık portlarında çalışan HTTP/HTTPS servislerini keşfeder.

**Örnek çıktı:**
```
🔍 Ağ taraması başlıyor...
✅ 5 cihaz bulundu:

📱 192.168.1.10
   🌐 http://192.168.1.10:80 - My Router
📱 192.168.1.20
   🌐 http://192.168.1.20:8080 - Apache Server

💾 Sonuçlar .localla_services.json dosyasına kaydedildi
```

### Belirli IP'de Portları Tara
```bash
./localla ports 192.168.1.1
```
Bu komut belirli bir IP adresindeki yaygın portları tarar.

**Örnek çıktı:**
```
🔍 192.168.1.1 taraması başlıyor...
✅ 2 açık port bulundu:

✓ Port 80 açık
  🌐 http://192.168.1.1:80 - Router Admin
✓ Port 443 açık
  🔒 https://192.168.1.1:443

💾 Sonuçlar .localla_services.json dosyasına kaydedildi
```

### Bulunmuş Servisleri Listele
```bash
./localla list
```
Kaydedilen tüm servisleri ve cihazları JSON formatında gösterir.

### Yardım
```bash
./localla help
```

## 🎯 Özellikler

- [x] Temel CLI yapısı
- [x] ARP taraması ile ağ keşfi
- [x] Port taraması
- [x] HTTP/HTTPS servis keşfi
- [x] JSON çıktı desteği
- [x] Tarama sonuçlarını kaydetme
- [x] Unit testler

## 📊 Çıktı Formatı

Tüm tarama sonuçları `.localla_services.json` dosyasına kaydedilir:

```json
{
  "timestamp": "2026-04-22T16:02:23+03:00",
  "devices": [
    {
      "ip": "192.168.1.10",
      "mac": "aa:bb:cc:dd:ee:ff",
      "ports": [80, 443, 8080]
    }
  ],
  "services": [
    {
      "ip": "192.168.1.10",
      "port": 80,
      "protocol": "http",
      "title": "My Router"
    }
  ]
}
```

## 🧪 Testler

Testleri çalıştırmak için:

```bash
go test ./cmd -v
```

Test kapsamı:
- Port listesi ayrıştırma
- Port açık olup olmadığını kontrolü
- IP adres işlemleri
- Veri yapıları (Device, Service, ScanResult)

## ⚙️ Teknik Detaylar

- **Paralel tarama**: Verimli için goroutine'ler kullanır (max 10-20 eşzamanlı bağlantı)
- **HTTP Timeout**: 3 saniye (hızlı tarama için)
- **Port Timeout**: 2 saniye (port kontrolü için)
- **Desteklenen Portlar**: 80, 443, 8000, 8080, 8443, 3000, 5000, 9000
- **JSON Desteği**: Tüm sonuçlar JSON olarak kaydedilir

MIT

## 👨‍💻 Katkıda Bulunma

Pull request'ler kabul edilir!

---

**Geliştirici**: UmutKavil
