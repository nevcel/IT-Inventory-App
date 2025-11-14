async function loadInventory() {
    const response = await fetch('http://localhost:8080/inventory');
    const data = await response.json();
    const tbody = document.querySelector('#inventoryTable tbody');
    tbody.innerHTML = '';
    data.forEach(item => {
        const row = `<tr>
            <td>${item.id}</td>
            <td>${item.type}</td>
            <td>${item.name}</td>
            <td>${item.date_added}</td>
            <td>${item.notes}</td>
        </tr>`;
        tbody.innerHTML += row;
    });
}