+++
author = "Wesley Deal"
title = "Scroll & middle-click with the TrackPoint"
date = 2017-08-06
updated = 2024-04-25
[taxonomies]
tags = ["software", "thinkpad"]
[extra]
shorttitle = "Scroll & middle-click with a TrackPoint"
color = "#E10020"
+++

<!--<style>
:root{
	--color-accent: #E10020;
	--blockquote-bg: linear-gradient(135deg, #E1002020, #E1002000);
	--color-hover: #E10020;
	--color-logo-bg: #E10020;
	--color-logo-fg: #fff;
}
</style>-->

On ThinkPads running Windows, the middle mouse button below the TrackPoint can be used as a scroll modifier for the trackpoint or a middle-click/MOUSE3 button, but not both at the same time. This page covers ways to resolve this.

## TPMiddle
TPMiddle.exe is a brilliant program by Gerhard Wiesinger that allows both functions at once by detecting whether or not the mouse is moved between press and release.

This tool has some issues on newer devices. Suggested replacements include [W10Wheel.NET](https://github.com/ykon/w10wheel.net).

### Downloads

* [tpmiddle-0.6.exe](tpmiddle-0.6.exe) - **Recommended**, especially for Windows 10 users, 10.5kb

* [tpmiddle-0.7.exe](tpmiddle-0.7.exe) - Incompatible with some newer systems, 11kb  
also available from the [author's site](https://www.wiesinger.com/opensource/TPMiddle/V0.7/Release/tpmiddle.exe)

### Installation

1. Make sure you have [Microsoft Visual C++ 2010 Redistributables](https://www.microsoft.com/en-us/download/details.aspx?id=26999) installed
2. Ensure your middle button is set to "Use for scrolling" under Control Panel > Mouse > ThinkPad
3. Download the file
4. Drop it into `%appdata%\Microsoft\Windows\Start Menu\Programs\Startup`
5. Log out and back in to start it

## Alternatives

An [AutoHotKey-based substitute](http://web.archive.org/web/20220125153630/http://forum.notebookreview.com/threads/ultranav-middle-click-button-scroll.423415/) was listed in the Notebook Review forums a few years ago.


### W10Wheel.NET
[W10Wheel.NET](https://github.com/ykon/w10wheel.net/releases/latest) by ykon ([wiki](https://github.com/ykon/w10wheel.net/wiki)) seems to be the best option. You'll need [.NET Framework 4.x](https://dotnet.microsoft.com/en-us/download/dotnet-framework/thank-you/net48-web-installer).

As per [/u/JakandClank on reddit](https://old.reddit.com/r/thinkpad/comments/q2cfup/tpmiddle_on_newer_thinkpads/):

>I've got a pretty decent workaround with w10wheel that works just about as well as tpmiddle did. Here's how I do it:
>
> 1. Download w10wheel and unzip it somewhere.
> 2. Follow the instructions on [this page](https://www.tenforums.com/tutorials/57690-create-elevated-shortcut-without-uac-prompt-windows-10-a.html) to create a UAC-less shortcut for the program. This will allow the program to run as administrator without having to allow it every time with the User Account Control prompt. You want to run it as admin so you can scroll on things like Task Manager and your clipboard history.
> 3. Copy the shortcut into this directory: `C:\Users\(username)\AppData\Roaming\Microsoft\Windows\Start Menu\Programs\Startup`
>
> Now the w10wheel program will be installed (by just unzipping it, there's no installer) and will launch at startup without an annoying UAC prompt.
>
> To make w10wheel work just like tpmiddle, I did this:
>
> 1. Go to the trackpoint settings by going to Start-> start typing "mouse settings" -> click on "Mouse settings".
> 2. Under "Synaptics Trackpoint" click on **Settings**.
> 3. Set the **Middle Button Action** to **Middle Button**.
> 4. Launch w10wheel if it's not running already.
> 5. Right click on the w10wheel tray icon.
> 6. Set **Trigger** to **MiddleDrag**. This makes it so that holding down the middle button will allow you to scroll. Setting it to "Middle" makes it a toggle and I don't recommend that.
> 7. Make sure that **Horizontal Scroll** is checked in the same menu.
> 8. In the **Set Number** menu, set **verticalThreshold** to 1 and **horizontalThreshold** to 4. These values are the "resistance" to your mouse movement. It maxes at 500 but really anything past 5 is hard on a trackpoint but YMMV. If you're gonna adjust it I'd recommend doing it a few ticks at a time.
> 9. *Optional* Enable **Cursor Change** in the menu. This gives you a visual indicator that you are scrolling. I am yet to find a way to change which cursor type it changes to, unfortunately. I don't really like the "resize window" cursor it gives you.
