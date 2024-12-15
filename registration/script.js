document.getElementById('reg').addEventListener('click', function(event) {
    window.location.href = 'auth.html';
});

let err = document.getElementById('error');

document.getElementById('registrationForm').addEventListener('submit', function(event) {
    event.preventDefault();
    const form = this; // 'this' refers to the form element

    // More robust way to get form data, handles potential missing fields gracefully
    const formData = {
        name: form.name.value || "",
        email: form.email.value || "",
        password: form.password.value || ""
    };

    console.log(formData); // Log formData, not FormData

    fetch('http://localhost:8081/newuser', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(formData)
    })
    .then(response => {
        if (response.status === 201) {
            // Успешная регистрация
            alert('Успешная регистрация');
            return response.json();
        } else if (response.status === 204) {
            // Пользователь уже зарегистрирован
            throw new Error('Пользователь уже зарегистрирован');
        } else if (response.status === 400) {
            // Ошибка валидации (например, пароль или имя слишком короткие)
            throw new Error('Пароль и имя пользователя должны быть больше 5 символов');
        } else {
            // Неизвестная ошибка
            throw new Error('Неизвестная ошибка');
        }
    })
    .then(data => {
        alert('Успешная регистрация');
    })
    .catch(error => {
        // Отображение ошибки на странице
        err.textContent = error.message;
        if (err.textContent == 'Пользователь уже зарегистрирован') {
            var isAdmin = confirm("Пользователь уже зарегистрирован\nХотите авторизоваться?");
            if (isAdmin) {
                window.location.href = 'auth.html';
            }
        }
        if (err.textContent == `Unexpected token 'O', "OK" is not valid JSON`) {
            err.textContent = ''
        }
        console.error('Ошибка регистрации:', error);
    });
});