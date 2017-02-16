/**
 * Created by zhibinpan on 21/12/2016.
 */

$(function (){

    $('#signInBtn').click(function (event) {

        var username = $('#inputUsername').val();
        var password = $('#inputPassword').val();

        $.ajax({
            type: "POST",
            url: '/api/v1/login',
            data: JSON.stringify({
                username: username,//'gftrader',
                password: password,//'A98D8B1134D34F6E161463F757139'
            }),
            success: function (data, textStatus) {
                console.log(textStatus);
                $('#loginForm').hide();
                $('#messageForm').show();
            },
            contentType: 'application/json',
            dataType: 'json'
        });

        event.preventDefault();
    });

    $('#target').submit(function (event) {
        var topic = $('#inputTopic').val();
        var text = $('#inputText').val();

        try {
            var obj = JSON.parse(text);
            if (typeof obj !== 'object') {
                alert('must provide message in JSON format');
            } else {
                $.ajax({
                    type: "POST",
                    url: '/api/v1/notification',
                    data: JSON.stringify({
                        topic: topic,
                        message: obj,
                        appId: 'gftrader',
                        appKey: 'A98D8B1134D34F6E161463F757139'
                    }),
                    success: function (data, textStatus) {
                        console.log(textStatus);
                    },
                    contentType: 'application/json',
                    dataType: 'json'
                });
            }
        }
        catch (err) {
            console.error(err);
            alert('must provide message in JSON format');
        }

        event.preventDefault();
    });
});

