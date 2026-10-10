// Checks the Qatar National Grid in app/index.html, with no browser.
//   node tools/check-qng.js
// The reference values in qng-reference.json came from outside this program:
//   k_pyproj           point scale factor from PROJ (pyproj), EPSG:2932's projection on International 1924
//   lat_kml, lon_kml   WGS84 from Excel to KML's own Qatar National Grid conversion
//   loop, link         traverses on true grid coordinates, with ground distances from the
//                      ellipsoid geodesic (pyproj) at 30 m height, and grid angles
// Exit code 1 if anything is off.
"use strict";
const fs = require("fs");
const path = require("path");

const html = fs.readFileSync(path.join(__dirname, "..", "app", "index.html"), "utf8");
const eng = html.match(/\/\* ENGINE START[\s\S]*?\*\/([\s\S]*?)\/\* ENGINE END \*\//);
const qng = html.match(/\/\* QNG START[\s\S]*?\*\/([\s\S]*?)\/\* QNG END \*\//);
if (!eng || !qng) { console.error("ENGINE or QNG markers not found in app/index.html"); process.exit(1); }
const SURV = new Function(eng[1] + "\nreturn SURV;")();
const QNG = new Function(qng[1] + "\nreturn QNG;")();
const ref = require("./qng-reference.json");

let failed = 0;
function check(name, cond, detail) {
  console.log((cond ? "ok   " : "FAIL ") + name + (detail ? "  (" + detail + ")" : ""));
  if (!cond) failed++;
}

// 1. The projection: point scale factor as PROJ gives it, and WGS84 as Excel to KML gives it
for (const [name, r] of Object.entries(ref.points)) {
  const k = QNG.pointScale(r.E, r.N), ll = QNG.toWGS84(r.E, r.N, 0);
  const dm = Math.hypot((ll[0] - r.lat_kml) * 111320, (ll[1] - r.lon_kml) * 111320 * Math.cos(ll[0] * Math.PI / 180));
  check("scale factor at " + name, Math.abs(k - r.k_pyproj) < 2e-9, "k " + k.toFixed(9) + ", PROJ " + r.k_pyproj.toFixed(9));
  check("WGS84 at " + name + " as Excel to KML", dm < 0.001, (dm * 1000).toFixed(4) + " mm apart");
}
check("Doha is in Qatar", QNG.inQatar(229500, 391000));
check("swapped coordinates are not", !QNG.inQatar(391000, 229500));

// 2. Traverses: ground distances times the combined scale factor land on the true grid coordinates
function travel(T, type) {
  const st = T.stations, n = st.length;
  const mid = st.reduce((a, s) => [a[0] + s[0] / n, a[1] + s[1] / n], [0, 0]);
  const csf = QNG.pointScale(mid[0], mid[1]) * QNG.heightFactor(mid[1], T.h);
  const az = (a, b) => SURV.azimuth(a[0], a[1], b[0], b[1]);
  const names = [type === "loop" ? "L0" : "K0"].concat(T.rows.map(r => r.name));
  const rows = names.map((nm, i) => i < names.length - 1
    ? { name: nm, angle: +T.rows[i].angle, dist: +T.rows[i].distance, dh: null }
    : { name: nm, angle: +T.closeAngle, dist: null, dh: null });
  const S = { name: names[0], E: st[0][0], N: st[0][1], H: null };
  const end = type === "loop" ? S : { name: names[names.length - 1], E: st[n - 1][0], N: st[n - 1][1], H: null };
  const endAz = type === "loop" ? az(st[0], T.bs) : az(st[n - 1], T.fs);
  return { csf, run: sc => SURV.compute({
    type, method: "bowditch", angleDir: "cw", scale: sc, tol: { angK: 5, ratio: 20000, hC: 8 }, sigma: { ang: 2, dA: 1, dB: 1 }, useH: false,
    startPt: S, endPt: type === "link" ? end : null, startOri: { mode: "point", az: az(st[0], T.bs) }, endOri: { mode: "point", az: endAz }, rows
  }) };
}
for (const [type, T] of [["loop", ref.loop], ["link", ref.link]]) {
  const t = travel(T, type), r = t.run(t.csf), r1 = t.run(1);
  const worst = res => Math.max(...T.stations.map((s, i) => Math.hypot(res.finalE[i] - s[0], res.finalN[i] - s[1])));
  check(type + " with the grid scale factor " + t.csf.toFixed(8), worst(r) < 0.003 && r.closure.fL < 0.003,
    "misclosure " + (r.closure.fL * 1000).toFixed(2) + " mm, worst station " + (worst(r) * 1000).toFixed(2) + " mm from truth");
  check(type + " without it is visibly off", worst(r1) > 0.01 || r1.closure.fL > 0.01,
    "misclosure " + (r1.closure.fL * 1000).toFixed(1) + " mm, worst station " + (worst(r1) * 1000).toFixed(1) + " mm");
}

if (failed) { console.error(failed + " check(s) failed"); process.exit(1); }
console.log("all Qatar National Grid checks passed");
