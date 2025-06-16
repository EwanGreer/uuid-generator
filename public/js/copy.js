function myFunction() {
  const copyText = document.getElementById("uuid").innerText;
  const notif = document.getElementById("copy-notification");

  navigator.clipboard
    .writeText(copyText)
    .then(() => {
      notif.classList.remove("hidden", "fade-out");
      notif.classList.add("bounce-in");

      setTimeout(() => {
        notif.classList.remove("bounce-in");
      }, 600);

      setTimeout(() => {
        notif.classList.add("fade-out");
      }, 1500);

      setTimeout(() => {
        notif.classList.add("hidden");
        notif.classList.remove("fade-out");
      }, 1750);
    })
    .catch((err) => {
      console.error("Failed to copy UUID:", err);
    });
}
