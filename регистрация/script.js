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
        document.getElementById('error-message').textContent = "Регистрация успешна!";
        // Redirect after successful registration
        // window.location.href = '/success'; 
    })
    .catch(error => {
        console.error('Ошибка регистрации:', error);
    });
});

