async function login() {
    const name = document.getElementById("name").value;
    const email = document.getElementById("email").value;
    const telegram = document.getElementById("telegram").value;
    const role = document.getElementById("role").value;

    if (!name || !email) {
        alert("Введите имя и email");
        return;
    }

    const res = await fetch("/users", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
            name,
            email,
            telegram_id: telegram,
            role
        })
    });

    const user = await res.json();

    localStorage.setItem("user", JSON.stringify(user));

    if (user.role === "admin") {
        window.location.href = "/admin";
    } else {
        window.location.href = "/user";
    }
}