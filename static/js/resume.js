
class HTTP {

    static sendJSON(method, url, body) {
        let payload;
        if (body !== null) {
            payload = JSON.stringify(body);
        }
        return $.ajax({
            type: method,
            url: url,
            data: payload,
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json'
            }
        })
    }
}


const profileId = parseInt(localStorage.getItem('profileID'))

Moon.component("resume", {
    template: `<div>
    <h2 class="ui dividing header"> create a new resume </h2>
<div m-if="hasMessage" class="ui transition message {{messageState}}">
    <div class="header">{{ messageText}}</div>
    <i class="close icon" m-on:click="destroyMessage"></i>
 </div>
    <form method="post" class="ui large form">
        <div class="field">
            <label>title</label>
            <input type="text" name="title" m-on:change="updateTitle(event)">
        </div>
        <button class="ui fluid large submit button" m-on:click.prevent="create(event)"> {{action}} </button>
    </form>
</div>
    `,
    data: function () {
        return {
            title: '',
            hasMessage: false,
            messageState: '',
            messageText: '',
            created: false,
            value: {},
            action: 'Create'
        }
    },
    methods: {
        updateTitle(e) {
            this.set('title', e.target.value)
        },
        message(state, text) {
            this.set('hasMessage', true)
            this.set('messageState', state)
            this.set('messageText', text)
        },
        destroyMessage() {
            this.set('hasMessage', false)
        },
        async create(e) {
            const title = this.get('title');
            const isCreated = this.get('created')
            let data;
            let err;
            if (isCreated) {
                const value = this.get('value')
                if (value.title === title) {
                    this.callMethod('message', ['success', 'saved'])
                    return null
                }
                await HTTP.sendJSON('POST', "/resume/update", {
                    id: value.id,
                    profileID: profileId,
                    title: title
                }).then((res) => {
                    data = res;
                }, (error) => {
                    err = error;
                })
            } else {
                await HTTP.sendJSON('POST', "/resume/new", {
                    profileID: profileId,
                    title: title
                }).then((res) => {
                    data = res;
                }, (error) => {
                    err = error;
                })
            }
            if (data) {
                this.callMethod('message', ['success', 'successful created'])
                this.set('created', true)
                this.set('value', data)
                this.set('action', 'Update')
            } else if (err) {
                this.callMethod('message', ['error', 'some fish'])
            }
        }
    },
})


const app = new Moon({
    el: "#app",
});