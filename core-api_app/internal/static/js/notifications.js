document.addEventListener('DOMContentLoaded', function() {
    // 1. Получаем элементы
    const badge = document.getElementById('notificationBadge');
    const notifList = document.getElementById('notificationList');
    const markAllBtn = document.getElementById('markAllAsRead');
    const deleteAllBtn = document.getElementById('deleteAllNotifications'); // Новая кнопка

    // 2. Инициализация WebSocket соединения
    let ws;
    try {
        const protocol = window.location.protocol === 'https:' ? 'wss://' : 'ws://';
        ws = new WebSocket(protocol + window.location.host + '/ws');
        
        ws.onopen = function() {
            console.log('WebSocket соединение установлено');
        };
        
        ws.onmessage = function(event) {
            try {
                const data = JSON.parse(event.data);
                addNotification(data.title, data.message);
            } catch (e) {
                console.error('Ошибка обработки сообщения:', e);
            }
        };
        
        ws.onclose = function() {
            console.log('WebSocket соединение закрыто, попытка переподключения...');
            setTimeout(() => {
                document.location.reload();
            }, 100000);
        };
        
        ws.onerror = function(error) {
            console.error('WebSocket ошибка:', error);
        };
    } catch (e) {
        console.error('Ошибка инициализации WebSocket:', e);
    }

    // 3. Загружаем уведомления из localStorage
    let notifications = [];
    try {
        const saved = localStorage.getItem('notifications');
        if (saved) notifications = JSON.parse(saved);
    } catch (e) {
        console.error('Ошибка загрузки:', e);
    }

    // 4. Сохраняем уведомления
    function save() {
        localStorage.setItem('notifications', JSON.stringify(notifications));
        console.log('💾 Сохранено:', notifications);
    }

    // 5. Обновляем бейдж
    function updateBadge() {
        if (!badge) return;
        const unread = notifications.filter(n => !n.read).length;
        badge.textContent = unread;
        badge.style.display = unread > 0 ? 'block' : 'none';
    }

    // 6. Отрисовываем список
    function renderList() {
        if (!notifList) return;
        
        notifList.innerHTML = notifications.length === 0
            ? `<div class="text-center py-4 text-muted">Нет уведомлений</div>`
            : notifications.map(notif => `
                <div class="list-group-item ${notif.read ? '' : 'bg-warning bg-opacity-10 border-start border-warning border-1'}">
                    <div class="d-flex justify-content-between align-items-start">
                        <div class="flex-grow-1" onclick="markNotificationAsRead(${notif.id})">
                            <h6 class="mb-1 ${notif.read ? '' : 'text-primary'}">${notif.title}</h6>
                            <p class="mb-1 small">${notif.message}</p>
                        </div>
                        <div class="d-flex flex-column align-items-end">
                            <small class="mb-1 ${notif.read ? 'text-muted' : 'text-primary'}">${notif.time}</small>
                            <button class="btn btn-sm btn-outline-danger" 
                                    onclick="deleteNotification(${notif.id}, event)">
                                <i class="bi bi-trash"></i>
                            </button>
                        </div>
                    </div>
                </div>
            `).join('');
    }

    // 7. Добавляем уведомление
    window.addNotification = function(title, message) {
        const newNotif = {
            id: Date.now(),
            title: title || 'Новое уведомление',
            message: message || '...',
            read: false,
            time: new Date().toLocaleTimeString()
        };
        
        notifications.unshift(newNotif);
        if (notifications.length > 50) {
            notifications.pop();
        }
        
        save();
        updateBadge();
        renderList();
        
        try {
            new Audio('/static/sounds/notification.mp3').play().catch(e => {});
        } catch (e) {}
    };

    // 8. Пометить как прочитанное
    window.markNotificationAsRead = function(id) {
        const notif = notifications.find(n => n.id === id);
        if (notif && !notif.read) {
            notif.read = true;
            save();
            updateBadge();
            renderList();
        }
    };

    // 9. Удалить уведомление
    window.deleteNotification = function(id, event) {
        event.stopPropagation(); // Предотвращаем всплытие события
        notifications = notifications.filter(n => n.id !== id);
        save();
        updateBadge();
        renderList();
    };

    // 10. Прочитать все
    if (markAllBtn) {
        markAllBtn.addEventListener('click', function() {
            notifications.forEach(n => n.read = true);
            save();
            updateBadge();
            renderList();
        });
    }

    // 11. Удалить все уведомления
    if (deleteAllBtn) {
        deleteAllBtn.addEventListener('click', function() {
            if (confirm('Вы уверены, что хотите удалить все уведомления?')) {
                notifications = [];
                save();
                updateBadge();
                renderList();
            }
        });
    }

    // 12. Первый запуск
    updateBadge();
    renderList();

    // 13. Очистка при закрытии страницы
    window.addEventListener('beforeunload', function() {
        if (ws && ws.readyState === WebSocket.OPEN) {
            ws.close();
        }
    });
});