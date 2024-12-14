let err_message = document.getElementById('err'); // Убедитесь, что элемент с id="err" существует в HTML

document.getElementById('reg').addEventListener('click', function(event) {
    window.location.href = 'title.html';
});

document.getElementById('form').addEventListener('submit', function(event) {
    event.preventDefault(); // Предотвращаем отправку формы по умолчанию

    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;

    // Отправляем POST-запрос на /login
    fetch('http://localhost:8081/login', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            email: email,
            password: password,
        }),
    })
    .then(response => {
        if (!response.ok) {
            // Обрабатываем ошибки
            if (response.status === 401) {
                // Если статус 401 (Unauthorized)
                throw new Error('ошибка входа'); // Выбрасываем ошибку с текстом "ошибка входа"
            } else {
                // Для других ошибок
                throw new Error(`Ошибка ${response.status}: ${response.statusText}`);
            }
        }
        return response.json(); // Если запрос успешный, возвращаем данные
    })
    .then(data => {
        // Обработка успешного ответа

        // Сохраняем JWT-токен (предположим, что он возвращается в поле "token")
        const token = data.token;

        // Сохраняем токен в локальном хранилище
        localStorage.setItem('jwtToken', token);

        // Выполняем GET-запрос на /restricted/user/:email с использованием JWT-токена
        return fetch(`http://localhost:8081/restricted/user/${email}`, {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${token}`, // Добавляем токен в заголовок
            },
        });
    })
    .then(response => {
        if (!response.ok) {
            // Обрабатываем ошибки GET-запроса
            throw new Error(`Ошибка ${response.status}: ${response.statusText}`);
        }
        return response.json(); // Возвращаем данные GET-запроса
    })
    .then(userData => {
        // Выводим результат GET-запроса в alert
        alert(`Успешный вход: ${JSON.stringify(userData)}`);    
        window.location.href = '../static/main.html'; // Перенаправление на другую страницу
    })
    .catch(error => {
        // Обработка ошибок
        if (error.message === 'ошибка входа') {
            // Если ошибка связана с 401
            err_message.textContent = 'ошибка входа'; // Изменяем текст сообщения
        } else {
            // Для других ошибок
            err_message.textContent = 'Произошла ошибка: ' + error.message;
        }
    });
});