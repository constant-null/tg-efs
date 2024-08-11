const appChoices = document.querySelectorAll('.approach-choice')

for (let i = 0; i < appChoices.length; i++) {
    appChoices[i].addEventListener("click", (e) => {
        const choice = e.target.dataset.choice;
        const dropbtn = e.currentTarget.parentElement.parentElement.querySelector('.dropbtn');
        dropbtn.firstChild.data = 'к'+choice;
        const approachInput = e.currentTarget.parentElement.parentElement.querySelector('input');
        approachInput.value = choice;
    })
}