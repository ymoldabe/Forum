var navLinks = document.querySelectorAll("nav a");
for (var i = 0; i < navLinks.length; i++) {
	var link = navLinks[i]
	if (link.getAttribute('href') == window.location.pathname) {
		link.classList.add("live");
		break;
	}
}

var btn1 = document.querySelector('#green');
var btn2 = document.querySelector('#red');

btn1.addEventListener('click', function() {
  
    if (btn2.classList.contains('red')) {
      btn2.classList.remove('red');
    } 
  this.classList.toggle('green');
  
});

btn2.addEventListener('click', function() {
  
    if (btn1.classList.contains('green')) {
      btn1.classList.remove('green');
    } 
  this.classList.toggle('red');
  
});


// coment

document.addEventListener("DOMContentLoaded", function () {
  const submitButton = document.getElementById("submit-btn");
  const errorMessage = document.getElementById("error-message");

  submitButton.addEventListener("click", function () {
    const nameInput = document.getElementById("name");
    const commentInput = document.getElementById("comment");

    if (nameInput.value.trim() === "" || commentInput.value.trim() === "") {
      errorMessage.style.display = "block";
    } else {
      errorMessage.style.display = "none";
      // Here you can implement the logic to actually submit the comment
    }
  });
});