const createBtn = document.getElementById("createBtn");
const analyticsBtn = document.getElementById("analyticsBtn");

createBtn.addEventListener("click", async () => {
    const url = document.getElementById("originalUrl").value.trim();
    const code = document.getElementById("customCode").value.trim();
    const resultDiv = document.getElementById("result");
    resultDiv.textContent = "Creating...";

    if (!url) {
        resultDiv.textContent = "URL is required!";
        return;
    }

    const payload = { url };
    if (code) payload.custom_code = code;

    try {
        const res = await fetch("http://localhost:8080/shorten", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(payload)
        });
        const data = await res.json();

        if (res.ok) {
            const shortUrl = `http://localhost:8080/s/${data.code}`;
            resultDiv.innerHTML = `Short link created: <a href="${shortUrl}" target="_blank">${shortUrl}</a>`;
        } else {
            resultDiv.textContent = data.error || "Error creating link";
        }
    } catch (err) {
        resultDiv.textContent = "Network error";
    }
});

analyticsBtn.addEventListener("click", async () => {
    const code = document.getElementById("analyticsCode").value.trim();
    const group = document.getElementById("groupSelect").value;
    const analyticsDiv = document.getElementById("analyticsResult");
    analyticsDiv.textContent = "Loading...";

    if (!code) {
        analyticsDiv.textContent = "Short code is required!";
        return;
    }

    try {
        const res = await fetch(`http://localhost:8080/analytics/${encodeURIComponent(code)}?group=${group}`);
        const data = await res.json();

        if (res.ok) {
            if (group === "all") {
                // показываем все визиты без группировки
                analyticsDiv.textContent = JSON.stringify(data, null, 2);
            } else {
                // показываем агрегированную статистику
                analyticsDiv.textContent = JSON.stringify(data, null, 2);
            }
        } else {
            analyticsDiv.textContent = data.error || "Error fetching analytics";
        }
    } catch (err) {
        analyticsDiv.textContent = "Network error";
    }
});
