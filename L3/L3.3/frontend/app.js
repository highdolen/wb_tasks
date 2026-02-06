const API = "http://localhost:8080";
let currentPage = 1;
const limit = 5;

async function fetchComments(parentId = null) {
    let url = `${API}/comments`;
    if (parentId !== null) url += `?parent=${parentId}`;
    else url += `?page=${currentPage}&limit=${limit}`;

    const res = await fetch(url);
    if (!res.ok) throw new Error("Ошибка запроса");
    return await res.json();
}

async function createComment(text, parentId = null) {
    await fetch(`${API}/comments`, {
        method: "POST",
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify({text, parent_id: parentId})
    });
}

async function deleteComment(id) {
    await fetch(`${API}/comments/${id}`, {method: "DELETE"});
}

function renderComment(comment, container, level = 0) {
    const div = document.createElement("div");
    div.className = "comment";
    div.style.marginLeft = `${level * 20}px`;

    let childrenHtml = `<div class="actions">
        <button data-reply="${comment.id}">Ответить</button>
        <button data-toggle="${comment.id}">+</button>
        <button data-delete="${comment.id}">Удалить</button>
    </div>`;

    div.innerHTML = `<div><b>#${comment.id}</b> ${comment.text}</div>${childrenHtml}`;
    container.appendChild(div);

    const childrenContainer = document.createElement("div");
    childrenContainer.id = `children-${comment.id}`;
    div.appendChild(childrenContainer);

    if (comment.Children && comment.Children.length) {
        comment.Children.forEach(child => renderComment(child, childrenContainer, level + 1));
    }
}

async function loadComments() {
    const container = document.getElementById("comments");
    container.innerHTML = "";

    try {
        const comments = await fetchComments();
        comments.forEach(c => renderComment(c, container));
        document.getElementById("pageInfo").innerText = currentPage;
    } catch (e) {
        container.innerText = "Нет комментариев";
    }
}

document.addEventListener("click", async e => {
    if (e.target.dataset.reply) {
        const text = prompt("Текст ответа:");
        if (text) {
            await createComment(text, Number(e.target.dataset.reply));
            loadComments();
        }
    }

    if (e.target.dataset.delete) {
        await deleteComment(Number(e.target.dataset.delete));
        loadComments();
    }

    if (e.target.dataset.toggle) {
        const parentId = Number(e.target.dataset.toggle);
        const childrenContainer = document.getElementById(`children-${parentId}`);
        if (childrenContainer.style.display === "none" || !childrenContainer.style.display) {
            const children = await fetchComments(parentId);
            childrenContainer.innerHTML = "";
            children.forEach(c => renderComment(c, childrenContainer, 1));
            childrenContainer.style.display = "block";
            e.target.innerText = "-";
        } else {
            childrenContainer.style.display = "none";
            e.target.innerText = "+";
        }
    }
});

document.getElementById("addRoot").onclick = async () => {
    const text = document.getElementById("newText").value;
    if (!text) return;
    await createComment(text);
    document.getElementById("newText").value = "";
    loadComments();
};

document.getElementById("searchBtn").onclick = async () => {
    const q = document.getElementById("searchInput").value;
    if (!q) return;
    const res = await fetch(`${API}/comments/search?q=${q}`);
    const data = await res.json();
    const container = document.getElementById("comments");
    container.innerHTML = "";
    data.forEach(c => renderComment(c, container));
};

document.getElementById("resetBtn").onclick = () => {
    document.getElementById("searchInput").value = "";
    loadComments();
};

document.getElementById("prevPage").onclick = () => {
    if (currentPage > 1) {
        currentPage--;
        loadComments();
    }
};
document.getElementById("nextPage").onclick = () => {
    currentPage++;
    loadComments();
};

loadComments();
