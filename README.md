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

1. **Set up.** In the panel at the top, type the **Station** you start on, its **Backsight**,
   the **Close point** and the **Start Easting** and **Start Northing**.
2. **The close point decides the kind of traverse.** The start station again: a *closed loop*.
   Another known point: a *closed link* (give its coordinates in **Settings**). Empty: an *open*
   traverse, with no check.
3. **List the points.** Each row is one set-up: at **Station**, sighting back to **Backsight**,
   the **Angle** and **Distance** to the new point, **Name**. Station and Backsight follow from the
   row above. Easting and Northing are worked out as you type; type an Easting and Northing
   instead and the angle and distance are worked back from them.
4. **Close.** Open **2 Close and adjust**, pick **Compass**, **Transit** or **Least squares** and
   press **Apply Adjustment**. The list then shows the adjusted coordinates. Your observations are
   never changed. **Generate Report** gives a printable report, a PDF and the adjusted
   coordinates as CSV.
5. **Save** (Ctrl+S) keeps the job as a file.

**Add point by click** turns on adding: each click on the plan adds a point there, with its
angle and distance worked back from the plan, and a dashed line shows the leg before you click.
In a closed traverse the new point goes in before the closing leg, and the closing leg is
measured on the plan too (you are asked first if it holds observed values). Esc or the button
stops. **Add row** adds an empty row; the buttons at the end of each row insert a row above it or
remove it.

**The plan.** Scroll to zoom, drag to move, double-click or **Fit** to see it all. Point at a
station for its coordinates; click it to find its row. The selected row's leg is highlighted.

**Undo and redo** (Ctrl+Z, Ctrl+Y, or the arrows in the title bar) take back any change: a
typed value, a removed row, a loaded file, an applied adjustment.

**Step 2** shows the result against the tolerance first, with the reason when it fails, then the
precision, error distance and direction, angular and level misclosure, area and perimeter, then
the adjusted values. Switch between **Compass**, **Transit** and **Least squares** at the top.

## Settings

| | |
|---|---|
| The Angle column holds | **Bearings** (whole-circle bearing of each line, the default) or **Angles turned from the backsight**, clockwise or anticlockwise, as a field book has them. With bearings there is no angular check: the misclosure is the linear one only. |
| Typed as | Decimal degrees, `DDD.MMSS` (`81.3127` is 81°31'27"), or gons. `81°31'27"` and `81 31 27` are read too. |
| Scale factor | Ground distances times this give grid distances. |
| Backsight | For angles: the backsight point's coordinates, or a bearing to it. |
| Close | For a closed link: the close point's coordinates. For angles: the closing orientation and the closing angle at the close point. |
| Tolerance | Primary, Secondary, Tertiary or your own: angular *K*″ × √n, linear precision 1 : *x*, levels *C* mm × √km. |
| Carry levels | Adds Ht diff and Level columns; the level misclosure is checked and spread. |
| Least squares | Standard errors of the angles and distances. |

**Load the sample loop** loads a primary loop on two pillars, observed with angles, levelled,
with two side shots. A row's **Side** can be *Point* (a traverse station) or *Side shot*: a point
taken from the same set-up as the row above it, worked out from the adjusted stations and never
part of the adjustment.

## Load traverse data

| File | |
|---|---|
| CSV, TXT (comma, tab, semicolon or space separated) | With a heading row, the columns are found by name: `Station, Backsight, Angle, Distance, Ht diff, Easting, Northing, Name, Description, Side`. Without one, the columns are read in the order `Station, Backsight, Angle, Distance, Easting, Northing, Name, Description`. A row with Easting and Northing but no angle or distance has them worked back from the coordinates. |
| JSON | A job saved with **Save**. |

Drag a file onto the window, or press **Load traverse data** (Ctrl+O). The Angle column is read as the Settings say: bearings unless you chose angles.

## Close and adjust

| Adjustment | |
|---|---|
| Compass (Bowditch) | The misclosure is spread in proportion to the length of each leg. |
| Transit | The misclosure in Easting and Northing is spread in proportion to each leg's ΔE and ΔN. |
| Least squares | Angles and distances are weighted by their standard errors and adjusted together, with residuals, a chi-square test and the 95% error ellipse of every station. Needs a closed traverse. |

The closure table gives the precision, the error distance and direction, the area and the
perimeter; the adjusted values give each station's coordinates, the adjusted bearing and
distance and the corrections ΔX, ΔY. Below them, the result against the tolerance.
**More results** opens the full checks, the angles, the least-squares statistics, the side
shots and the plan.

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
| Ctrl+Z / Ctrl+Y | Undo / redo |
| Enter / Shift+Enter | Next / previous row, same column |
| Esc | Stop adding points, close the settings or the report |

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
