(function () {
    var url = "http://127.0.0.1:8080/?rand=" + Math.random();
    var tags = ["site_pv", "site_uv", "page_pv", "page_uv"];
    var xhr = new XMLHttpRequest();
    xhr.open("POST", url, true);
    xhr.setRequestHeader("x-bsz-referer", window.location.href);
    xhr.onreadystatechange = function () {
        if (xhr.readyState === 4) {
            if (xhr.status === 200) {
                var res_1 = JSON.parse(xhr.responseText);
                if (res_1.success === true) {
                    tags.map(function (tag) {
                        var ele = document.getElementById("busuanzi_".concat(tag));
                        if (ele)
                            ele.innerHTML = res_1['data'][tag];
                    });
                }
            }
        }
    };
    xhr.send();
})();
