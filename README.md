# Library Management System

Go ile geliştirilmiş kütüphane yönetim sistemi.

## Özellikler

- Kullanıcı girişi
- Kitap yönetimi
- Üye yönetimi
- Ödünç kitap işlemleri
- Dashboard ekranı

## Kurulum

Projeyi indirdikten sonra bağımlılıkları yüklemek için:

```bash
go mod download

#çalıştırmak için
go run main.go
#Uygulama çalıştıktan sonra tarayıcıdan aç:
http://localhost:8080

# Proje Yapısı
# auth/: Kimlik doğrulama işlemleri
# database/: Veritabanı bağlantısı
# handlers/: Sayfa ve işlem kontrolleri
# models/: Veri modelleri
# templates/: HTML sayfaları
# static/: CSS ve statik dosyalar
# main.go: Uygulama başlangıç dosyası