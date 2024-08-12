const elements = {}
init()
fetch('/sheet_data?'+window.Telegram.WebApp.initData).then(
   async  (r) => fill(await r.json())
)

document.querySelector("#sheet-form").addEventListener("submit", (e) => {
    e.preventDefault()
    e.stopPropagation()
    let formData = new FormData(e.target);
    let object = {};
    formData.forEach((value, key) => object[key] = value);
    fetch('/sheet_data', {
        method: 'POST',
        body: JSON.stringify(object),
        headers: {
            "X-User-Data": window.Telegram.WebApp.initData
        }
    }).then((response) => {
        // do something with response here...
    });

})


function init() {
    const boundElements = document.querySelectorAll("[data-bind]")
    for (let i = 0; i < boundElements.length; i++) {
        if (!elements[boundElements[i].name]) {
            elements[boundElements[i].name] = []
        }
        elements[boundElements[i].name].push(boundElements[i])
    }
}

function fill(data) {
    for (const [key, value] of Object.entries(data)) {
        if (!elements[key]) continue;

        for (let i = 0; i < elements[key].length; i++) {
            switch (elements[key][i].dataset.bind) {
                case "value":
                    elements[key][i].value = value;
                    break;
                case "data":
                    elements[key][i].firstChild.data = value;
                    break
                default:
                    console.log(`some bullshit happened (data-bind = ${elements[key][i].dataset.bind})`)
            }
        }
    }
}