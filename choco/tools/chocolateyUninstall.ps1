$ErrorActionPreference = 'Stop'

$packageName = 'hopsule'
$installDir = Join-Path $env:ChocolateyInstall 'lib\hopsule\tools'

if (Test-Path $installDir) {
    Remove-Item -Path $installDir -Recurse -Force
    Write-Host "hopsule uninstalled successfully."
} else {
    Write-Host "hopsule was not found in expected location."
}
