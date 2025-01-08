document.addEventListener("DOMContentLoaded", async () => {
    const navbar = document.getElementById("navbar");
    try {
      const response = await fetch("../navbar.html");
      const navbarHTML = await response.text();
      navbar.innerHTML = navbarHTML;
    } catch (err) {
      console.error("Error loading navbar:", err);
    }
  });
  