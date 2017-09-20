<template>
<div m-if="showForm" class="row nice-box">
    <h2 class="ui dividing header"> create a new resume </h2>
    <flash m-if="hasFlash" status="{{flashStatus}}" text="{{flashText}}" ></flash>
    <form method="post" class="ui large {{store.state.resume.loading}} form">
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