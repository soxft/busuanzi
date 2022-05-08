const bsz_fetch = () => {
    let url: string = "http://127.0.0.1:8080/api?rand=" + Math.random().toFixed(6);

    let xhr: XMLHttpRequest = new XMLHttpRequest();
    xhr.open("GET", url, true);
    // post
    let referer: string = window.location.href;
    xhr.setRequestHeader("x-bsz-referer", referer);
    xhr.onreadystatechange = function () {
        if (xhr.readyState === 4) {
            if (xhr.status === 200) {
                let res = JSON.parse(xhr.responseText);
                if (res.success === true) {
                    ["site_pv", "site_uv", "page_pv", "page_uv"].map(tag => {
                        let ele = document.getElementById(`busuanzi_${tag}`);
                        if (ele) ele.innerHTML = res['data'][tag];
                    })
                }
            }
        }
    }
    xhr.send();
}

bsz_fetch()
