async function showHistory(itemId) {

    const history = await apiRequest(`/items/${itemId}/history`)

    if (!history) return

    let text = "Item History\n\n"

    history.forEach(h => {

        text += `Action: ${h.action}\n`
        text += `User: ${h.changed_by}\n`
        text += `Time: ${h.changed_at}\n`

        if (h.old_value && h.new_value) {

            const oldVal = JSON.parse(h.old_value)
            const newVal = JSON.parse(h.new_value)

            Object.keys(newVal).forEach(key => {

                if (oldVal[key] !== newVal[key]) {
                    text += `${key}: ${oldVal[key]} → ${newVal[key]}\n`
                }

            })
        }

        text += "\n"
    })

    alert(text)
}

function exportCSV() {

    const token = localStorage.getItem("token")

    fetch("http://localhost:8080/audit/export", {
        headers: {
            "Authorization": "Bearer " + token
        }
    })
    .then(res => res.blob())
    .then(blob => {

        const url = window.URL.createObjectURL(blob)

        const a = document.createElement("a")
        a.href = url
        a.download = "history.csv"

        document.body.appendChild(a)
        a.click()
        a.remove()
    })
}