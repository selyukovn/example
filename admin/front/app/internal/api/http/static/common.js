
function handleResponse(
    onSuccessJson,
    onErrorJson = function (data, response) {
        let msg = data.message || ("Ошибка: " + response.statusText)
        console.error(msg)
    },
    onErrorText = function (text) {
        let msg = "Ошибка: " + text
        console.error(msg)
    },
    onSuccessText = function (text, response) {
        console.info(text, response)
    },
) {
    return function (response) {
        if (response.ok) {
            console.info("response:", response)
        } else {
            console.error("response:", response)
        }

        const isJson = (response.headers.get("Content-Type") || "").includes("application/json")
        if (isJson) {
            response.json()
                .then(data => {
                    if (response.ok) {
                        onSuccessJson && onSuccessJson(data, response)
                    } else {
                        onErrorJson && onErrorJson(data, response)
                    }
                })
                .catch(err => {
                    console.error("json parse error:", err)
                })
        } else {
            response.text()
                .then(text => {
                    if (response.ok) {
                        onSuccessText && onSuccessText(text, response)
                    } else {
                        // Например, 503, 524, ... -- когда ответ генерируется не в приложении.
                        onErrorText && onErrorText(text, response)
                    }
                })
                .catch(err => {
                    // невозможный(?) случай
                    console.error("text parse error:", err)
                })
        }
    }
}
