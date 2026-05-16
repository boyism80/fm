param(
    [Parameter(Mandatory = $true)]
    [string]$TemplatePath,
    [Parameter(Mandatory = $true)]
    [int]$ChannelId,
    [Parameter(Mandatory = $true)]
    [string]$OutputPath
)

if (-not (Test-Path -LiteralPath $TemplatePath)) {
    Write-Error "Game template not found: $TemplatePath"
    exit 1
}

$basePort = 8485
$port = $basePort + $ChannelId
$content = Get-Content -LiteralPath $TemplatePath -Raw
$content = [regex]::Replace($content, '(?m)^channel_id:\s*\d+', "channel_id: $ChannelId")
$content = [regex]::Replace($content, '(?m)^port:\s*\d+', "port: $port")

$outDir = Split-Path -Parent $OutputPath
if ($outDir -and -not (Test-Path -LiteralPath $outDir)) {
    New-Item -ItemType Directory -Path $outDir -Force | Out-Null
}

$utf8NoBom = New-Object System.Text.UTF8Encoding $false
[System.IO.File]::WriteAllText($OutputPath, $content, $utf8NoBom)
Write-Host "Wrote $OutputPath (channel_id=$ChannelId port=$port)"
