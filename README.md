# Traverse Pro

**Close and adjust a survey traverse, on your own computer.**

List the field book, close the traverse on known control and adjust it by Compass (Bowditch),
Transit or least squares. The arithmetic runs entirely on your machine and works with no
internet connection.

<img src="app/icon-512.png" alt="" width="96">

[Download for Windows](https://github.com/Asjad-Mehmood/Traiverse/releases/latest) ·
[Report a problem](https://github.com/Asjad-Mehmood/Traiverse/issues)

---

## Quick start

1. **Set up.** Under *Start*, type the station you set up on and its Easting and Northing. Under
   *Backsight*, type the point you sighted back to and its coordinates, or a bearing to it.
2. **Choose the kind of traverse.** A *closed loop* comes back to its start station. A *closed
   link* ends on another known point: type it under *Close*. An *open* traverse has no check.
3. **List the field book.** Each row is one set-up: at **Station**, sighting back to
   **Backsight**, turn the **Angle** and measure the **Distance** to the new point, **Name**.
   The Station and Backsight of each row follow from the row above. Easting and Northing are
   worked out as you type.
4. **Close.** Under *Close*, type the closing angle: at the close point, from the last station
   onto the backsight (loop) or the foresight point (link). Choose *None* if it was not observed.
5. **Adjust.** Open **2 Close and adjust**, pick the adjustment, check the misclosures against
   their tolerance and press **Apply Adjustment**. The list then shows the adjusted coordinates.
   Your observations are never changed.
6. **Report and save.** **Generate Report** for a printable report (*Print or PDF*),
   **Coordinates CSV** for the adjusted points and **Save** (Ctrl+S) to keep the job.

New to the program? Press **Try the sample** for a levelled primary loop on two pillars.

## The list

| Column | |
|---|---|
| Side | *Traverse* for a leg of the traverse. *Side shot* for a point taken from the same set-up as the traverse row above it; side shots are not part of the adjustment and are worked out from the adjusted stations. |
| Station, Backsight | Where the instrument stood and what it sighted back to. Filled in for you. |
| Angle | The horizontal angle turned from the backsight. Type it as `DDD.MMSS` (`81.3127`), `81°31'27"`, `81 31 27`, decimal degrees or gons, as chosen under *Angles typed as*. Clockwise or anticlockwise under *Angles turned*. |
| Distance | Horizontal distance to the new point, in metres. Multiplied by the *Scale factor* for grid distances. |
| Ht diff | Height difference to the new point. Shown when *Carry levels* is ticked. |
| Easting, Northing, Level | Worked out. Green once an adjustment is applied. |
| Name, Description | The new point. In a closed traverse the last leg always ends on the close point. |

**Add point by click.** Press it, then click on the plan: each click adds a station there, with
its angle and distance worked back from the plan. In a closed traverse the new station goes in
before the closing leg, and the closing leg is measured on the plan too. Esc stops.

## Load traverse data

| File | |
|---|---|
| CSV, TXT (comma, tab, semicolon or space separated) | With a heading row, the columns are found by name: `Station, Backsight, Angle, Distance, Ht diff, Easting, Northing, Name, Description, Side`. Without one, the columns are read in the order `Station, Backsight, Angle, Distance, Easting, Northing, Name, Description`. A row with Easting and Northing but no angle or distance has them worked back from the coordinates. |
| JSON | A job saved with **Save**. |

Drag a file onto the window, or press **Load traverse data** (Ctrl+O).

## Close and adjust

| Adjustment | |
|---|---|
| Compass (Bowditch) | The misclosure is spread in proportion to the length of each leg. |
| Transit | The misclosure in Easting and Northing is spread in proportion to each leg's ΔE and ΔN. |
| Least squares | Angles and distances are weighted by their standard errors (σ angle, σ distance + ppm) and adjusted together. Shows residuals, standardised residuals (over 3 is flagged), a chi-square test and the 95% error ellipse of every station. Needs a closed traverse. |

The angular misclosure is spread equally over the observed angles before Compass and Transit.

**Checks**, against the tolerance of the class (*Primary*, *Secondary*, *Tertiary* or your own):

| Check | Allowed |
|---|---|
| Angular misclosure | *K* seconds × √(number of angles): 5″, 10″, 20″ |
| Linear precision | 1 : 20,000, 1 : 10,000, 1 : 5,000 |
| Level misclosure | *C* mm × √(km): 8, 12, 24 mm |

## Install

### Windows 10 and 11

1. Download **`TraversePro-Setup.exe`** from the
   [latest release](https://github.com/Asjad-Mehmood/Traiverse/releases/latest).
2. Run it. It installs for your user only and does not need administrator rights.
3. If SmartScreen says *Windows protected your PC*, press **More info**, then **Run anyway**.
   This appears for any program without a paid signing certificate.

**Updates.** Each time the program opens, it checks for a newer version. If one exists, press
**Yes**. It downloads the installer, verifies it and restarts on the new version.

**Uninstall** from *Settings → Apps*, like any other program.

**The program does not open, or opens a blank window.** Microsoft Edge WebView2 is missing or
damaged. Install it from
[microsoft.com/edge/webview2](https://developer.microsoft.com/microsoft-edge/webview2/), then
open the program again.

## Keyboard shortcuts

| Keys | Action |
|---|---|
| Ctrl+O | Load traverse data |
| Ctrl+S | Save the job |
| Esc | Stop adding points by click, close the report |

---

## For developers

### Project layout

```
app/index.html            the whole program in one self-contained HTML file
lite/                     the Windows program (Go + WebView2), installer and self-updater
tools/check-engine.js     checks the traverse engine against known answers (node)
.github/workflows/        check.yml on every push; release.yml builds and publishes on a version tag
docs/releases/            release notes, one file per tag
Run-in-app-window.bat     opens app/index.html in an app window, no install
```

Everything the program does lives in `app/index.html`, with no build step and no dependencies.
To try a change, open the file in a browser.

**The traverse engine** (`SURV`, between `ENGINE START` and `ENGINE END` in `app/index.html`) is
the same code as the traverse engine of Excel to KML (`kml-convertor-src`, `app/index.html`),
copied unchanged. Keep the two in step: a fix to the arithmetic in one belongs in the other.
`node tools/check-engine.js` runs it against traverses whose answers are known.

### Making a release

The same as Excel to KML:

1. Merge the change into `main`.
2. Publish a GitHub release on `main` with a tag higher than the last one, for example `v1.0.1`.
   A tag starts with `v` and has three numbers. Put the notes in `docs/releases/<tag>.md` first.
3. The **Build and Release** workflow then:
   - checks the traverse engine;
   - stamps the version into the program and the Windows file details;
   - builds `TraversePro.exe` and the installer `TraversePro-Setup.exe`;
   - writes `latest.json` with the version and the installer's SHA-256;
   - attaches the installer and `latest.json` to the release, with `docs/releases/<tag>.md` as
     the description when that file exists.

Installed copies read `latest.json` from the latest release of this repository. They update only
when its version is higher than their own, and only if the download matches the SHA-256.
**To roll back**, publish the older code again under a higher version number.

> **Releases must stay public.** This repository is public, so its releases are too, and the
> updater reads them at `github.com/Asjad-Mehmood/Traiverse/releases/latest/download/`. If the
> code is ever made private, do what Excel to KML does: publish every release to a separate
> public repository as well (a `RELEASES_TOKEN` secret and a second publish step), and point
> `updateBase` in `lite/update_windows.go` there in a release made while this one is still public.

### Building the Windows program by hand

You need [Go 1.22 or later](https://go.dev/dl/), and [NSIS](https://nsis.sourceforge.io/) for the
installer.

```bat
cd lite
copy ..\app\index.html index.html
go build -ldflags "-H=windowsgui -s -w -X main.appVersion=1.0.0" -o TraversePro.exe .
makensis /DVERSION=1.0.0 installer.nsi
```

The program serves the page to its own window at `http://127.0.0.1:47851`. It tries ports 47852
and 47853 if that port is busy, then any free port. Its WebView2 data, which holds the job you
are working on, is in `%AppData%\TraversePro`.

---

Developed by **Asjad Mehmood**.

## Licence

Copyright © 2026 Asjad Mehmood. All rights reserved. See [LICENSE](LICENSE).
