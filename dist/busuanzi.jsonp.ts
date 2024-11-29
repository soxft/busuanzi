(function () {
    let url: string = "http://127.0.0.1:8080/jsonp",
        tags: string[] = ["site_pv", "site_uv", "page_pv", "page_uv"],
        current: HTMLOrSVGScriptElement = document.currentScript,
        pjax: boolean = current.hasAttribute("pjax"),
        api: string = current.getAttribute("data-api") || url,
        prefix: string = current.getAttribute("data-prefix") || "busuanzi",
        style: string = current.getAttribute("data-style") || "default";

    let format = (num: number, style: string = 'default'): string => {
        if (style === "comma") return num.toLocaleString();
        if (style === "short") {
            const units = ["", "K", "M", "B", "T"];
            let index = Math.floor(Math.log10(num) / 3);
            num /= Math.pow(1000, index);
            return `${Math.round(num * 100) / 100}${units[index]}`;
        }
        return num.toString();
    };

    let bsz_send = () => {
        let callbackName = `BszCallback_${Math.round(100000 * Math.random())}_${Date.now().toString().slice(-5)}`;
        let script = document.createElement('script');
        script.src = `${api}?callback=${callbackName}`;
        script.type = "text/javascript";
        script.defer = !0;
        script.referrerPolicy = "no-referrer-when-downgrade";

        // @ts-ignore
        window[callbackName] = (res: {
            site_pv: number,
            site_uv: number,
            page_pv: number,
            page_uv: number,
            token: string | null,
        }) => {
            try {
                tags.map((tag: string) => {
                    let element = document.getElementById(`${prefix}_${tag}`);
                    if (element != null) { // @ts-ignore
                        element.innerHTML = format(res[tag], style);
                    }

                    let container = document.getElementById(`${prefix}_container_${tag}`);
                    if (container != null) container.style.display = "inline";
                });
            } catch (e) {
                console.error(e);
            } finally {
                document.body.removeChild(script);
                // @ts-ignore
                delete window[callbackName];
            }
        };

        document.body.appendChild(script);
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
})();