!function(){let t=["site_pv","site_uv","page_pv","page_uv"],e=document.currentScript,a=e.hasAttribute("pjax"),n=e.getAttribute("data-api")||"http://127.0.0.1:8080/api",i=e.getAttribute("data-prefix")||"busuanzi",r=e.getAttribute("data-style")||"default",o="bsz-id",s=()=>{let e=new XMLHttpRequest;e.open("POST",n,!0);let a=localStorage.getItem(o);null!=a&&e.setRequestHeader("Authorization","Bearer "+a),e.setRequestHeader("x-bsz-referer",window.location.href),e.onreadystatechange=function(){if(4===e.readyState&&200===e.status){let a=JSON.parse(e.responseText);if(!0===a.success){t.map((t=>{let e=document.getElementById(`${i}_${t}`);null!=e&&(e.innerHTML=((t,e="default")=>{if("comma"===e)return t.toLocaleString();if("short"===e){const e=["","K","M","B","T"];let a=Math.floor(Math.log10(t)/3);return t/=Math.pow(1e3,a),`${Math.round(100*t)/100}${e[a]}`}return t.toString()})(a.data[t],r));let n=document.getElementById(`${i}_container_${t}`);null!=n&&(n.style.display="inline")}));let n=e.getResponseHeader("Set-Bsz-Identity");null!=n&&""!=n&&localStorage.setItem(o,n)}}},e.send()};if(s(),a){let t=window.history.pushState;window.history.pushState=function(){t.apply(this,arguments),s()},window.addEventListener("popstate",(function(t){s()}),!1)}}();