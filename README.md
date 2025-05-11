# AutoSRT Backend

AutoSRT, video ve ses dosyaları için otomatik altyazı oluşturma servisidir. Bu repo, projenin backend kısmını içerir.

## 🚀 Teknolojiler

- **Go (1.21+)**: Ana programlama dili
- **Gin**: HTTP web framework
- **MongoDB**: Ana veritabanı
- **DynamoDB**: Session yönetimi için
- **AWS S3**: Dosya depolama
- **AWS Transcribe**: Konuşma tanıma servisi
- **Paddle**: Ödeme sistemi entegrasyonu
- **JWT & Session**: Kimlik doğrulama
- **Resend**: E-posta servisi

## 🏗️ Mimari

Proje Clean Architecture prensiplerine göre tasarlanmıştır:

```
├── api
│   ├── http
│   │   └── delivery    # HTTP handlers
│   ├── middleware      # Middleware fonksiyonları
│   └── route          # Route tanımlamaları
├── bootstrap          # Uygulama başlangıç konfigürasyonları
├── domain            # İş mantığı arayüzleri ve modeller
├── repository        # Veritabanı işlemleri
├── usecase          # İş mantığı implementasyonları
└── utils            # Yardımcı fonksiyonlar
```

## 🔑 Özellikler

- 🔐 JWT ve Session tabanlı kimlik doğrulama
- 📝 Otomatik altyazı oluşturma
- 💳 Paddle ile abonelik sistemi
- 📧 E-posta bildirimleri
- 🌐 Çoklu dil desteği
- 🎥 Video ve ses dosyası işleme
- ⚡ Yüksek performanslı işlem kuyruğu

## 🛠️ Kurulum

1. Gereksinimleri yükleyin:
   ```bash
   go mod download
   ```

2. `.env` dosyasını oluşturun:
   ```env
   MONGODB_URI=your_mongodb_uri
   AWS_ACCESS_KEY=your_aws_access_key
   AWS_SECRET_KEY=your_aws_secret_key
   JWT_SECRET=your_jwt_secret
   PADDLE_API_KEY=your_paddle_api_key
   RESEND_API_KEY=your_resend_api_key
   ```

3. Uygulamayı başlatın:
   ```bash
   go run main.go
   ```

## 📝 API Endpoints

### Kimlik Doğrulama
- `POST /api/auth/register`: Kullanıcı kaydı
- `POST /api/auth/login`: Giriş
- `POST /api/auth/logout`: Çıkış
- `GET /api/auth/me`: Kullanıcı bilgileri

### Altyazı İşlemleri
- `POST /api/srt/create`: Altyazı oluşturma
- `GET /api/srt/list`: Altyazı listesi
- `GET /api/srt/{id}`: Altyazı detayları
- `DELETE /api/srt/{id}`: Altyazı silme

### Ödeme İşlemleri
- `POST /api/paddle/checkout`: Ödeme başlatma
- `POST /api/paddle/webhook`: Paddle webhook handler

## 🤝 Katkıda Bulunma

1. Bu repo'yu fork edin
2. Feature branch'i oluşturun (`git checkout -b feature/amazing-feature`)
3. Değişikliklerinizi commit edin (`git commit -m 'feat: add amazing feature'`)
4. Branch'inizi push edin (`git push origin feature/amazing-feature`)
5. Pull Request oluşturun

## 📄 Lisans

Bu proje MIT lisansı altında lisanslanmıştır. Detaylar için [LICENSE](LICENSE) dosyasına bakın.

## 📞 İletişim

Alper Çelik - [GitHub](https://github.com/kwa0x2)
