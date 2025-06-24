+++
title = "Prune-Backups.ps1: Prune Backup Files with PowerShell"
date = 2024-05-14
[taxonomies]
tags = ["software", "powershell"]
[extra]
shorttitle = "Prune-Backups.ps1"
color = "#02235b"
+++

Do you have a program that dumps daily backups into a folder which would eat up your disk? Do you want to automatically delete some of those files but keep backups of certain ages? This is just the case for which I wrote `Prune-Backups.ps1`

```ps1
param (
	[string]$Path = (Get-Location),
	[int]$Hours = 0,
	[int]$Days = 0,
	[int]$Weeks = 0,
	[int]$Months = 0,
	[int]$Years = 0,
	[switch]$Summary = $true,
	[switch]$WhatIf
)

$now = Get-Date
$hourlyDates  = 0..-$Hours  | % {$now.addHours($_)}
$dailyDates   = 0..-$Days   | % {$now.addDays($_)}
$weeklyDates  = 0..-$Weeks  | % {$now.addDays($_*7)}
$monthlyDates = 0..-$Months | % {$now.addMonths($_)}
$yearlyDates  = 0..-$Years  | % {$now.addYears($_)}
$datesToKeep  = @() + $hourlyDates + $dailyDates + $monthlyDates + $weeklyDates + $yearlyDates | Sort -Unique

$files = Get-ChildItem -Path $Path -File

function Get-ClosestFile {
	param (
		[datetime]$referenceDate,
		[array]$files
	)

	foreach ($file in $files) {
		$diff = (New-TimeSpan -Start $file.LastWriteTime -End $referenceDate).Duration()

		if (!$smallestDiff -or $diff -lt $smallestDiff) {
			$smallestDiff = $diff
			$closestFile = $file
		}
	}
	
	return $closestFile
}

$filesToKeep = @()

foreach ($date in $datesToKeep) {
	$file = Get-ClosestFile -referenceDate $date -files $files
	if ($file -ne $null -and $filesToKeep -notcontains $file) {
		$filesToKeep += $file
	}
}

$filesToKeep = $filesToKeep | Sort LastWriteTime -Unique
$lengthDeleted = 0

$files | ? {$_ -notin $filesToKeep} | % {
	$lengthDeleted += $_.Length
	if ($WhatIf) {
		Write-Output "WhatIf: Deleting file '$($_.FullName)' [$($_.LastWriteTime)]"

	} else {
		Write-Output "Deleting file '$($_.FullName)' [$($_.LastWriteTime)]"
		Remove-Item $_ -Force
	}
}

Write-Output "Operation complete."
if ($summary){
	Write-Output "Retained files:"
	$filesToKeep | % { " [$($_.LastWriteTime)] $($_.FullName)" }

	$deleted   = [math]::round($lengthDeleted/1Mb, 3)
	$remaining = [math]::round(($filesToKeep | Measure -Sum Length).Sum/1Mb, 3)
	$percentDeleted = [math]::round($deleted/($remaining+$deleted), 3)*100
	Write-Output "Removed $($deleted)MiB ($($percentDeleted)%), kept $($remaining)MiB"
}
```

Let's say you want to keep a backup from every day of the last week, as well as copies from 2, 3, and 4 weeks ago, 1, 2, and 3 months ago, 1 and 2 years ago. You can try it from the folder in question with the `-WhatIf` flag:

```ps1
./Prune-Backups.ps1 -Days 7 -Weeks 4 -Months 3 -Years 2 -WhatIf
```

It will report back which files would be deleted, which would be kept, and how much space you'll save.

Now you can create a scheduled daily task, pointing to a specific path and removing the `-WhatIf` flag:

>**Program/script**: `C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe`  
>**Add arguments (optional)**: `-NoProfile -NonInteractive -ExecutionPolicy Bypass -File C:\Scripts\Prune-Backups.ps1 -Path "C:\ProgramData\Vendor\Backups\" -Days 7 -Weeks 4 -Months 3 -Years 2`

**Note:** No matter what, the most recent item in the folder will be kept. The default behavior of the script, if called with no parameters, is to keep only the most recent file.
