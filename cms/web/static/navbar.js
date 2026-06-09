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

  // ===== Navbar admin drag-and-drop =====
  var dragSrc = null;

  document.addEventListener('dragstart', function (e) {
    var row = e.target.closest('[draggable="true"]');
    if (!row || !row.closest('#nav-items-list')) return;
    dragSrc = row;
    row.style.opacity = '0.4';
    e.dataTransfer.effectAllowed = 'move';
  });

  document.addEventListener('dragend', function (e) {
    var row = e.target.closest('[draggable="true"]');
    if (row) row.style.opacity = '';
    dragSrc = null;
    document.querySelectorAll('.nav-item-row.drag-over').forEach(function (el) {
      el.classList.remove('drag-over');
    });
  });

  document.addEventListener('dragover', function (e) {
    var row = e.target.closest('[draggable="true"]');
    if (!row || row === dragSrc || !row.closest('#nav-items-list')) return;
    e.preventDefault();
    e.dataTransfer.dropEffect = 'move';
    document.querySelectorAll('.nav-item-row.drag-over').forEach(function (el) {
      el.classList.remove('drag-over');
    });
    row.classList.add('drag-over');
  });

  document.addEventListener('drop', function (e) {
    e.preventDefault();
    var row = e.target.closest('[draggable="true"]');
    if (!row || !dragSrc || row === dragSrc || !row.closest('#nav-items-list')) return;

    document.querySelectorAll('.nav-item-row.drag-over').forEach(function (el) {
      el.classList.remove('drag-over');
    });

    var parent = row.parentNode;
    if (Array.from(parent.children).indexOf(dragSrc) < Array.from(parent.children).indexOf(row)) {
      parent.insertBefore(dragSrc, row.nextSibling);
    } else {
      parent.insertBefore(dragSrc, row);
    }

    // re-index all order_, label_, href_, icon_, and delete values
    parent.querySelectorAll('[draggable="true"]').forEach(function (el, idx) {
      el.querySelectorAll('input, select, button').forEach(function (input) {
        var name = input.getAttribute('name');
        if (!name) return;
        var match = name.match(/^(order_|label_|href_|icon_)(\d+)$/);
        if (match) {
          input.setAttribute('name', match[1] + idx);
        }
        if (input.getAttribute('name') === 'delete') {
          input.value = idx;
        }
      });
      // update order-input value
      var orderInput = el.querySelector('.order-input');
      if (orderInput) orderInput.value = idx;
    });
  });

  document.addEventListener('dragenter', function (e) { e.preventDefault(); });
  document.addEventListener('dragexit', function (e) { e.preventDefault(); });
})();
