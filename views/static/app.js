// app.js
let currentUser = null;
let editItemId = null;

function getTodayLocalISO() {
    const now = new Date();
    const yyyy = now.getFullYear();
    const mm = String(now.getMonth() + 1).padStart(2, "0");
    const dd = String(now.getDate()).padStart(2, "0");
    return `${yyyy}-${mm}-${dd}`;
}

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
            const canDelete =
                currentUser &&
                (currentUser.role === "admin" || item.owner_id === currentUser.id);

            const deleteButtonHtml = canDelete
                ? `<button class="btn-primary" data-delete-id="${item.id}">Loeschen</button>`
                : `-`;
            const editButtonHtml = canDelete
                ? `<button class="btn-primary" data-edit-id="${item.id}">Edit</button>`
                : `-`;


            row.innerHTML = `
        <td>${item.id}</td>
        <td>${item.type}</td>
        <td>${item.name}</td>
        <td>${item.date_added}</td>
        <td>${item.date_edited || "-"}</td>
        <td>${item.notes || "-"}</td>
        <td>${editButtonHtml} ${deleteButtonHtml}</td>

      `;
            tbody.appendChild(row);
        });

        //Erweiterung: Gesamtartikel anzeigen
        const totalItemsElement = document.getElementById("totalItems");
        if (totalItemsElement) {
            totalItemsElement.textContent = items.length;
        }

        // Erweiterung: Verschiedene Types anzeigen
        const totalTypesElement = document.getElementById("totalTypes");
        if (totalTypesElement) {
            const typesSet = new Set(
                items
                    .map(x => (x.type || "").trim().toLowerCase())
                    .filter(x => x !== "")
            );
            totalTypesElement.textContent = typesSet.size;
        }

        // Erweiterung: Heute hinzugefügt anzeigen
        const todayAddedElement = document.getElementById("todayAdded");
        if (todayAddedElement) {
            const today = getTodayLocalISO();
            const countToday = items.filter(x => x.date_added === today).length;
            todayAddedElement.textContent = countToday;
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
            const canDelete =
                currentUser &&
                (currentUser.role === "admin" || item.owner_id === currentUser.id);

            const deleteButtonHtml = canDelete
                ? `<button class="btn-primary" data-delete-id="${item.id}">Loeschen</button>`
                : `-`;
            const editButtonHtml = canDelete
                ? `<button class="btn-primary" data-edit-id="${item.id}">Edit</button>`
                : `-`;


            row.innerHTML = `
                <td>${item.id}</td>
                <td>${item.type}</td>
                <td>${item.name}</td>
                <td>${item.date_added}</td>
                <td>${item.date_edited || "-"}</td>
                <td>${item.notes || "-"}</td>
                <td>${editButtonHtml} ${deleteButtonHtml}</td>


            `;
            tbody.appendChild(row);
        });

    } catch (err) {
        console.error("Fehler bei der Suche:", err);
    }
}



//Logout-Funktion
async function doLogout() {
    await fetch("/logout", { method: "POST" });
    window.location.href = "/login.html";
}

//User-Info Anzeigen
function showUserInfo(me) {
    const userInfo = document.getElementById("userInfo");
    if (!userInfo || !me) return;

    if (me.role === "admin") {
        userInfo.textContent = `Eingeloggt als: ${me.username} (ADMIN)`;
    } else {
        userInfo.textContent = `Eingeloggt als: ${me.username}`;
    }
}

async function deleteItem(id) {
    const res = await fetch(`/inventory/${id}`, { method: "DELETE" });

    if (res.ok) {
        loadInventory();
    } else {
        console.error("Loeschen fehlgeschlagen:", res.status);
    }
}

// Listener for Item-Search
document.addEventListener("DOMContentLoaded", async () => {
    const me = await checkAuth();
    if (!me) return;

    currentUser = me
    showUserInfo(me);
    loadInventory();
    // Item-Search bereich
    const searchInput = document.getElementById("searchInput");
    if (searchInput) {
        searchInput.addEventListener("input", function () {
            const query = this.value.trim();

            if (query === "") {
                loadInventory();
            } else {
                searchInventory(query);
            }
        });
    }



    //Item_Add Formular wird angezeigt
    const addItemBtn = document.getElementById("addItemBtn");
    const addItemForm = document.getElementById("addItemForm");

    if (addItemBtn && addItemForm) {
        addItemBtn.addEventListener("click", () => {
            if (addItemForm.style.display === "none") {
                addItemForm.style.display = "block";
            } else {
                addItemForm.style.display = "none";
            }
        });
    }


    //Neue Item wird an Json gesendet

    if (addItemForm) {
        addItemForm.addEventListener("submit", async function (e) {
            e.preventDefault();

            const item = {
                type: document.getElementById("itemType").value,
                name: document.getElementById("itemName").value,
                notes: document.getElementById("itemNotes").value
            };

            const isEdit = editItemId !== null;

            const url = isEdit ? `/inventory/${editItemId}` : "/inventory";
            const method = isEdit ? "PUT" : "POST";

            try {
                const response = await fetch(url, {
                    method: method,
                    headers: {
                        "Content-Type": "application/json"
                    },
                    body: JSON.stringify(item)
                });

                if (!response.ok) {
                    throw new Error(`Serverstatus: ${response.status}`);
                }

                editItemId = null;

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




    const logoutBtn = document.getElementById("logoutBtn");
    if (logoutBtn) {
        logoutBtn.addEventListener("click", async () => {
            await doLogout();
        });
    }
    const tableBody = document.getElementById("inventoryBody");
    if (tableBody) {
        tableBody.addEventListener("click", async (e) => {
            // sucht das naechste Element (Button) das data-delete-id oder data-edit-id hat
            const actionBtn = e.target.closest("[data-delete-id], [data-edit-id]");
            if (!actionBtn) return;

            // LOESCHEN
            if (actionBtn.dataset.deleteId) {
                await deleteItem(actionBtn.dataset.deleteId);
                return;
            }

            // EDIT
            if (actionBtn.dataset.editId) {
                const id = actionBtn.dataset.editId;

                const res = await fetch("/inventory");
                if (!res.ok) return;

                const data = await res.json();
                const items = Array.isArray(data) ? data : [data];

                const item = items.find(x => String(x.id) === String(id));

                if (!item) return;

                editItemId = item.id;

                const addItemForm = document.getElementById("addItemForm");
                addItemForm.style.display = "block";

                document.getElementById("itemType").value = item.type || "";
                document.getElementById("itemName").value = item.name || "";
                document.getElementById("itemNotes").value = item.notes || "";
            }
        });
    }

});
