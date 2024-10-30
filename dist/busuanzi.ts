(function (){
    // 在此处设置您的后端地址 如 https://example.com/api
    let url: string = "http://127.0.0.1:8080/api",
        tags: string[] = ["site_pv", "site_uv", "page_pv", "page_uv"],
        current: HTMLOrSVGScriptElement = document.currentScript,
        pjax: boolean = current.hasAttribute("pjax"),                          // 是否启用 pjax
        api: string = current.getAttribute("data-api") || url,                 // 自定义后端地址
        prefix: string = current.getAttribute("data-prefix") || "busuanzi",    // 自定义标签ID前缀
        style: string = current.getAttribute("data-style") || "default",       // 数字显示风格 default | comma | short
        storageName: string = "bsz-id";                                        // 本地存储名称

    let format = (num: number, style: string = 'default'): string => {
        switch (style) {
            case "comma":
                return num.toLocaleString();
            case "short": {
                let units = ["", "K", "M", "B", "T"];
                let index = 0;
                while (num >= 1000 && index < units.length - 1) {
                    num /= 1000;
                    index++;
                }
                return Math.round(num * 100) / 100 + units[index]; // 四舍五入到两位小数
            }
            default:
                return num.toString();
        }
    };

    let bsz_send = () => {
        let xhr: XMLHttpRequest = new XMLHttpRequest();
        xhr.open("POST", api, true);

        // set user identity
        let token: string | null = localStorage.getItem(storageName);
        if (token != null) xhr.setRequestHeader("Authorization", "Bearer " + token);
        xhr.setRequestHeader("Content-Type", "application/json");
        xhr.setRequestHeader("x-bsz-referer", window.location.href);
        xhr.onreadystatechange = function () {
            if (xhr.readyState === 4) {
                if (xhr.status === 200) {
                    let res: any = JSON.parse(xhr.responseText);
                    if (res.success === true) {
                        tags.map((tag: string) => {
                            let element = document.getElementById(`${prefix}_${tag}`);
                            if (element != null) element.innerHTML = format(res['data'][tag], style);

                            let container = document.getElementById(`${prefix}_container_${tag}`);
                            if (container != null) container.style.display = "inline";
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