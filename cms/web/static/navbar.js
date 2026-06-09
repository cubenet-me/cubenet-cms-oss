(function () {
  'use strict';

  var navbar = document.querySelector('.navbar');
  if (!navbar) return;

  var lastScroll = 0;
  var ticking = false;
  var threshold = 50;
  var isHidden = false;

  function update() {
    var currentScroll = window.scrollY;

    if (currentScroll <= 0) {
      if (isHidden) {
        navbar.classList.remove('navbar-hidden');
        isHidden = false;
      }
      lastScroll = currentScroll;
      ticking = false;
      return;
    }

    if (currentScroll > lastScroll + 5 && currentScroll > threshold) {
      if (!isHidden) {
        navbar.classList.add('navbar-hidden');
        isHidden = true;
      }
    } else if (currentScroll < lastScroll - 5) {
      if (isHidden) {
        navbar.classList.remove('navbar-hidden');
        isHidden = false;
      }
    }

    lastScroll = currentScroll;
    ticking = false;
  }

  window.addEventListener('scroll', function () {
    if (!ticking) {
      window.requestAnimationFrame(update);
      ticking = true;
    }
  }, { passive: true });
})();
