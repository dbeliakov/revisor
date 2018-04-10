<template>
  <div id="app">
    <div id="logo" style="margin: auto; text-align: center; font-family: 'Slabo 27px', serif;">
      <h1 style="color: #7777cc; font-size: 45pt;">Revisor</h1>
    </div>
    <div v-if="$auth.ready()">
        <router-view></router-view>
    </div>
    <div v-if="!$auth.ready()">
        Loading ...
    </div>
  </div>
</template>

<script>

export default {
  name: 'App'
}
</script>

<style>
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}

#logo {
  font-family: 'Slabo 27px', serif;
}
</style>
