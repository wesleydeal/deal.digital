if(localStorage.getItem("mode")=="dark")
    document.documentElement.classList.add('dark');

document.addEventListener('DOMContentLoaded', function(){ 
    document.body.innerHTML += '<a class="toggle" onclick="toggle()"></a>';
}, false);


function darkmode(){
    document.documentElement.classList.add('dark');
    localStorage.setItem("mode", "dark");
}

function nodark(){
    document.documentElement.classList.remove('dark');
    localStorage.setItem("mode", "light");
}

function toggle(){
    if(document.documentElement.classList.contains('dark')) nodark();
    else darkmode();
}