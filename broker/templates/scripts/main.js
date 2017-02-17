/**
 * Created by zhibinpan on 21/12/2016.
 */

$(function (){

    var username;
    var password;

    $('#signInBtn').click(function (event) {

        username = $('#inputUsername').val();
        password = $('#inputPassword').val();

        $.ajax({
            type: "POST",
            url: '/api/v1/login',
            data: JSON.stringify({
                username: username,
                password: password,
            }),
            success: function (data, textStatus) {
                console.log(textStatus);
                $('#loginForm').hide();
                $('#messageForm').show();
            },
            error: function (xhr, textStatus, error) {
                "use strict";
                alert(textStatus + ': ' + error);
            },
            contentType: 'application/json',
            dataType: 'json'
        });

        event.preventDefault();
    });

    $('#target').submit(function (event) {
        var topic = $('#inputTopic').val();
        var text = $('#inputText').val();

        $.ajax({
            type: "POST",
            url: '/api/v1/notification',
            data: JSON.stringify({
                topic: topic,
                message: {
                    'content': text
                },
                appId: username,
                appKey: password,
            }),
            success: function (data, textStatus) {
                console.log(textStatus);
            },
            contentType: 'application/json',
            dataType: 'json'
        });

        /*try {
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
                        appId: username,//'gftrader',
                        appKey: password,//'A98D8B1134D34F6E161463F757139'
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
        }*/

        event.preventDefault();
    });
});

