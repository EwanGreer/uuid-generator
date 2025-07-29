document.addEventListener("DOMContentLoaded", function () {
  const slider = document.getElementById("password-length-slider");
  const output = document.getElementById("password-length-slider-length");

  if (slider && output) {
    console.log("Slider value:", slider.value);
    output.textContent = slider.value;

    slider.addEventListener("input", function () {
      output.textContent = this.value;
    });
  } else {
    console.warn("Slider or output not found");
  }
});
