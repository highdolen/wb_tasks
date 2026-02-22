document.addEventListener("DOMContentLoaded", () => {
    document.getElementById("createEventBtn").addEventListener("click", createEvent);
    loadEvents();
});

async function createEvent() {
    const name = document.getElementById("eventName").value;
    const date = document.getElementById("eventDate").value;
    const totalSeats = parseInt(document.getElementById("totalSeats").value);
    const bookingTTL = parseInt(document.getElementById("bookingTTL").value);

    if (!name || !date || !totalSeats || !bookingTTL) {
        alert("Заполните все поля");
        return;
    }

    const res = await fetch("/events", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
            name,
            date,
            total_seats: totalSeats,
            booking_ttl_minutes: bookingTTL
        })
    });

    const data = await res.json();

    if (!res.ok) {
        alert(data.error);
        return;
    }

    alert("Событие создано");
    loadEvents();
}

async function loadEvents() {
    const res = await fetch("/events");
    const events = await res.json();

    const container = document.getElementById("events");
    container.innerHTML = "";

    events.forEach(ev => {
        const div = document.createElement("div");
        div.className = "event";

        div.innerHTML = `
            <b>${ev.name}</b><br>
            дата: ${new Date(ev.date).toLocaleString()}<br>
            свободно: ${ev.available_seats}/${ev.total_seats}<br>
            TTL: ${ev.booking_ttl / 60000000000} мин
        `;

        container.appendChild(div);
    });
}