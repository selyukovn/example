const EMAIL_STEP = document.getElementById('email-step');
const CODE_STEP = document.getElementById('code-step');
const RESEND_LINK = document.getElementById('resend-link');
const ERROR = document.getElementById('error');
const RELOAD_LINK = document.getElementById("reload-link")
const EXPIRE = document.getElementById("expire-at")
const CONFIRM_BTN = document.getElementById("confirm-button")
const CODE_INPUT = document.getElementById("code")

// ---------------------------------------------------------------------------------------------------------------------
// Request
// ---------------------------------------------------------------------------------------------------------------------

let signInId = '';

function requestCfm() {
    ERROR.textContent = ""

    const email = document.getElementById('email').value
    if (!email) {
        ERROR.textContent = 'Пожалуйста, введите email'
        return
    }

    fetch(PAGE_SETTINGS_UrlSignInRequest, {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify({
            email: email,
        })
    }).then(handleResponse(function (data) {
        signInId = data.sign_in_id
        retriesLeft = data.retries_left
        canRetryAt = new Date(data.can_retry_at)
        expireAt = new Date(data.expire_at)
        // show code step
        EMAIL_STEP.style.display = 'none'
        CODE_STEP.style.display = 'block'
        RESEND_LINK.style.display = retriesLeft > 0 ? 'block' : 'none'
        RELOAD_LINK.style.display = retriesLeft > 0 ? 'block' : 'none'
        // run expire timer
        function updateExpireTime() {
            EXPIRE.textContent = "Протухнет через: " + (function() {
                const diffMs = expireAt - Date.now()

                if (diffMs <= 0) {
                    return "00:00";
                }

                const totalSeconds = Math.floor(diffMs / 1000);
                const minutes = Math.floor(totalSeconds / 60);
                const seconds = totalSeconds % 60;

                const formattedMinutes = String(minutes).padStart(2, '0');
                const formattedSeconds = String(seconds).padStart(2, '0');

                return `${formattedMinutes}:${formattedSeconds}`;
            })()
        }
        updateExpireTime();
        const expireTimeIntervalId = setInterval(updateExpireTime, 1000)
        setTimeout(function() {
            clearInterval(expireTimeIntervalId)
            RESEND_LINK.style.display = "none"
            CONFIRM_BTN.disabled = true
            CODE_INPUT.disabled = true
            EXPIRE.textContent = "Протухло"
            ERROR.textContent = "Поезд ушел"
        }, expireAt - Date.now() + 1000)
    }, function (data, response) {
        ERROR.textContent = data.message || ("Ошибка: " + response.statusText)
    }, function (text) {
        ERROR.textContent = "Ошибка: " + text
    })).catch(err => {
        console.error(err)
        ERROR.textContent = "Ошибка сети: " + err.message
    })
}

// ---------------------------------------------------------------------------------------------------------------------
// Request Retry
// ---------------------------------------------------------------------------------------------------------------------

let retriesLeft = 0;
let canRetryAt = Date.now();
let expireAt = Date.now();

function requestCfmRetry() {
    ERROR.textContent = ""

    if (retriesLeft <= 0) {
        ERROR.textContent = "Повторы кончились"
        return false
    }

    if (Date.now() < canRetryAt) {
        ERROR.textContent = "Еще рано. Можно после " + (
            String(canRetryAt.getHours()).padStart(2, '0')
            + ":"
            + String(canRetryAt.getMinutes()).padStart(2, '0')
        )
        return false
    }

    fetch(PAGE_SETTINGS_UrlSignInRequestRetry, {
        method: 'PUT',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify({
            sign_in_id: signInId
        })
    }).then(handleResponse(function (data) {
        retriesLeft = data.retries_left
        canRetryAt = new Date(data.can_retry_at)
        RESEND_LINK.style.display = retriesLeft > 0 ? 'block' : 'none'
    }, function (data, response) {
        ERROR.textContent = data.message || ("Ошибка: " + response.statusText)
    }, function (text) {
        ERROR.textContent = "Ошибка: " + text
    })).catch(err => {
        console.error(err)
        ERROR.textContent = "Ошибка сети: " + err.message
    })

    return false
}

// ---------------------------------------------------------------------------------------------------------------------
// Confirm
// ---------------------------------------------------------------------------------------------------------------------

function confirm() {
    ERROR.textContent = ""

    const code = document.getElementById('code').value
    if (!code) {
        ERROR.textContent = "Пожалуйста, введите код"
        return false
    }

    fetch(PAGE_SETTINGS_UrlSignInConfirm, {
        method: 'PUT',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify({
            sign_in_id: signInId,
            code: code,
        })
    }).then(handleResponse(function (data) {
        if (data.is_passed) {
            window.location.href = data.redirect_url
        } else {
            ERROR.textContent = "Неверный код! Осталось попыток: " + data.attempts_left
        }
    }, function (data, response) {
        ERROR.textContent = data.message || ("Ошибка: " + response.statusText)
    }, function (text) {
        ERROR.textContent = "Ошибка: " + text
    })).catch(err => {
        console.error(err)
        ERROR.textContent = "Ошибка сети: " + err.message
    })
}

// ---------------------------------------------------------------------------------------------------------------------

document.getElementById("resend-link").addEventListener("click", e => {
    e.stopPropagation();
    requestCfmRetry();
    return false;
})

document.getElementById("confirm-button").addEventListener("click", e => {
    e.stopPropagation();
    confirm();
    return false;
})

document.getElementById("email-step-continue-button").addEventListener("click", e => {
    e.stopPropagation();
    requestCfm();
    return false;
})

document.getElementById("reload-link").addEventListener("click", e => {
    e.stopPropagation();
    document.location.reload();
    return false;
})

// ---------------------------------------------------------------------------------------------------------------------
