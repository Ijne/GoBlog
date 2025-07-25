document.addEventListener('DOMContentLoaded', function() {
    const registerForm = document.getElementById('registerForm');
    
    if (registerForm) {
        registerForm.addEventListener('submit', async function(e) {
            e.preventDefault();
            
            // Валидация паролей
            const password = document.getElementById('password').value;
            const confirmPassword = document.getElementById('confirm_password').value;
            
            if (password !== confirmPassword) {
                showError('Пароли не совпадают');
                return;
            }
            
            // Сбор данных формы
            const formData = {
                username: document.getElementById('username').value.trim(),
                email: document.getElementById('email').value.trim(),
                password: password
            };
            
            try {
                // Отправка данных
                const response = await fetch('/registration', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(formData)
                });
                
                const data = await response.json();
                
                if (response.ok) {
                    if (data.redirect) {
                        window.location.href = data.redirect;
                    } else {
                        showSuccess('Регистрация прошла успешно!');
                    }
                } else {
                    showError(data.error || 'Ошибка регистрации');
                }
            } catch (error) {
                console.error('Registration error:', error);
                showError('Сетевая ошибка. Попробуйте позже.');
            }
        });
    }
    
    function showError(message) {
        // Удаляем предыдущие сообщения об ошибках
        const existingError = document.querySelector('.form-error-message');
        if (existingError) existingError.remove();
        
        // Создаем элемент с ошибкой
        const errorElement = document.createElement('div');
        errorElement.className = 'alert alert-danger form-error-message mt-3';
        errorElement.textContent = message;
        
        // Вставляем перед кнопкой отправки
        const submitButton = document.querySelector('.register-submit-btn');
        submitButton.parentNode.insertBefore(errorElement, submitButton.nextSibling);
        
        // Анимация появления
        errorElement.style.opacity = '0';
        errorElement.style.transform = 'translateY(-10px)';
        errorElement.style.transition = 'all 0.3s ease';
        
        setTimeout(() => {
            errorElement.style.opacity = '1';
            errorElement.style.transform = 'translateY(0)';
        }, 10);
    }
    
    function showSuccess(message) {
        // Аналогично showError, но для успешных сообщений
        const existingMessage = document.querySelector('.form-success-message');
        if (existingMessage) existingMessage.remove();
        
        const successElement = document.createElement('div');
        successElement.className = 'alert alert-success form-success-message mt-3';
        successElement.textContent = message;
        
        const form = document.getElementById('registerForm');
        form.prepend(successElement);
        
        successElement.style.opacity = '0';
        successElement.style.transform = 'translateY(-10px)';
        successElement.style.transition = 'all 0.3s ease';
        
        setTimeout(() => {
            successElement.style.opacity = '1';
            successElement.style.transform = 'translateY(0)';
        }, 10);
    }
});