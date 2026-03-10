function renderItemsTable(items) {

    const table = document.getElementById("itemsTable")
    table.innerHTML = ""

    const role = getRole()

    items.forEach(item => {

        const row = document.createElement("tr")

        let actions = ""

        if (role === "admin") {

            actions = `
                <button onclick="editItem(${item.id})">Edit</button>
                <button onclick="deleteItem(${item.id})">Delete</button>
                <button onclick="showHistory(${item.id})">History</button>
            `
        }

        if (role === "manager") {

            actions = `
                <button onclick="editItem(${item.id})">Edit</button>
                <button onclick="showHistory(${item.id})">History</button>
            `
        }

        if (role === "viewer") {

            actions = ``
        }

        row.innerHTML = `
        <td>${item.id}</td>
        <td id="name-${item.id}">${item.name}</td>
        <td id="qty-${item.id}">${item.quantity}</td>
        <td id="actions-${item.id}">
            ${actions}
        </td>
        `

        table.appendChild(row)
    })
}