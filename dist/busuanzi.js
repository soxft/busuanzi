!function() {
    document.getElementById("busuanzi_container_site_pv").style.display = "none"
    document.getElementById("busuanzi_container_page_pv").style.display = "none"
    document.getElementById("busuanzi_container_site_uv").style.display = "none"
    var t = ["site_pv", "site_uv", "page_pv", "page_uv"]
      , e = document.currentScript
      , a = e.hasAttribute("pjax")
      , n = e.getAttribute("data-api") || "http://127.0.0.1:8080/api"
      , r = e.getAttribute("data-prefix") || "busuanzi"
      , i = "bsz-id"
      , s = function() {
        var e = new XMLHttpRequest;
        e.open("POST", n, !0);
        var a = localStorage.getItem(i);
        null != a && e.setRequestHeader("Authorization", "Bearer " + a),
        e.setRequestHeader("x-bsz-referer", window.location.href),
        e.onreadystatechange = function() {
            if (4 === e.readyState && 200 === e.status) {
                var a = JSON.parse(e.responseText);
                if (!0 === a.success) {
                    t.map((function(t) {
                        var e = document.getElementById("".concat(r, "_").concat(t));
                        e && (e.innerHTML = a.data[t])
                    }
                    ));
                    var n = e.getResponseHeader("Set-Bsz-Identity");
                    null != n && "" != n && localStorage.setItem(i, n)
                }
            }
        }
        ,
        e.send()
    };
    if (s(),
    a) {
        var o = window.history.pushState;
        window.history.pushState = function() {
            o.apply(this, arguments),
            s()
        }
        ,
        window.addEventListener("popstate", (function(t) {
            s()
        }
        ), !1)
    }
    document.getElementById("busuanzi_container_site_pv").style.display = "inline"
    document.getElementById("busuanzi_container_page_pv").style.display = "inline"
    document.getElementById("busuanzi_container_site_uv").style.display = "inline"
}();
