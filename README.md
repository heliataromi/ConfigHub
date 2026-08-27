# Auto V2Ray & Clash Config Collector

[![Update Configs](https://github.com/heliataromi/ConfigHub/actions/workflows/scraper.yml/badge.svg)](https://github.com/heliataromi/ConfigHub/actions/workflows/scraper.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

*Read this document in other languages:*
🇬🇧 **[English](#english)** | 🇮🇷 **[فارسی (Persian)](#فارسی)**

---

<a id="english"></a>
## English

An automated V2Ray and Clash (Mihomo) configuration scraper built in Go. This tool scrapes various Iranian Telegram channels and subscription links every couple of hours, parses the configurations, removes duplicates, resolves IPs to identify their geographical location, and outputs ready-to-use subscription links for both V2Ray and Clash Meta clients.

### Features
* **Hourly Updates:** Runs automatically every couple of hours via GitHub Actions.
* **Deduplication:** Parses actual URL parameters to eliminate duplicates.
* **GeoIP Recognition:** Automatically resolves IPs and domains to assign country codes and flags.
* **🚀 Clash Meta (Mihomo) Generation:** Full support for Clash Meta YAML with automated latency testing (`url-test`), failover (`fallback`), load balancing, clean Country Auto Groups, and Fake-IP DNS.
* **📦 Sing-box Generation:** Generates clean, modern JSON configuration profiles (`singbox.json` & `singbox_lite.json`) with auto latency testing (`urltest`), failover, load balancing, Fake-IP DNS, and Iranian domestic traffic bypass.
* **Categorized & Encoded:** Generates separate files for each protocol in both Normal and Base64 formats to support all clients.

### Subscription Links

Copy the link corresponding to your preferred protocol or client.
*(Base64 links are strictly recommended for iOS users and older clients).*

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

#### 🚀 Clash / Clash Meta (Mihomo) Subscriptions
> Designed with auto speed testing (`url-test`), failover (`fallback`), load balancing, and Country Auto Groups.

| Subscription | Format | Raw URL |
| :--- | :---: | :--- |
| **Clash Meta (All / Mixed)** | `YAML` | `https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/clash.yaml` |
| **Clash Meta Lite (Mobile)** | `YAML` | `https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/clash_lite.yaml` |

#### 📦 Sing-box Subscriptions
> Designed with automated URL testing (`urltest`), failover, load balancing, and Fake-IP split DNS.

| Subscription | Format | Raw URL |
| :--- | :---: | :--- |
| **Sing-box (All / Mixed)** | `JSON` | `https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/singbox.json` |
| **Sing-box Lite (Mobile)** | `JSON` | `https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/singbox_lite.json` |

#### Recommended Clients
* **Sing-box:** [Sing-box](https://github.com/SagerNet/sing-box) | [NekoBox](https://github.com/MatsuriDayo/NekoBoxForAndroid) | [Karing](https://github.com/KaringX/karing) | [Hiddify](https://github.com/hiddify/hiddify-next) | [Streisand](https://apps.apple.com/app/streisand/id6450534064)
* **Clash Meta / Mihomo:** [Clash Verge Rev](https://github.com/clash-verge-rev/clash-verge-rev) | [Mihomo Party](https://github.com/mihomo-party-org/mihomo-party) | [Clash Meta for Android](https://github.com/MetaCubeX/ClashMetaForAndroid) | [FlClash](https://github.com/chen08209/FlClash)
* **Android:** [v2rayNG](https://github.com/2dust/v2rayNG) | [NekoBox](https://github.com/MatsuriDayo/NekoBoxForAndroid)
* **Windows:** [v2rayN](https://github.com/2dust/v2rayN) | [NekoBox](https://github.com/qr243vbi/nekobox)
* **iOS:** [Stash](https://apps.apple.com/app/stash/id1596063349) | Shadowrocket | V2Box | Streisand

#### How to Use

##### 🌐 For V2Ray Clients (v2rayNG, v2rayN, NekoBox)
1. Copy one of the V2Ray subscription links from the table above.
2. Open your client (e.g. **v2rayNG**).
3. Go to `Subscription Group` -> `Add (+)` -> Paste the link in the `URL` field and save.
4. Click **Update Subscription** to fetch the latest configurations.

##### 🚀 For Clash Clients (Clash Verge Rev, Mihomo Party, FlClash)
1. Copy the **Clash Meta** URL (`clash.yaml` or `clash_lite.yaml`).
2. Open your client (e.g. **Clash Verge Rev**).
3. Navigate to **Profiles** -> Paste the subscription URL into the **Profile URL** box -> Click **Import**.
4. Select the downloaded profile to activate it.
5. In **Home / Proxies**, leave Group on `PROXY` and Proxy on **`⚡ AUTO (Fastest Node)`** for zero-maintenance automated speed testing and routing.

##### 📦 For Sing-box Clients (Sing-box, Karing, NekoBox, Hiddify)
1. Copy the **Sing-box** JSON URL (`singbox.json` or `singbox_lite.json`).
2. Open your client (e.g. **Karing** or **Sing-box**).
3. Go to **Profiles / Subscriptions** -> Add / Import from URL -> Paste the link.
4. Activate the profile. Under Outbounds / Selectors, choose **`⚡ AUTO (Fastest Node)`**.

> In Hope of a Free Internet.

---

<a id="فارسی"></a>
## فارسی (Persian)

یک ابزار خودکار برای جمع‌آوری کانفیگ‌های V2Ray و کلش (Clash Meta) که به زبان Go نوشته شده‌است. این ربات هر چند ساعت کانفیگ‌های جدید را از کانال‌های تلگرامی و اشتراک‌های ایرانی استخراج کرده، تکراری‌ها را حذف می‌کند، لوکیشن سرورها را تشخیص می‌دهد و لینک‌های اشتراک (Subscription) آماده را برای نرم‌افزارهای V2Ray و Clash تولید می‌کند.

### ویژگی‌ها
* **به‌روزرسانی ساعتی:** جمع‌آوری خودکار کانفیگ‌های جدید هر چند ساعت از طریق GitHub Actions.
* **حذف تکراری‌ها (Deduplication):** بررسی پارامترهای کانفیگ برای حذف کانفیگ‌های تکراری.
* **تشخیص لوکیشن:** استخراج IP و دامنه، و نمایش پرچم و کدهای کشورها.
* **🚀 پشتیبانی کامل از کلش (Clash Meta / Mihomo):** تولید خودکار کانفیگ کلش همراه با تست پینگ خودکار (`url-test`)، پشتیبان (`fallback`)، گروه‌های کشوری تفکیک‌شده و جلوگیری از شلوغی با پوشهٔ اختصاصی انتخاب دستی (`🎯 MANUAL`).
* **📦 پشتیبانی از سینگ‌باکس (Sing-box):** تولید پروفایل کامل و سبک JSON برای Sing-box به همراه تست پینگ خودکار (`urltest`)، فیلتر ترافیک مستقیم ایران و سیستم Fake-IP DNS.
* **دسته‌بندی و انکود Base64:** تفکیک کانفیگ‌ها بر اساس پروتکل و ارائهٔ نسخه Base64 برای سازگاری کامل با آیفون.

### لینک‌های اشتراک (سابسکریپشن)

لینک مربوط به پروتکل یا برنامهٔ دلخواه خود را کپی کرده و در برنامهٔ خود وارد کنید.
*(استفاده از لینک‌های Base64 برای کاربران آیفون و برنامه‌های قدیمی‌تر اکیداً توصیه می‌شود).*

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

#### 🚀 لینک‌های اشتراک مخصوص کلش (Clash Meta / Mihomo)
> دارای تست پینگ خودکار، انتخاب سریع‌ترین سرور و دسته‌بندی کشوری.

| نوع اشتراک | فرمت | لینک مستقیم (Raw URL) |
| :--- | :---: | :--- |
| **کلش ترکیبی (همه کانفیگ‌ها)** | `YAML` | `https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/clash.yaml` |
| **کلش لایت (مخصوص موبایل)** | `YAML` | `https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/clash_lite.yaml` |

#### 📦 لینک‌های اشتراک مخصوص سینگ‌باکس (Sing-box)
> دارای تست پینگ خودکار، انتخاب سریع‌ترین نود و روتینگ تفکیک‌شده ترافیک ایران.

| نوع اشتراک | فرمت | لینک مستقیم (Raw URL) |
| :--- | :---: | :--- |
| **سینگ‌باکس ترکیبی (همه کانفیگ‌ها)** | `JSON` | `https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/singbox.json` |
| **سینگ‌باکس لایت (مخصوص موبایل)** | `JSON` | `https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/singbox_lite.json` |

#### برنامه‌های پیشنهادی
* **سینگ‌باکس (Sing-box):** [Sing-box](https://github.com/SagerNet/sing-box) | [NekoBox](https://github.com/MatsuriDayo/NekoBoxForAndroid) | [Karing](https://github.com/KaringX/karing) | [Hiddify](https://github.com/hiddify/hiddify-next) | [Streisand](https://apps.apple.com/app/streisand/id6450534064)
* **کلش (Clash Meta / Mihomo):** [Clash Verge Rev](https://github.com/clash-verge-rev/clash-verge-rev) | [Mihomo Party](https://github.com/mihomo-party-org/mihomo-party) | [Clash Meta for Android](https://github.com/MetaCubeX/ClashMetaForAndroid) | [FlClash](https://github.com/chen08209/FlClash)
* **اندروید:** [v2rayNG](https://github.com/2dust/v2rayNG) | [NekoBox](https://github.com/MatsuriDayo/NekoBoxForAndroid)
* **ویندوز:** [v2rayN](https://github.com/2dust/v2rayN) | [NekoBox](https://github.com/qr243vbi/nekobox)
* **آیفون (iOS):** [Stash](https://apps.apple.com/app/stash/id1596063349) | Shadowrocket | V2Box | Streisand

#### آموزش استفاده

##### 🌐 مخصوص برنامه‌های V2Ray (مانند v2rayNG, v2rayN, NekoBox)
1. یکی از لینک‌های جدول V2Ray بالا را کپی کنید.
2. وارد نرم‌افزار خود (مثلاً **v2rayNG**) شوید.
3. از منوی کناری به بخش `Subscription Group` رفته و روی دکمه `+` کلیک کنید.
4. یک نام دلخواه بنویسید و لینک را در بخش `URL` جای‌گذاری (Paste) کنید و تیک تایید را بزنید.
5. در صفحهٔ اصلی برنامه، از منوی سه‌نقطه گزینهٔ **Update Subscription** را انتخاب کنید تا جدیدترین کانفیگ‌ها دریافت شوند.

##### 🚀 مخصوص برنامه‌های کلش (Clash Verge Rev, Mihomo Party, FlClash)
1. یکی از لینک‌های **Clash Meta** جدول بالا را کپی کنید.
2. وارد نرم‌افزار خود (مثلاً **Clash Verge Rev**) شوید.
3. از منوی کناری به بخش **Profiles** رفته و لینک را در کادر **Profile URL** قرار داده و **Import** را بزنید.
4. روی پروفایل جدید کلیک کنید تا انتخاب و فعال شود.
5. در بخش **Home** یا **Proxies**، گزینهٔ `PROXY` را روی **`⚡ AUTO (Fastest Node)`** بگذارید تا همیشه کم‌پینگ‌ترین و پایدارترین سرور به‌طور خودکار انتخاب شود.

##### 📦 مخصوص برنامه‌های سینگ‌باکس (Sing-box, Karing, NekoBox, Hiddify)
1. یکی از لینک‌های **Sing-box** جدول بالا را کپی کنید.
2. وارد نرم‌افزار خود (مثلاً **Karing** یا **Sing-box**) شوید.
3. به بخش **Profiles** یا **Subscriptions** رفته و لینک را از طریق گزینهٔ **Import from URL** اضافه کنید.
4. پروفایل را فعال کرده و گروه خروجی (Outbound) را روی **`⚡ AUTO (Fastest Node)`** بگذارید.

> به امید اینترنت آزاد.