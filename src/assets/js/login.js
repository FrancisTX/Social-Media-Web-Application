const switchers = [...document.querySelectorAll('.switcher')]

switchers.forEach(item => {
	item.addEventListener('click', function() {
		switchers.forEach(item => item.parentElement.classList.remove('is-active'))
		this.parentElement.classList.add('is-active')
	})
})

const imgDiv = document.querySelector('.profile-pic-div');
const img = document.querySelector('#photo');
const file = document.querySelector('#file');
const uploadBtn = document.querySelector('#uploadBtn');

imgDiv.addEventListener('mouseenter', function(){
    uploadBtn.style.display = "block";
});


imgDiv.addEventListener('mouseleave', function(){
    uploadBtn.style.display = "none";
});

file.addEventListener('change', function(){
    const choosedFile = this.files[0];

    if (choosedFile) {

        const reader = new FileReader(); 

        reader.addEventListener('load', function(){
            console.log(choosedFile)
            console.log(reader.result)
            img.setAttribute('src', reader.result);
        });

        reader.readAsDataURL(choosedFile);
    }
});