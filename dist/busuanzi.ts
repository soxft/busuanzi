(function (){
    let url: string = "http://127.0.0.1:8080/api?rand=" + Math.random().toFixed(6);
    let tags: string[] = ["site_pv", "site_uv", "page_pv", "page_uv"];

    let xhr: XMLHttpRequest = new XMLHttpRequest();
    xhr.open("GET", url, true);

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
})()