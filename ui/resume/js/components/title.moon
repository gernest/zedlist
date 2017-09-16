<template>
<div class="row nice-box">
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
        <button class="ui fluid large primary submit button" m-on:click.prevent="create(event)"> {{store.state.resume.action}} </button>
    </form>
</div>
</template>

<script>
const title= require('./title.js');
exports= title()
</script>