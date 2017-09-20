<template>
<div m-if="hasMessage">
<div class="ui transition message {{status}}">
    <div class="header">{{ text}}</div>
    <i class="close icon" m-on:click="destroyMessage"></i>
 </div>
 </div>
</template>
<script>
const flash= require('./flash.js');
exports= flash()
</script>