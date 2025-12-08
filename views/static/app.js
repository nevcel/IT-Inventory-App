// app.js

async function loadInventory() {
    try {
        const response = await fetch("/inventory");
        if (!response.ok) throw new Error(`Serverstatus: ${response.status}`);

        const data = await response.json();
        const tbody = document.getElementById("inventoryBody");



        tbody.innerHTML = ""; // Tabelle leeren$


        // Falls dein Server ein einzelnes Objekt liefert, wandle es in ein Array
        const items = Array.isArray(data) ? data : [data];

        items.forEach(item => {
            const row = document.createElement("tr");
            row.innerHTML = `
        <td>${item.id}</td>
        <td>${item.type}</td>
        <td>${item.name}</td>
        <td>${item.date_added}</td>
        <td>${item.date_removed || "-"}</td>
        <td>${item.notes || "-"}</td>
      `;
            tbody.appendChild(row);
        });
    } catch (err) {
        console.error("Fehler beim Laden des Inventars:", err);
        document.getElementById("inventoryBody").innerHTML =
            "<tr><td colspan='6' style='color:red'>Inventardaten konnten nicht geladen werden.</td></tr>";
    }
}

document.addEventListener("DOMContentLoaded", loadInventory);
