(function() {
  'use strict';

  function renderServer(server) {
    var title = server.name || server.slug || 'Server';
    var desc = server.description || 'Minecraft server';
    var version = server.version || '';
    var addr = server.address || server.slug + '.example.com';
    var online = server.status === 'online';
    var players = server.players != null ? server.players : '?';
    var maxPl = server.max_players != null ? server.max_players : '?';

    return '<a href="/server/' + server.slug + '" class="server-card-ref group">' +
      '<div class="server-card-ref-bg"></div>' +
      '<div class="server-card-ref-overlay"></div>' +
      '<div class="server-card-ref-violet"></div>' +
      '<div class="server-card-ref-hover"><span>Открыть</span></div>' +
      '<div class="server-card-ref-content">' +
        '<h3 class="server-card-ref-title">' + escapeHtml(title) + '</h3>' +
        '<p class="server-card-ref-desc">' + escapeHtml(desc) + '</p>' +
        '<div class="server-card-ref-tags">' +
          (online
            ? '<span class="server-card-ref-badge-online"><span class="server-card-ref-dot"></span> Онлайн</span>'
            : '<span class="server-card-ref-badge-online" style="border-color:rgba(255,120,120,0.18);background:rgba(255,120,120,0.10);color:#fecaca;"><span class="server-card-ref-dot" style="background:#ff7c7c;box-shadow:0 0 18px rgba(255,120,120,0.65);"></span> Офлайн</span>') +
          (version ? '<span class="server-card-ref-badge-version">' + escapeHtml(version) + '</span>' : '') +
        '</div>' +
      '</div>' +
    '</a>';
  }

  function escapeHtml(str) {
    var div = document.createElement('div');
    div.appendChild(document.createTextNode(str));
    return div.innerHTML;
  }

  function loadServerList(container) {
    if (!container) return;
    fetch('/api/v1/servers')
      .then(function(r) { return r.json(); })
      .then(function(servers) {
        if (!servers || servers.length === 0) {
          container.innerHTML = '<div class="glass rounded-3xl p-10 text-center"><p style="color:rgba(255,255,255,0.45);">Серверов пока нет</p></div>';
          return;
        }
        if (servers.length === 1) {
          container.className = 'grid gap-4 sm:grid-cols-2';
          container.innerHTML = renderServer(servers[0]);
          return;
        }
        container.className = 'home-server-marquee';
        var duration = Math.max(24, servers.length * 10);
        container.innerHTML =
          '<div class="home-server-track home-server-track-animated" style="--home-server-duration:' + duration + 's">' +
            '<div class="home-server-set">' +
              servers.map(function(s) { return renderServer(s); }).join('') +
            '</div>' +
            '<div class="home-server-set" aria-hidden="true">' +
              servers.map(function(s) { return renderServer(s); }).join('') +
            '</div>' +
          '</div>';
      })
      .catch(function() {
        container.innerHTML = '<div class="glass rounded-3xl p-10 text-center"><p style="color:rgba(255,255,255,0.45);">Не удалось загрузить серверы</p></div>';
      });
  }

  document.addEventListener('DOMContentLoaded', function() {
    var el = document.getElementById('server-list');
    loadServerList(el);
  });

  window.loadServerList = loadServerList;
})();
