// app.js

async function checkAuth() {
    const res = await fetch("/me");
    if (!res.ok) {
        window.location.href = "/login.html";
        return null;
    }
    return await res.json();
}

async function loadInventory() {
    try {
        const response = await fetch("/inventory");
        if (!response.ok) throw new Error(`Serverstatus: ${response.status}`);

        const data = await response.json();
        const tbody = document.getElementById("inventoryBody");



        tbody.innerHTML = ""; // Tabelle leeren


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

        //Erweiterung: Gesamtartikel anzeigen
        const totalItemsElement = document.getElementById("totalItems");
        if (totalItemsElement) {
            totalItemsElement.textContent = items.length;
        }

    } catch (err) {
        console.error("Fehler beim Laden des Inventars:", err);
        document.getElementById("inventoryBody").innerHTML =
            "<tr><td colspan='6' style='color:red'>Inventardaten konnten nicht geladen werden.</td></tr>";
    }
}

// Item search by Name
async function searchInventory(query) {
    try {
        const response = await fetch(`/items/search?q=${encodeURIComponent(query)}`);
        if (!response.ok) throw new Error(`Serverstatus: ${response.status}`);

        const data = await response.json();
        const tbody = document.getElementById("inventoryBody");
        tbody.innerHTML = "";

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
        console.error("Fehler bei der Suche:", err);
    }
}


document.addEventListener("DOMContentLoaded", loadInventory);

// Listener for Item-Search
document.addEventListener("DOMContentLoaded", async () => {
    const me = await checkAuth();
    if (!me) return;

    loadInventory();
    // Item-Search bereich
    const searchInput = document.getElementById("searchInput");
    if (!searchInput) return;

    searchInput.addEventListener("input", function () {
        const query = this.value.trim();

        if (query === "") {
            loadInventory();
        } else {
            searchInventory(query);
        }
    });

    //Item_Add Formular wird angezeigt
    const addItemBtn = document.getElementById("addItemBtn");
    const addItemForm = document.getElementById("addItemForm");

    if (!addItemBtn || !addItemForm) return;

    addItemBtn.addEventListener("click", () => {
        if (addItemForm.style.display === "none") {
            addItemForm.style.display = "block";
        } else {
            addItemForm.style.display = "none";
        }
    });

    //Neue Item wird an Json gesendet

    if (addItemForm) {
        addItemForm.addEventListener("submit", async function (e) {
            e.preventDefault();

            const item = {
                type: document.getElementById("itemType").value,
                name: document.getElementById("itemName").value,
                notes: document.getElementById("itemNotes").value
            };

            try {
                const response = await fetch("/inventory", {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json"
                    },
                    body: JSON.stringify(item)
                });

                if (!response.ok) {
                    throw new Error(`Serverstatus: ${response.status}`);
                }
                // Inventar neu laden
                loadInventory();
                // Formular leeren
                addItemForm.reset();
                // Formular schliessen
                addItemForm.style.display = "none";


            } catch (err) {
                console.error("Fehler beim Hinzufuegen des Artikels:", err);
            }
        });
    }


});
