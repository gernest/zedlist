$(document).ready(function(){
    $('.ui.dropdown').dropdown({
        on: 'hover'
    });
    $('.message .close').on('click',function(){
        $(this).closest('.message')
        .transition('fade');
    });
})


Moon.component("basic",{
    template: `<p>{{name}} hello , world </p>`,
    data: function(){
        return {
            name: "gernest"
        }
    }
});

const app=new Moon({
    el: "#app"
});