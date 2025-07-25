document.addEventListener('DOMContentLoaded', function() {
    const loginForm = document.getElementById('loginForm');
    
    if (loginForm) {
        loginForm.addEventListener('submit', async function(e) {
            e.preventDefault();
            
            // Сброс предыдущих ошибок
            resetErrors();
            
            // Сбор данных формы
            const email = document.getElementById('email').value.trim();
            const password = document.getElementById('password').value;
            
            // Валидация
            if (!validateEmail(email)) {
                showError('email', 'Введите корректный email');
                return;
            }
            
            // Показ состояния загрузки
            const submitBtn = loginForm.querySelector('.login-submit-btn');
            const originalBtnText = submitBtn.textContent;
            submitBtn.disabled = true;
            submitBtn.innerHTML = '<span class="spinner-border spinner-border-sm" role="status" aria-hidden="true"></span> Вход...';
            
            try {
                // Отправка данных
                const response = await fetch('/login', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({ email, password })
                });
                
                const data = await response.json();
                
                if (response.ok) {
                    if (data.redirect) {
                        window.location.href = data.redirect;
                    } else {
                        showSuccess('Вход выполнен успешно!');
                    }
                } else {
                    showFormError(data.error || 'Неверный email или пароль');
                }
            } catch (error) {
                console.error('Login error:', error);
                showFormError('Сетевая ошибка. Попробуйте позже.');
            } finally {
                // Восстановление кнопки
                submitBtn.disabled = false;
                submitBtn.textContent = originalBtnText;
            }
        });
    }
    
    function validateEmail(email) {
        const re = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        return re.test(email);
    }
    
    function showError(fieldId, message) {
        const field = document.getElementById(fieldId);
        const formGroup = field.closest('.form-group');
        
        formGroup.classList.add('has-error');
        
        let errorElement = formGroup.querySelector('.error-message');
        if (!errorElement) {
            errorElement = document.createElement('div');
            errorElement.className = 'error-message';
            formGroup.appendChild(errorElement);
        }
        
        errorElement.textContent = message;
    }
    
    function resetErrors() {
        document.querySelectorAll('.form-group').forEach(group => {
            group.classList.remove('has-error');
            const errorElement = group.querySelector('.error-message');
            if (errorElement) {
                errorElement.textContent = '';
            }
        });
    }
    
    function showFormError(message) {
        const errorElement = document.createElement('div');
        errorElement.className = 'alert alert-danger mt-3';
        errorElement.textContent = message;
        
        const form = document.getElementById('loginForm');
        form.prepend(errorElement);
        
        // Автоматическое скрытие через 5 секунд
        setTimeout(() => {
            errorElement.style.transition = 'opacity 0.5s ease';
            errorElement.style.opacity = '0';
            setTimeout(() => errorElement.remove(), 500);
        }, 5000);
    }
    
    function showSuccess(message) {
        const successElement = document.createElement('div');
        successElement.className = 'alert alert-success mt-3';
        successElement.textContent = message;
        
        const form = document.getElementById('loginForm');
        form.prepend(successElement);
        
        // Автоматическое скрытие через 3 секунды
        setTimeout(() => {
            successElement.style.transition = 'opacity 0.5s ease';
            successElement.style.opacity = '0';
            setTimeout(() => successElement.remove(), 500);
        }, 3000);
    }
});