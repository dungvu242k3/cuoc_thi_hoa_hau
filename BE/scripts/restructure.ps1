$ErrorActionPreference = "Stop"
$root = "c:\Users\dungv\cuoc_thi_hoa_hau\BE"
Set-Location $root

# Phase 1: Create new directories
$dirs = @(
    "graph\schema", "graph\resolver", "graph\generated", "graph\model", "graph\middleware", "graph\dataloader",
    "internal\cache", "internal\config", "internal\dao", "internal\database",
    "internal\ecode", "internal\handler", "internal\model", "internal\routers",
    "internal\server", "internal\service", "internal\types",
    "configs", "deployments", "scripts", "docs"
)
foreach ($d in $dirs) {
    $p = Join-Path $root $d
    if (-not (Test-Path $p)) { New-Item -ItemType Directory -Path $p -Force | Out-Null }
}
Write-Host "Directories created."

# Phase 2: Copy files to new locations
$copies = @(
    # graph (from internal/adapter/graph -> graph/)
    @{ Src = "internal\adapter\graph\schema\*"; Dst = "graph\schema\" },
    @{ Src = "internal\adapter\graph\generated\*"; Dst = "graph\generated\" },
    @{ Src = "internal\adapter\graph\model\*"; Dst = "graph\model\" },
    @{ Src = "internal\adapter\graph\resolver\*"; Dst = "graph\resolver\" },
    @{ Src = "internal\adapter\graph\middleware\*"; Dst = "graph\middleware\" },
    # cache
    @{ Src = "internal\adapter\cache\*"; Dst = "internal\cache\" },
    # handler
    @{ Src = "internal\adapter\handler\*"; Dst = "internal\handler\" },
    # dao (mongodb repos + connection)
    @{ Src = "internal\adapter\storage\mongodb\contestant_repo.go"; Dst = "internal\dao\contestant_repo.go" },
    @{ Src = "internal\adapter\storage\mongodb\user_repo.go"; Dst = "internal\dao\user_repo.go" },
    @{ Src = "internal\adapter\storage\mongodb\feedback_repo.go"; Dst = "internal\dao\feedback_repo.go" },
    @{ Src = "internal\adapter\storage\mongodb\score_repo.go"; Dst = "internal\dao\score_repo.go" },
    @{ Src = "internal\adapter\storage\mongodb\schedule_repo.go"; Dst = "internal\dao\schedule_repo.go" },
    # database (connection)
    @{ Src = "internal\adapter\storage\mongodb\connection.go"; Dst = "internal\database\connection.go" },
    # local storage -> internal/dao
    @{ Src = "internal\adapter\storage\local\local_storage.go"; Dst = "internal\dao\local_storage.go" },
    # model (domain)
    @{ Src = "internal\core\domain\contestant.go"; Dst = "internal\model\contestant.go" },
    @{ Src = "internal\core\domain\user.go"; Dst = "internal\model\user.go" },
    @{ Src = "internal\core\domain\feedback.go"; Dst = "internal\model\feedback.go" },
    @{ Src = "internal\core\domain\schedule.go"; Dst = "internal\model\schedule.go" },
    @{ Src = "internal\core\domain\score.go"; Dst = "internal\model\score.go" },
    @{ Src = "internal\core\domain\rbac.go"; Dst = "internal\model\rbac.go" },
    @{ Src = "internal\core\domain\rbac_test.go"; Dst = "internal\model\rbac_test.go" },
    # service
    @{ Src = "internal\core\service\auth_svc.go"; Dst = "internal\service\auth_svc.go" },
    @{ Src = "internal\core\service\auth_service.go"; Dst = "internal\service\auth_service.go" },
    @{ Src = "internal\core\service\contestant_svc.go"; Dst = "internal\service\contestant_svc.go" },
    @{ Src = "internal\core\service\contestant_service.go"; Dst = "internal\service\contestant_service.go" },
    @{ Src = "internal\core\service\feedback_svc.go"; Dst = "internal\service\feedback_svc.go" },
    @{ Src = "internal\core\service\schedule_svc.go"; Dst = "internal\service\schedule_svc.go" },
    @{ Src = "internal\core\service\scoring_svc.go"; Dst = "internal\service\scoring_svc.go" },
    # port -> internal/types (interfaces)
    @{ Src = "internal\core\port\auth_port.go"; Dst = "internal\types\auth_port.go" },
    @{ Src = "internal\core\port\contestant_port.go"; Dst = "internal\types\contestant_port.go" },
    @{ Src = "internal\core\port\feedback_port.go"; Dst = "internal\types\feedback_port.go" },
    @{ Src = "internal\core\port\infra.go"; Dst = "internal\types\infra.go" },
    @{ Src = "internal\core\port\score_port.go"; Dst = "internal\types\score_port.go" },
    @{ Src = "internal\core\port\schedule_port.go"; Dst = "internal\types\schedule_port.go" },
    @{ Src = "internal\core\port\security_port.go"; Dst = "internal\types\security_port.go" },
    @{ Src = "internal\core\port\storage.go"; Dst = "internal\types\storage.go" },
    @{ Src = "internal\core\port\user_port.go"; Dst = "internal\types\user_port.go" },
    @{ Src = "internal\core\port\probe_port.go"; Dst = "internal\types\probe_port.go" },
    @{ Src = "internal\core\port\repository.go"; Dst = "internal\types\repository.go" },
    @{ Src = "internal\core\port\service.go"; Dst = "internal\types\service_iface.go" },
    # internal/pkg -> pkg (validation, common)
    @{ Src = "internal\adapter\infra\security\bcrypt.go"; Dst = "internal\server\bcrypt.go" },
    @{ Src = "internal\adapter\infra\security\jwt.go"; Dst = "internal\server\jwt.go" },
    # deploy -> deployments
    @{ Src = "deploy\docker-compose.yml"; Dst = "deployments\docker-compose.yml" },
    # worker
    @{ Src = "internal\worker\task_handler.go"; Dst = "internal\service\task_handler.go" },
    # pkg validation
    @{ Src = "internal\pkg\validation\contestant.go"; Dst = "pkg\validation\contestant.go" },
    @{ Src = "internal\pkg\validation\pagination.go"; Dst = "pkg\validation\pagination.go" },
    @{ Src = "internal\pkg\validation\score.go"; Dst = "pkg\validation\score.go" },
    @{ Src = "internal\pkg\common\constants.go"; Dst = "pkg\common\constants.go" }
)

foreach ($c in $copies) {
    $srcPath = Join-Path $root $c.Src
    $dstPath = Join-Path $root $c.Dst
    $dstDir = Split-Path $dstPath -Parent
    if (-not (Test-Path $dstDir)) { New-Item -ItemType Directory -Path $dstDir -Force | Out-Null }
    
    if ($c.Src -like "*\*") {
        # Wildcard copy
        Copy-Item -Path $srcPath -Destination $dstDir -Force -ErrorAction SilentlyContinue
    } else {
        Copy-Item -Path $srcPath -Destination $dstPath -Force
    }
}
Write-Host "Files copied."

# Phase 3: Replace import paths and package names in ALL Go files under new directories
$replacements = @(
    # Import path replacements (order matters - most specific first)
    @{ Old = 'cuoc_thi_hoa_hau/internal/adapter/graph/generated'; New = 'cuoc_thi_hoa_hau/graph/generated' },
    @{ Old = 'cuoc_thi_hoa_hau/internal/adapter/graph/middleware'; New = 'cuoc_thi_hoa_hau/graph/middleware' },
    @{ Old = 'cuoc_thi_hoa_hau/internal/adapter/graph/model'; New = 'cuoc_thi_hoa_hau/graph/model' },
    @{ Old = 'cuoc_thi_hoa_hau/internal/adapter/graph/resolver'; New = 'cuoc_thi_hoa_hau/graph/resolver' },
    @{ Old = 'cuoc_thi_hoa_hau/internal/adapter/cache'; New = 'cuoc_thi_hoa_hau/internal/cache' },
    @{ Old = 'cuoc_thi_hoa_hau/internal/adapter/handler'; New = 'cuoc_thi_hoa_hau/internal/handler' },
    @{ Old = 'cuoc_thi_hoa_hau/internal/adapter/infra/security'; New = 'cuoc_thi_hoa_hau/internal/server' },
    @{ Old = 'cuoc_thi_hoa_hau/internal/adapter/storage/local'; New = 'cuoc_thi_hoa_hau/internal/dao' },
    @{ Old = 'cuoc_thi_hoa_hau/internal/adapter/storage/mongodb'; New = 'cuoc_thi_hoa_hau/internal/dao' },
    @{ Old = 'cuoc_thi_hoa_hau/internal/core/domain'; New = 'cuoc_thi_hoa_hau/internal/model' },
    @{ Old = 'cuoc_thi_hoa_hau/internal/core/port'; New = 'cuoc_thi_hoa_hau/internal/types' },
    @{ Old = 'cuoc_thi_hoa_hau/internal/core/service'; New = 'cuoc_thi_hoa_hau/internal/service' },
    @{ Old = 'cuoc_thi_hoa_hau/internal/pkg/validation'; New = 'cuoc_thi_hoa_hau/pkg/validation' },
    @{ Old = 'cuoc_thi_hoa_hau/internal/pkg/common'; New = 'cuoc_thi_hoa_hau/pkg/common' }
)

# Package name replacements (for files that changed package)
$pkgReplacements = @(
    @{ Dir = "internal\model"; Old = "package domain"; New = "package model" },
    @{ Dir = "internal\types"; Old = "package port"; New = "package types" },
    @{ Dir = "internal\dao"; Old = "package mongodb"; New = "package dao" },
    @{ Dir = "internal\dao"; Old = "package local"; New = "package dao" },
    @{ Dir = "internal\server"; Old = "package security"; New = "package server" }
)

# Find all Go files in new locations
$goFiles = Get-ChildItem -Path $root -Include "*.go" -Recurse | Where-Object {
    $rel = $_.FullName.Substring($root.Length + 1)
    # Only process files in new locations (avoid old locations)
    $rel -match "^(graph|internal\\(cache|config|dao|database|ecode|handler|model|routers|server|service|types)|pkg\\(validation|common)|cmd|configs)" -and
    $rel -notmatch "^internal\\adapter" -and
    $rel -notmatch "^internal\\core" -and
    $rel -notmatch "^internal\\pkg"
}

# Also include cmd/ files
$cmdFiles = Get-ChildItem -Path (Join-Path $root "cmd") -Include "*.go" -Recurse
$goFiles = @($goFiles) + @($cmdFiles) | Sort-Object -Property FullName -Unique

foreach ($file in $goFiles) {
    $content = Get-Content -Path $file.FullName -Raw -Encoding UTF8
    $changed = $false
    
    # Apply import replacements
    foreach ($r in $replacements) {
        if ($content -match [regex]::Escape($r.Old)) {
            $content = $content.Replace($r.Old, $r.New)
            $changed = $true
        }
    }
    
    # Apply package name replacements
    foreach ($pr in $pkgReplacements) {
        $dirPath = Join-Path $root $pr.Dir
        if ($file.DirectoryName -eq $dirPath) {
            if ($content -match [regex]::Escape($pr.Old)) {
                $content = $content.Replace($pr.Old, $pr.New)
                $changed = $true
            }
        }
    }
    
    if ($changed) {
        Set-Content -Path $file.FullName -Value $content -NoNewline -Encoding UTF8
        Write-Host "Updated: $($file.FullName.Substring($root.Length + 1))"
    }
}

Write-Host "`nPhase 3 complete: Imports updated."
Write-Host "Done! Run 'go build ./...' to verify."
