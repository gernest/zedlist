
function send(endpoint, payload) {
    return $.ajax({
        type: 'POST',
        url: endpoint,
        data: JSON.stringify(payload),
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json'
        }
    })
}

const profileId = parseInt(localStorage.getItem('profileID'))

Moon.component("resume", {
    template: `<div>
    <h2 class="ui dividing header"> create a new resume </h2>
<div m-if="hasMessage" class="ui message {{messageState}}">
    {{ messageText}}
 </div>
    <form method="post" class="ui large form">
        <div class="field">
            <label>title</label>
            <input type="text" name="title" m-on:change="updateTitle(event)">
        </div>
        <button class="ui fluid large submit button" m-on:click="create(event)"> create </button>
    </form>
</div>
    `,
    data: function () {
        return {
            title: '',
            hasMessage: false,
            messageState: '',
            messageText: '',

        }
    },
    methods: {
        updateTitle(e) {
            this.set('title', e.target.value)
        },
        async create(e) {
            e.preventDefault();
            const title = this.get('title');
            let data;
            let err;
            await send("/resume/new", {
                profileID: profileId,
                title: title
            }).then((res) => {
                data = res;
            }, (xreq, status, error) => {
                err = error;
            })
            if(data){
                this.set('hasMessage',true)
                this.set('messageState','success')
                this.set('messageText','successful created')
            }
        }
    },
})


const app = new Moon({
    el: "#app",
});