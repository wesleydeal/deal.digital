+++
title = 'sdx1.net'
template = 'raw.html'
page_template = 'raw.html'
+++
{% raw () %}
<!DOCTYPE html>
<head>
    <title>Reading SMART Data on Seagate Drives - sdx1.net</title>
    <link href="/archive/sdx1/style.css" rel="stylesheet">
    <script type="text/javascript" src="/archive/sdx1/lightswitch.js"></script>
	<meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1"> 
    <meta name="theme-color" content="#333">
	<meta name="description" content="Lenovo's Corporate Discount, which currently uses the code NJ*PERKSEPP, offers a significant discount on expensive ThinkPads.">
</head>
<body>
    <nav>
        <a href="../../..">sdx1.net</a> > <a href="../..">Info</a> > <a href="../../#hardware">Hardware</a> > <a>Reading SMART Data on Seagate Drives</a>
    </nav>
    <div class="content">
        <h1>Reading SMART Data on Seagate Drives</h1>
		<p class="meta">Updated 2018-10-05. Written by sdx1 on 2017-09-17.</p>
        <p>Seagate hard drives often report extremely high read and seek error rates in SMART data. After pronouncing many such drives defective, I did some research on the subject and found that Seagate drives use non-standard formatting for read and seak error rates, and every piece of software I've used reports these values incorrectly. While the easiest advice on this subject is to ignore these values in the absence of other drive issues, there is a way to read them.</p>
		<h2>The Process</h2>
		<p>The vast majority of programs that read SMART data rely on smartctl (<a href="https://sourceforge.net/projects/smartmontools/files/latest/download?source=files" target="_blank">Windows download</a>), which we'll use to read these values.</p>
		<p>Assuming that our disk is located at /dev/hda (which is typical for smartctl on Windows and older Linuxes), we can get the disk status using <code>smartctl -A /dev/hda</code>. The raw read error rate will be reported as something like <strong>130917967</strong> and the seek something like <strong>13219996990</strong>. This is clearly wrong, because the drive hasn't failed or reported any SMART errors, and the normalized value still shows this attribute to be within acceptable margins.</p>
		<p>This is because the read and seek error rates are actually recorded as 48-bit hexadecimal values, where the first 16 bits (4 hexadecmal digits) represent the total number of read or seek errors and the last 32 bits (8 hex digits) represent the total number of reads or seeks attempted. We can get these values with the command <code>smartctl -v 1,hex48 -v 7,hex48 -A /dev/hda</code>.</p>
		<p>Using this command, our read and seek error values become <strong>0x<u>0000</u>07cda64f</strong> and <strong>0x<u>0003</u>13f9253e</strong> respectively. The first four digits following the x (emphasized) are the total number of errors. As you can see, read and seek errors are only 0 and 3 respectively for this particular drive. And the total number of reads are 130917967 and 335095102 respectively, obtained by taking the last eight digits and converting from hexadecimal to decimal numbers.</p>
		<h2>Additional Reading</h2>
		<ul>
			<li>my primary source: <a href="http://www.users.on.net/~fzabkar/HDD/Seagate_SER_RRER_HEC.html">an article by <em>fzabkar</em></a>, <a href="archive.html">archived here.</a></li>
			<li>a <a href="http://sgros.blogspot.com/2013/01/seagate-disk-smart-values.html">good blog post on the subject by Stjepan Groš</a></li>
			<li>this <a href="https://superuser.com/questions/393257/brand-new-seagate-hdd-has-high-raw-read-error-rate">StackExchange SuperUser question regarding the subject</a></li>
		</ul>
    </div>
</body>
{% end %}