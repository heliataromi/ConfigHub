# Auto V2Ray & Clash Config Collector

[![Automatic Config Scraper](https://github.com/heliataromi/ConfigHub/actions/workflows/scraper.yml/badge.svg)](https://github.com/heliataromi/ConfigHub/actions/workflows/scraper.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

*Languages:* [English](#english) | [فارسی](#فارسی)

---

<a id="english"></a>
## English

An automated collector for V2Ray, Xray, Clash (Mihomo), and Sing-box configurations built in Go. It scrapes public Iranian Telegram channels and subscription links every few hours, validates and deduplicates the proxies, resolves GeoIP locations, and exports organized subscription links.

### Features
* **Protocols:** VLESS, VMess, Trojan, Shadowsocks, SSR, TUIC, Hysteria 2, AnyTLS, WireGuard, and Socks in standard and Base64 formats.
* **Validation:** Filters broken or unreachable configs and eliminates duplicates via node fingerprints.
* **GeoIP Tagging:** Resolves node hostnames/IPs to assign ISO country codes and flags.
* **Clash Meta & Sing-box:** Exports YAML and JSON profiles with `url-test` latency benchmarks, fallback routing, and Iranian direct traffic bypass.
* **Country Subscriptions:** Standalone subscription endpoints for 60+ countries in [COUNTRIES.md](COUNTRIES.md).

---

### Quick Guide

* **General / V2Ray & Xray:** Use [Mixed (All Protocols)](#v2ray--standard-subscriptions) in v2rayNG, v2rayN, INCY, or Karing.
* **iOS (iPhone / iPad):** Use [Mixed Base64](#v2ray--standard-subscriptions) or [Sing-box JSON](#sing-box-subscriptions) in INCY, Streisand, Karing, or V2Box.
* **Low-Memory Devices:** Use [Mixed Lite](#v2ray--standard-subscriptions).
* **Automated Benchmarking:** Use [Clash Meta](#clash-meta-mihomo-subscriptions) or [Sing-box](#sing-box-subscriptions).
* **Telegram:** Use [Telegram Proxies](#telegram-proxies).
* **Country Filters:** See [COUNTRIES.md](COUNTRIES.md).

---

### Quick Copy

**V2Ray Mixed (All Protocols):**
```text
https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/mixed.txt
```

**iOS Mixed (Base64):**
```text
https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/mixed_base64.txt
```

**V2Ray Mixed Lite:**
```text
https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/mixed_lite.txt
```

**Clash Meta (YAML):**
```text
https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/clash.yaml
```

**Sing-box (JSON):**
```text
https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/singbox.json
```

---

### Subscriptions

#### V2Ray / Standard Subscriptions

| Protocol | Normal Link (Standard) | Base64 Link (iOS / Legacy) |
| :--- | :---: | :---: |
| **Mixed (All)** | [Normal](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/mixed.txt) | [Base64](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/mixed_base64.txt) |
| **Mixed Lite** | [Normal](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/mixed_lite.txt) | [Base64](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/mixed_lite_base64.txt) |
| **VLESS** | [Normal](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/vless.txt) | [Base64](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/vless_base64.txt) |
| **VMess** | [Normal](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/vmess.txt) | [Base64](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/vmess_base64.txt) |
| **Trojan** | [Normal](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/trojan.txt) | [Base64](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/trojan_base64.txt) |
| **Shadowsocks** | [Normal](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/ss.txt) | [Base64](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/ss_base64.txt) |
| **ShadowsocksR** | [Normal](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/ssr.txt) | [Base64](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/ssr_base64.txt) |
| **TUIC** | [Normal](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/tuic.txt) | [Base64](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/tuic_base64.txt) |
| **Hysteria 2** | [Normal](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/hy2.txt) | [Base64](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/hy2_base64.txt) |
| **Hysteria** | [Normal](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/hysteria.txt) | [Base64](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/hysteria_base64.txt) |
| **AnyTLS** | [Normal](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/anytls.txt) | [Base64](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/anytls_base64.txt) |
| **WireGuard** | [Normal](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/wireguard.txt) | [Base64](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/wireguard_base64.txt) |
| **Socks** | [Normal](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/socks.txt) | [Base64](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/socks_base64.txt) |

Country-filtered subscriptions (Germany, US, UK, Netherlands, Turkey, etc.) are available in [COUNTRIES.md](COUNTRIES.md).

#### Telegram Proxies

Direct MTProto and SOCKS proxies for Telegram:

| File | Link |
| :--- | :--- |
| **Telegram Proxies** | [telegram.txt](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/telegram.txt) |

#### Clash Meta (Mihomo) Subscriptions

| Subscription | Raw URL |
| :--- | :--- |
| **Clash Meta (Mixed)** | `https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/clash.yaml` |
| **Clash Meta Lite** | `https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/clash_lite.yaml` |

*Import directly:* [Clash / FlClash](clash://install-config?url=https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/clash.yaml&name=ConfigHub) · [Stash](stash://install-config?url=https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/clash.yaml&name=ConfigHub)

#### Sing-box Subscriptions

| Subscription | Raw URL |
| :--- | :--- |
| **Sing-box (Mixed)** | `https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/singbox.json` |
| **Sing-box Lite** | `https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/singbox_lite.json` |

*Import directly:* [Sing-box / Karing](sing-box://import-remote-profile?url=https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/singbox.json#ConfigHub)

---

### Client Compatibility

| Client | Android | Windows | iOS | macOS | Linux |
| :--- | :---: | :---: | :---: | :---: | :---: |
| [**v2rayNG**](https://github.com/2dust/v2rayNG) *(V2Ray / Xray)* | ✅ | — | — | — | — |
| [**v2rayN**](https://github.com/2dust/v2rayN) *(V2Ray / Xray)* | — | ✅ | — | ✅ | ✅ |
| [**INCY**](https://github.com/Incy-App) *(Multi-Core)* | ✅ | ✅ | ✅ | ✅ | ✅ |
| [**Karing**](https://github.com/KaringX/karing) *(Multi-Core)* | ✅ | ✅ | ✅ | ✅ | ✅ |
| [**Streisand**](https://apps.apple.com/app/streisand/id6450534064) *(Sing-box)* | — | — | ✅ | ✅ | — |
| [**NekoBox**](https://github.com/MatsuriDayo/NekoBoxForAndroid) *(Sing-box / Xray)* | ✅ | ✅ | — | — | ✅ |
| [**Clash Verge Rev**](https://github.com/clash-verge-rev/clash-verge-rev) *(Clash Meta)* | — | ✅ | — | ✅ | ✅ |
| [**Mihomo Party**](https://github.com/mihomo-party-org/mihomo-party) *(Clash Meta)* | — | ✅ | — | ✅ | ✅ |
| [**FlClash**](https://github.com/chen08209/FlClash) *(Clash Meta)* | ✅ | ✅ | ✅ | ✅ | ✅ |
| [**Clash Meta (CMFA)**](https://github.com/MetaCubeX/ClashMetaForAndroid) *(Clash Meta)* | ✅ | — | — | — | — |
| [**Hiddify**](https://github.com/hiddify/hiddify-next) *(Sing-box)* | ✅ | ✅ | ✅ | ✅ | ✅ |
| [**V2Box**](https://apps.apple.com/app/v2box-v2ray-client/id6446814065) *(V2Ray)* | — | — | ✅ | ✅ | — |
| [**FoXray**](https://apps.apple.com/app/foxray/id6448898396) *(Xray)* | — | — | ✅ | ✅ | — |
| [**Stash**](https://apps.apple.com/app/stash/id1596063349) *(Clash)* | — | — | ✅ | ✅ | — |
| [**Shadowrocket**](https://apps.apple.com/app/shadowrocket/id932747118) *(Multi-Protocol)* | — | — | ✅ | ✅ | — |

---

### Setup Instructions

#### V2Ray & Xray (v2rayNG, v2rayN, INCY, NekoBox)
1. Copy the **Mixed** or protocol-specific URL from the table above.
2. In your client, navigate to `Subscription Group` / `Subscriptions` -> `Add`.
3. Paste the URL and save.
4. Select **Update Subscription**.

#### Clash Meta (Clash Verge Rev, Mihomo Party, FlClash)
1. Copy the **Clash Meta** URL (`clash.yaml` or `clash_lite.yaml`) or use the import link.
2. In your client, go to **Profiles** -> Paste URL -> **Import**.
3. Activate the profile.
4. Under **Proxies**, set `🔰 PROXY` to `⚡ AUTO (Fastest Node)`.

#### Sing-box (Sing-box, Karing, Hiddify)
1. Copy the **Sing-box** JSON URL (`singbox.json` or `singbox_lite.json`) or use the import link.
2. Go to **Profiles / Subscriptions** -> Import from URL.
3. Activate the profile and set the outbound selector to `⚡ AUTO (Fastest Node)`.

> In Hope of a Free Internet.

---

<a id="فارسی"></a>
## فارسی

<div dir="rtl">

ابزار جمع‌آوری و پالایش خودکار کانفیگ‌های V2Ray، Xray، کلش (Clash Meta) و Sing-box که با زبان Go نوشته شده است. این ابزار هر چند ساعت کانال‌های عمومی تلگرام و اشتراک‌های مختلف را بررسی کرده، کانفیگ‌های خراب و تکراری را حذف می‌کند، بر اساس IP موقعیت جغرافیایی و پرچم کشورها را مشخص کرده و لینک‌های اشتراک مرتب را خروجی می‌دهد.

### قابلیت‌ها
* **پروتکل‌ها:** پشتیبانی از VLESS، VMess، Trojan، Shadowsocks، SSR، TUIC، Hysteria 2، AnyTLS، WireGuard و Socks در دو فرمت معمولی و Base64.
* **اعتبارسنجی:** حذف سرورهای نامعتبر و قطع و پاکسازی موارد تکراری بر اساس اثر انگشت کانفیگ (Fingerprint).
* **تشخیص لوکیشن:** تفکیک IP و دامنه برای تعیین کد و پرچم کشورها.
* **کلش متا و سینگ‌باکس:** تولید پروفایل‌های YAML و JSON همراه با تست پینگ خودکار (`url-test`)، پشتیبان (`fallback`) و روتینگ مستقیم سایت‌های ایرانی.
* **اشتراک‌های کشوری:** لینک‌های تفکیک‌شده برای بیش از ۶۰ کشور در فایل [COUNTRIES.md](COUNTRIES.md).

---

### راهنمای انتخاب لینک

* **استفاده عمومی / V2Ray و Xray:** از [لینک ترکیبی (همه پروتکل‌ها)](#لینک‌های-اشتراک-v2ray-و-xray) در برنامه‌های INCY ،v2rayN ،v2rayNG یا Karing استفاده کنید.
* **کاربران iOS (آیفون و آیپد):** از [ترکیبی Base64](#لینک‌های-اشتراک-v2ray-و-xray) یا [پروفایل Sing-box](#لینک‌های-اشتراک-سینگ‌باکس) در برنامه‌های Karing ،Streisand ،INCY یا V2Box استفاده کنید.
* **دستگاه‌های ضعیف‌تر:** از [ترکیبی لایت (Lite)](#لینک‌های-اشتراک-v2ray-و-xray) استفاده کنید.
* **انتخاب خودکار کم‌پینگ‌ترین سرور:** از [کلش](#لینک‌های-اشتراک-کلش-متا-mihomo) یا [سینگ‌باکس](#لینک‌های-اشتراک-سینگ‌باکس) استفاده کنید.
* **تلگرام:** از [پروکسی‌های مستقیم تلگرام](#پروکسی‌های-تلگرام) استفاده کنید.
* **کانفیگ‌های یک کشور خاص:** به [COUNTRIES.md](COUNTRIES.md) مراجعه کنید.

---

### کپی سریع

**ترکیبی V2Ray (همه پروتکل‌ها):**
```text
https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/mixed.txt
```

**مخصوص آیفون (Base64):**
```text
https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/mixed_base64.txt
```

**ترکیبی لایت (کم‌حجم):**
```text
https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/mixed_lite.txt
```

**کلش متا (YAML):**
```text
https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/clash.yaml
```

**سینگ‌باکس (JSON):**
```text
https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/singbox.json
```

---

### لینک‌های اشتراک

#### لینک‌های اشتراک V2Ray و Xray

| پروتکل | لینک معمولی (Standard) | لینک Base64 (مخصوص iOS) |
| :--- | :---: | :---: |
| **ترکیبی (همه)** | [Normal](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/mixed.txt) | [Base64](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/mixed_base64.txt) |
| **ترکیبی لایت** | [Normal](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/mixed_lite.txt) | [Base64](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/mixed_lite_base64.txt) |
| **VLESS** | [Normal](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/vless.txt) | [Base64](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/vless_base64.txt) |
| **VMess** | [Normal](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/vmess.txt) | [Base64](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/vmess_base64.txt) |
| **Trojan** | [Normal](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/trojan.txt) | [Base64](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/trojan_base64.txt) |
| **Shadowsocks** | [Normal](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/ss.txt) | [Base64](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/ss_base64.txt) |
| **ShadowsocksR** | [Normal](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/ssr.txt) | [Base64](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/ssr_base64.txt) |
| **TUIC** | [Normal](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/tuic.txt) | [Base64](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/tuic_base64.txt) |
| **Hysteria 2** | [Normal](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/hy2.txt) | [Base64](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/hy2_base64.txt) |
| **Hysteria** | [Normal](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/hysteria.txt) | [Base64](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/hysteria_base64.txt) |
| **AnyTLS** | [Normal](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/anytls.txt) | [Base64](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/anytls_base64.txt) |
| **WireGuard** | [Normal](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/wireguard.txt) | [Base64](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/wireguard_base64.txt) |
| **Socks** | [Normal](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/socks.txt) | [Base64](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/socks_base64.txt) |

کانفیگ‌های تفکیک‌شده بر اساس کشور در [COUNTRIES.md](COUNTRIES.md) در دسترس هستند.

#### پروکسی‌های تلگرام

پروکسی‌های مستقیم MTProto و SOCKS برای تلگرام:

| فایل | لینک مستقیم |
| :--- | :--- |
| **پروکسی‌های تلگرام** | [telegram.txt](https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/telegram.txt) |

#### لینک‌های اشتراک کلش متا (Mihomo)

| نوع اشتراک | لینک مستقیم (Raw URL) |
| :--- | :--- |
| **کلش ترکیبی** | `https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/clash.yaml` |
| **کلش لایت** | `https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/clash_lite.yaml` |

*ورود مستقیم:* [Clash / FlClash](clash://install-config?url=https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/clash.yaml&name=ConfigHub) · [Stash](stash://install-config?url=https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/clash.yaml&name=ConfigHub)

#### لینک‌های اشتراک سینگ‌باکس

| نوع اشتراک | لینک مستقیم (Raw URL) |
| :--- | :--- |
| **سینگ‌باکس ترکیبی** | `https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/singbox.json` |
| **سینگ‌باکس لایت** | `https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/singbox_lite.json` |

*ورود مستقیم:* [Sing-box / Karing](sing-box://import-remote-profile?url=https://raw.githubusercontent.com/heliataromi/ConfigHub/subscription/singbox.json#ConfigHub)

---

### برنامه‌های پیشنهادی

| برنامه (Client) | اندروید | ویندوز | iOS | مک (macOS) | لینوکس |
| :--- | :---: | :---: | :---: | :---: | :---: |
| [**v2rayNG**](https://github.com/2dust/v2rayNG) *(V2Ray / Xray)* | ✅ | — | — | — | — |
| [**v2rayN**](https://github.com/2dust/v2rayN) *(V2Ray / Xray)* | — | ✅ | — | ✅ | ✅ |
| [**INCY**](https://github.com/Incy-App) *(Multi-Core)* | ✅ | ✅ | ✅ | ✅ | ✅ |
| [**Karing**](https://github.com/KaringX/karing) *(Multi-Core)* | ✅ | ✅ | ✅ | ✅ | ✅ |
| [**Streisand**](https://apps.apple.com/app/streisand/id6450534064) *(Sing-box)* | — | — | ✅ | ✅ | — |
| [**NekoBox**](https://github.com/MatsuriDayo/NekoBoxForAndroid) *(Sing-box / Xray)* | ✅ | ✅ | — | — | ✅ |
| [**Clash Verge Rev**](https://github.com/clash-verge-rev/clash-verge-rev) *(Clash Meta)* | — | ✅ | — | ✅ | ✅ |
| [**Mihomo Party**](https://github.com/mihomo-party-org/mihomo-party) *(Clash Meta)* | — | ✅ | — | ✅ | ✅ |
| [**FlClash**](https://github.com/chen08209/FlClash) *(Clash Meta)* | ✅ | ✅ | ✅ | ✅ | ✅ |
| [**Clash Meta (CMFA)**](https://github.com/MetaCubeX/ClashMetaForAndroid) *(Clash Meta)* | ✅ | — | — | — | — |
| [**Hiddify**](https://github.com/hiddify/hiddify-next) *(Sing-box)* | ✅ | ✅ | ✅ | ✅ | ✅ |
| [**V2Box**](https://apps.apple.com/app/v2box-v2ray-client/id6446814065) *(V2Ray)* | — | — | ✅ | ✅ | — |
| [**FoXray**](https://apps.apple.com/app/foxray/id6448898396) *(Xray)* | — | — | ✅ | ✅ | — |
| [**Stash**](https://apps.apple.com/app/stash/id1596063349) *(Clash)* | — | — | ✅ | ✅ | — |
| [**Shadowrocket**](https://apps.apple.com/app/shadowrocket/id932747118) *(چند پروتکله)* | — | — | ✅ | ✅ | — |

---

### راهنمای استفاده

#### برنامه‌های V2Ray و Xray (مانند Nekobox ،INCY ،v2rayN ،v2rayNG)
1. لینک اشتراک موردنظر را از جدول بالا کپی کنید.
2. در برنامه به بخش `Subscription Group` یا `Subscriptions` رفته و گزینهٔ افزودن (`+`) را بزنید.
3. لینک را در کادر `URL` جای‌گذاری و ذخیره کنید.
4. روی **Update Subscription** کلیک کنید تا سرورها بارگذاری شوند.

#### برنامه‌های کلش متا (مانند FlClash ،Mihomo Party ،Clash Verge Rev)
1. لینک `clash.yaml` را کپی کرده یا از لینک ورود مستقیم استفاده کنید.
2. در برنامه وارد بخش **Profiles** شده و لینک را در **Profile URL** وارد و **Import** کنید.
3. پروفایل را فعال کنید.
4. در بخش **Proxies**، گروه `🔰 PROXY` را روی `⚡ AUTO (Fastest Node)` قرار دهید.

#### برنامه‌های سینگ‌باکس (Hiddify ،Karing ،Sing-box)
1. لینک `singbox.json` را کپی کرده یا از لینک ورود مستقیم استفاده کنید.
2. در برنامه به بخش **Profiles / Subscriptions** رفته و گزینهٔ Import from URL را بزنید.
3. پروفایل را فعال کرده و خروجی را روی `⚡ AUTO (Fastest Node)` بگذارید.

> به امید اینترنت آزاد

</div>