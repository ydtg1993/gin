(function (w) {
    w.parseParams = function (uri, params) {
        let paramsArray = [];
        if (Object.keys(params).length === 0) return uri;
        Object.keys(params).forEach(key => params[key] && paramsArray.push(`${key}=${params[key]}`));
        if (uri.search(/\?/) === -1) {
            uri += `?${paramsArray.join('&')}`
        } else {
            uri += `&${paramsArray.join('&')}`
        }
        return uri
    };
    w.request = function (settings) {
        settings = Object.assign({
            url: '',
            method: "GET",
            header: {},
            data: {},
            callback: null,
            error_callback: null,
            timeout: 30000
        }, settings);
        let xhr = new XMLHttpRequest();
        if (settings.method === 'GET') settings.url = w.parseParams(settings.url, settings.data);
        xhr.open(settings.method, settings.url, true);
        xhr.timeout = settings.timeout;
        let token = '';
        if (document.querySelector('meta[name="csrf-token"]')) {
            token = document.querySelector('meta[name="csrf-token"]').getAttribute('content');
        }
        xhr.setRequestHeader("X-CSRF-TOKEN", token);
        if (settings.header) {
            for (let header in settings.header) {
                if (!settings.header.hasOwnProperty(header)) continue;
                xhr.setRequestHeader(header, settings.header[header]);
            }
        }
        if (settings.method === 'GET') {
            xhr.setRequestHeader("Content-type", "application/text;charset=UTF-8");
            xhr.responseType = "text";
            xhr.send(null);
        } else {
            if (settings.data instanceof FormData) {
                xhr.setRequestHeader("X-Requested-With", "XMLHttpRequest");
            } else {
                xhr.setRequestHeader("Content-Type", "application/json;charset=UTF-8");
                settings.data = JSON.stringify(settings.data);
            }
            xhr.responseType = "json";
            xhr.send(settings.data);
        }
        xhr.onload = function () {
            if (this.status >= 200 && this.status < 300) return;
            if (typeof settings.error_callback === 'function') settings.error_callback();
        };
        xhr.onreadystatechange = function () {
            if (xhr.readyState === xhr.DONE && xhr.status === 200) {
                let response = xhr.response;
                if (typeof settings.callback === 'function') settings.callback(response);
            }
        };
        xhr.onerror = function (e) {
            console.error(e);
            if (typeof settings.error_callback === 'function') settings.error_callback(e);
        };
    };

    w.search = class {
        constructor() {
            this.input_dom = document.querySelector("#search>input");
            this.button_dom = document.querySelector("#search>.search-button");
            this.button_dom.addEventListener('click',()=>{
                if(!this.input_dom.value)return;
                window.location.href = new URL(window.location.href).origin + "/search?keywords="+this.input_dom.value;
            });
        }
    };
    new w.search();
})(window);