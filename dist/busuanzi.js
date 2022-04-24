let xhr = new XMLHttpRequest();
xhr.open('GET', '//busuanzi.9420.ltd/api', true);
xhr.send(null);
xhr.onreadystatechange = function () {
    if (xhr.readyState === 4) {
        if (xhr.status === 200) {
            let data = JSON.parse(xhr.responseText);
            if (data.success === true) {
                document.getElementById("busuanzi_page_pv").innerText = data.data['page_pv'];
                document.getElementById("busuanzi_page_uv").innerText = data.data['page_uv'];
                document.getElementById("busuanzi_site_pv").innerText = data.data['site_pv'];
                document.getElementById("busuanzi_site_uv").innerText = data.data['site_uv'];
            }
        }
    }
};
