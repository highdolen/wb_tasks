async function loadItems() {

    const items = await apiRequest("/items")

    if (!items) return

    renderItemsTable(items)
}

async function createItem() {

    const role = getRole()

    if (role !== "admin") {
        alert("Only admin can create items")
        return
    }

    const name = document.getElementById("name").value
    const quantity = parseInt(document.getElementById("quantity").value)

    if (!name || !quantity) {
        alert("Fill name and quantity")
        return
    }

    await apiRequest("/items", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ name, quantity })
    })

    loadItems()
}
    

function editItem(id) {

    const role = getRole()

    if (role === "viewer") {
        alert("Viewer cannot edit items")
        return
    }

    const nameCell = document.getElementById(`name-${id}`)
    const qtyCell = document.getElementById(`qty-${id}`)
    const actions = document.getElementById(`actions-${id}`)

    const name = nameCell.innerText
    const qty = qtyCell.innerText

    nameCell.innerHTML = `<input id="edit-name-${id}" value="${name}">`
    qtyCell.innerHTML = `<input id="edit-qty-${id}" type="number" value="${qty}">`

    actions.innerHTML = `
        <button onclick="saveItem(${id})">Save</button>
        <button onclick="loadItems()">Cancel</button>
    `
}

async function saveItem(id) {

    await apiRequest(`/items/${id}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
            name: document.getElementById(`edit-name-${id}`).value,
            quantity: parseInt(document.getElementById(`edit-qty-${id}`).value)
        })
    })

    loadItems()
}

async function deleteItem(id) {

    const role = getRole()

    if (role !== "admin") {
        alert("Only admin can delete items")
        return
    }

    if (!confirm("Delete item?")) return

    await apiRequest(`/items/${id}`, {
        method: "DELETE"
    })

    loadItems()
}

document.addEventListener("DOMContentLoaded", () => {
    checkAuth()
    loadItems()
})