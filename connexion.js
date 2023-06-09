var PopupConnexion = document.getElementById('popupSeConnecter');
var PopupInscription = document.getElementById('popupSinscrire');
var BoutonConnexion = document.getElementById('boutonSeConnecter');
var BoutonInscription = document.getElementById('boutonSinscrire');
var wsh = document.getElementById('wsh');
var Vidéo = document.getElementById('background-video');
var FermerPopupConnexion = document.getElementById('boutonFermer');
var FermerPopupConnexion2 = document.getElementById('boutonFermer2');

PopupConnexion.style.display = 'none';
PopupInscription.style.display = 'none';

BoutonConnexion.addEventListener('click', function() {
    wsh.style.display = 'none';
    PopupInscription.style.display = 'none';
    Vidéo.style.marginTop = '-200px';
    PopupConnexion.style.display = 'block';
});

BoutonInscription.addEventListener('click', function() {
    wsh.style.display = 'none';
    PopupConnexion.style.display = 'none';
    Vidéo.style.marginTop = '-200px';
    PopupInscription.style.display = 'block';
});

FermerPopupConnexion.addEventListener('click', function() {
    wsh.style.display = 'flex';
    PopupConnexion.style.display = 'none';
    PopupInscription.style.display = 'none';
    Vidéo.style.marginTop = '0px';
});

FermerPopupConnexion2.addEventListener('click', function() {
    wsh.style.display = 'flex';
    PopupConnexion.style.display = 'none';
    PopupInscription.style.display = 'none';
    Vidéo.style.marginTop = '0px';
});

