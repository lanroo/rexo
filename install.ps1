# REXO installer for Windows.
#
# One-line usage (nothing else required):
#   irm https://raw.githubusercontent.com/lanroo/rexo/main/install.ps1 | iex
#
# It installs Scoop for you (user scope, no admin) if it is missing, then
# installs REXO through it. Installing via Scoop avoids the SmartScreen
# "unrecognized app" prompt you get from a raw .exe download.

$ErrorActionPreference = 'Stop'

function Write-Step($msg) { Write-Host "==> $msg" -ForegroundColor Cyan }

Write-Step "Installing REXO — Runtime for Execution & eXchange Orchestration"

# 1. Make sure Scoop is available (install it if not).
if (-not (Get-Command scoop -ErrorAction SilentlyContinue)) {
    Write-Step "Scoop not found. Installing Scoop (current user, no admin needed)..."
    try {
        Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser -Force
    } catch {
        Write-Host "Could not change execution policy automatically; continuing." -ForegroundColor Yellow
    }
    Invoke-RestMethod -Uri 'https://get.scoop.sh' | Invoke-Expression
} else {
    Write-Step "Scoop is already installed."
}

# 2. Install (or update) REXO from this repo's manifest.
$manifest = 'https://raw.githubusercontent.com/lanroo/rexo/main/scoop/rexo.json'
if (scoop list rexo 6>$null | Select-String -SimpleMatch 'rexo') {
    Write-Step "Updating REXO..."
    scoop update rexo
} else {
    Write-Step "Installing REXO..."
    scoop install $manifest
}

Write-Host ""
Write-Host "REXO is installed. Try it now:" -ForegroundColor Green
Write-Host "  rexo doctor"
Write-Host "  rexo init my-first-project"
