var bsz_fetch = function () {
    var url = "http://127.0.0.1:8080/api?rand=" + Math.random().toFixed(6);
    var xhr = new XMLHttpRequest();
    xhr.open("GET", url, true);
    var referer = window.location.href;
    xhr.setRequestHeader("x-bsz-referer", referer);
    xhr.onreadystatechange = function () {
        if (xhr.readyState === 4) {
            if (xhr.status === 200) {
                var res_1 = JSON.parse(xhr.responseText);
                if (res_1.success === true) {
                    ["site_pv", "site_uv", "page_pv", "page_uv"].map(function (tag) {
                        var ele = document.getElementById("busuanzi_".concat(tag));
                        if (ele)
                            ele.innerHTML = res_1['data'][tag];
                    });
                }
            }
        }
    };
    xhr.send();
};
bsz_fetch();
