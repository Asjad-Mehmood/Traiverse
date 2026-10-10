// Checks the traverse engine inside app/index.html, with no browser.
//   node tools/check-engine.js
// It takes the code between ENGINE START and ENGINE END, runs it, and works a few traverses
// whose answers are known. Exit code 1 if anything is off.
"use strict";
const fs = require("fs");
const path = require("path");

const html = fs.readFileSync(path.join(__dirname, "..", "app", "index.html"), "utf8");
const m = html.match(/\/\* ENGINE START[\s\S]*?\*\/([\s\S]*?)\/\* ENGINE END \*\//);
if (!m) { console.error("ENGINE START / ENGINE END markers not found in app/index.html"); process.exit(1); }
const SURV = new Function(m[1] + "\nreturn SURV;")();

let failed = 0;
function check(name, cond, detail) {
  console.log((cond ? "ok   " : "FAIL ") + name + (detail ? "  (" + detail + ")" : ""));
  if (!cond) failed++;
}
const near = (a, b, tol) => Math.abs(a - b) <= tol;

// 1. Angles as typed
check("DDD.MMSS", near(SURV.parseAngle("81.3127", "dms"), 81 + 31 / 60 + 27 / 3600, 1e-12));
check("D M S with symbols", near(SURV.parseAngle("81°31'27\"", "dms"), 81 + 31 / 60 + 27 / 3600, 1e-12));
check("gons", near(SURV.parseAngle("100", "gon"), 90, 1e-12));
check("bad minutes rejected", isNaN(SURV.parseAngle("81.6127", "dms")));

// 2. A perfect loop worked back from coordinates closes exactly, by every method
const P = [[1000, 1000], [1200, 1050], [1250, 1250], [1050, 1300], [950, 1150]];
const BS = [1000, 1500];
function loopFrom(pts, bs) {
  const n = pts.length, rows = [];
  const az = (a, b) => SURV.azimuth(a[0], a[1], b[0], b[1]);
  for (let i = 0; i < n; i++) {
    const back = i === 0 ? bs : pts[i - 1], fwd = pts[(i + 1) % n];
    rows.push({ name: "S" + i, angle: SURV.n360(az(pts[i], fwd) - az(pts[i], back)), dist: SURV.dist2D(pts[i][0], pts[i][1], fwd[0], fwd[1]), dh: 0 });
  }
  // closing angle at S0: from the last station back onto the backsight
  rows.push({ name: "S0", angle: SURV.n360(az(pts[0], bs) - az(pts[0], pts[n - 1])), dist: null, dh: null });
  return rows;
}
for (const method of ["bowditch", "transit", "lsq"]) {
  const t = {
    type: "loop", method, angleDir: "cw", scale: 1, tol: { angK: 10, ratio: 10000, hC: 12 }, sigma: { ang: 3, dA: 2, dB: 2 }, useH: false,
    startPt: { name: "S0", E: P[0][0], N: P[0][1], H: null }, endPt: null,
    startOri: { mode: "point", az: SURV.azimuth(P[0][0], P[0][1], BS[0], BS[1]) },
    endOri: { mode: "point", az: SURV.azimuth(P[0][0], P[0][1], BS[0], BS[1]) }, rows: loopFrom(P, BS)
  };
  const r = SURV.compute(t);
  const worst = Math.max(...P.map((p, i) => Math.hypot(r.finalE[i] - p[0], r.finalN[i] - p[1])));
  check("perfect loop, " + method, near(r.closure.angMis, 0, 1e-6) && r.closure.fL < 1e-9 && worst < 1e-6, "worst " + worst.toExponential(2) + " m");
}

// 3. The same loop with 30 mm added to one distance: Bowditch spreads the misclosure by length
{
  const rows = loopFrom(P, BS); rows[2].dist += 0.030;
  const t = {
    type: "loop", method: "bowditch", angleDir: "cw", scale: 1, tol: { angK: 10, ratio: 10000, hC: 12 }, sigma: { ang: 3, dA: 2, dB: 2 }, useH: false,
    startPt: { name: "S0", E: P[0][0], N: P[0][1], H: null }, endPt: null,
    startOri: { mode: "point", az: SURV.azimuth(P[0][0], P[0][1], BS[0], BS[1]) },
    endOri: { mode: "point", az: SURV.azimuth(P[0][0], P[0][1], BS[0], BS[1]) }, rows
  };
  const r = SURV.compute(t);
  check("30 mm blunder found", near(r.closure.fL, 0.030, 1e-6), "linear misclosure " + (r.closure.fL * 1000).toFixed(2) + " mm");
  const sumCE = r.cE.reduce((a, b) => a + b, 0), sumCN = r.cN.reduce((a, b) => a + b, 0);
  check("Bowditch corrections cancel the misclosure", near(sumCE, -r.closure.fE, 1e-9) && near(sumCN, -r.closure.fN, 1e-9));
  check("loop closes back on its start", near(r.finalE[r.finalE.length - 1], P[0][0], 1e-9) && near(r.finalN[r.finalN.length - 1], P[0][1], 1e-9));
}

// 4. The sample traverse in the program (same as Excel to KML's example): within primary tolerance
{
  const rows = [["CP-01", "81.3127", 316.630, 0.646], ["TS-01", "124.3129", 308.240, 0.584], ["TS-02", "125.5354", 308.364, 0.566],
    ["TS-03", "104.5917", 293.903, -0.802], ["TS-04", "129.2554", 294.528, -0.478], ["TS-05", "130.2740", 288.697, -0.523], ["CP-01", "23.1021", null, null]]
    .map(r => ({ name: r[0], angle: SURV.parseAngle(r[1], "dms"), dist: r[2], dh: r[3] }));
  const S = { name: "CP-01", E: 228450.312, N: 441203.845, H: 4.215 }, az = SURV.azimuth(S.E, S.N, 228912.640, 441655.118);
  for (const method of ["bowditch", "transit", "lsq"]) {
    const r = SURV.compute({
      type: "loop", method, angleDir: "cw", scale: 1, tol: { angK: 5, ratio: 20000, hC: 8 }, sigma: { ang: 2, dA: 1, dB: 1.5 }, useH: true,
      startPt: S, endPt: null, startOri: { mode: "point", az }, endOri: { mode: "point", az }, rows
    });
    const c = r.closure;
    check("sample loop, " + method, c.angPass && c.precPass && c.hPass,
      "angular " + c.angMis.toFixed(1) + "\", 1:" + Math.floor(c.precision) + ", level " + (c.fH * 1000).toFixed(1) + " mm");
  }
}

if (failed) { console.error(failed + " check(s) failed"); process.exit(1); }
console.log("all engine checks passed");
