function getToken() {
    return localStorage.getItem("token")
}

function getRole() {
    return localStorage.getItem("role")
}

function isAuth() {
    return !!getToken()
}

function logout() {
    localStorage.removeItem("token")
    localStorage.removeItem("role")
    window.location.href = "/pages/login.html"
}

function checkAuth() {

    const path = window.location.pathname

    if (path === "/" || path === "/index.html") {
        if (isAuth()) {
            window.location.href = "/pages/items.html"
        }
        return
    }

    if (path.includes("login.html")) {
        if (isAuth()) {
            window.location.href = "/pages/items.html"
        }
        return
    }

    if (!isAuth()) {
        window.location.href = "/pages/login.html"
    }

    const role = getRole()

    const createBox = document.getElementById("create-box")

    if (createBox && role !== "admin") {
        createBox.style.display = "none"
    }

    const exportBtn = document.querySelector("button[onclick='exportCSV()']")

    if (exportBtn && role === "viewer") {
        exportBtn.style.display = "none"
    }

    const auditBtn = document.getElementById("audit-btn")

    if (auditBtn && role === "viewer") {
        auditBtn.style.display = "none"
    }
}

document.addEventListener("DOMContentLoaded", checkAuth)