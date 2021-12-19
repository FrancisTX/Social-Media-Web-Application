const imgDiv = document.getElementById('preview')
const file = document.querySelector('#selectedFile');

file.addEventListener('change', function(){
    const choosedFile = this.files[0];

    if (choosedFile) {

        const reader = new FileReader(); 

        reader.addEventListener('load', function(){
          imgDiv.src = reader.result
          imgDiv.style.display = "block";
        });

        reader.readAsDataURL(choosedFile);
    }
});