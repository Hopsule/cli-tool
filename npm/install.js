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

const PLATFORM_MAP = {
  darwin: "darwin",
  linux: "linux",
  win32: "windows",
};

const ARCH_MAP = {
  x64: "amd64",
  arm64: "arm64",
};

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

function downloadFile(url) {
  return new Promise((resolve, reject) => {
    const follow = (url, redirects = 0) => {
      if (redirects > 10) {
        return reject(new Error("Too many redirects"));
      }

      https
        .get(url, { headers: { "User-Agent": "hopsule-npm-installer" } }, (res) => {
          if (res.statusCode >= 300 && res.statusCode < 400 && res.headers.location) {
            return follow(res.headers.location, redirects + 1);
          }

          if (res.statusCode !== 200) {
            return reject(
              new Error(`Download failed: HTTP ${res.statusCode} from ${url}`)
            );
          }

          const chunks = [];
          res.on("data", (chunk) => chunks.push(chunk));
          res.on("end", () => resolve(Buffer.concat(chunks)));
          res.on("error", reject);
        })
        .on("error", reject);
    };

    follow(url);
  });
}

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

async function install() {
  const platform = getPlatform();
  const arch = getArch();
  const url = getDownloadUrl(platform, arch);
  const binDir = path.join(__dirname, "bin");
  const binaryPath = getBinaryPath();

  // Skip if binary already exists with correct version
  if (fs.existsSync(binaryPath)) {
    try {
      const output = execSync(`"${binaryPath}" --version`, {
        encoding: "utf8",
        stdio: "pipe",
      });
      if (output.includes(VERSION)) {
        console.log(`hopsule v${VERSION} already installed.`);
        return;
      }
    } catch (_) {}
  }

  console.log(`Installing hopsule v${VERSION} for ${platform}/${arch}...`);
  console.log(`Downloading from: ${url}`);

  const buffer = await downloadFile(url);

  // Ensure bin directory exists
  if (!fs.existsSync(binDir)) {
    fs.mkdirSync(binDir, { recursive: true });
  }

  // Extract based on platform
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
      // Some archives nest in a directory
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

    console.log(`hopsule v${VERSION} installed successfully!`);
  } finally {
    try {
      fs.rmSync(tmpExtractDir, { recursive: true, force: true });
    } catch (_) {}
  }
}

install().catch((err) => {
  console.error(`Failed to install hopsule: ${err.message}`);
  console.error("");
  console.error("You can install manually from:");
  console.error(`  https://github.com/${REPO}/releases/tag/v${VERSION}`);
  console.error("");
  console.error("Or use Homebrew: brew install hopsule/tap/hopsule");
  process.exit(1);
});
