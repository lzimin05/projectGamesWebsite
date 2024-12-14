document.getElementById('reg').addEventListener('click', function(event) {
    window.location.href = 'auth.html';
})

let err = document.getElementById('error')

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
        let errorMessage = "Неизвестная ошибка";
        if (response.status!=201) {
            throw new Error(errorMessage)
        }
    })
    .then(data => {
        console.log(data)            
        window.location.href = 'auth.html';

        //document.getElementById('error-message').textContent = "Регистрация успешна!";
        // Redirect after successful registration
        // window.location.href = '/success'; 
    })
    .catch(error => {
        err.textContent = 'пароль и имя пользователя должен быть больше 5'
        console.error('Ошибка регистрации:', error);
    });
});

