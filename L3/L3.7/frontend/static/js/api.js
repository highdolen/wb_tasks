const API = "http://localhost:8080"

async function apiRequest(url, options = {}) {

    options.headers = options.headers || {}

    const token = localStorage.getItem("token")

    if (token) {
        options.headers["Authorization"] = "Bearer " + token
    }

    const response = await fetch(API + url, options)

    if (response.status === 401) {
        logout()
        return null
    }

    if (!response.ok) {
        const text = await response.text()
        console.error("API error:", text)
        return null
    }

    return response.json()
}