+++
title = 'sdx1.net'
template = 'raw.html'
page_template = 'raw.html'
+++
{% raw () %}
<!DOCTYPE html PUBLIC "-//W3C//DTD HTML 4.01 Transitional//EN">
<html>
<head>
  <meta content="text/html; charset=ISO-8859-1"
 http-equiv="content-type">
  <title>Seagate SER, RRER &amp; HEC</title>
  <meta content="Franc Zabkar" name="author">
  <meta content="Seagate SER, RRER &amp; HEC SMART attributes"
 name="description">
     <meta name="viewport" content="width=device-width, initial-scale=1"> 
</head>
<body style="color: rgb(0, 0, 0); background-color: rgb(255, 255, 204);"
 alink="#000099" link="#000099" vlink="#990099">
<p style="font-size: 20px; font-weight: 700">This page is an archive of <a href="http://www.users.on.net/~fzabkar/HDD/Seagate_SER_RRER_HEC.html">www.<wbr>users.on.net/<wbr>~fzabkar/<wbr>HDD/<wbr>Seagate_SER_RRER_HEC.<wbr>html</a> for <a href="./">sdx1.net/info/hardware/smart-seagate/</a></p><hr>
<p><big style="font-family: Britannic Bold;"><big><span
 style="font-weight: bold;">Seagate's Seek Error Rate,
Raw Read Error
Rate, and Hardware ECC Recovered SMART attributes</span></big></big><br>
</p>
<p>Seagate's <strong>Seek
Error Rate</strong>, <strong>Raw
Read Error
Rate</strong>, and <strong>Hardware
ECC Recovered</strong> SMART
attributes create a lot of anxiety amongst Seagate users. This is
because the raw values are typically very high, and the normalised
values (Current / Worst / Threshold) are usually quite low. Despite
this, the numbers in most cases are perfectly OK.<br>
<br>
The anxiety arises because we intuitively expect that the normalised
values should reflect a "health" score, with 100 being the ideal value.
Similarly, we would expect that the raw values should reflect an error
count, in which case a value of 0 would be most desirable. However,
Seagate calculates and applies these attribute values in a
counterintuitive way. <br>
<br>
In fact the normalised values of Seagate's Seek Error Rate, Raw Read
Error Rate, and Hardware ECC Recovered attributes are logarithmic, not
linear, and the raw values are sector counts or seek counts, not error
counts.<br>
<br>
Seagate's SMART documentation is not publicly available. The following
information has not been gleaned from any official source, but is based
on my own testing and observation, and on testing by others. Therefore
it may contain errors. <br>
<br>
<br>
<big style="font-family: Britannic Bold;"><strong>Seek
Error Rate</strong></big><br>
<br>
The raw value of each SMART attribute occupies <strong>48
bits</strong>.
Seagate's Seek Error Rate attribute consists of two parts -- a <strong>16-bit
count of seek errors</strong> in the
uppermost 4 nibbles, and a <strong>32-bit
</strong><strong>count
of seeks</strong> in the lowermost 8
nibbles. In
order to see these data, we will need a SMART utility that reports all
48 bits, preferably in hexadecimal. Two such utilities are HD Sentinel
and HDDScan.<br>
<br>
I believe the relationship between the raw and normalised values of the
SER attribute is given by ... <br>
<br>
<strong>normalised SER = -10 log
(lifetime seek errors / lifetime seeks)</strong>
<br>
<br>
In the above formula, if the drive has recorded no errors, then we
would still need to set the number of errors to 1, otherwise the result
would be indeterminate. <br>
<br>
The following table correlates the normalised SER against the actual
error rate:<br>
<br>
</p>
<pre> 90 &mdash; &lt;= 1 error per 1000 million seeks<br> 80 &mdash; &lt;= 1 error per 100 million<br> 70 &mdash; &lt;= 1 error per 10 million<br> 60 &mdash; &lt;= 1 error per million<br> 50 &mdash; 10 errors per million<br> 40 &mdash; 100 errors per million<br> 30 &mdash; 1000 errors per million<br> 20 &mdash; 10 errors per thousand</pre>
<p> <br>
A drive that has not yet recorded <strong>1
million seeks</strong>
will show <strong>100</strong>
and <strong>253</strong>
for the <strong>Current</strong>
and <strong>Worst</strong>
values. I believe this is because the data
are not considered to be statistically significant until the drive has
recorded 1 million seeks. When this target is reached, the values drop
to 60 and 60, assuming there have been no errors. <br>
<br>
By way of example, here are the SMART data for my 13GB Seagate HDD:<br>
<a target="_blank"
 href="http://www.users.on.net/%7Efzabkar/SmartUDM/13GB.RPT">http://www.users.on.net/~fzabkar/SmartUDM/13GB.RPT</a><br>
<br>
</p>
<pre> Attribute ID Threshold Value Worst Raw <br> ======================================================<br> Seek Error Rate 7 30 53 38 052E0E3000EC</pre>
<p> <br>
The number of <strong>lifetime
seek errors</strong> = <strong>0x052E</strong>
(uppermost 4 nibbles)<br>
<br>
The number of <strong>lifetime
seeks</strong> = <strong>0x0E3000EC</strong>
(lowermost 8 nibbles)<br>
<br>
Using Google's calculator ...<br>
<br>
0x052E = 1326<br>
0x0E3000EC = 238 026 988<br>
<br>
<a target="_blank"
 href="http://www.google.com/search?q=0x052E+in+decimal">http://www.google.com/search?q=0x052E+in+decimal</a><br>
<a target="_blank"
 href="http://www.google.com/search?q=0x0E3000EC+in+decimal">http://www.google.com/search?q=0x0E3000EC+in+decimal</a><br>
<br>
Applying the formula ...<br>
<br>
<strong>normalised SER = -10 log
(0x052E / 0x0E3000EC) </strong><br>
<br>
<a target="_blank"
 href="http://www.google.com/search?q=-10+log+%280x052E+/+0x0E3000EC%29">http://www.google.com/search?q=-10+log+(0x052E+/+0x0E3000EC)</a><br>
<br>
... we get a result of <strong>52.54</strong>.<br>
<br>
Here is a second example:<br>
<a target="_blank"
 href="http://www.users.on.net/%7Efzabkar/SmartUDM/120GB.RPT">http://www.users.on.net/~fzabkar/SmartUDM/120GB.RPT</a><br>
<br>
</p>
<pre> Attribute ID Threshold Value Worst Raw <br> ======================================================<br> Seek Error Rate 7 30 79 60 00000580A6AC</pre>
<p> </p>
<p>The above drive is in fact
error free. It has recorded <strong>0x0580A6AC</strong>
seeks (= 92 million) without error.<br>
<br>
Applying the formula ...<br>
<br>
<strong> normalised SER = -10 log
(1 / 0x0580A6AC) </strong><br>
<br>
... we get a result of <strong>79.65</strong><br>
<br>
Note that we have used 1 instead of 0 for the error count (because log
0 is indeterminate).</p>
<p><br>
<big style="font-family: Britannic Bold;"><strong></strong></big></p>
<p><big style="font-family: Britannic Bold;"><strong>Raw
Read Error Rate
and Hardware ECC Recovered</strong></big><span
 style="font-family: Britannic Bold;"><span style="font-weight: bold;"></span></span></p>
<p><span style="font-family: Britannic Bold;"><span
 style="font-weight: bold;"></span></span><br>
The raw values of the <strong>RRER</strong>
and <strong>HER</strong>
attributes represent a <strong>sector
count</strong>, not an error
count. This figure rolls over to 0 once the count reaches about <strong>250
million</strong>. I suspect that the
drive records the total number of
errors in each block of 250 million sectors, and then recalculates the
normalised values of each attribute accordingly. This means that RRER
and HER would be updated according to a <strong>rolling
average</strong>
rather than on a lifetime basis. I'm almost certain that the normalised
values are also logarithmic, but I'm not sure how they are calculated.
The above figure of 250 million sectors applies to the <strong>7200.11</strong>
and <strong>DiamondMax 22</strong>
models, but may not apply to all.<br>
<br>
While writing this article I came upon a Seagate document entitled "<strong>Diagnostic
Commands</strong>". It doesn't
discuss SMART attributes, but it refers
to "<strong>Error Recovery Usage
Rate</strong>" and defines it as ...<br>
<br>
<strong>Error Recovery Usage Rate
= </strong><br>
<br>
<strong>-log10 {(Number of sectors
in which controller invoked
specified error recovery scheme)/[(Number of sectors transferred) *
(512 bytes/sector) * (8 bits/byte)]}</strong><br>
<br>
This lends support for my Seek Error Rate formula, and suggests that
the RRER and HER attributes may be similarly calculated.<br>
<br>
In fact the document mentions (but does not discuss) 5 different error
recovery schemes:<br>
<br>
</p>
<ul>
  <li><strong>HARD</strong>
= multiple retries invoked and failed</li>
  <li><strong>FIRM</strong>
= multiple retries invoked </li>
  <li><strong>SOFT</strong>
= 5 retries invoked </li>
  <li><strong>OTF</strong>
= 1 retry invoked (<strong>On The
Fly</strong>)</li>
  <li><strong>RAW</strong>
= OTF ECC invoked </li>
</ul>
<p>
<br>
"<strong>On The Fly</strong>"
means that errored data is corrected
using the ECC bytes, without an additional access of the platters.<br>
<br>
Based on the abovementioned Error Recovery Usage Rate formula, I now
postulate that the normalised value of the Raw Read Error Rate
attribute could be calculated as follows:<br>
<br>
<strong> normalised RRER = -10 log
(number of errored sectors / total
bits transferred)</strong><br>
<br>
The total number of bits is ...<br>
<br>
<strong> (250 million sectors) x
(512 bytes/sector) x (8 bits/byte) =
1.024 x 10^12</strong><br>
<br>
It seems to me that it makes more sense to use a round figure, say <strong>10^12</strong>.<br>
<br>
If we now let the number of errors equal 0 (or 1), then we have ...<br>
<br>
<strong> max normalised RRER = -10
log (1 / 10^12) = 120</strong><br>
<br>
Similarly, if we let the number of errors equal 250 million (ie every
sector is errored), then we have ...<br>
<br>
<strong> min normalised RRER = -10
log (1 / 4096) = 36</strong><br>
<br>
Therefore, if my hypothesis is correct, we would expect that the
threshold value of the RRER attribute would be 36, and its maximum
possible value would be 120. In fact my Internet research tends to
confirm a maximum of <strong>120</strong>
for 7200.11 models, but the <strong>threshold</strong>
figure is <strong>6</strong>.<br>
<br>
FWIW, here are the numbers for my own Seagate drives:<br>
<br>
</p>
<pre> Attribute ID Threshold Value Worst Raw<br>===============================================================<br>Raw Read Error Rate 1 6 114 100 00000386EBBA (ST3320620A)<br>Raw Read Error Rate 1 6 64 62 00000AFD20E3 (ST3120026A)<br>Raw Read Error Rate 1 34 77 66 000007820F8F (ST340016A)<br>Raw Read Error Rate 1 0 79 78 00000753BA8E (ST313021A)<br>Hardware ECC recovered 195 0 100 63 00000C62F66E (ST3320620A)<br>Hardware ECC recovered 195 0 64 62 00000AFD20E3 (ST3120026A)<br>Hardware ECC recovered 195 0 77 66 000007820F8F (ST340016A)</pre>
<p> <br>
<a target="_blank"
 href="http://www.users.on.net/%7Efzabkar/SmartUDM/320GB.RPT">http://www.users.on.net/~fzabkar/SmartUDM/320GB.RPT</a><br>
<a target="_blank"
 href="http://www.users.on.net/%7Efzabkar/SmartUDM/120GB.RPT">http://www.users.on.net/~fzabkar/SmartUDM/120GB.RPT</a><br>
<a target="_blank"
 href="http://www.users.on.net/%7Efzabkar/SmartUDM/40GB.RPT">http://www.users.on.net/~fzabkar/SmartUDM/40GB.RPT</a><br>
<a target="_blank"
 href="http://www.users.on.net/%7Efzabkar/SmartUDM/13GB.RPT">http://www.users.on.net/~fzabkar/SmartUDM/13GB.RPT</a></p>
<br>
<span style="font-weight: bold;">Edit
#1</span>:
Some time ago a data recovery specialist provided the following
information to me:<br>
<br>
<span style="font-weight: bold;">Raw
Read Error Rate = 10 * log10(NumberOfSectorsTransferredToOrFromHost *
512 * 8 / (Number of sectors requiring retries))</span><br>
<br>
... where the factor of 512*8 is used to convert from sectors to bits.
The attribute value is only computed when the number of bits in the
transferred bits count is in the range <span style="font-weight: bold;">10^10
</span>to<span style="font-weight: bold;"> 10^12</span>.<br>
<br>
The formula is essentially the same as the aforementioned one. However,
it also produces an inconsistent threshold value, which is why I have
always felt uncomfortable with it. The formula also refers to sectors
being transferred "ToOr<span style="font-weight: bold;">From</span>Host"
which doesn't make sense for a read attribute.
<p>Nevertheless, if we ignore the
threshold
anomaly, then for each block of 10^12 bits read ...</p>
<span style="font-weight: bold;">Number
of sectors requiring retries = 10^ [(120 - normalised RRER) / 10]<br>
<br>
<br>
</span>The following table
correlates the normalised RRER against the
actual
error rate:<br>
<pre> 120 &mdash; &lt;=1 errored sector in 10^12 bits read<br> 110 &mdash; 10 errored sectors in 10^12 bits read<br> 100 &mdash; 100 errored sectors in 10^12 bits read<br> 90 &mdash; 1000 errored sectors in 10^12 bits read</pre>
<br>
<span style="font-weight: bold;">Edit
#2</span>:
<span style="font-weight: bold;">&nbsp;SandForce
SF-2000 SSD</span>s have very
similar error rate attributes.<br>
<br>
For example, their <span style="font-weight: bold;">Raw Read Error
Rate </span>is
defined as follows:<br>
<br>
&nbsp;Normalized Equation: <span style="font-weight: bold;">10log10[BitsRead
/
(ReadErrors + 1)]<br>
<br>
</span>&nbsp; &nbsp;SectorsRead = Number of sectors read<br>
&nbsp; &nbsp;SectorsToBits = 512 * 8<br>
&nbsp; &nbsp;BitsRead = SectorsRead * SectorsToBits<br>
<br>
&nbsp;Normalized Value Range:<br>
&nbsp; &nbsp; Best = 120<br>
&nbsp;&nbsp; &nbsp;Worst = 38<br>
<br>
According to a <a
 href="http://hddguardian.googlecode.com/svn/docs/Kingston%20SMART%20attributes%20details.pdf">Kingston
SF-2000 datasheet</a>, "this
Attribute reads &lsquo;120&rsquo; until a sample size between
10E10 and 10E12 is available to be tracked by this Attribute".<br>
<br>
<p><big style="font-family: Britannic Bold;"><strong>References</strong></big><br>
<br>
Here are several <strong>Usenet
discussions</strong> where I have
posted the results of my <strong>experiments</strong>:<br>
<br>
<strong>Seagate - SMART Raw Read
Error Rate test:</strong><br>
<a target="_blank"
 href="http://groups.google.com/group/comp.sys.ibm.pc.hardware.storage/browse_thread/thread/b6eb8aa2476f9cac/030c515959145d44#030c515959145d44">http://groups.google.com/group/comp.sys.ibm.pc.hardware.storage/browse_thread/thread/b6eb8aa2476f9cac/030c515959145d44#030c515959145d44</a><br>
<br>
<strong>SER, RRER, and HEC
discussion:</strong><br>
<a target="_blank"
 href="http://groups.google.com/group/comp.sys.ibm.pc.hardware.storage/browse_thread/thread/54b8ad6d34549e95/ae6ca014b3ff211a#ae6ca014b3ff211a">http://groups.google.com/group/comp.sys.ibm.pc.hardware.storage/browse_thread/thread/54b8ad6d34549e95/ae6ca014b3ff211a#ae6ca014b3ff211a</a><br>
<br>
<strong>Seek Error Rate discussion:</strong><br>
<a target="_blank"
 href="http://groups.google.com/group/comp.sys.ibm.pc.hardware.storage/browse_thread/thread/87001db5c567fb9a/63ccf100808bc3f6#63ccf100808bc3f6">http://groups.google.com/group/comp.sys.ibm.pc.hardware.storage/browse_thread/thread/87001db5c567fb9a/63ccf100808bc3f6#63ccf100808bc3f6</a><br>
<br>
<strong>A report from a Seagate
user regarding the RRER attribute:</strong><br>
<a target="_blank"
 href="http://forums.seagate.com/t5/Barracuda-XT-Barracuda-Barracuda/New-Maxtor-STM3500320AS-500GB-S-M-A-R-T-Problem/m-p/22276">http://forums.seagate.com/t5/Barracuda-XT-Barracuda-Barracuda/New-Maxtor-STM3500320AS-500GB-S-M-A-R-T-Problem/m-p/22276</a><br>
<br>
<strong>HD Sentinel (DOS / Windows
/ Linux): </strong><br>
<a target="_blank" href="http://www.hdsentinel.com/">http://www.hdsentinel.com/</a>
<br>
<br>
<strong>HDDScan for Windows: </strong><br>
<a target="_blank" href="http://hddscan.com/">http://hddscan.com/</a>
<br>
<br>
<strong>Explanation of SMART
attributes: </strong><br>
<a target="_blank" href="http://en.wikipedia.org/wiki/S.M.A.R.T.">http://en.wikipedia.org/wiki/S.M.A.R.T.</a><br>
<br>
<strong>Kingston&reg; SF-2000
Based SSD SMART Attributes:</strong><br>
<a target="_blank"
 href="http://hddguardian.googlecode.com/svn/docs/Kingston%20SMART%20attributes%20details.pdf">http://hddguardian.googlecode.com/svn/docs/Kingston%20SMART%20attributes%20details.pdf</a><br>
<br>
</p>
</body>
</html>
{% end %}
