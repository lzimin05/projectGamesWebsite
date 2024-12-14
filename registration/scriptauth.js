var flag = false //зашел ли пользователь

document.getElementById('reg').addEventListener('click', function(event) {
    window.location.href = 'title.html';
})

document.getElementById('form').addEventListener('submit', function(event) {
    event.preventDefault();
    const form = this; // 'this' refers to the form element
    var email = document.getElementById('email')
    var password = document.getElementById('password')
    // More robust way to get form data, handles potential missing fields gracefully
    const formData = {
        email: email.value || "",
        password: password.value || ""
    };

    console.log(formData); // Log formData, not FormData

    fetch('http://localhost:8081/login', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(formData)
    })
    .then(response => {
        if (!response.ok) {
            return response.json().then(data => {
                let errorMessage = "Неизвестная ошибка";
                // Improved error handling: check for different error structures
                if (data && data.message) {
                    errorMessage = data.message;
                } else if (data && data.errors) { // Handle potential validation errors from backend
                    errorMessage = Object.values(data.errors).join(', ');
                } else if (response.status === 409) {
                    errorMessage = "Email уже существует";
                } else if (response.status === 400) {
                    errorMessage = "Некорректные данные";
                } else {
                    errorMessage = `Ошибка ${response.status}: ${response.statusText}`;
                }
                document.getElementById('error-message').textContent = errorMessage;
                throw new Error(errorMessage);
            });
        }
        return response.json();
    })
    .then(data => {
        if (data === true) {
            flag = true
            //зашел!!!
            window.location.href = '../main.html';

        } 
        document.getElementById('error-message').textContent = "Регистрация успешна!";
        // Redirect after successful registration
        // window.location.href = '/success'; 
    })
    .catch(error => {
        document.getElementById('err').textContent = 'ошибка входа';
        console.error('Ошибка авторизации:', error);
    });
});




