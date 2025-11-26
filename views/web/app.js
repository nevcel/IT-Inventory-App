let inventory = [];
let allInventory = [];

// Icons für verschiedene Typen
const typeIcons = {
    'air': '✈️',
    'bürostuhl': '🪑',
    'drucker': '🖨️',
    'macbook': '💻',
    'messestand': '🎪',
    'promotion': '📦',
    'default': '📌'
};

const typeColors = {
    'air': '#4CAF50',
    'bürostuhl': '#FFC107',
    'drucker': '#2196F3',
    'macbook': '#9C27B0',
    'messestand': '#F44336',
    'promotion': '#FF9800',
    'default': '#607D8B'
};

function getIcon(type) {
    const key = type.toLowerCase();
    return typeIcons[key] || typeIcons.default;
}

function getColor(type) {
    const key = type.toLowerCase();
    return typeColors[key] || typeColors.default;
}

async function loadInventory() {
    try {
        document.getElementById('errorContainer').innerHTML = '';
        document.getElementById('inventoryList').innerHTML = '<div class="loading">Lade Inventar...</div>';

        const response = await fetch('http://localhost:8080/inventory');
        if (!response.ok) {
            throw new Error('Fehler beim Laden der Daten');
        }
        const data = await response.json();

        allInventory = data;
        inventory = data;
        renderInventory();

    } catch (error) {
        console.error('Fehler:', error);
        document.getElementById('errorContainer').innerHTML =
            `<div class="error-message">Inventar konnte nicht geladen werden: ${error.message}</div>`;
        document.getElementById('inventoryList').innerHTML = `
                    <div class="empty-state">
                        <div class="empty-state-icon">⚠️</div>
                        <div>Fehler beim Laden der Daten</div>
                    </div>
                `;
    }
}

function renderInventory(filter = '') {
    const listEl = document.getElementById('inventoryList');

    const filtered = allInventory.filter(item =>
        item.name.toLowerCase().includes(filter.toLowerCase()) ||
        item.type.toLowerCase().includes(filter.toLowerCase()) ||
        (item.notes && item.notes.toLowerCase().includes(filter.toLowerCase()))
    );

    if (filtered.length === 0) {
        listEl.innerHTML = `
                    <div class="empty-state">
                        <div class="empty-state-icon">📦</div>
                        <div>Keine Items gefunden</div>
                    </div>
                `;
        return;
    }

    listEl.innerHTML = filtered.map(item => `
                <div class="inventory-item" onclick="showDetail(${item.id})">
                    <div class="item-icon" style="background: ${getColor(item.type)}20;">
                        ${getIcon(item.type)}
                    </div>
                    <div class="item-info">
                        <div class="item-name">${item.name}</div>
                        <div class="item-details">${item.type} • ${item.notes || 'Keine Notizen'}</div>
                    </div>
                    <div class="item-arrow">›</div>
                </div>
            `).join('');
}

function showDetail(id) {
    const item = allInventory.find(i => i.id === id);
    if (!item) return;

    document.getElementById('detailContent').innerHTML = `
                <div class="detail-row">
                    <div class="detail-label">ID</div>
                    <div class="detail-value">${item.id}</div>
                </div>
                <div class="detail-row">
                    <div class="detail-label">Typ</div>
                    <div class="detail-value">${item.type}</div>
                </div>
                <div class="detail-row">
                    <div class="detail-label">Name</div>
                    <div class="detail-value">${item.name}</div>
                </div>
                <div class="detail-row">
                    <div class="detail-label">Hinzugefügt am</div>
                    <div class="detail-value">${item.date_added}</div>
                </div>
                <div class="detail-row">
                    <div class="detail-label">Notizen</div>
                    <div class="detail-value">${item.notes || '-'}</div>
                </div>
            `;

    document.getElementById('detailModal').classList.add('active');
}

// Event Listeners
document.getElementById('closeDetailBtn').addEventListener('click', () => {
    document.getElementById('detailModal').classList.remove('active');
});

document.getElementById('searchInput').addEventListener('input', (e) => {
    renderInventory(e.target.value);
});

// Modal außerhalb klicken zum Schließen
document.getElementById('detailModal').addEventListener('click', (e) => {
    if (e.target.id === 'detailModal') {
        document.getElementById('detailModal').classList.remove('active');
    }
});

// Initial laden
loadInventory();