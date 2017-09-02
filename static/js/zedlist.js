$(document).ready(function(){
    $('.ui.dropdown').dropdown({
        on: 'hover'
    });
    $('.message .close').on('click',function(){
        $(this).closest('.message')
        .transition('fade');
    });
})