
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
const state = {};
const bus = new Moon();

Moon.component("resume", {
    template: `<div class="row nice-box">
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
        <button class="ui fluid large primary submit button" m-on:click.prevent="create(event)"> {{action}} </button>
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
    hooks: {
        init() {
            const self = this;
            this.sendMessage = (status, text) => {
                self.set('hasMessage', true)
                self.set('messageState', status)
                self.set('messageText', text)
            };
            this.createdResume = (resume) => {
                self.set('created', true)
                self.set('value', resume)
                self.set('action', 'Update')
            };
            bus.on('createdBasic', async (payload) => {
                const value = self.get('value')
                try {
                    const data = await HTTP.sendJSON('POST', "/resume/update", {
                        id: value.id,
                        profileID: profileId,
                        basicID: payload.id
                    })
                    this.sendMessage('success', 'successful updated')
                    state.resume = data
                } catch (e) {
                    self.sendMessage('error', `failed to update basic details`)
                }
            })
        }
    },
    methods: {
        updateTitle(e) {
            this.set('title', e.target.value)
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
                    this.sendMessage('success', 'saved')
                    return null
                }
                try {
                    const data = await HTTP.sendJSON('POST', "/resume/update", {
                        id: value.id,
                        profileID: profileId,
                        title: title
                    })
                    this.sendMessage('success', 'successful updated')
                    state.resume = data;
                } catch (e) {
                    this.sendMessage('error', 'failed to update resume')
                }
            } else {
                try {
                    const data = await HTTP.sendJSON('POST', "/resume/new", {
                        profileID: profileId,
                        title: title
                    })
                    this.sendMessage('success', 'successful created')
                    this.createdResume(data)
                    state.resume = data;
                    bus.emit('createResume', data)
                } catch (e) {
                    this.sendMessage('error', 'failed to create resume')
                }
            }
        }
    },
})

Moon.component('basic', {
    template: `<div class="row nice-box">
    <h2 class="ui dividing header"> Basic details </h2>
<div m-if="hasMessage" class="ui transition message {{messageState}}">
    <div class="header">{{ messageText}}</div>
    <i class="close icon" m-on:click="destroyMessage"></i>
 </div>
    <form method="post" class="ui large form">
        <div class="field">
            <label>Name</label>
            <input type="text" name="name" m-on:change="updateField(event)">
        </div>
        <div class="field">
            <label>Label</label>
            <input type="text" name="label" m-on:change="updateField(event)">
        </div>
        <!-------
        <div class="field">
            <label>Picture</label>
            <input type="file" name="picture" m-on:change="updateField(event)">
        </div>
        ---->
        <div class="field">
            <label>Email</label>
            <input type="email" name="email" m-on:change="updateField(event)">
        </div>
        <div class="field">
            <label>Phone</label>
            <input type="tel" name="phone" m-on:change="updateField(event)">
        </div>
        <div class="field">
            <label>Summary</label>
            <input type="text" name="summary" m-on:change="updateField(event)">
        </div>
        <button class="ui fluid large primary submit button" m-on:click.prevent="create(event)"> {{action}} </button>
    </form>
</div>`,
    data: function () {
        return {
            fields: {
                name: '',
                label: '',
                picture: '',
                email: '',
                phone: '',
                summary: ''
            },
            hasMessage: false,
            messageState: '',
            messageText: '',
            created: false,
            value: {},
            action: 'Create'
        }
    },
    hooks: {
        init() {
            const self = this;
            this.sendMessage = (status, text) => {
                self.set('hasMessage', true)
                self.set('messageState', status)
                self.set('messageText', text)
            };
        }
    },
    methods: {
        updateField(e) {
            this.set(`fields.${e.target.name}`, e.target.value)
        },
        destroyMessage() {
            this.set('hasMessage', false)
        },
        async create(e) {
            const value = this.get("fields")
            try {
                const data = await HTTP.sendJSON('POST', "/resume/basic", value)
                bus.emit('createdBasic', data)
                this.sendMessage('success', 'successful created basic details')
            } catch (e) {
                this.sendMessage('error', 'failed creating basic details')
                console.log(e)
            }
        }
    }
})

Moon.component('app', {
    template: `<div>
    <resume></resume>
    <basic m-if="showBasic"></basic>
    </div>
    `,
    data: function () {
        return {
            showBasic: false,
        }
    },
    hooks: {
        init() {
            const self = this;
            bus.on("createResume", (payload) => {
                self.set('showBasic', true)
                state.resume = payload;
            })
        }
    }
})

const app = new Moon({
    el: "#app"
});