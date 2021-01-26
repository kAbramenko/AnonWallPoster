document.getElementById('feedback-form').addEventListener('submit', function (evt) {
    if (this.msg.value == '') {
        alert("Заполните поле!")
        return
    }
    var ctx = this;
    var link = document.getElementById('submit_button');
    var loading_bar = document.getElementById('load_bar');
    var old_display = link.style.display;
    var old_visible = link.style.visibility;
    link.style.visibility = 'hidden';
    link.style.display = 'none';
    loading_bar.style.display = 'block';
    var http = new XMLHttpRequest();
    evt.preventDefault();
    http.open("POST", "/api/post", true);
    var form = new FormData(document.forms.feedback)

    http.send(form);
    http.onreadystatechange = function () {
        if (http.readyState == 4 && http.status == 200) {
            alert('Ваше сообщение получено.\nНаши модераторы проверят Ваш материал в течении нескольких минут.\nБлагодарим за интерес к нашему сообществу!');
            ctx.msg.removeAttribute('value'); // очистить поле сообщения (две строки)
            ctx.msg.value = '';
            ctx.photo.removeAttribute('value')
            ctx.photo.value = '';
            link.style.visibility = old_visible;
            link.style.display = old_display;
            loading_bar.style.display = 'none';
        } else if (http.status != 200) {
            alert("Сервис временно не доступен")
            link.style.visibility = old_visible;
            link.style.display = old_display;
            loading_bar.style.display = 'none';
        }
    }
    http.onerror = function () {
        alert('Извините, данные не были переданы');
        link.style.visibility = old_visible;
        link.style.display = old_display;
        loading_bar.style.display = 'none';
    }
}, false);
