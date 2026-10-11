# Traverse & Level

**Professional Traverse and Level Adjustment Software.** *Precision in every point.*

Enter the coordinates you measure on site, close on a known point and get an adjusted, checked
result with a signed PDF report. Level books are reduced, checked and adjusted the same way.
Every job is saved to your Excel to KML account, ready to open, copy or share.

<img src="app/icon-512.png" alt="" width="96">

[Download for Windows](https://github.com/Asjad-Mehmood/Traiverse/releases/latest/download/TraverseLevel-Setup.exe) ·
[Open the web app](https://www.exceltokml.com/traverse/app) ·
[Product page and manual](https://www.exceltokml.com/traverse) ·
[Report a problem](https://github.com/Asjad-Mehmood/Traiverse/issues)

---

## What it does

- **Coordinate traverse.** Station, backsight, the points measured on site and the closing point.
  Closed or open, typed or pasted from Excel. Bowditch, Transit or Equal adjustment, checked
  against your precision 1 : N.
- **Levelling.** Rise & Fall or Height of Instrument, the arithmetic checks, misclosure against
  your allowable and adjusted RLs.
- **PDF reports** with your report details and boxes for Surveyed, Checked and Approved by.
- **Dashboard, sharing and a traverse simulator** for training.
- **Offline, dark mode and phone.** After the first login it also opens without internet.

The full step-by-step manual is on the [product page](https://www.exceltokml.com/traverse).

## The Windows program

The program is a window onto the web app at
[www.exceltokml.com/traverse/app](https://www.exceltokml.com/traverse/app), in its own Start
menu entry and desktop shortcut, so the program and the website always match:

- **Log in with your Excel to KML account.** The same account as the website, and the same jobs.
- **Jobs are saved to your account** and appear in the web app and on your phone too.
- **Offline.** Open it once with internet. After that it opens without a connection; jobs saved
  offline upload by themselves when the connection is back.
- **Updates itself.** Each time it opens, it checks for a newer version of the program.

### Install on Windows 10 and 11

1. Download **`TraverseLevel-Setup.exe`** from the
   [latest release](https://github.com/Asjad-Mehmood/Traiverse/releases/latest).
2. Run it. It installs for your user only and does not need administrator rights.
3. If SmartScreen says *Windows protected your PC*, press **More info**, then **Run anyway**.
   This appears for any program without a paid signing certificate.

**Updates.** When a newer version exists, the program asks; press **Yes**. It downloads the
installer, verifies it and restarts on the new version.

**Coming from Traverse Pro?** Traverse Pro 1.x offers this update by itself and becomes Traverse &
Level, in the same folder. Jobs you saved as files are not touched.

**Uninstall** from *Settings → Apps*, like any other program.

**The program does not open, or opens a blank window.** Microsoft Edge WebView2 is missing or
damaged. Install it from
[microsoft.com/edge/webview2](https://developer.microsoft.com/microsoft-edge/webview2/), then
open the program again.

**No installer?** `Run-in-app-window.bat` opens the web app in an Edge or Chrome app window.

---

## For developers

### Project layout

```
lite/main.go              the Windows program (Go + WebView2): a window on the web app
lite/offline.html         shown only if the very first start has no internet
lite/update_windows.go    the self-updater (reads latest.json from the latest release)
lite/installer.nsi        the installer (NSIS), per user, no administrator rights
.github/workflows/        check.yml on every push; release.yml builds and publishes on a version tag
docs/releases/            release notes, one file per tag
Run-in-app-window.bat     opens the web app in an app window, no install
```

The app itself (traverse, levelling, reports, accounts, the jobs API) lives with the website, in
`Asjad-Mehmood/KML-WEBSITE-` (`traverse-app/` and `app/api/traverse/`). A change there reaches
the Windows program as soon as the website is deployed, with no new release.
A release here is only needed when the program itself changes: the window, the updater or the
installer.

### Making a release

The same as Excel to KML:

1. Merge the change into `Main`.
2. Put the notes in `docs/releases/<tag>.md`, then publish a GitHub release on `Main` with a tag
   higher than the last one, for example `v2.0.2`. A tag starts with `v` and has three numbers.
3. The **Build and Release** workflow then:
   - stamps the version into the program and the Windows file details;
   - builds `TraverseLevel.exe` and the installer `TraverseLevel-Setup.exe`;
   - writes `latest.json` with the version and the installer's SHA-256;
   - attaches the installer (also as `TraversePro-Setup.exe`, the name Traverse Pro 1.x updates
     from) and `latest.json` to the release, with `docs/releases/<tag>.md` as the description.

Installed copies read `latest.json` from the latest release of this repository. They update only
when its version is higher than their own, and only if the download matches the SHA-256.
**To roll back**, publish the older code again under a higher version number.

> **Releases must stay public.** This repository is public, so its releases are too, and the
> updater and the website read them at
> `github.com/Asjad-Mehmood/Traiverse/releases/latest/download/`. If the code is ever made
> private, do what Excel to KML does: publish every release to a separate public repository as
> well (a `RELEASES_TOKEN` secret and a second publish step), and point `updateBase` in
> `lite/update_windows.go` there in a release made while this one is still public.

### Building the Windows program by hand

You need [Go 1.22 or later](https://go.dev/dl/), and [NSIS](https://nsis.sourceforge.io/) for the
installer.

```bat
cd lite
go build -ldflags "-H=windowsgui -s -w -X main.appVersion=2.0.1" -o TraverseLevel.exe .
makensis /DVERSION=2.0.1 installer.nsi
```

On Linux or macOS, build the same program with `GOOS=windows GOARCH=amd64 go build ...`, and the
installer with `makensis` from your package manager.
