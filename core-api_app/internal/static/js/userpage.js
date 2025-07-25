document.addEventListener('DOMContentLoaded', function() {
    const subscribeBtn = document.getElementById('subscribeBtn');
    
    if (!subscribeBtn) {
        console.error("Error: Subscribe button not found");
        return;
    }

    const userId = subscribeBtn.dataset.userId;
    let isSubscribed = subscribeBtn.dataset.initialSubscribed === 'true';

    function updateButton() {
        if (isSubscribed) {
            subscribeBtn.innerHTML = '<i class="bi bi-person-dash"></i> Отписаться';
            subscribeBtn.classList.remove('btn-outline-success');
            subscribeBtn.classList.add('btn-outline-danger');
        } else {
            subscribeBtn.innerHTML = '<i class="bi bi-person-plus"></i> Подписаться';
            subscribeBtn.classList.remove('btn-outline-danger');
            subscribeBtn.classList.add('btn-outline-success');
        }
    }

    subscribeBtn.addEventListener('click', function() {
        console.log("Button clicked, current state:", isSubscribed);
        
        const action = isSubscribed ? 'unsubscribe' : 'subscribe';
        const data = {
            target_id: parseInt(userId),
            action: action
        };

        console.log("Sending data:", data);

        fetch('/subscribe', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(data)
        })
        .then(response => {
            console.log("Response status:", response.status);
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.json();
        })
        .then(data => {
            console.log("Success:", data);
            isSubscribed = !isSubscribed;
            updateButton();
        })
        .catch(error => {
            console.error("Error:", error);
            alert("Ошибка: " + error.message);
        });
    });

    updateButton();
});