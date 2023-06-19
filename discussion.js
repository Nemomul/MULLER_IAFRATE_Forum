var submitMessageBtn = document.getElementById('submitMessageBtn');
var messageInput = document.getElementById('messageInput');
var discussion = document.getElementById('discussion');

messageInput.addEventListener('input', function() {
    this.style.height = 'auto';
    this.style.height = (this.scrollHeight) + 'px';
});

submitMessageBtn.addEventListener('click', function() {
    var messageContainer = document.createElement('div');
    messageContainer.className = 'message';
    messageContainer.style.marginBottom = "20px";

    var messageText = document.createElement('p');
    messageText.textContent = messageInput.value;
    messageContainer.appendChild(messageText);

    discussion.appendChild(messageContainer);

    messageInput.value = "";
    messageInput.style.height = '50px';
});
