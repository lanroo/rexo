# Installing REXO

REXO is a single self-contained binary. No Python, Node, Docker, or LLM account
is required to run it.

Pick your operating system, then choose **by file** (download + click) or **by
command line**. After any install, run `rexo doctor` to confirm it works.

## Downloads (latest release)

| OS | CPU | Download |
|---|---|---|
| Windows | x86-64 | [rexo_windows_amd64.zip](https://github.com/lanroo/rexo/releases/latest/download/rexo_windows_amd64.zip) |
| macOS | Apple Silicon (M1–M4) | [rexo_darwin_arm64.tar.gz](https://github.com/lanroo/rexo/releases/latest/download/rexo_darwin_arm64.tar.gz) |
| macOS | Intel | [rexo_darwin_amd64.tar.gz](https://github.com/lanroo/rexo/releases/latest/download/rexo_darwin_amd64.tar.gz) |
| Linux | x86-64 | [rexo_linux_amd64.tar.gz](https://github.com/lanroo/rexo/releases/latest/download/rexo_linux_amd64.tar.gz) |
| Linux | ARM64 | [rexo_linux_arm64.tar.gz](https://github.com/lanroo/rexo/releases/latest/download/rexo_linux_arm64.tar.gz) |

All releases (with `checksums.txt`) live on the
[Releases page](https://github.com/lanroo/rexo/releases).

---

## Windows

### Recommended: one-line install (no security warnings)

Open **PowerShell** and paste this single line:

```powershell
irm https://raw.githubusercontent.com/lanroo/rexo/main/install.ps1 | iex
```

That's it. The script installs [Scoop](https://scoop.sh) for you (user scope, no
admin) if you don't have it, then installs REXO through it. Installing via Scoop
avoids the SmartScreen "unrecognized app" prompt you get from a raw `.exe`.

Then run:

```powershell
rexo doctor
rexo init my-first-project
```

Update later with `scoop update rexo`.

> Prefer to do it by hand? `scoop install https://raw.githubusercontent.com/lanroo/rexo/main/scoop/rexo.json`

### By file (no command line)

1. Download [rexo_windows_amd64.zip](https://github.com/lanroo/rexo/releases/latest/download/rexo_windows_amd64.zip).
2. Right-click the `.zip` → **Extract All**. You get `rexo.exe`.
3. Open the folder, click the address bar, type `powershell`, press Enter.
4. Run:

   ```powershell
   .\rexo.exe doctor
   .\rexo.exe init my-first-project
   ```

To run `rexo` from anywhere, move `rexo.exe` into a folder on your `PATH`.

### By command line (PowerShell)

```powershell
Invoke-WebRequest -Uri https://github.com/lanroo/rexo/releases/latest/download/rexo_windows_amd64.zip -OutFile rexo.zip
Expand-Archive rexo.zip -DestinationPath rexo-bin -Force
.\rexo-bin\rexo.exe doctor
```

---

## macOS

### By file

1. Download the archive for your Mac —
   [Apple Silicon (M1–M4)](https://github.com/lanroo/rexo/releases/latest/download/rexo_darwin_arm64.tar.gz)
   or [Intel](https://github.com/lanroo/rexo/releases/latest/download/rexo_darwin_amd64.tar.gz).
   (Apple menu →  **About This Mac** tells you which one.)
2. Double-click the `.tar.gz` to extract `rexo`.
3. Open **Terminal**, `cd` into the folder, and run:

   ```bash
   ./rexo doctor
   ./rexo init my-first-project
   ```

> First run may warn "unidentified developer". Right-click `rexo` → **Open** →
> **Open**, or **System Settings → Privacy & Security → Open Anyway**.

### By command line (Terminal)

Apple Silicon — swap `arm64`→`amd64` on Intel Macs:

```bash
curl -L -o rexo.tar.gz https://github.com/lanroo/rexo/releases/latest/download/rexo_darwin_arm64.tar.gz
tar -xzf rexo.tar.gz
sudo mv rexo /usr/local/bin/rexo
rexo doctor
```

---

## Linux

### By file

1. Download [rexo_linux_amd64.tar.gz](https://github.com/lanroo/rexo/releases/latest/download/rexo_linux_amd64.tar.gz)
   (or [ARM64](https://github.com/lanroo/rexo/releases/latest/download/rexo_linux_arm64.tar.gz)).
2. Extract it (double-click, or `tar -xzf rexo_linux_amd64.tar.gz`).
3. In a terminal in that folder:

   ```bash
   ./rexo doctor
   ./rexo init my-first-project
   ```

### By command line

x86-64 — swap `amd64`→`arm64` for ARM machines:

```bash
curl -L -o rexo.tar.gz https://github.com/lanroo/rexo/releases/latest/download/rexo_linux_amd64.tar.gz
tar -xzf rexo.tar.gz
sudo mv rexo /usr/local/bin/rexo
rexo doctor
```

---

## From source (developers, any OS)

Requires **Go 1.24+** and **Git**.

```bash
git clone https://github.com/lanroo/rexo.git
cd rexo
go build -o rexo ./cmd/rexo   # rexo.exe on Windows
./rexo doctor
```

Or install straight onto your `PATH`:

```bash
go install github.com/lanroo/rexo/cmd/rexo@latest
```

The binary lands in `$(go env GOPATH)/bin` — make sure that is on your `PATH`.

---

## Windows security warnings (SmartScreen & Defender)

REXO's `.exe` is **not yet code-signed**, so on Windows you may hit one or both:

1. **"Windows protected your PC" (SmartScreen, blue/purple box)** — click
   **More info → Run anyway**.
2. **"...contains a virus or potentially unwanted software" (Defender)** — this
   is a **false positive**, common for new, unsigned Go programs. The binary is
   compiled by GitHub Actions from this public source and its SHA-256 matches
   `checksums.txt` (see below), so you can verify exactly what you are running.

If you trust the source and want to run it anyway:

- Open **Windows Security → Virus & threat protection → Protection history**,
  find the `rexo.exe` item, and choose **Allow / Restore**; or
- Add an exclusion for the folder REXO lives in
  (**Virus & threat protection → Manage settings → Exclusions**).

**If Defender keeps deleting the binary** (e.g. a Scoop install extracts it and
it disappears, so `rexo` reports "Could not create process"), add a one-time
folder exclusion from an **Administrator** PowerShell, then reinstall:

```powershell
# Scoop install — exclude the scoop folder (adjust the path to your user):
Add-MpPreference -ExclusionPath "$HOME\scoop"
scoop uninstall rexo
scoop install https://raw.githubusercontent.com/lanroo/rexo/main/scoop/rexo.json
rexo version
```

Adding the exclusion needs admin once; **running `rexo` itself never needs
admin**. The exclusion is a temporary measure until the false-positive report
clears or the binary is code-signed.

Only do this for binaries you obtained from the official
[Releases page](https://github.com/lanroo/rexo/releases) and whose checksum you
verified. Proper code signing is on the roadmap to remove this friction.

Developers can avoid the warning entirely by building from source or using
`go install github.com/lanroo/rexo/cmd/rexo@latest`.

## Verify the download (optional)

Each release ships `checksums.txt`. To confirm a file wasn't corrupted:

- **Linux/macOS**: `shasum -a 256 <file>` (or `sha256sum -c checksums.txt`)
- **Windows**: `Get-FileHash .\rexo_windows_amd64.zip -Algorithm SHA256`
  and compare with the matching line in `checksums.txt`.

## Uninstall

REXO writes nothing outside the folders you point it at. To remove it, delete
the `rexo` / `rexo.exe` binary (and any copy you placed on your `PATH`).
