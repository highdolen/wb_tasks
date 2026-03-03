// Базовый путь к API
const API = "/api";

// Состояние отображения таблицы транзакций
let tableVisible = false;



// Загрузка списка транзакций с учётом фильтров
async function loadItems() {
  const type = (document.getElementById("filterType")?.value || "").trim();
  const category = (document.getElementById("filterCategory")?.value || "").trim();
  const from = (document.getElementById("filterFrom")?.value || "").trim();
  const to = (document.getElementById("filterTo")?.value || "").trim();

  const params = [];

  if (type !== "") params.push(`type=${encodeURIComponent(type)}`);
  if (category !== "") params.push(`category=${encodeURIComponent(category)}`);
  if (from !== "") params.push(`from=${from}T00:00:00Z`);
  if (to !== "") params.push(`to=${to}T23:59:59Z`);

  let url = `${API}/items`;

  if (params.length > 0) {
    url += "?" + params.join("&");
  }

  const res = await fetch(url);
  const data = await res.json();

  const table = document.getElementById("itemsTable");
  table.innerHTML = "";

  data.forEach(item => {
    const date = new Date(item.created_at).toLocaleString();

    table.innerHTML += `
      <tr>
        <td>${item.id}</td>
        <td>${item.type}</td>
        <td>${item.amount}</td>
        <td>${item.category || ""}</td>
        <td>${date}</td>
        <td>
          <button onclick="deleteItem(${item.id})">Delete</button>
        </td>
      </tr>
    `;
  });
}


// Переключение отображения таблицы транзакций
function toggleTable() {
  const block = document.getElementById("itemsBlock");
  const btn = document.getElementById("toggleBtn");

  tableVisible = !tableVisible;
  block.style.display = tableVisible ? "block" : "none";
  btn.textContent = tableVisible ? "Скрыть транзакции" : "Показать транзакции";

  if (tableVisible) loadItems();
}


// Удаление транзакции
async function deleteItem(id) {
  await fetch(`${API}/items/${id}`, { method: "DELETE" });
  if (tableVisible) loadItems();
}


// Отправка формы создания новой транзакции
document.getElementById("itemForm").addEventListener("submit", async e => {
  e.preventDefault();

  const createdAtValue = document.getElementById("createdAt").value;

  const item = {
    type: document.getElementById("type").value,
    amount: parseFloat(document.getElementById("amount").value),
    category: document.getElementById("category").value,
    created_at: new Date(createdAtValue).toISOString()
  };

  await fetch(`${API}/items`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(item)
  });

  e.target.reset();
  setCurrentDateTime();

  if (tableVisible) loadItems();
});



// Загрузка аналитики и сгруппированных данных
async function loadAnalytics() {
  try {
    const from = document.getElementById("fromDate").value;
    const to = document.getElementById("toDate").value;
    const group = document.getElementById("groupBy").value;

    const res = await fetch(
      `${API}/analytics?from=${from}T00:00:00Z&to=${to}T23:59:59Z`
    );
    const summary = await res.json();

    const resGrouped = await fetch(
      `${API}/analytics/grouped?from=${from}T00:00:00Z&to=${to}T23:59:59Z&group_by=${group}`
    );
    const grouped = await resGrouped.json();

    const result = document.getElementById("analyticsResult");

    result.innerHTML = `
      <pre>${JSON.stringify({ summary, grouped }, null, 2)}</pre>
    `;
  } catch (err) {
    console.error(err);
    alert("Ошибка загрузки аналитики");
  }
}


// Применение фильтров к таблице транзакций
function applyFilters() {
  const block = document.getElementById("itemsBlock");
  const btn = document.getElementById("toggleBtn");

  if (!tableVisible) {
    tableVisible = true;
    block.style.display = "block";
    btn.textContent = "Скрыть транзакции";
  }

  loadItems();
}


// Экспорт аналитики в CSV
function exportCSV() {
  const from = document.getElementById("fromDate").value;
  const to = document.getElementById("toDate").value;

  window.open(
    `${API}/analytics/export/csv?from=${from}T00:00:00Z&to=${to}T23:59:59Z`
  );
}


// Установка дефолтного периода (с начала месяца)
function setDefaultDates() {
  const today = new Date();
  const start = new Date(today.getFullYear(), today.getMonth(), 1);

  document.getElementById("fromDate").value = start.toISOString().split("T")[0];
  document.getElementById("toDate").value = today.toISOString().split("T")[0];
}


// Установка текущей даты и времени в форму
function setCurrentDateTime() {
  const now = new Date();
  now.setMinutes(now.getMinutes() - now.getTimezoneOffset());
  document.getElementById("createdAt").value = now.toISOString().slice(0,16);
}


// Инициализация значений при загрузке страницы
setDefaultDates();
setCurrentDateTime();