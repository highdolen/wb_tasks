async function login() {

    const username = document.getElementById("nickname").value
    const role = document.getElementById("role").value

    if (!username || !role) {
        alert("Enter username and role")
        return
    }

    const res = await fetch("http://localhost:8080/login", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({
            username,
            role
        })
    })

    if (!res.ok) {
        alert("Login failed")
        return
    }

    const data = await res.json()

    localStorage.setItem("token", data.token)
    localStorage.setItem("role", role)

    window.location.href = "/pages/items.html"
}