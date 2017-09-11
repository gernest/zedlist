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