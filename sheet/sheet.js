console.log(window.Telegram.WebApp.initData)
fetch('/test?'+window.Telegram.WebApp.initData).then((r) => console.log(r))