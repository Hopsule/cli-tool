$ErrorActionPreference = 'Stop'

$packageName = 'hopsule'
$version = '0.9.8'

# Determine architecture
$arch = if ([Environment]::Is64BitOperatingSystem) {
    if ($env:PROCESSOR_ARCHITECTURE -eq 'ARM64') { 'arm64' } else { 'amd64' }
} else {
    throw "32-bit systems are not supported."
}

$url = "https://github.com/Hopsule/cli-tool/releases/download/v${version}/hopsule-windows-${arch}.zip"
$checksum = '' # Updated by CI during release
$checksumType = 'sha256'

$installDir = Join-Path $env:ChocolateyInstall 'lib\hopsule\tools'

$packageArgs = @{
    PackageName    = $packageName
    UnzipLocation  = $installDir
    Url64bit       = $url
    Checksum64     = $checksum
    ChecksumType64 = $checksumType
}

Install-ChocolateyZipPackage @packageArgs

# Add to PATH via shim
$binaryPath = Join-Path $installDir 'hopsule.exe'
if (Test-Path $binaryPath) {
    Write-Host "hopsule v${version} installed successfully!"
} else {
    throw "Installation failed: binary not found at ${binaryPath}"
}
