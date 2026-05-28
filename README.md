# ⚡ Auto V2Ray Config Collector

[![Update Configs](https://github.com/heliataromi/ConfigHub/actions/workflows/scraper.yml/badge.svg)](https://github.com/heliataromi/ConfigHub/actions/workflows/scraper.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

*Read this document in other languages:*
🌐 **[English](#english)** | 🇮🇷 **[فارسی (Persian)](#فارسی)**

---

<a id="english"></a>
## 🌐 English

An automated V2Ray configuration scraper built in Go. This tool scrapes various Iranian Telegram channels every hour, parses the configurations, removes duplicates, resolves IPs to identify their geographical location, and outputs ready-to-use subscription links.

### ✨ Features
* **🔄 Hourly Updates:** Runs automatically every hour via GitHub Actions.
* **🧠 Smart Deduplication:** Parses actual URL parameters to eliminate duplicates, even if the query strings are scrambled or re-ordered.
* **🌍 GeoIP & Tunnel Recognition:** Automatically resolves IPs and domains to assign country flags (e.g., `🇩🇪 DE`).
* **🗂 Categorized & Encoded:** Generates separate files for each protocol (VLESS, VMess, Trojan, SS) in both Normal and Base64 formats to support all clients.

### 🔗 Subscription Links

Copy the link corresponding to your preferred protocol and import it into your V2Ray client.
*(Base64 links are strictly recommended for iOS users and older clients).*

| Protocol |                                  Normal Link (Standard)                                   |                                     Base64 Link (iOS / Legacy)                                     |
| :--- |:-----------------------------------------------------------------------------------------:|:--------------------------------------------------------------------------------------------------:|
| 🚀 **Mixed (All)** | [Normal](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/mixed.txt)  | [Base64](https://raw.githubusercontent.com/YOUR_USERNAME/YOUR_REPO/subscription/mixed_base64.txt)  |
| 🛡️ **VLESS** | [Normal](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/vless.txt)  | [Base64](https://raw.githubusercontent.com/YOUR_USERNAME/YOUR_REPO/subscription/vless_base64.txt)  |
| 📦 **VMess** | [Normal](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/vmess.txt)  | [Base64](https://raw.githubusercontent.com/YOUR_USERNAME/YOUR_REPO/subscription/vmess_base64.txt)  |
| 🐎 **Trojan** | [Normal](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/trojan.txt) | [Base64](https://raw.githubusercontent.com/YOUR_USERNAME/YOUR_REPO/subscription/trojan_base64.txt) |
| 🥷🏻 **Shadowsocks** |   [Normal](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/ss.txt)   |   [Base64](https://raw.githubusercontent.com/YOUR_USERNAME/YOUR_REPO/subscription/ss_base64.txt)   |

#### 📱 Recommended Clients
* **Android:** [v2rayNG](https://github.com/2dust/v2rayNG) | [NekoBox](https://github.com/MatsuriDayo/NekoBoxForAndroid)
* **Windows:** [v2rayN](https://github.com/2dust/v2rayN) | [NekoBox](https://github.com/qr243vbi/nekobox)
* **iOS:** Shadowrocket | V2Box | Streisand

#### 🛠️ How to use
1. Copy one of the subscription links from the table above.
2. Open your V2Ray client.
3. Go to `Subscription Group` -> `Add` (or `+` icon).
4. Paste the link in the `URL` field and save.
5. Click **Update Subscription** to fetch the latest configurations.

> In the hope of free internet.

---

<a id="فارسی"></a>
## 🇮🇷 فارسی (Persian)

یک ابزار خودکار و هوشمند برای جمع‌آوری کانفیگ‌های V2Ray که به زبان Go نوشته شده‌است. این ربات هر ساعت کانفیگ‌های جدید را از کانال‌های تلگرامی ایرانی استخراج کرده، تکراری‌ها را حذف می‌کند، لوکیشن سرورها را تشخیص می‌دهد و لینک‌های اشتراک (Subscription) آماده را برای شما تولید می‌کند.

### ✨ ویژگی‌ها
* **🔄 به‌روزرسانی ساعتی:** جمع‌آوری خودکار کانفیگ‌های جدید هر ۱ ساعت از طریق GitHub Actions.
* **🧠 حذف تکراری‌ها (Deduplication):** بررسی پارامترهای کانفیگ برای حذف کانفیگ‌های تکراری.
* **🌍 تشخیص لوکیشن:** استخراج IP و دامنه، و نمایش پرچم کشورها (مثلاً `🇩🇪 DE`).
* **🗂 دسته‌بندی و انکود Base64:** تفکیک کانفیگ‌ها بر اساس پروتکل (VLESS, VMess, Trojan, SS) و ارائهٔ نسخه Base64 برای سازگاری کامل با آیفون.

### 🔗 لینک‌های اشتراک (سابسکریپشن)

لینک مربوط به پروتکل دلخواه خود را کپی کرده و در برنامهٔ V2Ray خود وارد کنید.
*(استفاده از لینک‌های Base64 برای کاربران آیفون و برنامه‌های قدیمی‌تر اکیداً توصیه می‌شود).*

| پروتکل |                                   لینک معمولی (Standard)                                    |                                      لینک Base64 (مخصوص iOS)                                       |
| :--- |:-------------------------------------------------------------------------------------------:|:--------------------------------------------------------------------------------------------------:|
| 🚀 **ترکیبی (همه)** | [معمولی](https://raw.githubusercontent.com/YOUR_USERNAME/YOUR_REPO/subscription/mixed.txt)  | [Base64](https://raw.githubusercontent.com/YOUR_USERNAME/YOUR_REPO/subscription/mixed_base64.txt)  |
| 🛡️ **VLESS** | [معمولی](https://raw.githubusercontent.com/YOUR_USERNAME/YOUR_REPO/subscription/vless.txt)  | [Base64](https://raw.githubusercontent.com/YOUR_USERNAME/YOUR_REPO/subscription/vless_base64.txt)  |
| 📦 **VMess** | [معمولی](https://raw.githubusercontent.com/YOUR_USERNAME/YOUR_REPO/subscription/vmess.txt)  | [Base64](https://raw.githubusercontent.com/YOUR_USERNAME/YOUR_REPO/subscription/vmess_base64.txt)  |
| 🐎 **Trojan** | [معمولی](https://raw.githubusercontent.com/YOUR_USERNAME/YOUR_REPO/subscription/trojan.txt) | [Base64](https://raw.githubusercontent.com/YOUR_USERNAME/YOUR_REPO/subscription/trojan_base64.txt) |
| 🥷🏻 **Shadowsocks** |   [معمولی](https://raw.githubusercontent.com/YOUR_USERNAME/YOUR_REPO/subscription/ss.txt)   |   [Base64](https://raw.githubusercontent.com/YOUR_USERNAME/YOUR_REPO/subscription/ss_base64.txt)   |

#### 📱 برنامه‌های پیشنهادی
* **اندروید:** [v2rayNG](https://github.com/2dust/v2rayNG) | [NekoBox](https://github.com/MatsuriDayo/NekoBoxForAndroid)
* **ویندوز:** [v2rayN](https://github.com/2dust/v2rayN) | [NekoBox](https://github.com/qr243vbi/nekobox)
* **آیفون (iOS):** Shadowrocket | V2Box | Streisand

#### 🛠️ آموزش استفاده
1. یکی از لینک‌های جدول بالا را کپی کنید.
2. وارد نرم‌افزار خود (مثلاً v2rayNG) شوید.
3. از منوی کناری به بخش `Subscription Group` رفته و روی دکمه `+` کلیک کنید.
4. یک نام دلخواه بنویسید و لینک را در بخش `URL` جای‌گذاری (Paste) کنید و تیک تایید را بزنید.
5. در صفحهٔ اصلی برنامه، از منوی سه‌نقطه گزینهٔ **Update Subscription** را انتخاب کنید تا جدیدترین کانفیگ‌ها دریافت شوند.

> به امید اینترنت آزاد.