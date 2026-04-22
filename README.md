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

### Belirli IP'de Portları Tara
```bash
./localla ports 192.168.1.1
```

### Bulunmuş Servisleri Listele
```bash
./localla list
```

### Yardım
```bash
./localla help
```

## 🎯 Özellikler (Planlanan)

- [x] Temel CLI yapısı
- [ ] ARP taraması ile ağ keşfi
- [ ] Port taraması
- [ ] HTTP/HTTPS servis keşfi
- [ ] JSON çıktı desteği
- [ ] Tarama sonuçlarını kaydetme

## 📝 Lisans

MIT

## 👨‍💻 Katkıda Bulunma

Pull request'ler kabul edilir!

---

**Geliştirici**: UmutKavil
