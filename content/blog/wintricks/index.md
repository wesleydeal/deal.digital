+++
date = 2025-03-04
updated = 2025-06-27
title = "Sick Tricks for making Windows 11 Usable"
draft = false
[extra]
shorttitle = "Sick Windows 11 Tricks"
color = "#08a1f7"
+++

Inspired by [Nathan Lineback](http://toastytech.com/guis/indexwindows.html).

## Installation

Get an ISO ([consumer editions here](https://www.microsoft.com/software-download/windows11)). Write it to your install media:

* **Current Windows users:** Use [Rufus](https://rufus.ie/en/) to write it to a flash drive, and ensure that you tell it to *remove requirement for an online Microsoft account* and *disable data collection* when prompted.
* **Advanced or Linux users:** Install [Ventoy](https://www.ventoy.net/en/index.html) to your flash drive. Copy your Windows ISO to it. [Build an autounattend.xml with this tool](https://schneegans.de/windows/unattend-generator/) then use [VentoyPlugson](https://ventoy.net/en/plugin_plugson.html) to point at the XML file on the [Autoinstall](https://ventoy.net/en/plugin_autoinstall.html) page.

Boot your flash drive and install to the preferred disk. Activate using your legitimate product key and not, for example, some kind of *Microsoft Activation Script* you pulled from GitHub :^)

## Essential software & configuration

The first thing you'll need is a copy of [O&O ShutUp10](https://www.oo-software.com/en/shutup10) (or [Chris Titus's Windows Utility](https://github.com/ChrisTitusTech/winutil)). Go through and disable all telemetry and ads (though you may want to leave some, like location services).

Windows includes lots of animations that slow you down (start menu, smoother cursor in office applications, smooth Page Up/Down), and they can all be turned off in Settings > Accessibility > Visual effects > Animation.

You'll need a browser that isn't the ad-ridden spyware that is Edge--try [LibreWolf](https://librewolf.net/) or [Brave](https://brave.com/) with crypto wallet disabled.

[MSEdgeRedirect](https://github.com/rcmaehl/MSEdgeRedirect?tab=readme-ov-file) is a must if you happen to use the built-in Start menu search or weather widgets but don't want them to open links in Edge.

```batch
winget install --id=rcmaehl.MSEdgeRedirect --id=Brave.Brave --id=OO.Software.Shutup10 -e --silent --accept-source-agreements --accept-package-agreements
```

## Package Management

The quickest way I've found to get most software for Windows these days is using WinGet. I can't make promises as to exactly how well-moderated the package sources are, but I've yet to encounter
true malware within it. Here's the commands you'll wind up using:

```batch
winget search librewolf
winget add LibreWolf.LibreWolf
winget remove LibreWolf.LibreWolf
winget upgrade LibreWolf.LibreWolf
winget upgrade --all
```

## UX Tweaks

[Windhawk](https://windhawk.net/) will help you reclaim screen real estate by setting taskbar height and icon size to 36 and 20 respectively (pair with button width 32 for best results).

[AltSnap](https://github.com/RamonUnch/AltSnap) will allow you to drag and snap windows around using the Alt (or if configured, Super) key so that you don't have to target the edge of the window just to resize or find the titlebar to move it. If you're used to Linux desktop environments, it's a little itchy to live without this feature.

[PowerToys](https://learn.microsoft.com/en-us/windows/powertoys/) has lots of goodies including a color picker, window snapping tools, bulk renaming, file preview, and a macOS Spotlight search clone.
