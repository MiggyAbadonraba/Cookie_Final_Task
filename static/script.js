document.addEventListener('DOMContentLoaded', () => {
  const form = document.getElementById('cookie-form');
  const list = document.getElementById('cookie-list');

  // Загрузка и отрисовка всех
  async function loadAll() {
    const res = await fetch('/api/cookies');
    const items = await res.json();
    list.innerHTML = '';
    items.forEach(renderItem);
  }

  // Отрисовать одну запись с кнопками
  function renderItem(c) {
    const li = document.createElement('li');
    li.innerHTML = `
      <strong>${c.name}</strong>: ${c.ingredients} ${c.description ? '— '+c.description : ''}
      <button data-id="${c.id}" class="edit">✎</button>
      <button data-id="${c.id}" class="delete">🗑</button>
    `;
    list.appendChild(li);

    li.querySelector('.delete').onclick = async () => {
      await fetch(`/api/cookies/${c.id}`, { method: 'DELETE' });
      loadAll();
    };
    li.querySelector('.edit').onclick = () => {
      form.id.value = c.id;
      form.name.value = c.name;
      form.ingredients.value = c.ingredients;
      form.description.value = c.description;
    };
  }

  // Обработка submit для create/update
  form.onsubmit = async e => {
    e.preventDefault();
    const data = {
      name: form.name.value,
      ingredients: form.ingredients.value,
      description: form.description.value,
    };
    const id = form.id.value;
    const url = id? `/api/cookies/${id}` : '/api/cookies';
    const method = id? 'PUT' : 'POST';

    await fetch(url, {
      method,
      headers: {'Content-Type':'application/json'},
      body: JSON.stringify(data)
    });
    form.reset();
    loadAll();
  };

// Начальная загрузка
  loadAll();

// shutdown
const shutdownBtn = document.getElementById('shutdown-btn');
shutdownBtn.addEventListener('click', async () => {
  shutdownBtn.disabled = true;
  shutdownBtn.textContent = 'Выключаем…';
  try {
    // шлём POST, получаем ответ
    const res = await fetch('/api/shutdown', { method: 'POST' });
    console.log('[INFO]', await res.text());
  } catch (err) {
    console.error('[ERROR] shutdown:', err);
  }
  // спустя 2 секунды перезагружаем страницу
  setTimeout(() => {
    location.reload();
  }, 1000);
});
});
