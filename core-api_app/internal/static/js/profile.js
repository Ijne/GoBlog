// Обработчик создания поста
document.getElementById('createPostForm')?.addEventListener('submit', async (e) => {
    e.preventDefault();
    
    const responseDiv = document.getElementById('postResponse');
    responseDiv.innerHTML = '';
    
    try {
        const response = await fetch('/news', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({
                title: document.getElementById('postTitle').value,
                body: document.getElementById('postContent').value
            })
        });

        if (response.ok) {
            // Очищаем форму
            document.getElementById('postTitle').value = '';
            document.getElementById('postContent').value = '';
            
            // Обновляем страницу
            window.location.reload();
        } else {
            const errorData = await response.json();
            throw new Error(errorData.error || 'Ошибка сервера');
        }
    } catch (error) {
        responseDiv.innerHTML = `
            <div class="alert alert-danger">
                Ошибка: ${error.message}
            </div>
        `;
    }
});

// Функция для инициализации обработчиков удаления постов
function initDeletePostHandlers() {
    document.querySelectorAll('.delete-post-form').forEach(form => {
        form.addEventListener('submit', async (e) => {
            e.preventDefault();
            const postId = form.dataset.id;
            const postElement = document.getElementById(`post-${postId}`);
            
            if (!confirm("Вы уверены, что хотите удалить этот пост?")) return;

            try {
                const response = await fetch(`/news?id=${postId}`, {
                    method: 'DELETE'
                });

                if (response.ok) {
                    // Плавное исчезновение поста
                    postElement.style.transition = 'opacity 0.3s, transform 0.3s';
                    postElement.style.opacity = '0';
                    postElement.style.transform = 'translateX(100px)';
                    
                    // Полное удаление из DOM после анимации
                    setTimeout(() => {
                        postElement.remove();
                        
                        // Проверяем остались ли посты
                        if (!document.querySelector('.post-card')) {
                            document.querySelector('.my-posts').innerHTML = `
                                <div class="empty-posts alert alert-light rounded shadow-sm text-center py-5">
                                    <i class="bi bi-info-circle-fill text-muted" style="font-size: 2.5rem;"></i>
                                    <h4 class="mt-3 text-muted">У вас пока нет постов</h4>
                                    <p class="text-muted">Создайте свой первый пост выше</p>
                                </div>`;
                        }
                        
                        // Обновляем счетчик постов в шапке
                        const countBadge = document.querySelector('.badge.bg-light.text-dark');
                        const currentCount = parseInt(countBadge.textContent.split(': ')[1]);
                        countBadge.innerHTML = `<i class="bi bi-postcard"></i> Постов: ${currentCount - 1}`;
                    }, 300);
                } else {
                    const error = await response.text();
                    alert(`Ошибка удаления: ${error}`);
                }
            } catch (error) {
                alert(`Сетевая ошибка: ${error.message}`);
            }
        });
    });
}

// Инициализация при загрузке страницы
document.addEventListener('DOMContentLoaded', () => {
    initDeletePostHandlers();
    
    // Если используется динамическая загрузка контента, нужно вызывать initDeletePostHandlers()
    // после каждого изменения DOM
});