var BtnDiscussion = document.getElementById('discussionBtn');
var BtnCategorie = document.getElementById('categorieBtn');
var BtnPost = document.getElementById('postBtn');
var PartieDiscussion = document.getElementById('partieDiscussion');
var PartieCategorie = document.getElementById('partieCategorie');
var PartiePost = document.getElementById('partiePost');

BtnDiscussion.addEventListener('click', function() {
    PartieDiscussion.style.display = "block";
    PartieCategorie.style.display = "none";
    PartiePost.style.display = "none";
});

BtnCategorie.addEventListener('click', function() {
    PartieDiscussion.style.display = "none";
    PartieCategorie.style.display = "block";
    PartiePost.style.display = "none";
});

BtnPost.addEventListener('click', function() {
    PartieDiscussion.style.display = "none";
    PartieCategorie.style.display = "none";
    PartiePost.style.display = "block";
}
);


/*partie popup pour créer une discussion*/

var btn = document.getElementById('btn');
var closePopup = document.getElementById('close-popup');
var createBoxBtn = document.getElementById('create-box');
var popup = document.getElementById('popup');
var titleInput = document.getElementById('title');
var imgInput = document.getElementById('imgParcourir');
var preview = document.getElementById('preview');
var container = document.getElementById('container');

btn.addEventListener('click', function() {
    popup.style.display = "block";
});

closePopup.addEventListener('click', function() {
    popup.style.display = "none";
});

imgInput.addEventListener('change', function() {
    preview.src = URL.createObjectURL(imgInput.files[0]);
});

createBoxBtn.addEventListener('click', function() {
    var link = document.createElement('a');
    link.href = 'discussion.html';  /*changer ca pour le go*/

    var box = document.createElement('div');
    box.className = 'box';

    var title = document.createElement('h2');
    title.textContent = titleInput.value;
    box.appendChild(title);

    var image = document.createElement('img');
    image.src = preview.src;
    box.appendChild(image);

    link.appendChild(box);
    container.appendChild(link);

    popup.style.display = "none";

});

