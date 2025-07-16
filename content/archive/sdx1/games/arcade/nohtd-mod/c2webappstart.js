jQuery(window).resize(function()
{if(window.c2resizestretchmode===1)
{window.c2resizestretchmode=2;var canvas=document.getElementById("c2canvas");window.c2oldcanvaswidth=canvas.width;window.c2oldcanvasheight=canvas.height;window.c2eventtime=Date.now();var w=jQuery(window).width();var h=jQuery(window).height();cr_sizeCanvas(w,h);}
else if(window.c2resizestretchmode===2)
{if(Date.now()>window.c2eventtime+ 50)
{window.c2resizestretchmode=0;cr_sizeCanvas(window.c2oldcanvaswidth,window.c2oldcanvasheight);}}});jQuery(document).ready(function()
{cr_createRuntime("c2canvas");});function onVisibilityChanged()
{if(document.hidden||document.mozHidden||document.webkitHidden||document.msHidden)
cr_setSuspended(true);else
cr_setSuspended(false);};document.addEventListener("visibilitychange",onVisibilityChanged,false);document.addEventListener("mozvisibilitychange",onVisibilityChanged,false);document.addEventListener("webkitvisibilitychange",onVisibilityChanged,false);document.addEventListener("msvisibilitychange",onVisibilityChanged,false);