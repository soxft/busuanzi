let scriptTag: HTMLScriptElement;

const url = "https://busuanzi.9420.ltd/api?callback=BusuanziCallback";
const tags:Array<any> = ["site_pv","site_uv", "page_pv", "page_uv"];

const fetchUrl = (url: string, callback: Function) => {
    let str: string = "BusuanziCallback_" + Math.floor(1099511627776 * Math.random()).toString();
    window[str] = function(callback: Function){
        return function(a) {
            try {
                callback(a);
                scriptTag.parentElement.removeChild(scriptTag)
            } catch (c) {}
        }
    }(callback)
    scriptTag = document.createElement("script");
    scriptTag.type = "text/javascript"
    scriptTag.defer = true;
    scriptTag.src = url.replace("BusuanziCallback",str);
    scriptTag.referrerPolicy = "no-referrer-when-downgrade"
    document.getElementsByTagName("head")[0].appendChild(scriptTag)
   }

fetchUrl(url, function(a) {
    tags.map(tag => {
        let ele = document.getElementById(`busuanzi_${tag}`);
        if (ele) {
            ele.innerHTML = a[tag];
        }
    })
})
