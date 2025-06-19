var slider = document.getElementById("password-length-slider");
var output = document.getElementById("password-length-slider-length");
output.innerHTML = slider.value; // Display the default slider value

// Update the current slider value (each time you drag the slider handle)
slider.oninput = function () {
  output.innerHTML = this.value;
};
