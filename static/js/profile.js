document.addEventListener('DOMContentLoaded', function() {
    updateUserInfo();
});

const BASE_URL = 'http://localhost:8081';

// Функция для получения сессии (id пользователя)
async function getSessionId() {
    const token = localStorage.getItem('jwtToken');
    if (!token) {
        window.location.href = '../registration/auth.html';
        console.error('Token not found');
        return null;
    }

    try {
        const response = await fetch(`${BASE_URL}/restricted/sesion`, {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${token}`,
            },
        });

        if (response.ok) {
            const sessionId = await response.text();
            console.log('Session ID:', sessionId);
            return sessionId;
        } else {
            console.error('Failed to fetch session ID:', response.statusText);
            window.location.href = '../registration/auth.html';
            return null;
        }
    } catch (error) {
        console.error('Error fetching session ID:', error);
        return null;
    }
}

// Функция для получения почты пользователя по id
async function getUserEmailById(id) {
    const token = localStorage.getItem('jwtToken');
    if (!token) {
        window.location.href = '../registration/auth.html';
        console.error('Token not found');
        return null;
    }

    try {
        const response = await fetch(`${BASE_URL}/restricted/email/${id}`, {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${token}`,
            },
        });

        if (response.ok) {
            const email = await response.text(); // Получаем почту как строку
            //email =  // Удаляем первый и последний символ
            console.log('User Email:', email.slice(1, -2));
            return email.slice(1, -2);
        } else {
            console.error('Failed to fetch user email:', response.statusText);
            return null;
        }
    } catch (error) {
        console.error('Error fetching user email:', error);
        return null;
    }
}

// Функция для получения имени пользователя по id
async function getUserNameById(id) {
    const token = localStorage.getItem('jwtToken');
    if (!token) {
        console.error('Token not found');
        return null;
    }

    try {
        const response = await fetch(`${BASE_URL}/restricted/user/${id}`, {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${token}`,
            },
        });

        if (response.ok) {
            const name = await response.text(); // Получаем имя как строку
            console.log('User Name:', name.slice(1, -2));
            return name.slice(1, -2);; 
        } else {
            console.error('Failed to fetch user name:', response.statusText);
            return null;
        }
    } catch (error) {
        console.error('Error fetching user name:', error);
        return null;
    }
}

// Функция для обновления данных пользователя на странице
async function updateUserInfo() {
    const sessionId = await getSessionId();
    if (sessionId) {
        const email = await getUserEmailById(sessionId);
        const name = await getUserNameById(sessionId);

        if (email && name) {
            // Обновляем имя пользователя на странице
            const usernameElement = document.getElementById('username');
            if (usernameElement) {
                usernameElement.textContent = name;
            } else {
                console.error('Элемент с id="username" не найден');
            }

            // Обновляем почту пользователя на странице
            const emailElement = document.getElementById('email');
            if (emailElement) {
                emailElement.textContent = email;
            } else {
                console.error('Элемент с id="email" не найден');
            }
        }
    }
}

// Обработчик для кнопки редактирования имени пользователя
document.getElementById('edit-username').addEventListener('click', async function() {
    const newUsername = prompt('Введите новое имя пользователя (не менее 5 символов):');
    if (newUsername && newUsername.trim().length >= 5) {
        // Проверяем длину имени
        if (newUsername.trim().length > 30) {
            alert('Имя пользователя должно быть не более 30 символов!');
            return;
        }

        // Получаем id пользователя из сессии
        const sessionId = await getSessionId();
        if (!sessionId) {
            console.error('Не удалось получить id пользователя из сессии');
            return;
        }

        // Отправляем PUT-запрос на сервер
        const token = localStorage.getItem('jwtToken');
        if (!token) {
            window.location.href = '../registration/auth.html';
            console.error('Token not found');
            return;
        }

        try {
            const response = await fetch(`${BASE_URL}/restricted/user/${sessionId}`, {
                method: 'PUT',
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    newname: newUsername.trim(), // Отправляем новое имя в формате JSON
                }),
            });

            if (response.ok) {
                // Обновляем интерфейс
                const usernameElement = document.getElementById('username');
                if (usernameElement) {
                    usernameElement.textContent = newUsername.trim();
                }
                alert('Имя пользователя успешно обновлено!');
            } else {
                const errorData = await response.text();
                console.error('Ошибка обновления имени:', errorData);
                alert('Ошибка обновления имени. Попробуйте еще раз.');
            }
        } catch (error) {
            console.error('Ошибка при отправке запроса:', error);
            alert('Произошла ошибка. Попробуйте еще раз.');
        }
    } else if (newUsername && newUsername.trim().length < 5) {
        alert('Имя пользователя должно быть не менее 5 символов!');
    } else {
        alert('Имя пользователя не может быть пустым!');
    }
});

// Обработчик для кнопки выхода
document.getElementById('logout-button').addEventListener('click', function() {
    alert('Вы вышли из системы.');
    // Получаем токен из локального хранилища
    const token = localStorage.getItem('jwtToken');
    // Проверяем, есть ли токен
    if (token) {
        // Отправляем POST-запрос на /logout
        fetch(`${BASE_URL}/logout`, {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${token}`, // Добавляем токен в заголовок
            },
        })
        .then(response => {
            if (!response.ok) {
                throw new Error(`Ошибка выхода: ${response.status} ${response.statusText}`);
            }
            return response.json(); // Если запрос успешен, возвращаем данные
        })
        .then(data => {
            console.log('Выход выполнен успешно:', data);
        })
        .catch(error => {
            console.error('Ошибка выхода:', error.message);
        })
        .finally(() => {
            // Удаляем токен из локального хранилища
            localStorage.removeItem('jwtToken');

            // Перенаправляем пользователя на страницу входа
            window.location.href = '../registration/auth.html';
            console.log("click menu!");
        });
    } else {
        // Если токена нет, просто перенаправляем на страницу входа
        window.location.href = '../registration/auth.html';
        console.log("click menu!");
    }
    // Здесь можно добавить логику для перенаправления на страницу входа
});

// Обработчик для кнопки возврата на главную страницу
document.getElementById('back-button').addEventListener('click', function() {
    window.location.href = 'main.html'; // Перенаправление на главную страницу
});