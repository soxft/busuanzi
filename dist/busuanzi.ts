(function (){
    // 在此处设置您的后端地址 如 https://example.com/api
    let url: string = "http://127.0.0.1:8080/api",
        tags: string[] = ["site_pv", "site_uv", "page_pv", "page_uv"],
        current: HTMLOrSVGScriptElement = document.currentScript,
        pjax: boolean = current.hasAttribute("pjax"),
        api: string = current.getAttribute("data-api") || url,
        storageName: string = "busuanzi-identity";

    let bsz_send: Function = function () {
        let xhr: XMLHttpRequest = new XMLHttpRequest();
        xhr.open("POST", api, true);

        // set user identity
        let token = localStorage.getItem(storageName);
        if (token != null) xhr.setRequestHeader("Authorization", "Bearer " + token);
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

                        let setIdentity = xhr.getResponseHeader("Set-Bsz-Identity")
                        if (setIdentity != null && setIdentity != "") localStorage.setItem(storageName, setIdentity);
                    }
                }
            }
        }
        xhr.send();
    };
    bsz_send();

    if (!!pjax) {
        let history_pushState: Function = window.history.pushState;
        window.history.pushState = function () {
            history_pushState.apply(this, arguments);
            bsz_send();
        };

        window.addEventListener("popstate", function (_e: PopStateEvent) {
            bsz_send();
        }, false);
    }
})()