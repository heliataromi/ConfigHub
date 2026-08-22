# Auto V2Ray Config Collector

[![Update Configs](https://github.com/heliataromi/ConfigHub/actions/workflows/scraper.yml/badge.svg)](https://github.com/heliataromi/ConfigHub/actions/workflows/scraper.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

*Read this document in other languages:*
🇬🇧 **[English](#english)** | 🇮🇷 **[فارسی (Persian)](#فارسی)**

---

<a id="english"></a>
## English

An automated V2Ray configuration scraper built in Go. This tool scrapes various Iranian Telegram channels and subscription links every couple of hours, parses the configurations, removes duplicates, resolves IPs to identify their geographical location, and outputs ready-to-use subscription links.

### Features
* **Hourly Updates:** Runs automatically every couple of hours via GitHub Actions.
* **Deduplication:** Parses actual URL parameters to eliminate duplicates.
* **GeoIP Recognition:** Automatically resolves IPs and domains to assign country codes.
* **Categorized & Encoded:** Generates separate files for each protocol in both Normal and Base64 formats to support all clients.

### Subscription Links

Copy the link corresponding to your preferred protocol and import it into your V2Ray or Clash client.
*(Base64 links are strictly recommended for iOS users and older clients).*

#### 🚀 Clash / Clash Meta (Mihomo) Subscriptions
> Designed with auto speed testing (`url-test`), failover (`fallback`), load balancing, and Country Auto Groups.

| Type | Link |
| :--- | :---: |
| **Clash Meta (All)** | [Clash YAML](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/clash.yaml) |
| **Clash Meta Lite (Mobile)** | [Clash YAML](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/clash_lite.yaml) |

#### 🌐 V2Ray / Standard Subscriptions

| Protocol |                                  Normal Link (Standard)                                   |                                     Base64 Link (iOS / Legacy)                                     |
| :--- |:-----------------------------------------------------------------------------------------:|:--------------------------------------------------------------------------------------------------:|
| **Mixed (All)** | [Normal](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/mixed.txt)  | [Base64](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/mixed_base64.txt)  |
| **Mixed Lite (Mobile)** | [Normal](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/mixed_lite.txt)  | [Base64](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/mixed_lite_base64.txt)  |
| **VLESS** | [Normal](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/vless.txt)  | [Base64](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/vless_base64.txt)  |
| **VMess** | [Normal](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/vmess.txt)  | [Base64](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/vmess_base64.txt)  |
| **Trojan** | [Normal](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/trojan.txt) | [Base64](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/trojan_base64.txt) |
| **Shadowsocks** |   [Normal](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/ss.txt)   |   [Base64](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/ss_base64.txt)   |
| **ShadowsocksR**|   [Normal](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/ssr.txt)  |   [Base64](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/ssr_base64.txt)  |
| **TUIC** |   [Normal](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/tuic.txt)  |   [Base64](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/tuic_base64.txt)  |
| **Hysteria 2** |   [Normal](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/hy2.txt)  |   [Base64](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/hy2_base64.txt)  |
| **Hysteria** |   [Normal](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/hysteria.txt)|   [Base64](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/hysteria_base64.txt)|
| **Socks** |   [Normal](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/socks.txt) |   [Base64](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/socks_base64.txt) |
| **WireGuard** | [Normal](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/wireguard.txt)| [Base64](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/wireguard_base64.txt)|

#### Recommended Clients
* **Clash Meta / Mihomo:** [Clash Verge Rev](https://github.com/clash-verge-rev/clash-verge-rev) | [Mihomo Party](https://github.com/mihomo-party-org/mihomo-party) | [Clash Meta for Android](https://github.com/MetaCubeX/ClashMetaForAndroid) | [FlClash](https://github.com/chen08209/FlClash)
* **Android:** [v2rayNG](https://github.com/2dust/v2rayNG) | [NekoBox](https://github.com/MatsuriDayo/NekoBoxForAndroid)
* **Windows:** [v2rayN](https://github.com/2dust/v2rayN) | [NekoBox](https://github.com/qr243vbi/nekobox)
* **iOS:** [Stash](https://apps.apple.com/app/stash/id1596063349) | Shadowrocket | V2Box | Streisand

#### How to use
1. Copy one of the subscription links from the tables above.
2. Open your V2Ray or Clash client.
3. Import the subscription URL as a Remote Profile / Subscription.
4. Click **Update Subscription** to fetch the latest configurations.

> In Hope of a Free Internet.

---

<a id="فارسی"></a>
## فارسی (Persian)

یک ابزار خودکار برای جمع‌آوری کانفیگ‌های V2Ray که به زبان Go نوشته شده‌است. این ربات هر چند ساعت کانفیگ‌های جدید را از کانال‌های تلگرامی و اشتراک‌های ایرانی استخراج کرده، تکراری‌ها را حذف می‌کند، لوکیشن سرورها را تشخیص می‌دهد و لینک‌های اشتراک (Subscription) آماده (شامل V2Ray و Clash Meta) را برای شما تولید می‌کند.

### ویژگی‌ها
* **به‌روزرسانی ساعتی:** جمع‌آوری خودکار کانفیگ‌های جدید هر چند ساعت از طریق GitHub Actions.
* **حذف تکراری‌ها (Deduplication):** بررسی پارامترهای کانفیگ برای حذف کانفیگ‌های تکراری.
* **تشخیص لوکیشن:** استخراج IP و دامنه، و نمایش کدهای کشورها.
* **پشتیبانی از کلش (Clash Meta):** تولید خودکار کانفیگ کلش همراه با گروه‌های هوشمند تست پینگ خودکار (`url-test`)، پشتیبان (`fallback`) و گروه‌های تفکیک‌شده بر اساس کشور.
* **دسته‌بندی و انکود Base64:** تفکیک کانفیگ‌ها بر اساس پروتکل و ارائهٔ نسخه Base64 برای سازگاری کامل با آیفون.

### لینک‌های اشتراک (سابسکریپشن)

لینک مربوط به پروتکل یا برنامهٔ دلخواه خود را کپی کرده و در برنامهٔ خود وارد کنید.
*(استفاده از لینک‌های Base64 برای کاربران آیفون و برنامه‌های قدیمی‌تر اکیداً توصیه می‌شود).*

#### 🚀 لینک‌های اشتراک مخصوص کلش (Clash Meta / Mihomo)
> دارای تست پینگ خودکار، انتخاب سریع‌ترین سرور و دسته‌بندی کشوری.

| نوع | لینک اشتراک |
| :--- | :---: |
| **کلش ترکیبی (همه کانفیگ‌ها)** | [لینک کلش](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/clash.yaml) |
| **کلش لایت (مخصوص موبایل)** | [لینک کلش](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/clash_lite.yaml) |

#### 🌐 لینک‌های استاندارد V2Ray

| پروتکل |                                   لینک معمولی (Standard)                                    |                                      لینک Base64 (مخصوص iOS)                                       |
| :--- |:-------------------------------------------------------------------------------------------:|:--------------------------------------------------------------------------------------------------:|
| **ترکیبی (همه)** | [معمولی](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/mixed.txt)  | [Base64](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/mixed_base64.txt)  |
| **ترکیبی لایت (موبایل)** | [معمولی](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/mixed_lite.txt)  | [Base64](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/mixed_lite_base64.txt)  |
| **VLESS** | [معمولی](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/vless.txt)  | [Base64](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/vless_base64.txt)  |
| **VMess** | [معمولی](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/vmess.txt)  | [Base64](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/vmess_base64.txt)  |
| **Trojan** | [معمولی](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/trojan.txt) | [Base64](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/trojan_base64.txt) |
| **Shadowsocks** |   [معمولی](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/ss.txt)   |   [Base64](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/ss_base64.txt)   |
| **ShadowsocksR**|   [معمولی](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/ssr.txt)  |   [Base64](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/ssr_base64.txt)  |
| **TUIC** |   [معمولی](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/tuic.txt)  |   [Base64](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/tuic_base64.txt)  |
| **Hysteria 2** |   [معمولی](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/hy2.txt)  |   [Base64](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/hy2_base64.txt)  |
| **Hysteria** |   [معمولی](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/hysteria.txt)|   [Base64](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/hysteria_base64.txt)|
| **Socks** |   [معمولی](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/socks.txt) |   [Base64](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/socks_base64.txt) |
| **WireGuard** | [معمولی](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/wireguard.txt)| [Base64](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/wireguard_base64.txt)|

#### برنامه‌های پیشنهادی
* **کلش (Clash Meta / Mihomo):** [Clash Verge Rev](https://github.com/clash-verge-rev/clash-verge-rev) | [Mihomo Party](https://github.com/mihomo-party-org/mihomo-party) | [Clash Meta for Android](https://github.com/MetaCubeX/ClashMetaForAndroid) | [FlClash](https://github.com/chen08209/FlClash)
* **اندروید:** [v2rayNG](https://github.com/2dust/v2rayNG) | [NekoBox](https://github.com/MatsuriDayo/NekoBoxForAndroid)
* **ویندوز:** [v2rayN](https://github.com/2dust/v2rayN) | [NekoBox](https://github.com/qr243vbi/nekobox)
* **آیفون (iOS):** [Stash](https://apps.apple.com/app/stash/id1596063349) | Shadowrocket | V2Box | Streisand

#### آموزش استفاده
1. یکی از لینک‌های جدول بالا را کپی کنید.
2. وارد نرم‌افزار خود (مثلاً Clash Verge یا v2rayNG) شوید.
3. لینک را در بخش Profiles / Subscriptions وارد و ذخیره کنید.
4. روی دکمهٔ **Update** کلیک کنید تا جدیدترین کانفیگ‌ها دریافت شوند.

> به امید اینترنت آزاد.