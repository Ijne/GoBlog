document.addEventListener('DOMContentLoaded', function() {
    const notificationBadge = document.getElementById('notificationBadge');
    const notificationList = document.getElementById('notificationList');
    const markAllAsReadBtn = document.getElementById('markAllAsRead');
    
    let notifications = [];
    
    // Инициализация WebSocket соединения
    const wsScheme = window.location.protocol === "https:" ? "wss://" : "ws://";
    const wsPath = wsScheme + window.location.host + "/ws";
    const socket = new WebSocket(wsPath);
    
    socket.onmessage = function(event) {
        const data = JSON.parse(event.data);
        
        // Добавляем новое уведомление
        notifications.unshift({
            id: Date.now(),
            title: data.title || 'Новое уведомление',
            message: data.message,
            read: false,
            time: new Date().toLocaleTimeString()
        });
        
        updateBadge();
        renderNotifications();
    };
    
    socket.onclose = function(event) {
        console.log('WebSocket соединение закрыто', event);
    };
    
    socket.onerror = function(error) {
        console.log('WebSocket ошибка:', error);
    };
    
    // Обновление бейджа
    function updateBadge() {
        const unreadCount = notifications.filter(n => !n.read).length;
        notificationBadge.textContent = unreadCount;
        notificationBadge.style.display = unreadCount > 0 ? 'block' : 'none';
    }
    
    // Отображение уведомлений
    function renderNotifications() {
        notificationList.innerHTML = '';
        
        if (notifications.length === 0) {
            notificationList.innerHTML = `
                <div class="text-center py-4 text-muted">
                    <i class="bi bi-bell-slash" style="font-size: 1.5rem;"></i>
                    <p class="mt-2">Нет уведомлений</p>
                </div>
            `;
            return;
        }
        
        notifications.forEach(notification => {
            const notificationItem = document.createElement('div');
            notificationItem.className = `list-group-item ${notification.read ? '' : 'active'}`;
            notificationItem.innerHTML = `
                <div class="d-flex justify-content-between align-items-start">
                    <div>
                        <div class="fw-bold">${notification.title}</div>
                        <small class="text-muted">${notification.message}</small>
                    </div>
                    <small class="text-muted">${notification.time}</small>
                </div>
            `;
            
            notificationItem.addEventListener('click', () => {
                if (!notification.read) {
                    notification.read = true;
                    updateBadge();
                    notificationItem.classList.remove('active');
                }
            });
            
            notificationList.appendChild(notificationItem);
        });
    }
    
    // Пометить все как прочитанные
    markAllAsReadBtn.addEventListener('click', () => {
        notifications = notifications.map(n => ({...n, read: true}));
        updateBadge();
        renderNotifications();
    });
    
    // Инициализация
    updateBadge();
    renderNotifications();
}); 