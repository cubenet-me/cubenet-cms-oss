(function() {
  'use strict';

  function renderServer(server) {
    const online = server.status === 'online';
    const playerStr = server.players != null ? server.players : '?';
    const maxStr = server.max_players != null ? server.max_players : '?';
    return `<div class="server-card">
      <div class="server-card-top neu-card">
        <div style="position:relative;z-index:2;display:flex;flex-direction:column;justify-content:space-between;height:100%;">
          <div style="display:flex;align-items:flex-start;justify-content:space-between;">
            <div>
              <div class="server-card-badge">${server.version || '1.20.4'}</div>
              <h3 style="font-size:20px;font-weight:700;margin-top:8px;color:#f7f8ff;">${server.name || server.slug}</h3>
            </div>
            <div class="status-dot${online ? '' : ' is-offline'}"><span class="outer"></span><span class="inner"></span></div>
          </div>
          <div style="margin-top:16px;display:flex;gap:20px;font-size:13px;color:rgba(210,213,231,0.7);">
            <div style="display:flex;align-items:center;gap:6px;">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path stroke-linecap="round" stroke-linejoin="round" d="M15.75 6a3.75 3.75 0 11-7.5 0 3.75 3.75 0 017.5 0zM4.501 20.118a7.5 7.5 0 0114.998 0A17.933 17.933 0 0112 21.75c-2.676 0-5.216-.584-7.499-1.632z"/></svg>
              <span>${playerStr}/${maxStr}</span>
            </div>
            <div style="display:flex;align-items:center;gap:6px;">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path stroke-linecap="round" stroke-linejoin="round" d="M2.25 15a4.5 4.5 0 004.5 4.5H18a3.75 3.75 0 001.332-7.257 3 3 0 00-3.758-3.848 5.25 5.25 0 00-10.233 2.33A4.502 4.502 0 002.25 15z"/></svg>
              <span>${server.address || server.slug + '.example.com'}</span>
            </div>
          </div>
          ${server.description ? '<p style="margin-top:10px;font-size:13px;color:rgba(210,213,231,0.6);">' + server.description + '</p>' : ''}
        </div>
      </div>
    </div>`;
  }

  function renderModPreview(mods) {
    if (!mods || mods.length === 0) return '';
    const visible = mods.slice(0, 3);
    const remainder = mods.length - 3;
    return visible.map(m => `<span class="server-card-badge" style="font-size:10px;padding:2px 8px;">${m}</span>`).join('') +
      (remainder > 0 ? `<span class="server-card-badge" style="font-size:10px;padding:2px 8px;">+${remainder}</span>` : '');
  }

  function loadServerList(container) {
    if (!container) return;
    fetch('/api/v1/servers')
      .then(function(r) { return r.json(); })
      .then(function(servers) {
        if (servers.length === 0) {
          container.innerHTML = '<div class="glass text-center p-8"><p class="text-muted">Серверов пока нет</p></div>';
          return;
        }
        container.innerHTML = servers.map(renderServer).join('');
      })
      .catch(function() {
        container.innerHTML = '<div class="glass text-center p-8"><p class="text-muted">Не удалось загрузить серверы</p></div>';
      });
  }

  document.addEventListener('DOMContentLoaded', function() {
    var el = document.getElementById('server-list');
    loadServerList(el);
  });

  window.loadServerList = loadServerList;
})();
