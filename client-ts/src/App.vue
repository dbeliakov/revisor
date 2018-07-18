<template>
  <div id="app">
    <div v-if="$auth.ready()">
      <div v-if="!$auth.authenticated()" id="logo" style="margin: auto; margin-top: 40px; text-align: center;">
        <h1 style="color: #7777cc; font-size: 45pt; font-family: 'Slabo 27px', serif;">Revisor</h1>
      </div>

      <div v-if="$auth.authenticated()" class="ui secondary pointing menu" style="margin-top: 10px;">
        <span style="margin: 5px 20px; text-align: center; font-family: 'Slabo 27px', serif; font-size: 30pt; color: #7777cc;">Revisor</span>
        <router-link :to="{name: 'OutReviews'}" active-class="active" class="item">
          Исходящие
          <!--<div class="ui green label">1</div>-->
        </router-link>
        <router-link :to="{name: 'InReviews'}" active-class="active" class="item">
          Входящие
          <!--<div class="ui green label">4</div>-->
        </router-link>
        <router-link :to="{name: 'NewReview'}" active-class="active" class="item">
          Новое ревью
        </router-link>
        <div class="right menu">
          <router-link :to="{name: 'Profile'}" active-class="active" class="item">
            {{ $auth.user().first_name }} {{ $auth.user().last_name }} ({{ $auth.user().username }})
          </router-link>
          <a class="ui item" v-on:click="logout()" href="javascript:void(0);">
            Выйти
          </a>
        </div>
      </div>

      <router-view></router-view>
    </div>

    <div v-if="!$auth.ready()">
      Loading...
    </div>
  </div>
</template>

<script lang="ts">
import {Component, Vue, Watch} from 'vue-property-decorator';

@Component
export default class App extends Vue {
  public logout() {
    this.$auth.logout();
    this.$router.push({name: 'login'});
  }
}
</script>

<style lang="scss">
@import '~semantic-ui-css/semantic.min.css';

#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  /*margin-top: 60px;*/
}
</style>
