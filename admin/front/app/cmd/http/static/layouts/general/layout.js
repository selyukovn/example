document.addEventListener('DOMContentLoaded', function() {
    // Обработчик выхода
    document.querySelector('.sign-out-btn').addEventListener('click', function(e) {
        e.preventDefault();
        if (confirm('Вы уверены, что хотите выйти?')) {
            fetch(e.target.attributes['data-url'].value, {
                method: 'DELETE',
                headers: {'Content-Type': 'application/json'},
            }).then(handleResponse(function (data) {
                document.location.href = data.redirect_url
            })).catch(err => {
                console.error(err)
                ERROR.textContent = "Ошибка сети: " + err.message
            })
        }
    });
});
