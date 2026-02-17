#!/usr/bin/env node

"use strict";

const os = require("os");
const fs = require("fs");
const path = require("path");
const https = require("https");
const { execSync } = require("child_process");

const PACKAGE = require("./package.json");
const VERSION = PACKAGE.version;
const REPO = "Hopsule/cli-tool";
const BINARY_NAME = "hopsule";
const DOWNLOAD_TIMEOUT_MS = 120_000;

const PLATFORM_MAP = {
  darwin: "darwin",
  linux: "linux",
  win32: "windows",
};

const ARCH_MAP = {
  x64: "amd64",
  arm64: "arm64",
};

// ── Pretty output helpers ───────────────────────────────────────────

const CYAN = "\x1b[36m";
const GREEN = "\x1b[32m";
const DIM = "\x1b[2m";
const BOLD = "\x1b[1m";
const RESET = "\x1b[0m";
const RED = "\x1b[31m";
const YELLOW = "\x1b[33m";

function formatBytes(bytes) {
  if (bytes < 1024) return `${bytes} B`;
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
}

function formatSpeed(bytesPerSec) {
  if (bytesPerSec < 1024) return `${bytesPerSec.toFixed(0)} B/s`;
  if (bytesPerSec < 1024 * 1024) return `${(bytesPerSec / 1024).toFixed(1)} KB/s`;
  return `${(bytesPerSec / (1024 * 1024)).toFixed(1)} MB/s`;
}

function drawProgressBar(percent, width = 30) {
  const filled = Math.round(width * percent);
  const empty = width - filled;
  const bar = "█".repeat(filled) + "░".repeat(empty);
  return bar;
}

function clearLine() {
  if (process.stderr.isTTY) {
    process.stderr.write("\r\x1b[K");
  }
}

// ── Platform detection ──────────────────────────────────────────────

function getPlatform() {
  const platform = PLATFORM_MAP[os.platform()];
  if (!platform) {
    throw new Error(
      `Unsupported platform: ${os.platform()}. ` +
        `Supported: ${Object.keys(PLATFORM_MAP).join(", ")}`
    );
  }
  return platform;
}

function getArch() {
  const arch = ARCH_MAP[os.arch()];
  if (!arch) {
    throw new Error(
      `Unsupported architecture: ${os.arch()}. ` +
        `Supported: ${Object.keys(ARCH_MAP).join(", ")}`
    );
  }
  return arch;
}

function getDownloadUrl(platform, arch) {
  const ext = platform === "windows" ? "zip" : "tar.gz";
  return `https://github.com/${REPO}/releases/download/v${VERSION}/${BINARY_NAME}-${platform}-${arch}.${ext}`;
}

function getBinaryPath() {
  const binDir = path.join(__dirname, "bin");
  const ext = os.platform() === "win32" ? ".exe" : "";
  return path.join(binDir, `${BINARY_NAME}${ext}`);
}

// ── Download with progress ──────────────────────────────────────────

function downloadFile(url) {
  return new Promise((resolve, reject) => {
    const startTime = Date.now();
    let timer = null;

    const cleanup = () => {
      if (timer) clearTimeout(timer);
    };

    const follow = (currentUrl, redirects = 0) => {
      if (redirects > 10) {
        cleanup();
        return reject(new Error("Too many redirects"));
      }

      if (timer) clearTimeout(timer);
      timer = setTimeout(() => {
        reject(new Error(`Download timed out after ${DOWNLOAD_TIMEOUT_MS / 1000}s`));
      }, DOWNLOAD_TIMEOUT_MS);

      const req = https
        .get(currentUrl, { timeout: 30000, headers: { "User-Agent": "hopsule-npm-installer" } }, (res) => {
          if (res.statusCode >= 300 && res.statusCode < 400 && res.headers.location) {
            res.resume();
            return follow(res.headers.location, redirects + 1);
          }

          if (res.statusCode !== 200) {
            res.resume();
            cleanup();
            return reject(
              new Error(`Download failed: HTTP ${res.statusCode} from ${currentUrl}`)
            );
          }

          const totalSize = parseInt(res.headers["content-length"], 10) || 0;
          let downloaded = 0;
          const chunks = [];
          let lastUpdate = Date.now();

          res.on("data", (chunk) => {
            chunks.push(chunk);
            downloaded += chunk.length;

            const now = Date.now();
            if (now - lastUpdate < 100) return;
            lastUpdate = now;

            const elapsed = (now - startTime) / 1000;
            const speed = downloaded / elapsed;
            const percent = totalSize > 0 ? downloaded / totalSize : 0;

            if (process.stderr.isTTY) {
              clearLine();
              if (totalSize > 0) {
                const bar = drawProgressBar(percent);
                const pct = (percent * 100).toFixed(0).padStart(3);
                process.stderr.write(
                  `  ${CYAN}${bar}${RESET} ${pct}%  ${DIM}${formatBytes(downloaded)}/${formatBytes(totalSize)}${RESET}  ${DIM}${formatSpeed(speed)}${RESET}`
                );
              } else {
                const spinner = ["⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"];
                const frame = spinner[Math.floor(elapsed * 4) % spinner.length];
                process.stderr.write(
                  `  ${CYAN}${frame}${RESET} Downloading...  ${DIM}${formatBytes(downloaded)}${RESET}  ${DIM}${formatSpeed(speed)}${RESET}`
                );
              }
            }
          });

          res.on("end", () => {
            cleanup();
            if (process.stderr.isTTY) {
              clearLine();
              const elapsed = (Date.now() - startTime) / 1000;
              process.stderr.write(
                `  ${GREEN}${drawProgressBar(1)}${RESET} 100%  ${DIM}${formatBytes(downloaded)} in ${elapsed.toFixed(1)}s${RESET}\n`
              );
            }
            resolve(Buffer.concat(chunks));
          });

          res.on("error", (err) => {
            cleanup();
            reject(err);
          });
        })
        .on("error", (err) => {
          cleanup();
          reject(err);
        });

      req.on("timeout", () => {
        req.destroy();
        cleanup();
        reject(new Error(`Connection timed out (attempt to ${currentUrl.substring(0, 60)}...)`));
      });
    };

    follow(url);
  });
}

// ── Extraction ──────────────────────────────────────────────────────

async function extractTarGz(buffer, destDir) {
  const tmpFile = path.join(os.tmpdir(), `hopsule-${Date.now()}.tar.gz`);
  fs.writeFileSync(tmpFile, buffer);

  try {
    execSync(`tar xzf "${tmpFile}" -C "${destDir}"`, { stdio: "pipe" });
  } finally {
    try {
      fs.unlinkSync(tmpFile);
    } catch (_) {}
  }
}

async function extractZip(buffer, destDir) {
  const tmpFile = path.join(os.tmpdir(), `hopsule-${Date.now()}.zip`);
  fs.writeFileSync(tmpFile, buffer);

  try {
    if (os.platform() === "win32") {
      execSync(
        `powershell -Command "Expand-Archive -Path '${tmpFile}' -DestinationPath '${destDir}' -Force"`,
        { stdio: "pipe" }
      );
    } else {
      execSync(`unzip -o "${tmpFile}" -d "${destDir}"`, { stdio: "pipe" });
    }
  } finally {
    try {
      fs.unlinkSync(tmpFile);
    } catch (_) {}
  }
}

// ── Main install ────────────────────────────────────────────────────

async function install() {
  const platform = getPlatform();
  const arch = getArch();
  const url = getDownloadUrl(platform, arch);
  const binDir = path.join(__dirname, "bin");
  const binaryPath = getBinaryPath();

  console.log("");
  console.log(`  ${BOLD}${CYAN}Hopsule${RESET} ${DIM}v${VERSION}${RESET}`);
  console.log(`  ${DIM}Decision & Memory Layer for AI teams${RESET}`);
  console.log("");

  // Skip if binary already exists with correct version
  if (fs.existsSync(binaryPath)) {
    try {
      const output = execSync(`"${binaryPath}" --version`, {
        encoding: "utf8",
        stdio: "pipe",
      });
      if (output.includes(VERSION)) {
        console.log(`  ${GREEN}✓${RESET} Already installed ${DIM}(v${VERSION})${RESET}`);
        console.log("");
        return;
      }
    } catch (_) {}
  }

  console.log(`  ${DIM}Platform:${RESET} ${platform}/${arch}`);
  console.log(`  ${DIM}Source:${RESET}   github.com/${REPO}`);
  console.log("");

  const buffer = await downloadFile(url);

  // Ensure bin directory exists
  if (!fs.existsSync(binDir)) {
    fs.mkdirSync(binDir, { recursive: true });
  }

  // Extract
  process.stderr.write(`  ${DIM}Extracting...${RESET}`);
  const tmpExtractDir = path.join(os.tmpdir(), `hopsule-extract-${Date.now()}`);
  fs.mkdirSync(tmpExtractDir, { recursive: true });

  try {
    if (platform === "windows") {
      await extractZip(buffer, tmpExtractDir);
    } else {
      await extractTarGz(buffer, tmpExtractDir);
    }

    // Find and move the binary
    const ext = platform === "windows" ? ".exe" : "";
    const extractedBinary = path.join(tmpExtractDir, `${BINARY_NAME}${ext}`);

    if (!fs.existsSync(extractedBinary)) {
      const files = fs.readdirSync(tmpExtractDir);
      let found = false;
      for (const f of files) {
        const nested = path.join(tmpExtractDir, f, `${BINARY_NAME}${ext}`);
        if (fs.existsSync(nested)) {
          fs.copyFileSync(nested, binaryPath);
          found = true;
          break;
        }
      }
      if (!found) {
        throw new Error(
          `Binary not found in archive. Contents: ${files.join(", ")}`
        );
      }
    } else {
      fs.copyFileSync(extractedBinary, binaryPath);
    }

    // Make executable on unix
    if (platform !== "windows") {
      fs.chmodSync(binaryPath, 0o755);
    }

    clearLine();
    console.log(`  ${GREEN}✓${RESET} Extracted`);
    console.log("");
    console.log(`  ${GREEN}${BOLD}✓ hopsule v${VERSION} installed successfully!${RESET}`);
    console.log("");
    console.log(`  ${DIM}Get started:${RESET}`);
    console.log(`    ${CYAN}hopsule login${RESET}      ${DIM}Sign in to your account${RESET}`);
    console.log(`    ${CYAN}hopsule init${RESET}       ${DIM}Initialize a project${RESET}`);
    console.log(`    ${CYAN}hopsule${RESET}            ${DIM}Open interactive TUI${RESET}`);
    console.log("");
  } finally {
    try {
      fs.rmSync(tmpExtractDir, { recursive: true, force: true });
    } catch (_) {}
  }
}

install().catch((err) => {
  clearLine();
  console.error("");
  console.error(`  ${RED}${BOLD}✗ Failed to install hopsule${RESET}`);
  console.error(`  ${RED}${err.message}${RESET}`);
  console.error("");
  console.error(`  ${YELLOW}Install manually:${RESET}`);
  console.error(`    ${DIM}https://github.com/${REPO}/releases/tag/v${VERSION}${RESET}`);
  console.error("");
  console.error(`  ${YELLOW}Or use Homebrew:${RESET}`);
  console.error(`    ${CYAN}brew install hopsule/tap/hopsule${RESET}`);
  console.error("");
  process.exit(1);
});
