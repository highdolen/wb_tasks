const API_URL = "http://localhost:8080";

// Создание уведомления
async function createNotification() {
    const data = {
        channel: document.getElementById("channel").value,
        recipient: document.getElementById("recipient").value,
        subject: document.getElementById("subject").value,
        message: document.getElementById("message").value,
        send_at: new Date(document.getElementById("sendAt").value).toISOString()
    };

    const res = await fetch(`${API_URL}/notify`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data)
    });

    if (!res.ok) {
        alert("Ошибка при создании уведомления");
        return;
    }

    const result = await res.json();
    alert(`Уведомление создано\nID: ${result.id}`);
}

// Получить статус уведомления
async function getStatus() {
    const id = document.getElementById("notificationId").value;
    if (!id) {
        alert("Введите ID уведомления");
        return;
    }

    const res = await fetch(`${API_URL}/notify/${id}`);
    if (!res.ok) {
        document.getElementById("statusResult").textContent = "Не найдено";
        return;
    }

    const data = await res.json();
    document.getElementById("statusResult").textContent =
        JSON.stringify(data, null, 2);
}

// Удаление (отмена) уведомления
async function deleteNotification() {
    const id = document.getElementById("deleteId").value;
    if (!id) {
        alert("Введите ID уведомления");
        return;
    }

    if (!confirm("Точно отменить уведомление?")) return;

    const res = await fetch(`${API_URL}/notify/${id}`, {
        method: "DELETE"
    });

    if (!res.ok) {
        alert("Ошибка при отмене уведомления");
        return;
    }

    alert("Уведомление отменено");
}
