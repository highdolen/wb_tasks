async function filterHistory() {

    const user = document.getElementById("user").value
    const action = document.getElementById("action").value
    const from = document.getElementById("from").value
    const to = document.getElementById("to").value

    const params = new URLSearchParams()

    if (user) params.append("user", user)
    if (action) params.append("action", action)
    if (from) params.append("from", from)
    if (to) params.append("to", to)

    const history = await apiRequest("/audit/filter?" + params.toString())

    if (!history) return

    renderAuditTable(history)
}

function renderAuditTable(history) {

    const table = document.getElementById("auditTable")
    table.innerHTML = ""

    history.forEach(h => {

        const row = document.createElement("tr")

        row.innerHTML = `
            <td>${h.id}</td>
            <td>${h.item_id}</td>
            <td>${h.action}</td>
            <td>${h.changed_by}</td>
            <td>${h.changed_at}</td>
        `

        table.appendChild(row)
    })
}

function goBack() {
    window.location.href = "/pages/items.html"
}

document.addEventListener("DOMContentLoaded", () => {
    checkAuth()
    filterHistory()
})