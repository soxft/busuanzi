var scriptTag;
var url = "https://busuanzi.9420.ltd/api?callback=BusuanziCallback";
var tags = ["site_pv", "site_uv", "page_pv", "page_uv"];
var fetchUrl = function (url, callback) {
    var str = "BusuanziCallback_" + Math.floor(1099511627776 * Math.random()).toString();
    window[str] = function (callback) {
        return function (a) {
            try {
                callback(a);
                scriptTag.parentElement.removeChild(scriptTag);
            }
            catch (c) { }
        };
    }(callback);
    scriptTag = document.createElement("script");
    scriptTag.type = "text/javascript";
    scriptTag.defer = true;
    scriptTag.src = url.replace("BusuanziCallback", str);
    scriptTag.referrerPolicy = "no-referrer-when-downgrade";
    document.getElementsByTagName("head")[0].appendChild(scriptTag);
};
fetchUrl(url, function (a) {
    tags.map(function (tag) {
        var ele = document.getElementById("busuanzi_".concat(tag));
        if (ele) {
            ele.innerHTML = a[tag];
        }
    });
});
