+++
author = "Wesley Deal"
title = "Reading & Interpreting SMART Read/Seek Error Data on Seagate Drives"
date = 2017-09-17
updated = 2024-05-02
[taxonomies]
tags = ["hardware", "software"]
[extra]
shorttitle = "Interpret Seagate read/seek error rates"
+++

Seagate hard drives often report extremely high read and seek error rates in SMART data. After pronouncing a string of drives defective based on these fields, I did some digging and found that Seagate drives use non-standard formatting for these values. Every piece of disk-monitoring software I had used reports these values incorrectly. Incidentally, this was good for our bottom line at the computer shop where we sold quite a few customers SSD upgrades partly on the basis of these mistaken readings. While the easiest advice on this subject is to ignore these values in the absence of other drive issues, there is a way to read them.

## The Process

The vast majority of programs that read SMART data rely on smartctl ([Windows download](https://sourceforge.net/projects/smartmontools/files/latest/download?source=files)), which we'll use to read these values.

Assuming that our disk is located at /dev/hda (which is typical for smartctl on Windows and older Linuxes), we can get the disk status using `smartctl -A /dev/hda`. The raw read error rate will be reported as something like **130917967** and the seek something like **13219996990**. This is clearly wrong, because the drive hasn't failed or reported any SMART errors, and the normalized value still shows this attribute to be within acceptable margins.

This is because the read and seek error rates are actually recorded as 48-bit hexadecimal values, where the first 16 bits (4 hex digits) represent the total number of read or seek errors and the last 32 bits (8 hex digits) represent the total number of reads or seeks attempted. We can get these values with the command `smartctl -v 1,hex48 -v 7,hex48 -A /dev/hda`.

Using this command, our read and seek error values become **0x<u>0000</u>07cda64f** and **0x<u>0003</u>13f9253e** respectively. The first four digits following the x (emphasized) are the total number of errors. As you can see, read and seek errors are only 0 and 3 respectively for this particular drive. And the total number of reads are 130917967 and 335095102 respectively, obtained by taking the last eight digits and converting from hexadecimal to decimal numbers.

## Additional Reading

- my primary source: [an article by *fzabkar*](http://web.archive.org/web/20211208195204/http://www.users.on.net/~fzabkar/HDD/Seagate_SER_RRER_HEC.html)

- a [good blog post on the subject by Stjepan Groš](http://sgros.blogspot.com/2013/01/seagate-disk-smart-values.html)

- this [StackExchange SuperUser question](https://superuser.com/questions/393257/brand-new-seagate-hdd-has-high-raw-read-error-rate)
