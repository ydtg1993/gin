document.addEventListener("DOMContentLoaded", function() {
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

        w.Pagination = class Pagination {
            constructor(selector, totalPages, currentPage) {
                if (selector instanceof HTMLElement) {
                    this.DOM = selector;
                } else {
                    this.DOM = document.querySelector(selector);
                }
                this.totalPages = totalPages;
                this.currentPage = currentPage;
            }

            visiblePageCount(page){
                this.visiblePagesCount = page - 2
            }

            bindEvent(func){
                if (typeof func !== 'function')return this;
                this.callback = func;
                return this;
            }

            format(url){
                this.url = url;
                return this;
            }

            make() {
                if(this.totalPages <= 1)return;
                if(this.currentPage > this.totalPages || this.currentPage < 1)this.currentPage = 1;
                let paginationContainer = document.createElement("ul");
                paginationContainer.className = "dlp-pagination";
                this.DOM.append(paginationContainer);

                // Render Previous button
                this.renderButton(paginationContainer, '‹', this.currentPage > 1 ? this.currentPage - 1 : 1);

                // Render first page
                this.renderButton(paginationContainer, '1', 1);

                // Render ellipsis if needed before the middle pages
                if (this.currentPage > 4) {
                    let tmp = this.url;
                    this.url = undefined;
                    this.renderButton(paginationContainer, "...", -1);
                    this.url = tmp;
                }

                // Calculate a larger range of pages around the current page
                let visiblePagesCount = 7; // Adjust this value based on your preference
                if (this.visiblePagesCount) visiblePagesCount = this.visiblePagesCount;
                let start = Math.max(2, this.currentPage - Math.floor(visiblePagesCount / 2));
                let end = Math.min(this.totalPages - 1, start + visiblePagesCount - 1);

                // Render pages around current page
                for (let i = start; i <= end; i++) {
                    this.renderButton(paginationContainer, i.toString(), i);
                }

                // Render ellipsis if needed after the middle pages
                if (this.currentPage < this.totalPages - 3 && end < this.totalPages - 1) {
                    let tmp = this.url;
                    this.url = undefined;
                    this.renderButton(paginationContainer, "...", -1);
                    this.url = tmp;
                }

                // Render last page
                this.renderButton(paginationContainer, this.totalPages.toString(), this.totalPages);

                // Render Next button
                this.renderButton(paginationContainer, '›', this.currentPage < this.totalPages ? this.currentPage + 1 : this.totalPages);
            }

            renderButton(container, text, page) {
                const p = document.createElement('li');
                if (page === -1){
                    p.className = "dlp-pagination-page skip";
                }else if(this.currentPage === page && text !== '‹' && text !== '›'){
                    p.className = "dlp-pagination-page current";
                }else {
                    p.className = "dlp-pagination-page";
                }
                const a = document.createElement("a");
                if(this.url && page !== -1) {
                    a.setAttribute('href', this.url.replace(":page", page));
                }else {
                    a.setAttribute('href',"javascript:void(0);");
                }
                a.innerText = text;
                if (typeof this.callback === 'function'){
                    a.addEventListener('click', ()=>this.callback(page));
                }
                p.append(a);
                container.appendChild(p);
            }
        };

        w.search = class {
            constructor() {
                this.input_dom = document.querySelector("#search>input");
                this.button_dom = document.querySelector("#search>.search-button");
                this.button_dom.addEventListener('click',()=>{
                    if(!this.input_dom.value)return;
                    this.input_dom.value = this.input_dom.value.replace(/[!@#$%^&*()_+\{\}:“<>?,.\/;'\[\]\\|`~"\'【】！，。、]+/g, '');
                    this.input_dom.value = this.input_dom.value.substring(0,64);
                    window.location.href = new URL(window.location.href).origin + "/search?keywords="+this.input_dom.value;
                });
            }
        };
        new w.search();
    })(window);
});