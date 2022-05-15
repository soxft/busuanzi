(function (){
    let url: string = "http://127.0.0.1:8080/api?rand=" + Math.random();
    let tags: string[] = ["site_pv", "site_uv", "page_pv", "page_uv"];

    let bsz_send: Function = function () {
        let xhr: XMLHttpRequest = new XMLHttpRequest();
        xhr.open("POST", url, true);

        xhr.setRequestHeader("x-bsz-referer", window.location.href);
        xhr.onreadystatechange = function () {
            if (xhr.readyState === 4) {
                if (xhr.status === 200) {
                    let res: any = JSON.parse(xhr.responseText);
                    if (res.success === true) {
                        tags.map((tag: string) => {
                            let ele = document.getElementById(`busuanzi_${tag}`);
                            if (ele) ele.innerHTML = res['data'][tag];
                        })
                    }
                }
            }
        }
        xhr.send();
    };

    let history_pushState: Function = window.history.pushState;
    window.history.pushState = function () {
        history_pushState.apply(this, arguments);
        bsz_send();
    };

    window.addEventListener("popstate", function (_e: PopStateEvent) {
        bsz_send();
    }, false);
    bsz_send();
})()