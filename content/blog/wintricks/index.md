+++
date = 2025-03-04
updated = 2025-03-04
title = "Sick Windows Tricks"
draft = true
+++

Or, how I make Windows 11 borderline usable. Inspired by [Nathan Lineback](http://toastytech.com/guis/indexwindows.html).

## Installation

Get an ISO ([consumer editions here](https://www.microsoft.com/software-download/windows11)). Use [Rufus](https://rufus.ie/en/) to write it to a flash drive, and ensure that you tell it to *remove requirement for an online Microsoft account* and *disable data collection* when prompted.

Boot your flash drive and install to the preferred disk. Activate using your legitimate product key :^)

## Essential software & configuration

The first thing you'll need is a copy of [O&O ShutUp10](https://www.oo-software.com/en/shutup10) (or [Chris Titus's Windows Utility](https://github.com/ChrisTitusTech/winutil)). Go through and disable all telemetry and ads (though you may want to leave some, like location services).

Windows includes lots of animations that slow you down (start menu, smoother cursor in office applications, smooth Page Up/Down), and they can all be turned off in Settings > Accessibility > Visual effects > Animation.

You'll need a browser that isn't the ad-ridden spyware that is Edge--try [LibreWolf](https://librewolf.net/) or [Brave](https://brave.com/) with crypto wallet disabled.

## UX Tweaks

[Windhawk] will help you reclaim screen real estate by setting *taskbar height and icon size* to 36 and 20 respectively (pair with button width 32 for best results).

[AltSnap](https://github.com/RamonUnch/AltSnap) will allow you to drag and snap windows around using the Alt (or if configured, Super) key so that you don't have to target the edge of the window just to resize or find the titlebar to move it. It takes a little getting used to, but it makes things a good bit quicker.

[PowerToys](https://learn.microsoft.com/en-us/windows/powertoys/) has lots of goodies including a color picker, window snapping tools, bulk renaming, file preview, and a macOS Spotlight search clone.