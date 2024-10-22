function togglePropertyMenu() {
  const propertyMenuButton = document.getElementById('property-menu-button');
  const propertyMenu = document.getElementById('property-menu');
  const propertyArrow = document.getElementById('property-arrow');

  propertyMenuButton.addEventListener('click', () => {
    propertyMenu.classList.toggle('hidden');
    propertyArrow.classList.toggle('rotate-180');
  });

  window.addEventListener('click', function (e) {
    if (
      !propertyMenuButton.contains(e.target) &&
      !propertyMenu.contains(e.target)
    ) {
      propertyMenu.classList.add('hidden');
      propertyArrow.classList.remove('rotate-180');
    }
  });
}

function toggleMobileMenu() {
  const mobileMenuButton = document.getElementById('mobile-menu-button');
  const closeMobileMenuButton = document.getElementById('close-mobile-menu');
  const mobileMenu = document.getElementById('mobile-menu');

  mobileMenuButton.addEventListener('click', () => {
    mobileMenu.classList.remove('hidden');
  });

  closeMobileMenuButton.addEventListener('click', () => {
    mobileMenu.classList.add('hidden');
  });
}

window.onload = function () {
  togglePropertyMenu();
  toggleMobileMenu();
};
