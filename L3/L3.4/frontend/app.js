const uploadBtn = document.getElementById('uploadBtn');
const fileInput = document.getElementById('fileInput');
const imagesContainer = document.getElementById('imagesContainer');

// Список изображений с их статусом
let images = [];

// Функция обновления списка изображений
function renderImages() {
    imagesContainer.innerHTML = '';
    images.forEach(img => {
        const card = document.createElement('div');
        card.className = 'image-card';

        const imageEl = document.createElement('img');
        imageEl.src = img.status === 'processing' ? '' : `/image/${img.id}`;
        imageEl.alt = img.id;

        const statusEl = document.createElement('p');
        statusEl.textContent = img.status === 'processing' ? 'В обработке...' : 'Готово';

        const getBtn = document.createElement('button');
        getBtn.textContent = 'Скачать';
        getBtn.disabled = img.status === 'processing';
        getBtn.onclick = () => {
            window.open(`/image/${img.id}`, '_blank');
        };

        const deleteBtn = document.createElement('button');
        deleteBtn.textContent = 'Удалить';
        deleteBtn.onclick = () => {
            fetch(`/image/${img.id}`, { method: 'DELETE' })
                .then(res => {
                    if(res.ok){
                        images = images.filter(i => i.id !== img.id);
                        renderImages();
                    }
                });
        };

        card.appendChild(imageEl);
        card.appendChild(statusEl);
        card.appendChild(getBtn);
        card.appendChild(deleteBtn);
        imagesContainer.appendChild(card);
    });
}

// Обновление статусов изображений каждые 2 секунды
setInterval(() => {
    images.forEach((img, idx) => {
        if(img.status === 'processing') {
            fetch(`/image/${img.id}`, { method: 'GET' })
                .then(res => {
                    if(res.status === 200){
                        images[idx].status = 'done';
                        renderImages();
                    }
                });
        }
    });
}, 2000);

// Загрузка изображения
uploadBtn.onclick = () => {
    const file = fileInput.files[0];
    if(!file) return alert('Выберите файл');

    const formData = new FormData();
    formData.append('file', file);

    fetch('/upload', { method: 'POST', body: formData })
        .then(res => res.json())
        .then(data => {
            images.push({ id: data.id, status: 'processing' });
            renderImages();
        });
};
