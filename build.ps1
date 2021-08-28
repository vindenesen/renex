$now = Get-Date -UFormat "%Y-%m-%d_%T"
$sha1 = (git rev-parse HEAD).Trim()
$release = ""

if ($args.Length -ne 1) {
    Write-Error "Missing version info"
    exit 1
}

$version = ($args[0]).Split(".")
if ($version.Length -ne 4) {
    Write-Error "Wrong version info, needs four levels: 1.2.3.4"
    exit 1
}

$release = $args[0]

$json = @"
{
    "FixedFileInfo": {
        "FileVersion": {
            "Major": $($version[0]),
            "Minor": $($version[1]),
            "Patch": $($version[2]),
            "Build": $($version[3])
        },
        "ProductVersion": {
            "Major": $($version[0]),
            "Minor": $($version[1]),
            "Patch": $($version[2]),
            "Build": $($version[3])
        },
        "FileFlagsMask": "3f",
        "FileFlags ": "00",
        "FileOS": "040004",
        "FileType": "01",
        "FileSubType": "00"
    },
    "StringFileInfo": {
        "Comments": "",
        "CompanyName": "Trond Vindenes",
        "FileDescription": "Bulk rename using regular expressions",
        "FileVersion": "",
        "InternalName": "renex",
        "LegalCopyright": "Trond Vindenes",
        "LegalTrademarks": "",
        "OriginalFilename": "renex.exe",
        "PrivateBuild": "",
        "ProductName": "renex",
        "ProductVersion": "v$($release)",
        "SpecialBuild": ""
    },
    "VarFileInfo": {
        "Translation": {
            "LangID": "0409",
            "CharsetID": "04B0"
        }
    },
    "IconPath": "",
    "ManifestPath": ""
}
"@

$json | Out-File "versioninfo.json" -Encoding Ascii -Force
.\goversioninfo.exe
go build -ldflags "-X main.sha1Ver=$sha1 -X main.buildTime=$now -X main.release=$release -s -w"
Remove-Item "versioninfo.json"
Remove-Item "resource.syso"