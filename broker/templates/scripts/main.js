/**
 * Created by zhibinpan on 21/12/2016.
 */

$(function (){
    $('#target').submit(function (event) {
        var topic = $('#inputTopic').val();
        var text = $('#inputText').val();

        $.ajax({
            type: "POST",
            url: '/api/v1/notification',
            data: JSON.stringify({topic: topic, message: text}),
            success: function(data, textStatus) {
                console.log(textStatus + ': ' + data);
            },
            contentType: 'application/json',
            dataType: 'json'
        });

        event.preventDefault();
    });
});

