function forge {
    param(
        [Parameter(Mandatory)]
        [string]$Url,

        [Parameter()]
        [switch]$Rebuild  # Optional switch to force rebuild
    )

    $image  = "video-forge"
    $output = Join-Path $PWD "output"

    # If rebuild requested, remove old image first
    if ($Rebuild) {
        Write-Host "🔨 Forcing image rebuild..." -ForegroundColor Yellow
        docker rmi $image -Force > $null 2>&1
    }

    # Check if build is needed (not present or removed above)
    docker image inspect $image > $null 2>&1
    if ($LASTEXITCODE -ne 0) {
        Write-Host "📦 Building Docker image..." -ForegroundColor Cyan
        docker build -t $image .
        if ($LASTEXITCODE -ne 0) {
            throw "❌ Docker build failed"
        }
    }

    if (-not (Test-Path $output)) {
        New-Item -ItemType Directory -Force -Path $output | Out-Null
    }

    Write-Host "🚀 Running pipeline..." -ForegroundColor Green
    docker run --rm `
        -v "${output}:/app/output" `
        --add-host=host.docker.internal:host-gateway `
        $image $Url
}