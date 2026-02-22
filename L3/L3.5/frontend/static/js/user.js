let currentUser = JSON.parse(localStorage.getItem("user"));
let eventsMap = {};
let myBookings = {};

// загрузка событий и броней пользователя после загрузки страницы
document.addEventListener("DOMContentLoaded", () => {
    loadEvents();
    loadUserBookings();
});

// загрузка событий
async function loadEvents() {
    try {
        const res = await fetch("/events");
        const events = await res.json();

        const container = document.getElementById("events");
        container.innerHTML = "";

        events.forEach(ev => {
            eventsMap[ev.id] = ev;

            const div = document.createElement("div");
            div.innerHTML = `
                <b>${ev.name}</b><br>
                Дата: ${new Date(ev.date).toLocaleString()}<br>
                Свободно мест: ${ev.available_seats}<br>
                <button onclick="bookEvent(${ev.id})">Забронировать</button>
                <hr>
            `;
            container.appendChild(div);
        });
    } catch (err) {
        console.error(err);
        alert("Ошибка загрузки событий");
    }
}

// загрузка броней пользователя с сервера
async function loadUserBookings() {
    try {
        const res = await fetch(`/users/${currentUser.id}/bookings`);

        // если вдруг сервер вернул не 200 — просто считаем, что броней нет
        if (!res.ok) {
            console.warn("Пока нет броней или пользователь только создан");
            myBookings = {};
            renderBookings();
            return;
        }

        const bookings = await res.json();

        myBookings = {};

        bookings.forEach(b => {
            if (b.status === "pending" || b.status === "confirmed") {
                myBookings[b.id] = {
                    eventID: b.event_id,
                    eventName: b.event_name,
                    paid: b.status === "confirmed"
                };
            }
        });

        renderBookings();
    } catch (err) {
        console.error("Ошибка загрузки броней:", err);

        // просто показываем пустой список без алерта
        myBookings = {};
        renderBookings();
    }
}

// бронирование
async function bookEvent(eventID) {
    try {
        const res = await fetch(`/events/${eventID}/book`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({
                email: currentUser.email,
                name: currentUser.name,
                telegram_id: currentUser.telegram_id,
                role: currentUser.role
            })
        });

        const data = await res.json();
        if (!res.ok) {
            alert(data.error);
            return;
        }

        const bookingID = data.booking_id;
        const eventName = eventsMap[eventID].name;

        myBookings[bookingID] = {
            eventID,
            eventName,
            paid: false
        };

        renderBookings();
    } catch (err) {
        console.error(err);
        alert("Ошибка бронирования");
    }
}

// отображение броней
function renderBookings() {
    const container = document.getElementById("bookings");
    container.innerHTML = "";

    Object.entries(myBookings).forEach(([bookingID, booking]) => {
        const div = document.createElement("div");

        div.innerHTML = `
            Бронь #${bookingID}<br>
            Событие: <b>${booking.eventName}</b><br>
            Статус: ${booking.paid ? "Оплачено ✅" : "Ожидает оплаты"}<br>
            ${
                !booking.paid
                    ? `<button onclick="confirmBooking(${bookingID})">Оплатить</button>`
                    : ""
            }
            <hr>
        `;

        container.appendChild(div);
    });
}

// подтверждение брони
async function confirmBooking(bookingID) {
    try {
        const res = await fetch(`/bookings/${bookingID}/confirm`, {
            method: "POST"
        });

        const data = await res.json();
        if (!res.ok) {
            alert(data.error || "Ошибка подтверждения");
            return;
        }

        // после успешной оплаты бронь помечаем как paid
        if (myBookings[bookingID]) {
            myBookings[bookingID].paid = true;
            renderBookings();
        }
    } catch (err) {
        console.error(err);
        alert("Ошибка подтверждения брони");
    }
}